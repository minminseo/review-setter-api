package user

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id                string
	emailSearchKey    string
	encryptedEmail    string
	encryptedPassword string
	timezone          string
	themeColor        string
	language          string
	verifiedAt        *time.Time
}

func NewUser(
	id string, // ID生成はユースケースに任せる
	email string,
	password string,
	timezone string,
	themeColor string,
	language string,
	cryptoService *CryptoService,
	searchKey string,

) (*User, error) {

	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if err := validateTimezone(timezone); err != nil {
		return nil, err
	}
	if err := validateThemeColor(themeColor); err != nil {
		return nil, err
	}
	if err := validateLanguage(language); err != nil {
		return nil, err
	}

	encryptedEmail, err := cryptoService.Encrypt(email)
	if err != nil {
		return nil, err
	}
	encryptedPassword := encrypt(password)

	u := &User{
		id:                id,
		emailSearchKey:    searchKey,
		encryptedEmail:    encryptedEmail,
		encryptedPassword: encryptedPassword,
		timezone:          timezone,
		themeColor:        themeColor,
		language:          language,
		verifiedAt:        nil, // 初期値は未認証

	}

	return u, nil
}

// 認証用（認証に必要なフィールドのみ）
func ReconstructUserForAuth(
	id string,
	emailSearchKey string,
	encryptedEmail string,
	encryptedPassword string,
	themeColor string,
	language string,
	verifiedAt *time.Time,
) (*User, error) {
	u := &User{
		id:                id,
		emailSearchKey:    emailSearchKey,
		encryptedEmail:    encryptedEmail,
		encryptedPassword: encryptedPassword,
		timezone:          "", // 認証時は不要
		themeColor:        themeColor,
		language:          language,
		verifiedAt:        verifiedAt,
	}
	return u, nil
}

// 設定取得用（パスワード不要）
func ReconstructUserForSettings(
	id string,
	encryptedEmail string,
	timezone string,
	themeColor string,
	language string,
	verifiedAt *time.Time,
) (*User, error) {
	u := &User{
		id:             id,
		encryptedEmail: encryptedEmail,
		timezone:       timezone,
		themeColor:     themeColor,
		language:       language,
		verifiedAt:     verifiedAt,
	}
	return u, nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) EmailSearchKey() string {
	return u.emailSearchKey
}

func (u *User) EncryptedEmail() string {
	return u.encryptedEmail
}

func (u *User) EncryptedPassword() string {
	return u.encryptedPassword
}

func (u *User) Timezone() string {
	return u.timezone
}

func (u *User) ThemeColor() string {
	return u.themeColor
}

func (u *User) Language() string {
	return u.language
}

func (u *User) VerifiedAt() *time.Time {
	return u.verifiedAt
}

func (u *User) UpdateSetting(
	email string,
	timezone string,
	themeColor string,
	language string,
	cryptoService *CryptoService,
	searchKey string,
) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validateTimezone(timezone); err != nil {
		return err
	}
	if err := validateThemeColor(themeColor); err != nil {
		return err
	}
	if err := validateLanguage(language); err != nil {
		return err
	}

	encryptedEmail, err := cryptoService.Encrypt(email)
	if err != nil {
		return err
	}

	u.emailSearchKey = searchKey
	u.encryptedEmail = encryptedEmail
	u.timezone = timezone
	u.themeColor = themeColor
	u.language = language

	return nil
}

func (u *User) GetEmail(cryptoService *CryptoService) (string, error) {
	return cryptoService.Decrypt(u.encryptedEmail)
}

func (u *User) UpdatePassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	encryptedPassword := encrypt(password)
	u.encryptedPassword = encryptedPassword

	return nil
}

