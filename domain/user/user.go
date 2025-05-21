package user

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type timeZoneName string
type themeColorName string
type languageName string

type User struct {
	id         string
	email      string
	password   string
	timezone   timeZoneName
	themeColor themeColorName
	language   languageName
}

func NewUser(
	id string, // ID生成はユースケースに任せる
	email string,
	password string,
	timezone timeZoneName,
	themeColor themeColorName,
	language languageName,
) (*User, error) {
	u := &User{
		id:         id,
		email:      email,
		password:   password,
		timezone:   timezone,
		themeColor: themeColor,
		language:   language,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

const (
	// タイムゾーン
	TimeZoneTokyo timeZoneName = "Asia/Tokyo"
	// TimeZoneLondon TimeZoneName = "Europe/London"

	// テーマカラー
	ThemeColorDark  themeColorName = "dark"
	ThemeColorLight themeColorName = "light"

	// 言語
	LanguageJa languageName = "ja"
	// LanguageEn languageName = "en"
)

var allowedTimeZones = map[timeZoneName]struct{}{
	TimeZoneTokyo: {},
	// TimeZoneLondon: {},
}
var allowedThemeColors = map[themeColorName]struct{}{
	ThemeColorDark:  {},
	ThemeColorLight: {},
}
var allowedLanguages = map[languageName]struct{}{
	LanguageJa: {},
	// LanguageEn: {},
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(
			&u.email,
			validation.Required.Error("メールアドレスは必須です"),
			validation.RuneLength(1, 255).Error("メールアドレスは1〜255文字です"),
			is.Email.Error("メールアドレスを入力して下さい"),
		),
		validation.Field(
			&u.password,
			validation.Required.Error("パスワードは必須です"),
			validation.RuneLength(6, 0).Error("パスワードは6文字以上です"),
		),
		validation.Field(
			&u.timezone,
			validation.Required.Error("タイムゾーンは必須です"),
			validation.RuneLength(1, 64).Error("65文字以上のタイムゾーンは対応していません"),
			validation.By(func(value interface{}) error {
				tz, _ := value.(timeZoneName)
				if _, ok := allowedTimeZones[tz]; !ok {
					return errors.New("タイムゾーンの値が不正です")
				}
				return nil
			}),
		),
		// テーマカラー
		validation.Field(
			&u.themeColor,
			validation.Required.Error("テーマカラーは必須です"),
			validation.By(func(value interface{}) error {
				thmclr, _ := value.(themeColorName)
				if _, ok := allowedThemeColors[thmclr]; !ok {
					return errors.New("テーマカラーは'dark'または'light'で指定してください")
				}
				return nil
			}),
		),
		// 言語
		validation.Field(
			&u.language,
			validation.Required.Error("言語は必須です"),
			validation.RuneLength(1, 5).Error("5文字以上の言語は対応していません"),
			validation.By(func(value interface{}) error {
				lng, _ := value.(languageName)
				if _, ok := allowedLanguages[lng]; !ok {
					return errors.New("言語タグの値が不正です")
				}
				return nil
			}),
		),
	)
}
