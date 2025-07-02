package settings

type SettingConfigReq struct {
	UserID              int64            `json:"user_id"`
	PlanTitle           string           `json:"planTitle" binding:"required"`
	IconVisibility      int              `json:"iconVisibility,omitempty" binding:"required,oneof=0 1"`
	InsuranceVisibility int              `json:"insuranceVisibility,omitempty" binding:"required,oneof=0 1"`
	SelectButton        int              `json:"selectButton,omitempty" binding:"required,oneof=0 1"`
	AddonTitle          string           `json:"addonTitle,omitempty" binding:"required,max=50"`
	EnabledDescription  string           `json:"enabledDescription" binding:"required,max=200"`
	DisabledDescription string           `json:"disabledDescription" binding:"required,max=200"`
	FooterText          string           `json:"footerText"`
	FooterUrl           string           `json:"footerUrl"`
	OptInColor          string           `json:"optInColor" binding:"required,max=50"`
	OptOutColor         string           `json:"optOutColor" binding:"required,max=50"`
	PricingType         int              `json:"pricingType,omitempty" binding:"required,oneof=0 1"`
	PriceSelect         []PriceSelectReq `json:"priceSelect" binding:"required,dive"`
	TiersSelect         []TierSelectReq  `json:"tiersSelect" binding:"required,dive"`
	RestValuePrice      string           `json:"restValuePrice" binding:"required"`
	ProductTypeInput    string           `json:"productTypeInput"`
	SelectedCollections []string         `json:"selectedCollections"`
	Icons               []IconReq        `json:"icons" binding:"required,dive"`
}

type CartSettingVO struct {
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
	// 其他金额
	OtherMoney float64 `json:"other_money"`
	// 购物车状态 0 关闭 1 打开
	ShowCart int `json:"show_cart,omitempty"`
	// 购物车图标 0 关闭 1 打开
	ShowCartIcon int `json:"show_cart_icon,omitempty"`
	// 购物车图标 0 滑动 1 勾选
	SelectButton int `json:"select_button,omitempty"`
	// 产品type
	ProductType string `json:"product_type"`
	// 产品选中集合
	ProductCollection []string         `json:"product_collection"`
	PricingType       int              `json:"pricing_type,omitempty"`
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

type CartPublicVO struct {
	AddonTitle   string           `json:"addon_title"`              // 保险标题
	EnabledDesc  string           `json:"enabled_desc"`             // 按钮打开文案
	DisabledDesc string           `json:"disabled_desc"`            // 按钮关闭文案
	FootText     string           `json:"foot_text"`                // 保险底部
	FootURL      string           `json:"foot_url"`                 // 保险跳转
	InColor      string           `json:"in_color"`                 // 打开颜色
	OutColor     string           `json:"out_color"`                // 关闭颜色
	OtherMoney   float64          `json:"other_money"`              // 其他金额
	ShowCartIcon int              `json:"show_cart_icon,omitempty"` // 购物车图标 0 关闭 1 打开
	SelectButton int              `json:"select_button,omitempty"`  // 购物车图标 0 滑动 1 勾选
	PriceSelect  []PriceSelectReq `json:"price_select"`
	TiersSelect  []TierSelectReq  `json:"tiers_select"`
	Icon         string           `json:"icon"`
	Variants     map[string]int64 `json:"variants"`
	ProductId    int64            `json:"product_id"`
	MoneyFormat  string           `json:"money_format"`
	PricingType  int              `json:"pricing_type"`
}
