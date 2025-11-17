package util

import (
	"net/smtp"
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
			"Content-Type: text/plain; charset=\"UTF-8\";\n\n" +
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
