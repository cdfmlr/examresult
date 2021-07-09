package model

import (
	"errors"

	"gorm.io/gorm"
)

// ExamResult 考试结果，关联到某个学生
type ExamResult struct {
	gorm.Model

	StudentID uint   // 学生 ID：gorm.Model 给的 primerKey
	Exam      string // 考试名称
	Result    string // 考试结果：成绩

	Noticed bool // 提醒过了

	Student Student `gorm:"foreignKey:StudentID"`
}

// FirstOrCreateExamResult 保存一条新成绩记录 || 更新已存在的
// 调用参数必须包含非空的 Exam、Result 字段、以及从库中查询出来的 Student （有 gorm ID 的）：
//   FirstOrCreateExamResult(ExamResult{
//		 Exam:      "not_empty",
//		 Result:    "not_empty",
//		 Student:   Student{
//			 Model:    gorm.Model{ID: not_zero},
//		 },
//	  })
// 返回查询到的，或新建的 ExamResult
//
// 如果已经存在，**不会更新**原来的值。
func FirstOrCreateExamResult(result ExamResult) (ExamResult, error) {
	exam := ExamResult{}

	if result.Student.ID == 0 {
		return exam, ExamStudentError
	}

	// refer: https://gorm.io/docs/advanced_query.html#FirstOrCreate
	res := DB.Where(&ExamResult{ // 学生 ID + 考试名 唯一确定一个考试成绩
		StudentID: result.Student.ID,
		Exam:      result.Exam,
		// Result:    result.Result,
	}).Assign(result).Preload("Student").FirstOrCreate(&exam)

	return exam, res.Error
}

// GetExamsOfStudent 获取某个学生的全部成绩
func GetExamsOfStudent(student Student) ([]ExamResult, error) {
	var exams []ExamResult

	if student.ID == 0 {
		return exams, ExamStudentError
	}

	result := DB.Where(&ExamResult{
		StudentID: student.ID,
	}).Find(&exams)

	return exams, result.Error
}

var (
	ExamStudentError = errors.New("model: query ExamResult with unrecorded Student")
)
