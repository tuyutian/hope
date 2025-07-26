package users

// UserSetting 用户自定义设置表
type UserSetting struct {
	Id         int64  `xorm:"pk autoincr 'id' comment('ID')"`
	UserId     int64  `xorm:"notnull 'user_id' comment('用户id')"`
	Name       string `xorm:"notnull varchar(255) 'name' comment('自定义设置键')"`
	Value      string `xorm:"notnull text 'value' comment('配置值(JSON格式)')"`
	CreateTime int64  `xorm:"created 'create_time' bigint(20) default 0 notnull comment('创建时间')" json:"create_time"`
	UpdateTime int64  `xorm:"updated 'update_time' bigint(20) default 0 notnull comment('最近修改时间')" json:"update_time"`
}

// TableName 设置 UserSetting 对应的表名
func (u *UserSetting) TableName() string {
	return "user_setting"
}

var DefaultDashboardGuideStep = map[string]bool{
	"enabled":            false,
	"setting_protension": false,
	"setup_widget":       false,
	"how_work":           false,
	"choose":             false,
}

const (
	DashboardGuideStep = "dashboard_guide_step"
	DashboardGuideHide = "dashboard_guide_hide"
)
