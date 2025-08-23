package settings

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/domain/entity/apps"
	cartEntity "backend/internal/domain/entity/settings"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	appRepo "backend/internal/domain/repo/apps"
	cartSettingRepo "backend/internal/domain/repo/carts"
	"backend/internal/domain/repo/products"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type CartSettingService struct {
	cartSettingRepo  cartSettingRepo.CartSettingRepository
	userRepo         userRepo.UserRepository
	variantRepo      products.VariantRepository
	productRepo      products.ProductRepository
	subscriptionRepo userRepo.UserSubscriptionRepository
	shopGraphqlRepo  shopifyRepo.ShopGraphqlRepository
	appAuthRepo      appRepo.AppAuthRepository
}

func NewCartSettingService(repos *providers.Repositories) *CartSettingService {
	return &CartSettingService{
		cartSettingRepo:  repos.CartSettingRepo,
		userRepo:         repos.UserRepo,
		variantRepo:      repos.VariantRepo,
		productRepo:      repos.ProductRepo,
		subscriptionRepo: repos.UserSubscriptionRepo,
		shopGraphqlRepo:  repos.ShopGraphqlRepo,
		appAuthRepo:      repos.AppAuthRepo,
	}
}

func (s *CartSettingService) GetCart(ctx context.Context, uid int64) (cartEntity.CartSettingData, error) {
	// 查询购物车设置
	cartSetting, err := s.cartSettingRepo.First(ctx, uid)
	if err != nil {
		logger.Error(ctx, "get-cart-db异常", "Err:", err.Error())
		return cartEntity.CartSettingData{}, err
	}

	// 如果没有找到相关的购物车设置，返回空结构体
	if cartSetting == nil {
		return cartEntity.CartSettingData{}, nil
	}

	// 将 ProductCollection 解析为数组，默认空数组
	var collectionArr []cartEntity.CollectionItem
	if cartSetting.ProductCollection != "" {
		err = json.Unmarshal([]byte(cartSetting.ProductCollection), &collectionArr)
		if err != nil {
			logger.Error(ctx, "Unmarshal collection fail:"+err.Error())
		}
	}

	// 解析 Icons
	var icons []cartEntity.IconReq
	if err := json.Unmarshal([]byte(cartSetting.IconUrl), &icons); err != nil {
		logger.Error(ctx, "get-cart 解析 IconUrl 失败", "Err:", err.Error())
		return cartEntity.CartSettingData{}, fmt.Errorf("解析 IconUrl 失败: %w", err)
	}

	// 解析 PricingSelect
	var prices []cartEntity.PriceSelectReq
	if err := json.Unmarshal([]byte(cartSetting.PricingSelect), &prices); err != nil {
		logger.Error(ctx, "get-cart 解析 PricingSelect 失败", "Err:", err.Error())
		return cartEntity.CartSettingData{}, fmt.Errorf("解析 PricingSelect 失败: %w", err)
	}

	// 解析 TiersSelect
	var tiers []cartEntity.TierSelectReq
	if err := json.Unmarshal([]byte(cartSetting.TiersSelect), &tiers); err != nil {
		logger.Error(ctx, "get-cart 解析 TiersSelect 失败", "Err:", err.Error())
		return cartEntity.CartSettingData{}, fmt.Errorf("解析 TiersSelect 失败: %w", err)
	}
	var inCollection bool
	if cartSetting.InCollection != 0 {
		inCollection = true
	} else {
		inCollection = false
	}
	// 返回购物车设置结构体
	return cartEntity.CartSettingData{
		PlanTitle:         cartSetting.PlanTitle,
		AddonTitle:        cartSetting.AddonTitle,
		EnabledDesc:       cartSetting.EnabledDesc,
		DisabledDesc:      cartSetting.DisabledDesc,
		FootText:          cartSetting.FootText,
		FootURL:           cartSetting.FootUrl,
		InColor:           cartSetting.InColor,
		OutColor:          cartSetting.OutColor,
		ShowCart:          cartSetting.ShowCart,
		ShowCartIcon:      cartSetting.ShowCartIcon,
		Icons:             icons,
		SelectButton:      cartSetting.SelectButton,
		InCollection:      inCollection,
		ProductCollection: collectionArr,
		PriceSelect:       prices,
		TiersSelect:       tiers,
		PricingType:       cartSetting.PricingType,
		PricingRule:       cartSetting.PricingRule,
		OutPrice:          cartSetting.OutSelectPrice,
		OutTier:           cartSetting.OutSelectTier,
		AllPrice:          cartSetting.AllPriceSet,
		AllTiers:          cartSetting.AllTiersSet,
		FulfillmentRule:   cartSetting.FulfillmentRule,
		CSS:               cartSetting.CSS,
	}, nil
}

