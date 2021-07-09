package server

import (
	"examresult/config"
	"github.com/go-gomail/gomail"
	log "github.com/sirupsen/logrus"
)

// Email 是简单的邮件，只有收件人、主体、正文（不能附件）
type Email struct {
	To      string // receiver@addr.com
	Subject string
	Body    string // text/html
}

// EmailServer 邮件发送服务器接口
type EmailServer interface {
	DialTest() error        // 测试登陆
	Send(email Email) error // 发送邮件
}

// SMTP is the parameters to connect to a SMTP server.
// SMTP 就是 SMTP 配置，和 config.ConfSMTP 一摸一样，这里就直接取个别名了
type SMTP = config.ConfSMTP

// emailServer offers simple email service
type emailServer struct {
	SMTP
}

// DialTest 测试登陆
//
// MAIN: 这样应该在 main 里测试一下
func (e emailServer) DialTest() error {
	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	dial, err := d.Dial()
	if err != nil {
		return err
	}
	_ = dial.Close()
	return nil
}

// Send 发送邮件
func (e emailServer) Send(email Email) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.Username)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	return d.DialAndSend(m)
}

// GlobalEmailServer 全局的邮件发送服务器
// 直接调用 GlobalEmailServer.Send() 就可以发邮件了
var GlobalEmailServer EmailServer

func initEmail() {
	log.WithField("SMTP", *config.SMTP).Info("init GlobalEmailServer")
	GlobalEmailServer = emailServer{*config.SMTP}
}
