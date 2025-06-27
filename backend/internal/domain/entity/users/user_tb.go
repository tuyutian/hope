package users

type User struct {
	Id              int64  `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	Name            string `xorm:"'name' varchar(255) notnull comment('shopify-name')" json:"name"`
	Domain          string `xorm:"'domain' varchar(255) notnull default '' comment('my shopify Domain网站域名')" json:"domain"`
	PlanDisplayName string `xorm:"'plan_display_name' varchar(40) notnull default '' comment('shopify套餐版本')" json:"plan_display_name"`
	AccessToken     string `xorm:"'access_token' varchar(512) notnull default '' comment('shopify-token')" json:"access_token"`
	UserToken       string `xorm:"'user_token' varchar(255) notnull default '' comment('用户token')" json:"user_token"`
	Pwd             string `xorm:"'pwd' varchar(255) notnull default '' comment('密码')" json:"pwd"`
	Level           int    `xorm:"'level' tinyint(2) default 0 notnull comment('app内部套餐等级')" json:"level"`
	Email           string `xorm:"'email' varchar(100) notnull default '' comment('邮箱')" json:"email"`
	CountryName     string `xorm:"'country_name' varchar(25) notnull default '' comment('国家名称')" json:"country_name"`
	CountryCode     string `xorm:"'country_code' varchar(25) notnull default '' comment('国家简码')" json:"country_code"`
	City            string `xorm:"'city' varchar(25) notnull default '' comment('城市')" json:"city"`
	TrialTime       int64  `xorm:"'trial_time' bigint(20) default 0 notnull comment('试用时间')" json:"trial_time"`
	CurrencyCode    string `xorm:"'currency_code' varchar(10) notnull default '' comment('货币简码')" json:"currency_code"`
	MoneyFormat     string `xorm:"'money_format' varchar(20) notnull default '' comment('货币单位符号')" json:"money_format"`
	LastLogin       int64  `xorm:"'last_login' bigint(20) default 0 notnull" json:"last_login"`
	IsDel           int    `xorm:"'is_del' tinyint(1) default 1 notnull comment('删除状态 1 正常 2 卸载 3关店')" json:"is_del"`
	PublishId       string `xorm:"'publish_id' varchar(100) notnull default '' comment('店铺publish_id')" json:"publish_id"`
	Steps           string `xorm:"'steps' varchar(500) notnull default '' comment('新手引导')" json:"steps"`
	Collection      string `xorm:"'collection' text comment('用户集合列表')" json:"collection"`
	UnInstallTime   int64  `xorm:"'uninstall_time' bigint(20) default 0 notnull comment('卸载时间')" json:"uninstall_time"`
	CreateTime      int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime      int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}
