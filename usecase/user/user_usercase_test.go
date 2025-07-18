package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

func TestUserUsecase_SignUp_StrictOrderAndCount(t *testing.T) {
	testEmail := "test@example.com"
	testPassword := "password123"
	testSearchKey := "search_key"
	testID := "test-id"

	dto := CreateUserInput{
		Email:      testEmail,
		Password:   testPassword,
		Timezone:   "Asia/Tokyo",
		ThemeColor: "dark",
		Language:   "ja",
	}

	tests := []struct {
		name     string
		dto      CreateUserInput
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator)
		wantErr  bool
	}{
		{
			name: "新規ユーザー作成成功",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(nil, errors.New("not found")).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),

					mockUserRepo.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),

					mockEmailSender.EXPECT().
						SendVerificationEmail(dto.Language, testEmail, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "既存ユーザー（認証済み）でエラー",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				existingUser := &userDomain.User{
					ID:         testID,
					VerifiedAt: &time.Time{},
				}
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(existingUser, nil).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "既存ユーザー（未認証）で認証コード再送信成功",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				existingUser := &userDomain.User{
					ID:         testID,
					VerifiedAt: nil,
				}
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(existingUser, nil).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),

					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),

					mockUserRepo.EXPECT().
						UpdatePassword(gomock.Any(), testID, gomock.Any()).
						Return(nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						DeleteByUserID(gomock.Any(), testID).
						Return(nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),

					mockEmailSender.EXPECT().
						SendVerificationEmail(dto.Language, testEmail, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "トランザクション失敗",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(nil, errors.New("not found")).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						Return(errors.New("transaction failed")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator)
			result, err := usecase.SignUp(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("SignUp() result should not be nil when no error expected")
			}
		})
	}
}

func TestUserUsecase_VerifyEmail(t *testing.T) {
	testEmail := "test@example.com"
	testCode := "123456"
	testSearchKey := "search_key"
	testID := "test-id"
	testToken := "jwt-token"

	dto := VerifyEmailInput{
		Email: testEmail,
		Code:  testCode,
	}

	tests := []struct {
		name     string
		dto      VerifyEmailInput
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator)
		wantErr  bool
	}{
		{
			name: "認証成功",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(&userDomain.User{ID: testID, ThemeColor: "dark", Language: "ja"}, nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						FindByUserID(gomock.Any(), testID).
						Return(&userDomain.EmailVerification{
							ID:        "verification-id",
							UserID:    testID,
							CodeHash:  userDomain.HashVerificationCodeForTest(testCode),
							ExpiresAt: time.Now().Add(10 * time.Minute),
						}, nil).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),

					mockUserRepo.EXPECT().
						UpdateVerifiedAt(gomock.Any(), gomock.Any(), testID).
						Return(nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						Delete(gomock.Any(), "verification-id").
						Return(nil).
						Times(1),

					mockTokenGenerator.EXPECT().
						GenerateToken(testID).
						Return(testToken, nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "ユーザーが見つからない",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "既に認証済み",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(&userDomain.User{ID: testID, VerifiedAt: &time.Time{}}, nil).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "認証情報が見つからない",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(&userDomain.User{ID: testID}, nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						FindByUserID(gomock.Any(), testID).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "トランザクション失敗",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(&userDomain.User{ID: testID}, nil).
						Times(1),

					mockEmailVerificationRepo.EXPECT().
						FindByUserID(gomock.Any(), testID).
						Return(&userDomain.EmailVerification{
							ID:        "verification-id",
							CodeHash:  userDomain.HashVerificationCodeForTest(testCode),
							ExpiresAt: time.Now().Add(10 * time.Minute),
						}, nil).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						Return(errors.New("transaction failed")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator)
			result, err := usecase.VerifyEmail(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("VerifyEmail() result should not be nil when no error expected")
			}
		})
	}
}

func TestUserUsecase_LogIn(t *testing.T) {
	testEmail := "test@example.com"
	testPassword := "password123"
	testSearchKey := "search_key"
	testID := "test-id"
	testToken := "jwt-token"

	dto := LoginUserInput{
		Email:    testEmail,
		Password: testPassword,
	}

	tests := []struct {
		name     string
		dto      LoginUserInput
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator)
		wantErr  bool
	}{
		{
			name: "ログイン成功",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
				user := &userDomain.User{
					ID:                testID,
					VerifiedAt:        &time.Time{},
					ThemeColor:        "dark",
					Language:          "ja",
					EncryptedPassword: string(hashedPassword),
				}

				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(user, nil).
						Times(1),

					mockTokenGenerator.EXPECT().
						GenerateToken(testID).
						Return(testToken, nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "ユーザーが見つからない",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "メールアドレスが認証されていない",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				user := &userDomain.User{
					ID:         testID,
					VerifiedAt: nil,
				}
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(user, nil).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "トークン生成失敗",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
				user := &userDomain.User{
					ID:                testID,
					VerifiedAt:        &time.Time{},
					EncryptedPassword: string(hashedPassword),
				}
				gomock.InOrder(
					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						FindByEmailSearchKey(gomock.Any(), testSearchKey).
						Return(user, nil).
						Times(1),

					mockTokenGenerator.EXPECT().
						GenerateToken(testID).
						Return("", errors.New("token generation failed")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator)
			result, err := usecase.LogIn(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("LogIn() result should not be nil when no error expected")
			}
		})
	}
}

func TestUserUsecase_GetUserSetting(t *testing.T) {
	testID := "test-id"
	testEmail := "test@example.com"

	tests := []struct {
		name     string
		userID   string
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator, *userDomain.CryptoService)
		want     *GetUserOutput
		wantErr  bool
	}{
		{
			name:   "ユーザー設定取得成功",
			userID: testID,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator, mockCryptoService *userDomain.CryptoService) {
				encryptedEmail, _ := mockCryptoService.Encrypt(testEmail)
				user := &userDomain.User{
					ID:             testID,
					Timezone:       "Asia/Tokyo",
					ThemeColor:     "dark",
					Language:       "ja",
					EncryptedEmail: encryptedEmail,
				}

				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(user, nil).
						Times(1),
				)
			},
			want: &GetUserOutput{
				Email:      testEmail,
				Timezone:   "Asia/Tokyo",
				ThemeColor: "dark",
				Language:   "ja",
			},
			wantErr: false,
		},
		{
			name:   "ユーザーが見つからない",
			userID: testID,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator, mockCryptoService *userDomain.CryptoService) {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator, mockCryptoService)
			result, err := usecase.GetUserSetting(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSetting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("GetUserSetting() result should not be nil when no error expected")
			}
		})
	}
}

