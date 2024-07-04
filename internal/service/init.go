package service

import (
	"micros/internal/db"
)

//Clean 清理数据库
func Clean() {
	db := db.Instance()
	db.Migrator().DropTable(&Student{})
}

//Migrate 初始化数据库
func Migrate() {
	db := db.Instance()
	// 自动迁移模式
	db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&Student{})
}

//Init 初始化
func Init() {
	Migrate()
}