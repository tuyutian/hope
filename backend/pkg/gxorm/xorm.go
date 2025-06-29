// Package gxorm for xorm engine database connection
package gxorm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	xLog "xorm.io/xorm/log"
)

var (
	// EngineNotExist engine not found
	EngineNotExist = errors.New("current db engine not exist")

	// EngineNameEmpty engine name is empty
	EngineNameEmpty = errors.New("current engine name is empty")
)

// DbBaseConf 数据库基本配置
// 单个数据库实例连接，如果需要读写分离，请使用 EngineGroupConf 初始化master/slaves
type DbBaseConf struct {
	Ip        string
	Port      int
	User      string
	Password  string
	Database  string
	Charset   string // 字符集 utf8mb4 支持表情符号
	Collation string // 整理字符集 utf8mb4_unicode_ci
	ParseTime bool   // 是否格式化时间
	Loc       string // 时区字符串 Local,PRC

	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout
}

// DbConf mysql连接信息
// parseTime=true changes the output type of DATE and DATETIME
// values to time.Time instead of []byte / string
// The date or datetime like 0000-00-00 00:00:00 is converted
// into zero value of time.Time.
type DbConf struct {
	DbBaseConf

	UsePool      bool // 当前db实例是否采用db连接池,默认不采用，如采用请求配置该参数
	MaxIdleConns int  // 设置连接池的空闲数大小
	MaxOpenConns int  // 最大open connection个数

	// sets the maximum amount of time a connection may be reused.
	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	MaxLifetime time.Duration

	ShowSql bool      // 是否输出sql，输出句柄是logger
	Logger  io.Writer // sql日志输出interface
}

// 每个数据库连接pool就是一个db引擎
var engineMap = map[string]*xorm.Engine{}

