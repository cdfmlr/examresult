package model

import (
	"errors"
	"gorm.io/gorm"
)

// Student 是学生模型
// 这里指的学生是一个拥有 学校、学号、密码，有资格登录教务系统(jwxt.{School}.edu.cn)的对象。
// 学校、学号 (School 和 Sid) 两个字段唯一确定一名学生
type Student struct {
	gorm.Model
	School   string `gorm:"not null"`
	Sid      string `gorm:"not null"`
	Password string `gorm:"not null"`

	// 用来通知的，还可以写更多种通知方式，也就需要更多种地址
	Email string
}

// SaveStudent 保存新学生 || 更新已存在的学生
// 返回更新后的学生
func SaveStudent(student Student) (Student, error) {
	s := Student{}

	if student.School == "" || student.Sid == "" {
		return s, StudentRecordError
	}

	// 1. Find a record with Where conditions
	//    or create one if not found
	// 2. Update the record with Assign attributes,
	//    regardless found or not
	// 3. put the result record to &s
	//
	// Refer: https://gorm.io/docs/advanced_query.html#FirstOrCreate
	res := DB.Where(&Student{
		School: student.School,
		Sid:    student.Sid,
	}).Assign(student).FirstOrCreate(&s)

	return s, res.Error
}

// GetAllStudents 获取全部学生
func GetAllStudents() ([]Student, error) {
	var students []Student
	result := DB.Find(&students)
	return students, result.Error
}

var (
	StudentRecordError = errors.New("model: Student missing School or Sid")
)
