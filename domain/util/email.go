package util

import (
	"bytes"
	"cachacariaapi/domain/entities"
	"html/template"
	"net/smtp"
	"time"
)

type EmailConfig struct {
	SMTPHost string `toml:"host"`
	SMTPPort string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}

func NewEmailConfig(
	SMTPHost string, SMTPPort string,
	Username string,
	Password string,
	From string,
) *EmailConfig {
	return &EmailConfig{
		SMTPHost: SMTPHost,
		SMTPPort: SMTPPort,
		Username: Username,
		Password: Password,
		From:     From,
	}
}

// SendEmail sends a plain text email using SMTP
func SendEmail(cfg EmailConfig, to []string, subject, body string) error {
	msg := []byte(
		"Subject: " + subject + "\n" +
			"MIME-Version: 1.0;\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
			body,
	)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)

	return smtp.SendMail(
		cfg.SMTPHost+":"+cfg.SMTPPort,
		auth,
		cfg.From,
		to,
		msg,
	)
}

func RenderTemplate(path string, data any) (string, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendAccountCreatedEmail(cfg EmailConfig, to []string, user entities.User) error {
	html, err := RenderTemplate("templates/account_created.gohtml", map[string]any{
		"Name":  user.Email,
		"Email": user.Email,
		"Year":  time.Now().Year(),
	})
	if err != nil {
		return err
	}

	return SendEmail(cfg, to, "Conta criada com sucesso", html)
}

func SendPasswordChangedEmail(cfg EmailConfig, to []string, user entities.User) error {
	html, err := RenderTemplate("templates/password_changed.gohtml", map[string]any{
		"Name":  user.Email,
		"Email": user.Email,
		"Year":  time.Now().Year(),
	})
	if err != nil {
		return err
	}

	return SendEmail(cfg, to, "Senha alterada com sucesso", html)
}
