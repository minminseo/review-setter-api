package usecase

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
)

type CreateUserInput struct {
	Email      string
	Password   string
	Timezone   string
	ThemeColor string
	Language   string
}

type createUserOutput struct {
	ID    string
	Email string
}

type loginUserInput struct {
	Email    string
	Password string
}

type loginUserOutput struct {
	Token      string
	ThemeColor string
	Language   string
}

type getUserOutput struct {
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type updateUserInput struct {
	ID         string
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type updateUserOutput struct {
	Email      string
	Timezone   string
	ThemeColor string
	Language   string
}

type userUsecase struct {
	userRepo userDomain.UserRepository
}

func NewUserUsecase(userRepo userDomain.UserRepository) IUserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (uu *userUsecase) SignUp(dto CreateUserInput) (*createUserOutput, error) {
	id := uuid.NewString()

	newUser, err := userDomain.NewUser(id, dto.Email, dto.Password, dto.Timezone, dto.ThemeColor, dto.Language)
	if err != nil {
		return nil, err
	}

	err = uu.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	resUser := &createUserOutput{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return resUser, nil
}

func (uu *userUsecase) Login(dto loginUserInput) (*loginUserOutput, error) {
	user, err := uu.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	err = user.IsValidPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// JWTトークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}

	resUser := &loginUserOutput{
		Token:      tokenString,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}
	return resUser, nil
}

func (uu *userUsecase) GetUserSetting(userID string) (*getUserOutput, error) {
	user, err := uu.userRepo.GetSettingByID(userID)
	if err != nil {
		return nil, err
	}
	resUser := &getUserOutput{
		Email:      user.Email,
		Timezone:   user.Timezone,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}
	return resUser, nil
}

func (uu *userUsecase) UpdateSetting(user updateUserInput) (*updateUserOutput, error) {
	targetUser, err := uu.userRepo.GetSettingByID(user.ID)
	if err != nil {
		return nil, err
	}

	err = targetUser.Set(user.Email, user.Timezone, user.ThemeColor, user.Language)
	if err != nil {
		return nil, err
	}

	err = uu.userRepo.Update(targetUser)
	if err != nil {
		return nil, err
	}

	resUser := &updateUserOutput{
		Email:      targetUser.Email,
		Timezone:   targetUser.Timezone,
		ThemeColor: targetUser.ThemeColor,
		Language:   targetUser.Language,
	}

	return resUser, nil
}
