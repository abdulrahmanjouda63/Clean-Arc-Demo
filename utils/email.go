package utils

import (
	"gopkg.in/gomail.v2"
)

// SendEmail sends an email using the provided SMTP configuration.
func SendEmail(to, subject, body, smtpUser, smtpPass, smtpHost string, smtpPort int) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return d.DialAndSend(m)
}
