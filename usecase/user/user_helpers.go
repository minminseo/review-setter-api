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
func (uu *userUsecase) sendVerificationEmail(language, toEmail, code string) error {
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
func (uu *userUsecase) resendVerification(ctx context.Context, userID string, dto CreateUserInput) (*CreateUserOutput, error) {
	newUser, err := userDomain.NewUser(userID, dto.Email, dto.Password, dto.Timezone, dto.ThemeColor, dto.Language)
	if err != nil {
		return nil, err
	}

	err = uu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := uu.userRepo.Update(ctx, newUser); err != nil {
			return err
		}
		if err := uu.userRepo.UpdatePassword(ctx, newUser.ID, newUser.EncryptedPassword); err != nil {
			return err
		}

		// 古い認証コードを削除
		if err := uu.emailVerificationRepo.DeleteByUserID(ctx, newUser.ID); err != nil {
		}

		verificationID := uuid.NewString()
		verification, code, err := userDomain.NewEmailVerification(verificationID, newUser.ID)
		if err != nil {
			return err
		}

		if err := uu.emailVerificationRepo.Create(ctx, verification); err != nil {
			return err
		}

		// メール送信処理
		if err := uu.sendVerificationEmail(newUser.Language, newUser.Email, code); err != nil {
			fmt.Printf("警告: %s への認証メールの再送信に失敗しました: %v\n", newUser.Email, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:    newUser.ID,
		Email: newUser.Email,
	}, nil
}
