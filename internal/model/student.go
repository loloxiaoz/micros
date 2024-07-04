package model

import "gorm.io/gorm"

// Student 实体
type Student struct {
	gorm.Model
	Name string `gorm:"index"`
	Age  int    `json:"age"`
}
