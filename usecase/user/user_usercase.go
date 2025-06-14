package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

type userUsecase struct {
	userRepo              userDomain.UserRepository
	emailVerificationRepo userDomain.EmailVerificationRepository
	transactionManager    transaction.ITransactionManager
}

func NewUserUsecase(
	userRepo userDomain.UserRepository,
	emailVerificationRepo userDomain.EmailVerificationRepository,
	transactionManager transaction.ITransactionManager,
) IUserUsecase {
	return &userUsecase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		transactionManager:    transactionManager,
	}
}

func (uu *userUsecase) SignUp(ctx context.Context, dto CreateUserInput) (*CreateUserOutput, error) {
	// 既存ユーザーかチェック
	existingUser, err := uu.userRepo.FindByEmail(ctx, dto.Email)
	if err == nil && existingUser != nil {
		// 認証済みならエラー
		if existingUser.IsVerified() {
			return nil, errors.New("このメールアドレスは既に使用されています")
		}
		// 未認証なら、情報を更新して認証コードを再送信
		return uu.resendVerification(ctx, existingUser, dto.Password)
	}

	id := uuid.NewString()
	newUser, err := userDomain.NewUser(id, dto.Email, dto.Password, dto.Timezone, dto.ThemeColor, dto.Language)
	if err != nil {
		return nil, err
	}

	// ユーザーと認証情報をトランザクションで保存
	err = uu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := uu.userRepo.Create(ctx, newUser); err != nil {
			return err
		}

		verificationID := uuid.NewString()
		verification, code, err := userDomain.NewEmailVerification(verificationID, newUser.ID)
		if err != nil {
			return err
		}

		if err := uu.emailVerificationRepo.Create(ctx, verification); err != nil {
			return err
		}

		// --- メール送信処理 ---
		if err := uu.sendVerificationEmail(newUser.Email, code); err != nil {
			fmt.Printf("警告: %s への認証メール送信に失敗しました: %v\n", newUser.Email, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resUser := &CreateUserOutput{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return resUser, nil
}

// VerifyEmail は認証コードを検証し、ユーザーを有効化します。
func (uu *userUsecase) VerifyEmail(ctx context.Context, dto VerifyEmailInput) (*LoginUserOutput, error) {
	user, err := uu.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	if user.IsVerified() {
		return nil, errors.New("既に認証済みです")
	}

	verification, err := uu.emailVerificationRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, errors.New("認証情報が見つかりません")
	}

	if verification.IsExpired() {
		return nil, errors.New("認証コードの有効期限が切れています")
	}

	if !verification.ValidateCode(dto.Code) {
		return nil, errors.New("認証コードが正しくありません")
	}

	// ユーザーを認証済みに更新し、認証情報を削除
	err = uu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		user.SetVerified()
		if err := uu.userRepo.Update(ctx, user); err != nil {
			return err
		}
		if err := uu.emailVerificationRepo.Delete(ctx, verification.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 認証成功後、JWTを発行してログインさせる
	return uu.createLoginResponse(user)
}

func (uu *userUsecase) LogIn(ctx context.Context, dto LoginUserInput) (*LoginUserOutput, error) {
	user, err := uu.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	if !user.IsVerified() {
		return nil, errors.New("メールアドレスが認証されていません")
	}

	err = user.IsValidPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	return uu.createLoginResponse(user)

}

func (uu *userUsecase) GetUserSetting(ctx context.Context, userID string) (*GetUserOutput, error) {
	user, err := uu.userRepo.GetSettingByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resUser := &GetUserOutput{
		Email:      user.Email,
		Timezone:   user.Timezone,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}
	return resUser, nil
}

func (uu *userUsecase) UpdateSetting(ctx context.Context, user UpdateUserInput) (*UpdateUserOutput, error) {
	targetUser, err := uu.userRepo.GetSettingByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	err = targetUser.Set(user.Email, user.Timezone, user.ThemeColor, user.Language)
	if err != nil {
		return nil, err
	}

	err = uu.userRepo.Update(ctx, targetUser)
	if err != nil {
		return nil, err
	}

	resUser := &UpdateUserOutput{
		Email:      targetUser.Email,
		Timezone:   targetUser.Timezone,
		ThemeColor: targetUser.ThemeColor,
		Language:   targetUser.Language,
	}

	return resUser, nil
}

func (uu *userUsecase) UpdatePassword(ctx context.Context, userID, password string) error {
	user, err := uu.userRepo.GetSettingByID(ctx, userID)
	if err != nil {
		return err
	}
	err = user.SetPassword(password)
	if err != nil {
		return err
	}
	err = uu.userRepo.UpdatePassword(ctx, userID, user.EncryptedPassword)
	if err != nil {
		return err
	}
	return nil
}
