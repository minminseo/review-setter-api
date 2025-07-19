package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
)

func TestTransactionManager(t *testing.T) {
	if testing.Short() {
		t.Skip("統合テストをスキップします")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	transactionManager := NewTransactionManager(testDBPool)
	userRepo := NewUserRepository()
	verificationRepo := NewEmailVerificationRepository()
	ctx := GetTestContext()

	t.Run("正常系:トランザクションが完了し保存できること", func(t *testing.T) {
		user1 := &userDomain.User{
			ID:                uuid.New().String(),
			EmailSearchKey:    "user1@example.com",
			EncryptedEmail:    "encrypted_user1_email",
			EncryptedPassword: "encrypted_password",
			Timezone:          "UTC",
			ThemeColor:        "light",
			Language:          "en",
		}
		user2 := &userDomain.User{
			ID:                uuid.New().String(),
			EmailSearchKey:    "user2@example.com",
			EncryptedEmail:    "encrypted_user2_email",
			EncryptedPassword: "encrypted_password",
			Timezone:          "UTC",
			ThemeColor:        "light",
			Language:          "en",
		}

		err := transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			if err := userRepo.Create(ctx, user1); err != nil {
				t.Fatalf("userRepo.Create failed: %v", err)
			}
			if err := userRepo.Create(ctx, user2); err != nil {
				t.Fatalf("userRepo.Create failed: %v", err)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("トランザクション内でエラーが発生: %v", err)
		}

		want1, _ := userRepo.FindByEmailSearchKey(ctx, user1.EmailSearchKey)
		want2, _ := userRepo.FindByEmailSearchKey(ctx, user2.EmailSearchKey)

		if diff := cmp.Diff(want1.ID, user1.ID); diff != "" {
			t.Errorf("FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(want2.ID, user2.ID); diff != "" {
			t.Errorf("FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("異常系:トランザクション内でエラーが発生した際はロールバックされること", func(t *testing.T) {
		user1 := &userDomain.User{
			ID:                uuid.New().String(),
			EmailSearchKey:    "rollback1@example.com",
			EncryptedEmail:    "encrypted_rollback1_email",
			EncryptedPassword: "encrypted_password",
			Timezone:          "UTC",
			ThemeColor:        "light",
			Language:          "en",
		}
		user2 := &userDomain.User{
			ID:                uuid.New().String(),
			EmailSearchKey:    "rollback2@example.com",
			EncryptedEmail:    "encrypted_rollback2_email",
			EncryptedPassword: "encrypted_password",
			Timezone:          "UTC",
			ThemeColor:        "light",
			Language:          "en",
		}

		_ = transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			if err := userRepo.Create(ctx, user1); err != nil {
				t.Fatalf("userRepo.Create failed: %v", err)
			}
			if err := userRepo.Create(ctx, user2); err != nil {
				t.Fatalf("userRepo.Create failed: %v", err)
			}
			err := errorRepositoryOperation(ctx, user1)
			if err != nil {
				return err
			}
			return nil
		})

		want1, _ := userRepo.FindByEmailSearchKey(ctx, user1.EmailSearchKey)
		if want1 != nil {
			t.Errorf("user1が保存されている(ロールバックされていない)")
		}
		want2, _ := userRepo.FindByEmailSearchKey(ctx, user2.EmailSearchKey)
		if want2 != nil {
			t.Errorf("user2が保存されている(ロールバックされていない)")
		}
	})

	t.Run("正常系:複数種類のエンティティを同一トランザクションで保存できること", func(t *testing.T) {
		userID := uuid.New().String()
		user := &userDomain.User{
			ID:                userID,
			EmailSearchKey:    "multi@example.com",
			EncryptedEmail:    "encrypted_multi_email",
			EncryptedPassword: "encrypted_password",
			Timezone:          "UTC",
			ThemeColor:        "light",
			Language:          "en",
		}
		verification := &userDomain.EmailVerification{
			ID:        uuid.New().String(),
			UserID:    userID,
			CodeHash:  "hashed_code_123",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}

		err := transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			if err := userRepo.Create(ctx, user); err != nil {
				return err
			}
			if err := verificationRepo.Create(ctx, verification); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			t.Fatalf("トランザクション内でエラーが発生: %v", err)
		}

		savedUser, _ := userRepo.FindByEmailSearchKey(ctx, user.EmailSearchKey)
		savedVerification, _ := verificationRepo.FindByUserID(ctx, userID)

		if diff := cmp.Diff(savedUser.ID, user.ID); diff != "" {
			t.Errorf("User FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(savedVerification.UserID, verification.UserID); diff != "" {
			t.Errorf("Verification FindByUserID() mismatch (-want +got):\n%s", diff)
		}
	})
}

func errorRepositoryOperation(ctx context.Context, u *userDomain.User) error {
	return errors.New("明示的なエラー")
}
