package usecase

type IUserUsecase interface {
	SignUp(user CreateUserInput) (*createUserOutput, error)
	Login(user loginUserInput) (*loginUserOutput, error)
	GetUserSetting(userID string) (*getUserOutput, error)
	UpdateSetting(user updateUserInput) (*updateUserOutput, error)
}
