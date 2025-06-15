package user

type SignUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}

type VerifyEmailResponse struct {
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}

type CsrfTokenResponse struct {
	CsrfToken string `json:"csrf_token"`
}

type GetUserSettingResponse struct {
	Email      string `json:"email"`
	Timezone   string `json:"timezone"`
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}

type UpdateUserSettingResponse struct {
	Email      string `json:"email"`
	Timezone   string `json:"timezone"`
	ThemeColor string `json:"theme_color"`
	Language   string `json:"language"`
}
