package user

import "context"

type IUserUsecase interface {
	SignUp(ctx context.Context, user CreateUserInput) (*CreateUserOutput, error)
	LogIn(ctx context.Context, user LoginUserInput) (*LoginUserOutput, error)
	GetUserSetting(ctx context.Context, userID string) (*GetUserOutput, error)
	UpdateSetting(ctx context.Context, user UpdateUserInput) (*UpdateUserOutput, error)
	UpdatePassword(ctx context.Context, userID, password string) error
	VerifyEmail(ctx context.Context, input VerifyEmailInput) (*LoginUserOutput, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, code, newPassword string) error
}

type iEmailSender interface {
	SendVerificationEmail(ctx context.Context, language, toEmail, code string) error
}

type iTokenGenerator interface {
	GenerateToken(userID string) (string, error)
}
