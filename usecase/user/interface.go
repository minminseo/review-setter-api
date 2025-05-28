package user

type IUserUsecase interface {
	SignUp(user CreateUserInput) (*CreateUserOutput, error)
	LogIn(user LoginUserInput) (*LoginUserOutput, error)
	GetUserSetting(userID string) (*GetUserOutput, error)
	UpdateSetting(user UpdateUserInput) (*UpdateUserOutput, error)
	UpdatePassword(userID, password string) error
}
