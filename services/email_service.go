package services

import (
	"fmt"
	"temp/global"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// EmailService handles email sending operations
type EmailService struct {
	enabled  bool
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
}

// NewEmailService creates a new email service
func NewEmailService(enabled bool, host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		enabled:  enabled,
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendEmail sends a generic email
func (s *EmailService) SendEmail(to, subject, body string) error {
	if !s.enabled {
		global.Logger.Info("Email sending is disabled in config")
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.username, s.password)
	if err := d.DialAndSend(m); err != nil {
		global.Logger.Error("Failed to send email", zap.Error(err))
		return err
	}
	global.Logger.Info("Email sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Our Service!"
	body := fmt.Sprintf(`
        <html>
        <body>
            <h1>Welcome, %s!</h1>
            <p>Thank you for registering. We're excited to have you on board.</p>
        </body>
        </html>
    `, name)
	return s.SendEmail(to, subject, body)
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(to, verificationToken string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf(`
        <html>
        <body>
            <h1>Email Verification</h1>
            <p>Please verify your email address by clicking the link below:</p>
            <p><a href="http://localhost:8080/verify-email?token=%s">Verify Email</a></p>
        </body>
        </html>
    `, verificationToken)
	return s.SendEmail(to, subject, body)
}
