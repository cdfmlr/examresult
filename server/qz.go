package server

import (
	"examresult/model"
	"github.com/cdfmlr/qzgo"
	"time"
)

// QzClient 实现和强智教务系统的通信。
// 提供查课表的 QueryWeekCourses 方法。
//
// 请务必使用 NewQzClient 来构造 QzClient 实例。
//
// 由于底层的 qzgo.Client 登录会不定时过期，
// 所以不宜长时间保持一个 QzClient 对象。
type QzClient struct {
	qzgo.Client

	student model.Student

	Auth    *qzgo.AuthUserRespBody
	Current *qzgo.GetCurrentTimeRespBody
}

// NewQzClient 新建一个客户端，
// 完成登录、获取教务时间。
// 参数：
//  - student 是用来登录的学生实例
func NewQzClient(student model.Student) (*QzClient, error) {
	client := &QzClient{
		Client: qzgo.Client{
			School: student.School,
			Xh:     student.Sid,
			Pwd:    student.Password,
		},
		student: student,
	}

	// Login
	authResp, err := client.AuthUser()
	client.Auth = authResp
	if err != nil {
		return client, err
	}

	// Query Current jw Time
	currentResp, err := client.GetCurrentTime(time.Now().Format("2006-01-02"))
	client.Current = currentResp

	return client, err
}

// GetExamResults 获取学生**当前学期**成绩
func (c *QzClient) GetExamResults() ([]model.ExamResult, error) {
	cjcxResult, err := c.GetCjcx(c.Xh, c.Current.Xnxqh)
	if err != nil {
		return nil, err
	}

	var exams []model.ExamResult
	for _, cj := range cjcxResult.Result {
		exams = append(exams, model.ExamResult{
			StudentID: c.student.ID,
			Exam:      cj.Kcmc,
			Result:    cj.Zcj,
			Student:   c.student,
		})
	}

	return exams, nil
}
