package user

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
)

// 認証コードをEメールで送信
func (uu *userUsecase) sendVerificationEmail(toEmail, code string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 認証情報
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// メール本文の作成
	subject := "Subject: Review Setter 認証コード\r\n"
	body := fmt.Sprintf("あなたの認証コードは %s です。\r\n有効期限は10分です。\r\n", code)
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

// ログイン成功時のレスポンスとJWTを生成
func (uu *userUsecase) createLoginResponse(user *userDomain.User) (*LoginUserOutput, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{
		Token:      tokenString,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}, nil
}

// 未認証ユーザーに認証コードを再送信
func (uu *userUsecase) resendVerification(ctx context.Context, user *userDomain.User, password string) (*CreateUserOutput, error) {
	// パスワードを更新
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	err := uu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := uu.userRepo.Update(ctx, user); err != nil {
			return err
		}

		// 古い認証コードを削除
		if err := uu.emailVerificationRepo.DeleteByUserID(ctx, user.ID); err != nil {
		}

		verificationID := uuid.NewString()
		verification, code, err := userDomain.NewEmailVerification(verificationID, user.ID)
		if err != nil {
			return err
		}

		if err := uu.emailVerificationRepo.Create(ctx, verification); err != nil {
			return err
		}

		// メール送信処理
		if err := uu.sendVerificationEmail(user.Email, code); err != nil {
			fmt.Printf("警告: %s への認証メールの再送信に失敗しました: %v\n", user.Email, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
