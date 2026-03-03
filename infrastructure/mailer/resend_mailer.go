package mailer

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

type ResendEmailSender struct {
	client *resend.Client
	from   string
}

func NewResendEmailSender() *ResendEmailSender {

	apikey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apikey)

	from := os.Getenv("RESEND_FROM_EMAIL")
	return &ResendEmailSender{
		client: client,
		from:   from,
	}
}

func (s *ResendEmailSender) SendVerificationEmail(language, toEmail, code string) error {
	var subject, htmlBody string

	switch language {
	case "ja":
		subject = "Review Setter 認証コード"
		htmlBody = fmt.Sprintf("あなたの認証コードは %s です。\r\n有効期限は10分です。\r\n", code)
	default: // 現状はja以外はenのみ
		subject = "Review Setter Verification Code"
		htmlBody = fmt.Sprintf("Your verification code is %s.\r\nIt is valid for 10 minutes.\r\n", code)
	}

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{toEmail},
		Html:    htmlBody,
		Subject: subject,
	}
	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("メールの送信に失敗しました: %w", err)
	}

	return nil
}