func TestUserUsecase_UpdateSetting(t *testing.T) {
	testID := "test-id"
	testEmail := "test@example.com"
	testSearchKey := "search_key"

	dto := UpdateUserInput{
		ID:         testID,
		Email:      testEmail,
		Timezone:   "Asia/Tokyo",
		ThemeColor: "light",
		Language:   "en",
	}

	tests := []struct {
		name     string
		dto      UpdateUserInput
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator, *userDomain.CryptoService)
		wantErr  bool
	}{
		{
			name: "ユーザー設定更新成功",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator, mockCryptoService *userDomain.CryptoService) {
				encryptedEmail, _ := mockCryptoService.Encrypt("old@example.com")
				user := &userDomain.User{
					ID:             testID,
					Timezone:       "Asia/Tokyo",
					ThemeColor:     "dark",
					Language:       "ja",
					EncryptedEmail: encryptedEmail,
				}

				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(user, nil).
						Times(1),

					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "ユーザーが見つからない",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator, mockCryptoService *userDomain.CryptoService) {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "更新処理失敗",
			dto:  dto,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator, mockCryptoService *userDomain.CryptoService) {
				encryptedEmail, _ := mockCryptoService.Encrypt("old@example.com")
				user := &userDomain.User{
					ID:             testID,
					Timezone:       "Asia/Tokyo",
					ThemeColor:     "dark",
					Language:       "ja",
					EncryptedEmail: encryptedEmail,
				}

				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(user, nil).
						Times(1),

					mockHasher.EXPECT().
						GenerateSearchKey(testEmail).
						Return(testSearchKey).
						Times(1),

					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(errors.New("update failed")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator, mockCryptoService)
			result, err := usecase.UpdateSetting(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSetting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("UpdateSetting() result should not be nil when no error expected")
			}
		})
	}
}

func TestUserUsecase_UpdatePassword(t *testing.T) {
	testID := "test-id"
	testPassword := "newpassword123"

	tests := []struct {
		name     string
		userID   string
		password string
		mockFunc func(*userDomain.MockUserRepository, *userDomain.MockEmailVerificationRepository, *transaction.MockITransactionManager, *userDomain.MockIHasher, *MockiEmailSender, *MockiTokenGenerator)
		wantErr  bool
	}{
		{
			name:     "パスワード更新成功",
			userID:   testID,
			password: testPassword,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				user := &userDomain.User{ID: testID}
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(user, nil).
						Times(1),

					mockUserRepo.EXPECT().
						UpdatePassword(gomock.Any(), testID, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name:     "ユーザーが見つからない",
			userID:   testID,
			password: testPassword,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(nil, errors.New("not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name:     "パスワード更新失敗",
			userID:   testID,
			password: testPassword,
			mockFunc: func(mockUserRepo *userDomain.MockUserRepository, mockEmailVerificationRepo *userDomain.MockEmailVerificationRepository, mockTransactionManager *transaction.MockITransactionManager, mockHasher *userDomain.MockIHasher, mockEmailSender *MockiEmailSender, mockTokenGenerator *MockiTokenGenerator) {
				user := &userDomain.User{ID: testID}
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetSettingByID(gomock.Any(), testID).
						Return(user, nil).
						Times(1),

					mockUserRepo.EXPECT().
						UpdatePassword(gomock.Any(), testID, gomock.Any()).
						Return(errors.New("update failed")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := userDomain.NewMockUserRepository(ctrl)
			mockEmailVerificationRepo := userDomain.NewMockEmailVerificationRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockHasher := userDomain.NewMockIHasher(ctrl)
			mockEmailSender := NewMockiEmailSender(ctrl)
			mockTokenGenerator := NewMockiTokenGenerator(ctrl)
			mockCryptoService, _ := userDomain.NewCryptoService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

			usecase := NewUserUsecase(
				mockUserRepo,
				mockEmailVerificationRepo,
				mockTransactionManager,
				mockCryptoService,
				mockHasher,
				mockEmailSender,
				mockTokenGenerator,
			)

			tt.mockFunc(mockUserRepo, mockEmailVerificationRepo, mockTransactionManager, mockHasher, mockEmailSender, mockTokenGenerator)
			err := usecase.UpdatePassword(context.Background(), tt.userID, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
