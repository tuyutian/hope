package entity

import (
	"time"
)

// UserCartSetting  保险用户基础配置表
type UserCartSetting struct {
	Id                int     `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid               int     `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	PlanTitle         string  `xorm:"'plan_title' varchar(100) comment('保险标题(内部)')" json:"plan_title"`
	AddonTitle        string  `xorm:"'addon_title' varchar(100) comment('保险标题')" json:"addon_title"`
	EnabledDesc       string  `xorm:"'enabled_desc' varchar(200) comment('按钮打开文案')" json:"enabled_desc"`
	DisabledDesc      string  `xorm:"'disabled_desc' varchar(200) comment('按钮关闭文案')" json:"disabled_desc"`
	FootText          string  `xorm:"'foot_text' varchar(100) comment('保险底部')" json:"foot_text"`
	FootUrl           string  `xorm:"'foot_url' varchar(100) comment('保险跳转')" json:"foot_url"`
	InColor           string  `xorm:"'in_color' varchar(50) comment('打开颜色')" json:"in_color"`
	OutColor          string  `xorm:"'out_color' varchar(50) comment('关闭颜色')" json:"out_color"`
	OtherMoney        float64 `xorm:"'other_money' decimal(10,2) default 0.00 comment('其他金额')" json:"other_money"`
	ShowCart          *int    `xorm:"'show_cart' tinyint(1) default 0 notnull comment('购物车状态 0 关闭 1 打开')" json:"show_cart"`
	ShowCartIcon      *int    `xorm:"'show_cart_icon' tinyint(1) default 0 notnull comment('购物车图标 0 关闭 1 打开')" json:"show_cart_icon"`
	IconUrl           string  `xorm:"'icon_url' text comment('选中url(json)')" json:"icon_url"`
	SelectButton      *int    `xorm:"'select_button' tinyint(1) default 0 notnull comment('购物车图标 0 滑动 1 勾选')" json:"select_button"`
	ProductType       string  `xorm:"'product_type' varchar(100) comment('产品type')" json:"product_type"`
	ProductCollection string  `xorm:"'product_collection' varchar(100) comment('产品选中集合')" json:"product_collection"`
	PricingType       *int    `xorm:"'pricing_type' tinyint(1) default 0 notnull comment('购物车图标 0 金额 1百分比')" json:"pricing_type"`
	PricingSelect     string  `xorm:"'pricing_select' text comment('金额计算范围')" json:"pricing_select"`
	TiersSelect       string  `xorm:"'tiers_select' text comment('百分比计算范围')" json:"tiers_select"`
	CreateTime        int64   `xorm:"'create_time' int(11) notnull comment('创建时间')" json:"create_time"`
	UpdateTime        int64   `xorm:"'update_time' int(11) notnull comment('修改时间')" json:"update_time"`
}

//func (m *InUserSetting) TableName() string {
//	return "in_user_setting"
//}

func (s *UserCartSetting) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	s.CreateTime = now
	s.UpdateTime = now
}

func (s *UserCartSetting) BeforeUpdate() {
	now := time.Now().Unix()
	s.UpdateTime = now
}
