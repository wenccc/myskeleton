package email

import (
	"github.com/wenccc/myskeleton/config"
	"sync"
)

type Mailer struct {
	Driver
}

var (
	once         sync.Once
	internalMail *Mailer
)

func NewMail() *Mailer {
	once.Do(func() {
		internalMail = &Mailer{Driver: &SMTP{}}
	})
	return internalMail
}

func (mailer *Mailer) Send(email Email) bool {
	return mailer.Driver.Send(email, config.GetStringMapString("mail.smtp"))
}
