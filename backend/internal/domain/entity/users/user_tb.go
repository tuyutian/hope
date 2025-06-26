package users

import (
	"time"
)

type User struct {
	Id              int    `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Name            string `xorm:"'name' varchar(255) notnull comment('shopify-name')" json:"name"`
	Domain          string `xorm:"'domain' varchar(255) comment('my shopify Domain网站域名')" json:"domain"`
	PlanDisplayName string `xorm:"'plan_display_name' varchar(40) comment('shopify套餐版本')" json:"plan_display_name"`
	AccessToken     string `xorm:"'access_token' varchar(255) comment('shopify-token')" json:"access_token"`
	UserToken       string `xorm:"'user_token' varchar(255) comment('用户token')" json:"user_token"`
	Pwd             string `xorm:"'pwd' varchar(255) comment('密码')" json:"pwd"`
	Level           int    `xorm:"'level' tinyint(2) default 0 notnull comment('app内部套餐等级')" json:"level"`
	Email           string `xorm:"'email' varchar(100) comment('邮箱')" json:"email"`
	CountryName     string `xorm:"'country_name' varchar(25) comment('国家名称')" json:"country_name"`
	CountryCode     string `xorm:"'country_code' varchar(25) comment('国家简码')" json:"country_code"`
	City            string `xorm:"'city' varchar(25) comment('城市')" json:"city"`
	TrialTime       int64  `xorm:"'trial_time' int(11) default 0 notnull comment('试用时间')" json:"trial_time"`
	CurrencyCode    string `xorm:"'currency_code' varchar(255) comment('货币简码')" json:"currency_code"`
	MoneyFormat     string `xorm:"'money_format' varchar(20) comment('货币单位符号')" json:"money_format"`
	LastLogin       int64  `xorm:"'last_login' int(11) default 0 notnull" json:"last_login"`
	IsDel           int    `xorm:"'is_del' tinyint(1) default 1 notnull comment('删除状态 1 正常 2 卸载 3关店')" json:"is_del"`
	PublishId       string `xorm:"'publish_id' varchar(100) comment('店铺publish_id')" json:"publish_id"`
	Steps           string `xorm:"'steps' varchar(255) comment('新手引导')" json:"steps"`
	Collection      string `xorm:"'collection' text comment('用户集合列表')" json:"collection"`
	UnInstallTime   int64  `xorm:"'uninstall_time' int(11) default 0 notnull comment('卸载时间')" json:"uninstall_time"`
	CreateTime      int64  `xorm:"'create_time' int(11) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime      int64  `xorm:"'update_time' int(11) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

//func (m *InUser) TableName() string {
//	return "user"
//}

func (u *User) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	u.CreateTime = now
	u.UpdateTime = now
}

func (u *User) BeforeUpdate() {
	now := time.Now().Unix()
	u.UpdateTime = now
}
