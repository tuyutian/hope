package apps

// AppDefinition App定义表
type AppDefinition struct {
	Id          int64  `xorm:"pk autoincr 'id' comment('ID')"`
	AppId       string `xorm:"notnull varchar(50) 'app_id' comment('App唯一标识')"`
	Name        string `xorm:"notnull varchar(100) 'name' comment('App名称')"`
	Description string `xorm:"text 'description' comment('App描述')"`
	IconUrl     string `xorm:"varchar(255) default '' 'icon_url' comment('App图标')"`
	CallbackUrl string `xorm:"varchar(255) default '' 'callback_url' comment('回调URL')"`
	ApiKey      string `xorm:"notnull varchar(100) 'api_key' comment('API Key')"`
	ApiSecret   string `xorm:"notnull varchar(100) 'api_secret' comment('API Secret')"`
	Scopes      string `xorm:"notnull text 'scopes' comment('授权域')"`
	Status      int8   `xorm:"notnull tinyint(1) default 1 'status' comment('状态 1:启用 0:禁用')"`
	CreateTime  int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime  int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

// TableName 设置 AppDefinition 对应的表名
func (a *AppDefinition) TableName() string {
	return "app_definition"
}

// AppConfig App配置表
type AppConfig struct {
	Id          int64  `xorm:"pk autoincr 'id' comment('ID')"`
	AppId       string `xorm:"notnull varchar(50) 'app_id' comment('App标识')"`
	ConfigKey   string `xorm:"notnull varchar(100) 'config_key' comment('配置键')"`
	ConfigValue string `xorm:"notnull text 'config_value' comment('配置值')"`
	Description string `xorm:"varchar(255) default '' 'description' comment('描述')"`
	CreateTime  int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime  int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

// TableName 设置 AppConfig 对应的表名
func (a *AppConfig) TableName() string {
	return "app_config"
}
