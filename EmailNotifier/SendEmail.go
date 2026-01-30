// Package emailnotifier provides functions for sending email through email
package emailnotifier

import (
	"fmt"
	"os"

	gomail "gopkg.in/gomail.v2"
)

// EmailService holds SMTP server configuration
type EmailService struct {
	Addr   string // smtp.gmail.com
	Port   int    // 587
	User   string // SMTP username
	Pass   string // SMTP password
	From   string // Default sender
	Dialer *gomail.Dialer
}

// NewEmailService constructs an EmailService using environment variables
func NewEmailService() (*EmailService, error) {

	host := os.Getenv("SMTP_ADDRESS") // e.g. smtp.gmail.com
	user := os.Getenv("APP_EMAIL")
	pass := os.Getenv("EMAIL_AUTHPASSWORD")
	from := os.Getenv("APP_EMAIL")

	if host == "" || user == "" || pass == "" || from == "" {
		return nil, fmt.Errorf("missing required SMTP environment variables")
	}

	d := gomail.NewDialer(host, 587, user, pass)

	// // 🔥 Force IPv4 (important for Gmail)
	// d.NewDialer = &net.Dialer{
	// 	Timeout:   10 * time.Second,
	// 	DualStack: false,
	// }

	return &EmailService{
		Addr:   host,
		Port:   587,
		User:   user,
		Pass:   pass,
		From:   from,
		Dialer: d,
	}, nil
}

// SendNotification sends an email with the given recipient and message body
func (e *EmailService) SendNotification(to string, body string) error {

	if to == "" {
		return fmt.Errorf("no recipient specified")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello From AlShifa Platform OTP For Your OnBoarding")
	m.SetBody("text/plain", body)

	return e.Dialer.DialAndSend(m)
}
