package email

import (
	"fmt"
	jmail "github.com/jordan-wright/email"
	"github.com/wenccc/myskeleton/logger"
	"net/smtp"
)

type SMTP struct{}

func (s *SMTP) Send(email Email, config map[string]string) bool {

	jm := jmail.NewEmail()
	jm.From = fmt.Sprintf("%s <%s>", email.From.Name, email.From.Address)
	jm.To = email.To
	jm.Bcc = email.Bcc
	jm.Cc = email.Cc
	jm.Subject = email.Subject
	jm.HTML = email.HTML
	jm.Text = email.Text

	logger.DebugJSON("邮件", "准备发送邮件", jm)

	err := jm.Send(
		fmt.Sprintf("%s:%s", config["host"], config["port"]),
		smtp.PlainAuth("", config["username"], config["password"], config["host"]),
	)
	if err != nil {
		logger.DebugJSON("邮件", "发送邮件", err.Error())
		return false
	}
	logger.DebugString("发送邮件", "发件成功", "")
	return true
}
