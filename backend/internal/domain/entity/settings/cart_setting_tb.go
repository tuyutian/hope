package settings

// UserCartSetting  保险用户基础配置表
type UserCartSetting struct {
	Id                int64   `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	UserID            int64   `xorm:"'user_id' bigint(20) notnull comment('用户id')" json:"user_id"`
	PlanTitle         string  `xorm:"'plan_title' varchar(100) notnull default '' comment('保险标题(内部)')" json:"plan_title"`
	AddonTitle        string  `xorm:"'addon_title' varchar(100) notnull default '' comment('保险标题')" json:"addon_title"`
	EnabledDesc       string  `xorm:"'enabled_desc' varchar(200) notnull default '' comment('按钮打开文案')" json:"enabled_desc"`
	DisabledDesc      string  `xorm:"'disabled_desc' varchar(200) notnull default '' comment('按钮关闭文案')" json:"disabled_desc"`
	FootText          string  `xorm:"'foot_text' varchar(100) notnull default '' comment('保险底部')" json:"foot_text"`
	FootUrl           string  `xorm:"'foot_url' varchar(255) notnull default '' comment('保险跳转')" json:"foot_url"`
	InColor           string  `xorm:"'in_color' varchar(50) notnull default '' comment('打开颜色')" json:"in_color"`
	OutColor          string  `xorm:"'out_color' varchar(50) notnull default '' comment('关闭颜色')" json:"out_color"`
	OtherMoney        float64 `xorm:"'other_money' decimal(12,2) notnull default 0.00 comment('其他金额')" json:"other_money"`
	ShowCart          int     `xorm:"'show_cart' tinyint(1) default 0 notnull comment('购物车状态 0 关闭 1 打开')" json:"show_cart"`
	ShowCartIcon      int     `xorm:"'show_cart_icon' tinyint(1) default 0 notnull comment('购物车图标 0 关闭 1 打开')" json:"show_cart_icon"`
	IconUrl           string  `xorm:"'icon_url' text comment('选中url(json)')" json:"icon_url"`
	SelectButton      int     `xorm:"'select_button' tinyint(1) default 0 notnull comment('购物车图标 0 滑动 1 勾选')" json:"select_button"`
	InCollection      int     `xorm:"'in_collection' tinyint(1) default 0 not null comment('是否启用集合筛选 0 关闭 1 打开')" json:"in_collection"`
	ProductCollection string  `xorm:"'product_collection' varchar(100) notnull default '' comment('产品选中集合')" json:"product_collection"`
	PricingType       int     `xorm:"'pricing_type' tinyint(1) default 0 notnull comment('购物车图标 0 金额 1百分比')" json:"pricing_type"`
	PricingRule       int     `xorm:"'pricing_rule' tinyint(1) default 0 notnull comment('金额计算方式 0 统一设置 1单独设置')" json:"pricing_rule"`
	PricingSelect     string  `xorm:"'pricing_select' text comment('金额计算范围')" json:"pricing_select"`
	TiersSelect       string  `xorm:"'tiers_select' text comment('百分比计算范围')" json:"tiers_select"`
	AllTiersSet       float64 `xorm:"'all_tiers_set' decimal(12,2) notnull default 0.00 comment('所有订单适用固定百分比') " json:"all_tiers_set"`
	AllPriceSet       float64 `xorm:"'all_price_set' decimal(12,2) notnull default 0.00 comment('所有订单适用固定金额') " json:"all_price_set"`
	FulfillmentRule   int     `xorm:"'fulfillment_rule' tinyint(1) default 0 notnull comment('在订单处于哪个发货阶段才计算保险佣金(0,1,2 分别代表第一个发货完成，全都发货完成，付费后就算)')" json:"fulfillment_rule"`
	CSS               string  `xorm:"'css' text comment('css样式自定义')" json:"css"`
	CreateTime        int64   `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime        int64   `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
