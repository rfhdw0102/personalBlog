package utils

import (
	"blogs/pkg/config"
	"fmt"
	"net/smtp"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	emailCfg := config.Get().Email
	// 构建认证信息
	auth := smtp.PlainAuth(
		"",
		emailCfg.Username,
		emailCfg.Password,
		emailCfg.Host,
	)

	// 构建邮件内容
	message := fmt.Sprintf("From: %s\r\n", emailCfg.Username)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// 构建SMTP服务器地址
	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)

	// 发送邮件
	err := smtp.SendMail(
		addr,
		auth,
		emailCfg.Username,
		[]string{to},
		[]byte(message),
	)

	return err
}