// SetCartSetting  打开设置
func (s *CartSettingService) SetCartSetting(ctx context.Context, req cartEntity.SettingConfigReq) error {
	// 操作购物车设置
	cartSetting, err := s.cartSettingRepo.First(ctx, req.UserID)
	if err != nil {
		logger.Error(ctx, "set-cart-db异常", "Err:", err.Error())
		return err
	}

	// 转成 JSON
	jsonIcon, err := json.Marshal(req.Icons)

	if err != nil {
		logger.Error(ctx, "set-cart 购物车图片json异常", "Err:", err.Error())
		return err
	}
	iconStr := string(jsonIcon)

	// 转成 JSON
	jsonPrice, err := json.Marshal(req.PriceSelect)

	if err != nil {
		logger.Error(ctx, "set-cart 价格json异常", "Err:", err.Error())
		return err
	}

	priceStr := string(jsonPrice)

	// 转成 JSON
	jsonTiers, err := json.Marshal(req.TiersSelect)

	if err != nil {
		logger.Error(ctx, "set-cart 百分比json异常", "Err:", err.Error())
		return err
	}

	tiersStr := string(jsonTiers)
	productCollection, err := json.Marshal(req.SelectedCollections)
	if err != nil {
		return err
	}
	var inCollection int
	if req.InCollection {
		inCollection = 1
	} else {
		inCollection = 0
	}
	subscribe, _ := s.subscriptionRepo.GetActiveSubscription(ctx, req.UserID)
	if subscribe == nil {
		req.ProtectifyVisibility = 0
	}
	userCartSetting := cartEntity.UserCartSetting{
		PlanTitle:         req.PlanTitle,
		AddonTitle:        req.AddonTitle,
		EnabledDesc:       req.EnabledDescription,
		DisabledDesc:      req.DisabledDescription,
		FootText:          req.FooterText,
		FootUrl:           req.FooterUrl,
		InColor:           req.OptInColor,
		OutColor:          req.OptOutColor,
		ShowCart:          req.ProtectifyVisibility,
		ShowCartIcon:      req.IconVisibility,
		IconUrl:           iconStr,
		SelectButton:      req.SelectButton,
		InCollection:      inCollection,
		ProductCollection: string(productCollection),
		PricingSelect:     priceStr,
		TiersSelect:       tiersStr,
		PricingType:       req.PricingType,
		PricingRule:       req.PricingRule,
		OutSelectPrice:    utils.ParseMoneyFloat(req.OutPrice),
		OutSelectTier:     utils.ParseMoneyFloat(req.OutTier),
		AllPriceSet:       utils.ParseMoneyFloat(req.AllPrice),
		AllTiersSet:       utils.ParseMoneyFloat(req.AllTiers),
		FulfillmentRule:   req.FulfillmentRule,
		CSS:               req.CSS,
	}
	needOpenCartPlugin := 0
	if cartSetting == nil {
		// 创建购物车
		userCartSetting.UserID = req.UserID
		_, err = s.cartSettingRepo.Create(ctx, &userCartSetting)
		if userCartSetting.ShowCart == 1 {
			needOpenCartPlugin = 1
		}
	} else {
		// 更新购物车
		userCartSetting.Id = cartSetting.Id
		err = s.cartSettingRepo.Update(ctx, &userCartSetting)
		if cartSetting.ShowCart == 0 && userCartSetting.ShowCart == 1 {
			needOpenCartPlugin = 1
		}
		if cartSetting.ShowCart == 1 && userCartSetting.ShowCart == 0 {
			needOpenCartPlugin = 2
		}
	}

	if err != nil {
		logger.Error(ctx, "set-cart-db(2)异常", "Err:", err.Error())
		return err
	}
	if needOpenCartPlugin > 0 {
		// When needOpenCartPlugin == 1, enable cart; when == 2, disable cart via Shopify app metafield
		client := ctx.Value(ctxkeys.ShopifyGraphqlClient).(*shopify_graphql.GraphqlClient)
		appData := ctx.Value(ctxkeys.AppData).(*apps.AppData)
		// Get current app installation to obtain ownerId for app metafields
		s.shopGraphqlRepo.WithClient(client)
		appAuth, err := s.appAuthRepo.GetByUserAndApp(ctx, req.UserID, appData.AppID)
		if err != nil {
			logger.Error(ctx, "appAuth fetch fail:"+err.Error())
			return err
		}
		if appAuth == nil || appAuth.InstallationId == 0 {
			logger.Error(ctx, "set-cart app安装ID为空")
			return fmt.Errorf("app安装ID为空")
		}

		// Determine value to set based on needOpenCartPlugin
		cartEnable := "false"
		if needOpenCartPlugin == 1 {
			cartEnable = "true"
		}

		_, err = s.shopGraphqlRepo.MetafieldSet(ctx, fmt.Sprintf("gid://shopify/AppInstallation/%d", appAuth.InstallationId), shopifyEntity.MetafieldConditionalNs, shopifyEntity.MetafieldTypeBoolean, "cart_enable", cartEnable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CartSettingService) GetPublicCart(ctx context.Context, appId string, shop string) (*cartEntity.CartPublicData, error) {
	// 获取uid
	user, err := s.userRepo.FirstByShop(ctx, appId, shop)

	if err != nil {
		logger.Error(ctx, "public-cart db异常", "Err:", err.Error())
		return nil, err
	}

	if user == nil || user.IsDel > 0 {
		logger.Error(ctx, "public-cart 用户不存在或卸载", "shop:", shop)
		return nil, fmt.Errorf("user not found")
	}

	// 查询购物车设置
	cartSetting, err := s.cartSettingRepo.First(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "public-cart db异常", "Err:", err.Error())
		return nil, err
	}

	// 如果没有找到相关的购物车设置，返回空结构体
	if cartSetting == nil || cartSetting.ShowCart == 0 {
		logger.Warn(ctx, "public-cart 配置没打开", "uid:", user.ID, "shop:", shop)
		return nil, nil
	}

	// 解析 Icons
	var icons []cartEntity.IconReq
	if err := json.Unmarshal([]byte(cartSetting.IconUrl), &icons); err != nil {
		logger.Error(ctx, "public-cart 解析 IconUrl 失败", "Err:", err.Error())
		return nil, fmt.Errorf("解析 IconUrl 失败: %w", err)
	}

	var iconSelect string
	for _, icon := range icons {
		if icon.Selected == true {
			iconSelect = icon.Src
			break
		}
	}
	var prices []cartEntity.PriceSelectReq
	var tiers []cartEntity.TierSelectReq
	if cartSetting.PricingType == 0 {
		// 解析 PricingSelect
		if err := json.Unmarshal([]byte(cartSetting.PricingSelect), &prices); err != nil {
			logger.Error(ctx, "get-cart 解析 PricingSelect 失败", "Err:", err.Error())
			return nil, fmt.Errorf("解析 PricingSelect 失败: %w", err)
		}
	} else {
		// 解析 TiersSelect
		if err := json.Unmarshal([]byte(cartSetting.TiersSelect), &tiers); err != nil {
			logger.Error(ctx, "get public 解析 TiersSelect 失败", "Err:", err.Error())
			return nil, fmt.Errorf("解析 TiersSelect 失败: %w", err)
		}
	}

	// 获取shopify 变体列表
	variants, productID, err := s.variantRepo.GetVariantConfig(ctx, user.ID)

	if err != nil {
		logger.Error(ctx, "public-cart db(2)异常", "Err:", err.Error())
		return nil, err
	}

	if len(variants) == 0 {
		logger.Warn(ctx, "public-cart 无产品数据", "uid:", user.ID, "shop:", shop)
		return nil, fmt.Errorf("public-cart 无产品数据") // 返回空 map，表示没有数据
	}

	// 返回购物车设置结构体
	return &cartEntity.CartPublicData{
		AddonTitle:     cartSetting.AddonTitle,
		EnabledDesc:    cartSetting.EnabledDesc,
		DisabledDesc:   cartSetting.DisabledDesc,
		FootText:       cartSetting.FootText,
		FootURL:        cartSetting.FootUrl,
		InColor:        cartSetting.InColor,
		OutColor:       cartSetting.OutColor,
		ShowCartIcon:   cartSetting.ShowCartIcon,
		Icon:           iconSelect,
		SelectButton:   cartSetting.SelectButton,
		PriceSelect:    prices,
		TiersSelect:    tiers,
		OutSelectPrice: cartSetting.OutSelectPrice,
		OutSelectTier:  cartSetting.OutSelectTier,
		AllPriceSet:    cartSetting.AllPriceSet,
		AllTiersSet:    cartSetting.AllTiersSet,
		Variants:       variants,
		ProductId:      productID,
		MoneyFormat:    user.MoneyFormat,
		PricingType:    cartSetting.PricingType,
	}, nil
}
