package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"strconv"

	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	To       string
	Subject  string
	Body     string
	HTMLBody string
	IsHTML   bool
}

func NewMail(to, subject, body string) *Mail {
	return &Mail{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}
}

func NewHTMLMail(to, subject, htmlBody string) *Mail {
	return &Mail{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
		IsHTML:   true,
	}
}

func NewMailFromTemplate(to, subject, templatePath string, data interface{}) (*Mail, error) {
	htmlContent, err := parseTemplate(templatePath, data)
	if err != nil {
		return nil, err
	}

	return &Mail{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlContent,
		IsHTML:   true,
	}, nil
}

func parseTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (m *Mail) Send() error {
	host := app.GetEnv("MAIL_SERVER", "smtp.gmail.com")
	portStr := app.GetEnv("MAIL_PORT", "587")
	username := app.GetEnv("MAIL_USERNAME", "")
	password := app.GetEnv("MAIL_PASSWORD", "")

	if host == "" || portStr == "" || username == "" || password == "" {
		return errors.New("SMTP configuration is not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("invalid SMTP_PORT value")
	}

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mail := gomail.NewMessage()

	// Validation
	if m.To == "" || m.Subject == "" {
		return errors.New("to and subject fields cannot be empty")
	}

	if !m.IsHTML && m.Body == "" {
		return errors.New("body field cannot be empty")
	}

	if m.IsHTML && m.HTMLBody == "" {
		return errors.New("HTML body field cannot be empty")
	}

	mail.SetHeader("From", username)
	mail.SetHeader("To", m.To)
	mail.SetHeader("Subject", m.Subject)

	if m.IsHTML {
		mail.SetBody("text/html", m.HTMLBody)
		if m.Body != "" {
			mail.AddAlternative("text/plain", m.Body)
		}
	} else {
		mail.SetBody("text/plain", m.Body)
	}

	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
