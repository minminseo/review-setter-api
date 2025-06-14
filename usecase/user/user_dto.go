package user

type CreateUserInput struct {
	Email      string
	Password   string
	Timezone   string
	ThemeColor string
	Language   string
}

type CreateUserOutput struct {
	ID    string
	Email string
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token      string
	ThemeColor string
	Language   string
}

type GetUserOutput struct {
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type UpdateUserInput struct {
	ID         string
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type UpdateUserOutput struct {
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type VerifyEmailInput struct {
	Email string
	Code  string
}
