package services

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
)

type MailService struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewMailService() *MailService {
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = os.Getenv("SMTP_USER")
	}
	return &MailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     from,
	}
}

func (m *MailService) SendMail(to []string, subject, body string) error {
	if m.host == "" {
		return fmt.Errorf("mail: SMTP_HOST is not set")
	}
	if len(to) == 0 {
		return fmt.Errorf("mail: no recipients")
	}

	addr := net.JoinHostPort(m.host, m.port)

	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.password, m.host)
	}

	msg := m.buildMessage(to, subject, body)
	if err := smtp.SendMail(addr, auth, m.from, to, msg); err != nil {
		return fmt.Errorf("mail: send failed: %w", err)
	}
	return nil
}

func (m *MailService) buildMessage(to []string, subject, body string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", m.from)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(to, ", "))
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return []byte(b.String())
}
