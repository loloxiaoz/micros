package service

import (
	"micros/internal/db"
	"micros/internal/model"
	"micros/pkg/toolkit"
)

// CreateStudent 创建学生
func CreateStudent(student *model.Student) error {
	student.ID = toolkit.GenID(1)
	result := db.Instance().Create(&student)
	return result.Error
}
