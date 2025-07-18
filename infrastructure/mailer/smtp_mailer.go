package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPEmailSender struct{}

func NewSMTPEmailSender() *SMTPEmailSender {
	return &SMTPEmailSender{}
}

func (s *SMTPEmailSender) SendVerificationEmail(language, toEmail, code string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 認証情報
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// メール本文の作成（言語によって切り替え）
	var subject, body string
	switch language {
	case "ja":
		subject = "Subject: Review Setter 認証コード\r\n"
		body = fmt.Sprintf("あなたの認証コードは %s です。\r\n有効期限は10分です。\r\n", code)
	default: // 現状はja以外はenのみ
		subject = "Subject: Review Setter Verification Code\r\n"
		body = fmt.Sprintf("Your verification code is %s.\r\nIt is valid for 10 minutes.\r\n", code)
	}
	msg := []byte("From: " + from + "\r\n" +
		"To: " + toEmail + "\r\n" +
		subject +
		"\r\n" +
		body)

	// Eメール送信
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("メールの送信に失敗しました: %w", err)
	}
	return nil
}
