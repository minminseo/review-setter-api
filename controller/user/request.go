package user

type sighUpUserRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Timezone   string `json:"timezone"`
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}

type logInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email      string `json:"email"`
	Timezone   string `json:"timezone"`
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}

type updatePasswordRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
