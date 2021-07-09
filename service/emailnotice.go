package service

import (
	"examresult/model"
	"examresult/server"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// emailNoticer 发邮件通知成绩
type emailNoticer struct{}

func (e emailNoticer) Notice(result model.ExamResult) bool {
	err := server.GlobalEmailServer.Send(server.Email{
		To:      result.Student.Email,
		Subject: "出成绩了",
		Body: fmt.Sprintf(
			"Hi %s,\n出成绩了!<br/><br/>\n"+
				"<em>%s</em>: <i>%s</i>.<br/><br/>\n\n"+
				"<hr/>Go https://jwxt.%s.edu.cn for details.",
			result.Student.Sid,
			result.Exam,
			result.Result,
			result.Student.School,
		),
	})
	if err != nil {
		log.WithError(err).Warn("service EmailNoticer: send failed")
		return false
	}
	return true
}

var EmailNoticer *emailNoticer

func initEmailNoticer() {
	log.Info("init EmailNoticer")
	EmailNoticer = &emailNoticer{}
	NoticeHandler.AddNoticer(EmailNoticer)
}
