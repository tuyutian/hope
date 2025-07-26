package apps

const (
	UserAppAuthTable = "user_app_auth"
)

// UserAppAuth 用户对应用授权表
type UserAppAuth struct {
	Id             int64  `xorm:"pk autoincr 'id' comment('ID')"`
	UserId         int64  `xorm:"notnull 'user_id' comment('用户ID')"`
	Shop           string `xorm:"notnull varchar(100) default '' 'shop' comment('my shopify Domain网站域名')"`
	AppId          string `xorm:"notnull varchar(50) 'app_id' comment('App标识')"`
	AuthToken      string `xorm:"varchar(255) default '' 'auth_token' comment('授权token')"`
	RefreshToken   string `xorm:"varchar(255) default '' 'refresh_token' comment('刷新token')"`
	TokenExpiresAt int64  `xorm:"default 0 'token_expires_at' comment('token过期时间')"`
	Scopes         string `xorm:"notnull text 'scopes' comment('授权域')"`
	Status         int8   `xorm:"notnull tinyint(1) default 1 'status' comment('状态 1:有效 0:已撤销')"`
	CreateTime     int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime     int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

// TableName 设置 UserAppAuth 对应的表名
func (u *UserAppAuth) TableName() string {
	return UserAppAuthTable
}
