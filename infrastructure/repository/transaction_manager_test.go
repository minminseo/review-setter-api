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
		user1, _ := userDomain.ReconstructUserForAuth(
			uuid.New().String(),
			"user1@example.com",
			"encrypted_user1_email",
			"encrypted_password",
			"light",
			"en",
			nil,
		)
		user2, _ := userDomain.ReconstructUserForAuth(
			uuid.New().String(),
			"user2@example.com",
			"encrypted_user2_email",
			"encrypted_password",
			"light",
			"en",
			nil,
		)

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

		want1, _ := userRepo.FindByEmailSearchKey(ctx, user1.EmailSearchKey())
		want2, _ := userRepo.FindByEmailSearchKey(ctx, user2.EmailSearchKey())

		if diff := cmp.Diff(want1.ID(), user1.ID()); diff != "" {
			t.Errorf("FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(want2.ID(), user2.ID()); diff != "" {
			t.Errorf("FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("異常系:トランザクション内でエラーが発生した際はロールバックされること", func(t *testing.T) {
		user1, _ := userDomain.ReconstructUserForAuth(
			uuid.New().String(),
			"rollback1@example.com",
			"encrypted_rollback1_email",
			"encrypted_password",
			"light",
			"en",
			nil,
		)
		user2, _ := userDomain.ReconstructUserForAuth(
			uuid.New().String(),
			"rollback2@example.com",
			"encrypted_rollback2_email",
			"encrypted_password",
			"light",
			"en",
			nil,
		)

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

		want1, _ := userRepo.FindByEmailSearchKey(ctx, user1.EmailSearchKey())
		if want1 != nil {
			t.Errorf("user1が保存されている(ロールバックされていない)")
		}
		want2, _ := userRepo.FindByEmailSearchKey(ctx, user2.EmailSearchKey())
		if want2 != nil {
			t.Errorf("user2が保存されている(ロールバックされていない)")
		}
	})

	t.Run("正常系:複数種類のエンティティを同一トランザクションで保存できること", func(t *testing.T) {
		userID := uuid.New().String()
		user, _ := userDomain.ReconstructUserForAuth(
			userID,
			"multi@example.com",
			"encrypted_multi_email",
			"encrypted_password",
			"light",
			"en",
			nil,
		)
		verification, _ := userDomain.ReconstructEmailVerification(
			uuid.New().String(),
			userID,
			"hashed_code_123",
			time.Now().Add(15 * time.Minute),
		)

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

		savedUser, _ := userRepo.FindByEmailSearchKey(ctx, user.EmailSearchKey())
		savedVerification, _ := verificationRepo.FindByUserID(ctx, userID)

		if diff := cmp.Diff(savedUser.ID(), user.ID()); diff != "" {
			t.Errorf("User FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(savedVerification.UserID(), verification.UserID()); diff != "" {
			t.Errorf("Verification FindByUserID() mismatch (-want +got):\n%s", diff)
		}
	})
}

func errorRepositoryOperation(ctx context.Context, u *userDomain.User) error {
	return errors.New("明示的なエラー")
}