const (
	// タイムゾーン
	TimeZoneTokyo      string = "Asia/Tokyo"
	TimeZoneLondon     string = "Europe/London"
	TimeZoneUTC        string = "UTC"
	TimeZoneParis      string = "Europe/Paris"
	TimeZoneMoscow     string = "Europe/Moscow"
	TimeZoneDubai      string = "Asia/Dubai"
	TimeZoneKolkata    string = "Asia/Kolkata"
	TimeZoneShanghai   string = "Asia/Shanghai"
	TimeZoneSydney     string = "Australia/Sydney"
	TimeZoneAuckland   string = "Pacific/Auckland"
	TimeZoneNewYork    string = "America/New_York"
	TimeZoneChicago    string = "America/Chicago"
	TimeZoneDenver     string = "America/Denver"
	TimeZoneLosAngeles string = "America/Los_Angeles"
	TimeZoneHonolulu   string = "Pacific/Honolulu"
	TimeZoneSaoPaulo   string = "America/Sao_Paulo"
	TimeZoneSantiago   string = "America/Santiago"

	// テーマカラー
	ThemeColorDark  string = "dark"
	ThemeColorLight string = "light"

	// 言語
	LanguageJa string = "ja"
	LanguageEn string = "en"
)

var allowedTimeZones = map[string]struct{}{
	TimeZoneTokyo:      {},
	TimeZoneLondon:     {},
	TimeZoneUTC:        {},
	TimeZoneParis:      {},
	TimeZoneMoscow:     {},
	TimeZoneDubai:      {},
	TimeZoneKolkata:    {},
	TimeZoneShanghai:   {},
	TimeZoneSydney:     {},
	TimeZoneAuckland:   {},
	TimeZoneNewYork:    {},
	TimeZoneChicago:    {},
	TimeZoneDenver:     {},
	TimeZoneLosAngeles: {},
	TimeZoneHonolulu:   {},
	TimeZoneSaoPaulo:   {},
	TimeZoneSantiago:   {},
}
var allowedThemeColors = map[string]struct{}{
	ThemeColorDark:  {},
	ThemeColorLight: {},
}
var allowedLanguages = map[string]struct{}{
	LanguageJa: {},
	LanguageEn: {},
}

// パスワードハッシュ化
func encrypt(plainText string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// パスワード検証
func (user *User) IsValidPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.encryptedPassword), []byte(password))
	if err != nil {
		return errors.New("パスワードが一致しません")
	}
	return nil
}

func validateEmail(email string) error {
	return validation.Validate(
		email,
		validation.Required.Error("メールアドレスは必須です"),
		validation.RuneLength(7, 254).Error("メールアドレスは1〜254文字です"),
		is.Email.Error("メールアドレスを入力して下さい"),
	)
}

func validatePassword(password string) error {
	return validation.Validate(
		password,
		validation.Required.Error("パスワードは必須です"),
		validation.RuneLength(6, 0).Error("パスワードは6文字以上です"),
	)
}

func validateTimezone(timezone string) error {
	return validation.Validate(
		timezone,
		validation.Required.Error("タイムゾーンは必須です"),
		validation.RuneLength(1, 64).Error("65文字以上のタイムゾーンは対応していません"),
		validation.By(func(value interface{}) error {
			tz, _ := value.(string)
			if _, ok := allowedTimeZones[tz]; !ok {
				return errors.New("タイムゾーンの値が不正です")
			}
			return nil
		}),
	)
}

func validateThemeColor(themeColor string) error {
	return validation.Validate(
		themeColor,
		validation.Required.Error("テーマカラーは必須です"),
		validation.By(func(value interface{}) error {
			thmclr, _ := value.(string)
			if _, ok := allowedThemeColors[thmclr]; !ok {
				return errors.New("テーマカラーは'dark'または'light'で指定してください")
			}
			return nil
		}),
	)
}

func validateLanguage(language string) error {
	return validation.Validate(
		language,
		validation.Required.Error("言語は必須です"),
		validation.RuneLength(1, 5).Error("5文字以上の言語は対応していません"),
		validation.By(func(value interface{}) error {
			lng, _ := value.(string)
			if _, ok := allowedLanguages[lng]; !ok {
				return errors.New("言語タグの値が不正です")
			}
			return nil
		}),
	)
}

// 認証済みかを確認
func (u *User) IsVerified() bool {
	return u.verifiedAt != nil
}

type IHasher interface {
	GenerateSearchKey(email string) string
}
