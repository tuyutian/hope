package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/application/users"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	userEntity "backend/internal/domain/entity/users"
	"backend/internal/domain/repo/jwtauth"
	"backend/internal/infras/config"
	"backend/internal/providers"
	"backend/pkg/crypto/bcrypt"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

// AuthWare auth 中间件
type AuthWare struct {
	userService  *users.UserService
	jwtRepo      jwtauth.JWTRepository
	customCrypto bcrypt.BCrypto
	aesCrypto    bcrypt.BCrypto
	shopifyConf  config.Shopify
}

// CookieClaims cookie 中的登录信息
type CookieClaims struct {
	UserID int64 `json:"userid,string"`
	SubID  int64 `json:"subid,string"`
}

// NewAuthWare JWT 中间件
func NewAuthWare(userService *users.UserService, repos *providers.Repositories, shopifyConf config.Shopify) *AuthWare {
	return &AuthWare{
		userService: userService,
		jwtRepo:     repos.JwtRepo,
		aesCrypto:   repos.AesCrypto,
		shopifyConf: shopifyConf,
	}
}

var (
	errJwtTokenEmpty = errors.New("jwt token is empty")
	errBadRequest    = errors.New("bad request")
	errSignInvalid   = errors.New("sign invalid")
	errUserNoLogin   = errors.New("user has no login")
)

// CheckLogin 验证登录是否合法
// 如果合法将解析 token，并将 claims 写入 context
func (auth *AuthWare) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			claims *jwt.BizClaims
			ctx    = c.Request.Context()
		)
		if claims, err = auth.checkLogin(c); err != nil {
			// record request logger
			logger.Warn(ctx, "failed to auth error", "err", err)
			statusCode := auth.parseError(err)
			// abort
			c.AbortWithStatusJSON(statusCode, gin.H{
				"code":    statusCode,
				"message": http.StatusText(statusCode),
			})
			return
		}
		if strings.HasPrefix(claims.Dest, "https://") {
			c.Header("Content-Security-Policy", "frame-ancestors "+claims.Dest+" https://admin.shopify.com;")
		}
		ctx = context.WithValue(ctx, ctxkeys.BizClaims, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// CheckAdmin 验证是否为超管
func (auth *AuthWare) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims, ok := ctx.Value(ctxkeys.BizClaims).(*jwt.BizClaims)
		if !ok {
			// record request logger
			logger.Warn(ctx, "failed to get claims from context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": http.StatusText(http.StatusForbidden),
			})
			return
		}

		if claims.AdminID <= 0 {
			// record request logger
			logger.Warn(ctx, "user is not admin", "claims", claims)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": http.StatusText(http.StatusForbidden),
			})
			return
		}

		c.Next()
	}
}

// parseError 解析 err 并返回http status code
func (auth *AuthWare) parseError(err error) int {
	// token 过期
	if errors.Is(err, jwt.ErrTokenExpired) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return http.StatusUnauthorized
	}

	// 所有的登录方式都校验不通过
	if errors.Is(err, errUserNoLogin) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, message.ErrInvalidAccount) {
		return http.StatusNotFound
	}

	return http.StatusBadRequest
}

func (auth *AuthWare) checkLogin(c *gin.Context) (*jwt.BizClaims, error) {
	ctx := c.Request.Context()

	claims, err := auth.checkJwt(c)
	if err != nil {
		if auth.skipErr(err) {
			logger.Warn(ctx, "skip login error", "trace_error", err)
			// 登录认证不通过，就直接返回错误
			return nil, errUserNoLogin
		}
		// 登录认证解析失败的错误类型
		return nil, err
	}

	// claims 解析成功
	if claims != nil {
		return claims, nil
	}

	// 以上登录认证都不通过，就直接返回错误
	return nil, errUserNoLogin
}

func (auth *AuthWare) skipErr(err error) bool {
	// 如果是jwt token为空
	if errors.Is(err, errJwtTokenEmpty) || errors.Is(err, errSignInvalid) {
		return true
	}

	return false
}

