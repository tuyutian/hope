package users

const (
	RedactHistoryTable = "redact_history"
)

// RedactHistory Redact历史记录表(最小化记录)
type RedactHistory struct {
	Id         int64  `xorm:"pk autoincr 'id' comment('ID')"`
	AppId      string `xorm:"notnull varchar(50) 'app_id' comment('App标识')"`
	Shop       string `xorm:"notnull varchar(100) 'shop' comment('Shop域名')"`
	RedactTime int64  `xorm:"notnull 'redact_time' comment('Redact处理时间')"`
	CreateTime int64  `xorm:"created notnull 'create_time' comment('创建时间')"`
}

func (r *RedactHistory) TableName() string {
	return RedactHistoryTable
}