// InitDB new a db engine
// mysql charset查看
// mysql> show character set where charset="utf8mb4";
// +---------+---------------+--------------------+--------+
// | Charset | Description   | Default collation  | Maxlen |
// +---------+---------------+--------------------+--------+
// | utf8mb4 | UTF-8 Unicode | utf8mb4_general_ci |      4 |
// +---------+---------------+--------------------+--------+
// 1 row in set (0.00 sec)
func (conf *DbBaseConf) InitDB() (*xorm.Engine, error) {
	if conf.Ip == "" {
		conf.Ip = "127.0.0.1"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	// 默认字符序，定义了字符的比较规则
	if conf.Collation == "" {
		conf.Collation = "utf8mb4_general_ci"
	}

	if conf.Loc == "" {
		conf.Loc = "Local"
	}

	if conf.Timeout == 0 {
		conf.Timeout = 10 * time.Second
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 5 * time.Second
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 5 * time.Second
	}

	// mysql connection time loc
	loc, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return nil, err
	}

	// mysql config
	mysqlConf := mysql.Config{
		User:   conf.User,
		Passwd: conf.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
		DBName: conf.Database,
		// Connection parameters
		Params: map[string]string{
			"charset": conf.Charset,
		},
		Collation:            conf.Collation,
		Loc:                  loc,               // Location for time.Time values
		Timeout:              conf.Timeout,      // Dial timeout
		ReadTimeout:          conf.ReadTimeout,  // I/O read timeout
		WriteTimeout:         conf.WriteTimeout, // I/O write timeout
		AllowNativePasswords: true,              // Allows the native password authentication method
		ParseTime:            conf.ParseTime,    // Parse time values to time.Time
	}

	// 连接实例对象，并非立即连接db,用的时候才会真正的建立连接
	var db *xorm.Engine
	db, err = xorm.NewEngine("mysql", mysqlConf.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewEngine create a db engine
// 如果配置上有显示sql执行时间和采用pool机制，就会建立db连接池
func (conf *DbConf) NewEngine() (*xorm.Engine, error) {
	db, err := conf.InitDB()
	if err != nil {
		return nil, err
	}

	if conf.ShowSql {
		if conf.Logger == nil {
			conf.Logger = os.Stdout
		}

		dbLogger := xLog.NewSimpleLogger(conf.Logger)
		dbLogger.ShowSQL(true)
		db.SetLogger(dbLogger)
	}

	// 设置连接池
	if conf.UsePool {
		db.SetMaxIdleConns(conf.MaxIdleConns) // 设置连接池的空闲数大小
		db.SetMaxOpenConns(conf.MaxOpenConns) // 设置最大打开连接数
	}

	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	if conf.MaxLifetime > 0 {
		db.SetConnMaxLifetime(conf.MaxLifetime)
	}

	return db, nil
}

// SetEngineName 给当前数据库指定engineName
// 一般用在多个db 数据库连接引擎的时候，可以给当前的db engine设置一个name
// 这样业务上游层，就可以通过 GetEngine(name)获得当前db engine
func (conf *DbConf) SetEngineName(name string) error {
	if name == "" {
		return EngineNameEmpty
	}

	// 初始化db句柄
	db, err := conf.NewEngine()
	if err != nil {
		return fmt.Errorf("current %s db engine init error: %s", name, err.Error())
	}

	engineMap[name] = db
	return nil
}

// ShortConnect 短连接设置，一般用于短连接服务的数据库句柄
func (conf *DbConf) ShortConnect() (*xorm.Engine, error) {
	conf.UsePool = false
	return conf.NewEngine()
}

// GetEngineByName 从db engine中获取一个数据库连接句柄
// 根据数据库连接句柄name获取指定的连接句柄
func GetEngineByName(name string) (*xorm.Engine, error) {
	if _, ok := engineMap[name]; ok {
		return engineMap[name], nil
	}

	return nil, EngineNotExist
}

// CloseAllDb 由于xorm db.Close()是关闭当前连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllDb() {
	for name, db := range engineMap {
		if err := db.Close(); err != nil {
			fmt.Println("close db error: ", err.Error())
			continue
		}

		delete(engineMap, name) // 销毁连接句柄标识
	}
}

// CloseDbByName 关闭指定name的db engine
func CloseDbByName(name string) error {
	if _, ok := engineMap[name]; ok {
		if err := engineMap[name].Close(); err != nil {
			fmt.Println("close db error: ", err.Error())
			return err
		}

		delete(engineMap, name)
	}

	return EngineNotExist
}

// ======================多个引擎组设置==================

var engineGroupMap = map[string]*xorm.EngineGroup{}

// EngineGroupConf 读写分离引擎配置
type EngineGroupConf struct {
	Master DbBaseConf
	Slaves []DbBaseConf

	UsePool bool // 是否采用db连接池,默认不采用，如采用请求配置该参数
	// the following configuration is for the configuration on each instance of master and slave
	// not the overall configuration of the engine group.
	// 下面的配置对于每个实例的配置，并非整个引擎组的配置
	MaxIdleConns int // 设置连接池的空闲数大小
	MaxOpenConns int // 最大open connection个数
	// sets the maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are reused forever.
	MaxLifetime time.Duration

	ShowSql bool      // 是否输出sql，输出句柄是logger
	Logger  io.Writer // sql日志输出interface
}

// NewEngineGroup 创建读写分离的引擎组，附带一些拓展配置
// 这里可以采用功能模式，方便后面对引擎组句柄进行拓展
// 默认采用连接池方式建立连接
func (conf *EngineGroupConf) NewEngineGroup(policies ...xorm.GroupPolicy) (*xorm.EngineGroup, error) {
	master, err := conf.Master.InitDB()
	if err != nil {
		return nil, err
	}

	slaveLen := len(conf.Slaves)
	if slaveLen == 0 {
		return nil, errors.New("slave db conf is empty")
	}

	slaves := make([]*xorm.Engine, 0, slaveLen)
	slaveErrs := make([]error, 0, slaveLen)

	var db *xorm.Engine
	for k := range conf.Slaves {
		db, err = conf.Slaves[k].InitDB()
		if err != nil {
			slaveErrs = append(slaveErrs, err)
			continue
		}

		slaves = append(slaves, db)
	}

	if len(slaveErrs) > 0 {
		fmt.Println("init slaves error: ", slaveErrs)

		errMsgSlice := make([]string, 0, len(slaveErrs))
		for _, sErr := range slaveErrs {
			errMsgSlice = append(errMsgSlice, sErr.Error())
		}

		return nil, errors.New("init slaves error: " + strings.Join(errMsgSlice, " "))
	}

	var eg *xorm.EngineGroup
	if len(policies) > 0 {
		eg, err = xorm.NewEngineGroup(master, slaves, policies...)
	} else {
		eg, err = xorm.NewEngineGroup(master, slaves)
	}

	if err != nil {
		return nil, err
	}

	if conf.ShowSql {
		if conf.Logger == nil {
			conf.Logger = os.Stdout
		}

		// xorm v1.0.x以上版本使用xorm log包 NewSimpleLogger方法
		dbLogger := xLog.NewSimpleLogger(conf.Logger)
		dbLogger.ShowSQL(true)
		eg.SetLogger(dbLogger)
	}

	if conf.UsePool {
		eg.SetMaxIdleConns(conf.MaxIdleConns) // 最大db空闲数
		eg.SetMaxOpenConns(conf.MaxOpenConns) // db最大连接数,小于数据库配置的最大连接
	}

	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	if conf.MaxLifetime > 0 {
		eg.SetConnMaxLifetime(conf.MaxLifetime)
	}

	return eg, nil
}

// SetEngineGroupName 给db engine group设置名字
func (conf *EngineGroupConf) SetEngineGroupName(name string, policies ...xorm.GroupPolicy) error {
	if name == "" {
		return EngineNameEmpty
	}

	// 初始化db句柄
	var eg *xorm.EngineGroup
	var err error
	if len(policies) > 0 {
		eg, err = conf.NewEngineGroup(policies...)
	} else {
		eg, err = conf.NewEngineGroup()
	}

	if err != nil {
		return fmt.Errorf("current %s db engine group init error: %s", name, err.Error())
	}

	engineGroupMap[name] = eg

	return nil
}

// ======================引擎组辅助方法================

// CloseAllEngineGroup 关闭当前引擎组连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllEngineGroup() {
	for name, db := range engineGroupMap {
		if err := db.Close(); err != nil {
			fmt.Println("close all db engine group error: ", err.Error())
			continue
		}

		delete(engineGroupMap, name) // 销毁连接句柄标识
	}
}

// CloseEngineGroupByName 关闭指定name的db engine group
func CloseEngineGroupByName(name string) error {
	if _, ok := engineGroupMap[name]; ok {
		if err := engineGroupMap[name].Close(); err != nil {
			fmt.Println("close db engine group error: ", err.Error())
			return err
		}

		delete(engineGroupMap, name)
	}

	return EngineNotExist
}

// GetEngineGroupName 从引擎组中获得一个db engine group
func GetEngineGroupName(name string) (*xorm.EngineGroup, error) {
	if _, ok := engineGroupMap[name]; ok {
		return engineGroupMap[name], nil
	}

	return nil, EngineNotExist
}