func (auth *AuthWare) checkJwt(c *gin.Context) (*jwt.BizClaims, error) {
	var (
		token = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		ctx   = c.Request.Context()
	)
	if len(token) == 0 {
		return nil, errJwtTokenEmpty
	}

	// 包含 bcrypt.EncryptedPrefix 前缀，为新版 JWT 生成的 token
	if strings.HasPrefix(token, bcrypt.EncryptedPrefix) {
		// strings.TrimPrefix 底层使用的是切片截取
		token = token[len(bcrypt.EncryptedPrefix):]
	}
	// token 校验失败
	claims, err := auth.jwtRepo.Verify(ctx, token)
	if err != nil {
		return nil, err
	}
	// 去掉 shopify 域名前的 https
	if claims.Dest != "" && strings.HasPrefix(claims.Dest, "https://") {
		claims.Dest = strings.TrimPrefix(claims.Dest, "https://")
	}
	// 验证claims是否合法
	err = auth.checkClaims(ctx, token, claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (auth *AuthWare) checkClaims(ctx context.Context, token string, claims *jwt.BizClaims) error {
	// 主账号登录
	if claims.UserID > 0 {
		if err := auth.checkUser(ctx, claims); err != nil {
			return err
		}
		return nil
	}

	// 超管账户登录
	if claims.AdminID > 0 {
		if err := auth.checkManage(ctx, claims); err != nil {
			return err
		}

		return nil
	}

	// shopify登录
	if len(claims.Dest) > 0 {
		if err := auth.checkShop(ctx, token, claims); err != nil {
			return err
		}
		return nil
	}

	return message.ErrInvalidAccount
}

var accessTokenRelPath = "admin/oauth/access_token"

func (auth *AuthWare) checkUser(ctx context.Context, claims *jwt.BizClaims) error {
	// 检查主账号是否存在
	_, err := auth.userService.GetLoginUserFromID(ctx, claims.UserID)
	return err
}

func (auth *AuthWare) checkShop(ctx context.Context, token string, claims *jwt.BizClaims) error {
	// 检查主账号是否存在
	user, err := auth.userService.GetLoginUserFromShop(ctx, claims.Dest)
	if user != nil {
		claims.UserID = user.ID
		return nil
	}
	if err != nil {
		return err
	}
	url := "https://" + claims.Dest + accessTokenRelPath
	data := struct {
		ClientId           string `json:"client_id"`
		ClientSecret       string `json:"client_secret"`
		GrantType          string `json:"grant_type"`
		SubjectToken       string `json:"subject_token"`
		SubjectTokenType   string `json:"subject_token_type"`
		RequestedTokenType string `json:"requested_token_type"`
	}{
		ClientId:           auth.shopifyConf.AppKey,
		ClientSecret:       auth.shopifyConf.AppSecret,
		GrantType:          "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectToken:       token,
		SubjectTokenType:   "urn:ietf:params:oauth:token-type:id_token",
		RequestedTokenType: "urn:shopify:params:oauth:token-type:offline-access-token",
	}
	client := utils.NewHTTPClient()
	sessionToken := new(shopifyEntity.Token)
	err = client.PostJSON(ctx, url, &data, &sessionToken)
	if err != nil {
		return err
	}
	user, err = auth.sessionStore(ctx, sessionToken)
	if err != nil {
		return err
	}
	if user != nil {
		claims.UserID = user.ID
		return nil
	}
	return message.ErrInvalidAccount
}

// sessionStore   session 授权部分
func (auth *AuthWare) sessionStore(ctx context.Context, token *shopifyEntity.Token) (*userEntity.User, error) {
	user, err := auth.userService.AuthFromSession(ctx, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (auth *AuthWare) checkManage(ctx context.Context, claims *jwt.BizClaims) error {
	// 检查超管用户是否存在
	if _, err := auth.userService.GetLoginAdminFromID(ctx, claims.AdminID); err != nil {
		return err
	}

	return nil
}
