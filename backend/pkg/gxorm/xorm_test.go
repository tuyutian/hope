package gxorm

import (
	"log"
	"testing"
	"time"

	"xorm.io/xorm"
)

/**
* sql:
* CREATE DATABASE IF NOT EXISTS test default charset utf8mb4;
* create table user (id int primary key auto_increment,name varchar(200),age tinyint) engine=innodb;
* 模拟数据插入
* mysql> insert into user (name) values("xiaoming");
   Query OK, 1 row affected (0.11 sec)

   mysql> insert into user (name) values("hello");
   Query OK, 1 row affected (0.04 sec)
*/

type myUser struct {
	Id   int    `xorm:"pk autoincr"` // 定义的字段属性，要用空格隔开
	Name string `xorm:"varchar(50)"`
	User string `xorm:"varchar(50)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestXORM(t *testing.T) {
	var e *xorm.Engine
	log.Println(e == nil)

	dbConf := &DbConf{
		DbBaseConf: DbBaseConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "root123456",
			Database:  "test",
			ParseTime: true,
		},

		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ShowSql:      true,
	}

	// 设置数据库连接对象，并非真正连接，只有在用的时候才会建立连接
	db, err := dbConf.NewEngine()
	if db == nil || err != nil {
		log.Println("db error")
		return
	}

	// 关闭数据库连接
	defer db.Close()

	log.Println("====get user===")
	user := &myUser{}
	has, err := db.Where("id = ?", 1).Get(user)
	log.Println(has, err)
	log.Println("user info: ", user.Id, user.Name)
}

/*
=== RUN   TestXORM
2025/03/21 10:25:04 true
2025/03/21 10:25:04 ====get user===
[xorm] [info]  2025/03/21 10:25:04.244163 [SQL]
SELECT `id`, `name`, `user` FROM `user` WHERE (id = ?) LIMIT 1 [1] - 1.77225ms
2025/03/21 10:25:04 true <nil>
2025/03/21 10:25:04 user info:  1 heige
--- PASS: TestXORM (0.00s)
PASS
*/

func TestWriteRead(t *testing.T) {
	// 测试读写分离场景
	rwConf := &EngineGroupConf{
		Master: DbBaseConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "root123456",
			Database:  "test",
			ParseTime: true,
		},
		Slaves: []DbBaseConf{
			DbBaseConf{
				Ip:        "127.0.0.1",
				Port:      3306,
				User:      "test1",
				Password:  "root123456",
				Database:  "test",
				ParseTime: true,
			},
			DbBaseConf{
				Ip:        "127.0.0.1",
				Port:      3306,
				User:      "test2",
				Password:  "root123456",
				Database:  "test",
				ParseTime: true,
			},
		},
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ShowSql:      true,
		MaxLifetime:  200 * time.Second,
	}

	eg, err := rwConf.NewEngineGroup()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	userInfo := &myUser{}
	has, err := eg.Where("id = ?", 1).Get(userInfo)
	log.Println("get id = 1 of userInfo: ", has, err)

	log.Println("=======engine select=========")
	user2 := &myUser{}
	has, err = eg.Where("id = ?", 3).Get(user2)
	log.Println(has, err)
	log.Println(user2)

	// 采用读写分离实现数据插入
	user4 := &myUser{
		Name: "xiaoxiao",
		User: "xixi",
	}

	// 插入单条数据，多条数据请用Insert(user3,user4,user5)
	affectedNum, err := eg.InsertOne(user4)
	log.Println("affected num: ", affectedNum)
	log.Println("insert id: ", user4.Id)
	log.Println("err: ", err)

	log.Println("get on slave to query")
	user5 := &myUser{}
	log.Println(eg.Slave().Where("id = ?", 4).Get(user5))
}
