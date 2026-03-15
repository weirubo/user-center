package email

import (
	"fmt"
	"net/smtp"

	"user-center/internal/conf"
)

func SendEmail(cfg *conf.SMTP, to, subject, body string) error {
	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n", cfg.From, to, subject)
	msg := []byte(header + body)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	err := smtp.SendMail(addr, auth, cfg.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func SendVerifyCode(cfg *conf.SMTP, to, code string) error {
	subject := "验证码"
	body := fmt.Sprintf("您的验证码是：%s，5分钟内有效。\n\n如非本人操作，请忽略。", code)
	return SendEmail(cfg, to, subject, body)
}
