package user

import "context"

type IUserUsecase interface {
	SignUp(ctx context.Context, user CreateUserInput) (*CreateUserOutput, error)
	LogIn(ctx context.Context, user LoginUserInput) (*LoginUserOutput, error)
	GetUserSetting(ctx context.Context, userID string) (*GetUserOutput, error)
	UpdateSetting(ctx context.Context, user UpdateUserInput) (*UpdateUserOutput, error)
	UpdatePassword(ctx context.Context, userID, password string) error
}
