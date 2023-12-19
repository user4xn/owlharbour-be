package helper

import (
	"simpel-api/pkg/util"

	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	from := util.GetEnv("MAIL_USERNAME", "fallback")
	password := util.GetEnv("MAIL_PASSWORD", "fallback")
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
