package db

import (
	"sync"
	"time"

	"micros/internal/config"
	"micros/internal/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	once sync.Once
	ins  *gorm.DB
	conf *config.DB
)

// Instance 数据库实例
func Instance() *gorm.DB {
	return ins
}

func assemble(c *config.DB) string {
	dsn := c.User + ":" + c.Password
	dsn += "@tcp(" + c.Host + ")"
	dsn += "/" + c.Database
	dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}

// Init 初始化数据库
func Init(c *config.DB) *gorm.DB {
	once.Do(func() {
		dsn := assemble(c)
		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			log.Logger().Fatal("数据库无法连接" + err.Error())
		}
		sqlDB, err := conn.DB()
		if err != nil {
			log.Logger().Fatal("数据库无法连接" + err.Error())
		}
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime))
		ins = conn
	})
	return ins
}
