package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

type userUsecase struct {
	userRepo              userDomain.UserRepository
	emailVerificationRepo userDomain.EmailVerificationRepository
	transactionManager    transaction.ITransactionManager
	cryptoService         *userDomain.CryptoService
	hasher                userDomain.IHasher
	emailSender           iEmailSender
	tokenGenerator        iTokenGenerator
}

func NewUserUsecase(
	userRepo userDomain.UserRepository,
	emailVerificationRepo userDomain.EmailVerificationRepository,
	transactionManager transaction.ITransactionManager,
	cryptoService *userDomain.CryptoService,
	hasher userDomain.IHasher,
	emailSender iEmailSender,
	tokenGenerator iTokenGenerator,
) IUserUsecase {

	return &userUsecase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		transactionManager:    transactionManager,
		cryptoService:         cryptoService,
		hasher:                hasher,
		emailSender:           emailSender,
		tokenGenerator:        tokenGenerator,
	}
}

func (uu *userUsecase) SignUp(ctx context.Context, dto CreateUserInput) (*CreateUserOutput, error) {
	// 検索キーを生成して既存ユーザーかチェック
	searchKey := uu.hasher.GenerateSearchKey(dto.Email)
	existingUser, err := uu.userRepo.FindByEmailSearchKey(ctx, searchKey)
	if err == nil && existingUser != nil {
		// 認証済みならエラー
		if existingUser.IsVerified() {
			return nil, errors.New("このメールアドレスは既に使用されています")
		}

		// 未認証なら認証コードを再送信
		newUser, err := userDomain.NewUser(existingUser.ID, dto.Email, dto.Password, dto.Timezone, dto.ThemeColor, dto.Language, uu.cryptoService, searchKey)
		if err != nil {
			return nil, err
		}
		err = uu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			if err := uu.userRepo.Update(ctx, newUser); err != nil {
				return err
			}
			if err := uu.userRepo.UpdatePassword(ctx, newUser.ID, newUser.EncryptedPassword); err != nil {
				return err
			}

			// 古い認証コードを削除
			if err := uu.emailVerificationRepo.DeleteByUserID(ctx, newUser.ID); err != nil {
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

			// メール送信処理
			if err := uu.emailSender.SendVerificationEmail(newUser.Language, dto.Email, code); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		result := &CreateUserOutput{
			ID:    newUser.ID,
			Email: dto.Email,
		}
		return result, nil
	}

	id := uuid.NewString()
	newUser, err := userDomain.NewUser(id, dto.Email, dto.Password, dto.Timezone, dto.ThemeColor, dto.Language, uu.cryptoService, searchKey)
	if err != nil {
		return nil, err
	}
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

		if err := uu.emailSender.SendVerificationEmail(newUser.Language, dto.Email, code); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	decryptedEmail, err := newUser.GetEmail(uu.cryptoService)
	if err != nil {
		return nil, err
	}
	resUser := &CreateUserOutput{
		ID:    newUser.ID,
		Email: decryptedEmail,
	}

	return resUser, nil
}

// VerifyEmail は認証コードを検証し、ユーザーを有効化します。
func (uu *userUsecase) VerifyEmail(ctx context.Context, dto VerifyEmailInput) (*LoginUserOutput, error) {
	searchKey := uu.hasher.GenerateSearchKey(dto.Email)
	user, err := uu.userRepo.FindByEmailSearchKey(ctx, searchKey)
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
		if err := uu.userRepo.UpdateVerifiedAt(ctx, user.VerifiedAt, user.ID); err != nil {
			return err
		}
		if err := uu.emailVerificationRepo.DeleteByUserID(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 認証成功後、JWTを発行してログインさせる
	tokenString, err := uu.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("トークンの生成に失敗しました")
	}
	result := &LoginUserOutput{
		Token:      tokenString,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}
	return result, nil
}

func (uu *userUsecase) LogIn(ctx context.Context, dto LoginUserInput) (*LoginUserOutput, error) {
	searchKey := uu.hasher.GenerateSearchKey(dto.Email)
	user, err := uu.userRepo.FindByEmailSearchKey(ctx, searchKey)
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

	tokenString, err := uu.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("トークンの生成に失敗しました")
	}

	result := &LoginUserOutput{
		Token:      tokenString,
		ThemeColor: user.ThemeColor,
		Language:   user.Language,
	}
	return result, nil
}

func (uu *userUsecase) GetUserSetting(ctx context.Context, userID string) (*GetUserOutput, error) {
	user, err := uu.userRepo.GetSettingByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	email, err := user.GetEmail(uu.cryptoService)
	if err != nil {
		return nil, err
	}

	resUser := &GetUserOutput{
		Email:      email,
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

	searchKey := uu.hasher.GenerateSearchKey(user.Email)

	err = targetUser.Set(user.Email, user.Timezone, user.ThemeColor, user.Language, uu.cryptoService, searchKey)
	if err != nil {
		return nil, err
	}

	err = uu.userRepo.Update(ctx, targetUser)
	if err != nil {
		return nil, err
	}

	decryptedEmail, err := targetUser.GetEmail(uu.cryptoService)
	if err != nil {
		return nil, err
	}

	resUser := &UpdateUserOutput{
		Email:      decryptedEmail,
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
