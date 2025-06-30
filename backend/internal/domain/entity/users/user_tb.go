package users

// User 用户表
type User struct {
	ID              int64  `xorm:"pk autoincr 'id' comment('ID')"`
	AppId           string `xorm:"notnull varchar(50) 'app_id' comment('App标识')"`
	Name            string `xorm:"notnull varchar(255) 'name' comment('shopify-name')"`
	Shop            string `xorm:"notnull varchar(100) 'shop' comment('shopify域名')"`
	RealDomain      string `xorm:"varchar(100) default '' 'real_domain' comment('网站真实域名')"`
	PlanDisplayName string `xorm:"varchar(40) default '' 'plan_display_name' comment('shopify套餐版本')"`
	AccessToken     string `xorm:"varchar(512) default '' 'access_token' comment('shopify-token')"`
	Password        string `xorm:"varchar(255) default '' 'password' comment('密码')"`
	Plans           int    `xorm:"int(11) default 0 'plans' comment('app套餐id')"`
	Email           string `xorm:"varchar(100) default '' 'email' comment('邮箱')"`
	Phone           string `xorm:"varchar(20) default '' 'phone' comment('电话号码')"`
	CountryName     string `xorm:"varchar(50) default '' 'country_name' comment('国家名称')"`
	CountryCode     string `xorm:"varchar(5) default '' 'country_code' comment('国家简码')"`
	City            string `xorm:"varchar(50) default '' 'city' comment('城市')"`
	FreeTrialDays   int8   `xorm:"tinyint(4) default 0 'free_trial_days' comment('试用天数')"`
	TrialTime       int64  `xorm:"default 0 'trial_time' comment('试用时间')"`
	CurrencyCode    string `xorm:"varchar(10) default '' 'currency_code' comment('货币简码')"`
	Timezone        int    `xorm:"int(11) default 0 'timezone' comment('The shop's time zone offset expressed as a number of minutes.')"`
	MoneyFormat     string `xorm:"varchar(20) default '' 'money_format' comment('货币单位符号')"`
	LastLogin       int64  `xorm:"default 0 'last_login' comment('最后登录时间')"`
	IsDel           int8   `xorm:"tinyint(1) default 0 'is_del' comment('删除状态 0正常 1已删除')"`
	PublishId       int64  `xorm:"bigint(20) default '' 'publish_id' comment('店铺publish_id')"`
	InstallTime     int64  `xorm:"bigint(20) default 0 'install_time' comment('安装时间')"`
	UninstallTime   int64  `xorm:"bigint(20) default 0 'uninstall_time' comment('卸载时间')"`
	CreateTime      int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime      int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

// TableName 设置 User 对应的表名
func (u *User) TableName() string {
	return "user"
}
