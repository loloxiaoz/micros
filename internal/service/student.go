package service

import (
	"micros/internal/db"
	"micros/pkg/toolkit"

	"gorm.io/gorm"
)

// Student 实体
type Student struct {
	gorm.Model
	Name string `gorm:"index"`
	Age int `json:"age"`
}

//CreateStudent 创建学生
func CreateStudent(student *Student) error {
	student.ID = toolkit.GenID(1)
	result := db.Instance().Create(&student) 
	return result.Error
}