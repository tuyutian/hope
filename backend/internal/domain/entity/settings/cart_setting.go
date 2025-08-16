package settings

type CollectionItem struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

type SettingConfigReq struct {
	UserID               int64            `json:"user_id,omitempty"`
	PlanTitle            string           `json:"planTitle" binding:"required"`
	IconVisibility       int              `json:"iconVisibility" binding:"oneof=0 1"`
	ProtectifyVisibility int              `json:"protectifyVisibility" binding:"oneof=0 1"`
	SelectButton         int              `json:"selectButton" binding:"oneof=0 1"`
	AddonTitle           string           `json:"addonTitle" binding:"required,max=50"`
	EnabledDescription   string           `json:"enabledDescription" binding:"required,max=200"`
	DisabledDescription  string           `json:"disabledDescription" binding:"required,max=200"`
	FooterText           string           `json:"footerText"`
	FooterUrl            string           `json:"footerUrl"`
	OptInColor           string           `json:"optInColor" binding:"required,max=50"`
	OptOutColor          string           `json:"optOutColor" binding:"required,max=50"`
	PricingType          int              `json:"pricingType" binding:"oneof=0 1"`
	PricingRule          int              `json:"pricingRule" binding:"oneof=0 1"`
	PriceSelect          []PriceSelectReq `json:"priceSelect" binding:"required,dive"`
	TiersSelect          []TierSelectReq  `json:"tiersSelect" binding:"required,dive"`
	OutPrice             string           `json:"outPrice" binding:"required"`
	OutTier              string           `json:"outTier" binding:"required"`
	AllTiers             string           `json:"allTiers"`
	AllPrice             string           `json:"allPrice"`
	InCollection         bool             `json:"onlyInCollection"`
	SelectedCollections  []CollectionItem `json:"selectedCollections"`
	Icons                []IconReq        `json:"icons" binding:"required,dive"`
	FulfillmentRule      int              `json:"fulfillmentRule" binding:"oneof=0 1 2"`
	CSS                  string           `json:"css"`
}

type CartSettingData struct {
	// 保险标题(内部)
	PlanTitle string `json:"plan_title"`
	// 保险标题
	AddonTitle string `json:"addon_title" binding:"required,oneof=0 1"`
	// 按钮打开文案
	EnabledDesc string `json:"enabled_desc"`
	// 按钮关闭文案
	DisabledDesc string `json:"disabled_desc"`
	// 保险底部
	FootText string `json:"foot_text"`
	// 保险跳转
	FootURL string `json:"foot_url"`
	// 打开颜色
	InColor string `json:"in_color"`
	// 关闭颜色
	OutColor string `json:"out_color"`
	// 购物车状态 0 关闭 1 打开
	ShowCart int `json:"show_cart"`
	// 购物车图标 0 关闭 1 打开
	ShowCartIcon int `json:"show_cart_icon"`
	// 购物车图标 0 滑动 1 勾选
	SelectButton int `json:"select_button"`
	// 产品type
	ProductType     string  `json:"product_type"`
	AllTiers        float64 `json:"all_tiers"`
	AllPrice        float64 `json:"all_price"`
	OutPrice        float64 `json:"out_price"`
	OutTier         float64 `json:"out_tier"`
	FulfillmentRule int     `json:"fulfillment_rule"`
	CSS             string  `json:"css"`
	// 产品选中集合
	ProductCollection []CollectionItem `json:"product_collection"`
	InCollection      bool             `json:"in_collection"`
	PricingType       int              `json:"pricing_type"`
	PricingRule       int              `json:"pricing_rule"`
	PriceSelect       []PriceSelectReq `json:"price_select"`
	TiersSelect       []TierSelectReq  `json:"tiers_select"`
	Icons             []IconReq        `json:"icons"`
}

type PriceSelectReq struct {
	Min   string `json:"min" binding:"required"`
	Max   string `json:"max" binding:"required"`
	Price string `json:"price" binding:"required"`
}

type TierSelectReq struct {
	Min        string `json:"min" binding:"required"`
	Max        string `json:"max" binding:"required"`
	Percentage string `json:"percentage" binding:"required"`
}

type IconReq struct {
	Id       int64  `json:"id" binding:"required"`
	Src      string `json:"src" binding:"required"`
	Selected bool   `json:"selected"`
}

type CartPublicData struct {
	AddonTitle     string           `json:"addon_title"`              // 保险标题
	EnabledDesc    string           `json:"enabled_desc"`             // 按钮打开文案
	DisabledDesc   string           `json:"disabled_desc"`            // 按钮关闭文案
	FootText       string           `json:"foot_text"`                // 保险底部
	FootURL        string           `json:"foot_url"`                 // 保险跳转
	InColor        string           `json:"in_color"`                 // 打开颜色
	OutColor       string           `json:"out_color"`                // 关闭颜色
	ShowCartIcon   int              `json:"show_cart_icon,omitempty"` // 购物车图标 0 关闭 1 打开
	SelectButton   int              `json:"select_button,omitempty"`  // 购物车图标 0 滑动 1 勾选
	PriceSelect    []PriceSelectReq `json:"price_select"`
	TiersSelect    []TierSelectReq  `json:"tiers_select"`
	OutSelectPrice float64          `json:"out_select_price"`
	OutSelectTier  float64          `json:"out_select_tier"`
	AllTiersSet    float64          `json:"all_tiers_set"`
	AllPriceSet    float64          `json:"all_price_set"`
	Icon           string           `json:"icon"`
	Variants       map[string]int64 `json:"variants"`
	ProductId      int64            `json:"product_id"`
	MoneyFormat    string           `json:"money_format"`
	PricingType    int              `json:"pricing_type"`
}
