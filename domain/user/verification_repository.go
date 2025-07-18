package user

import "context"

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *EmailVerification) error
	FindByUserID(ctx context.Context, userID string) (*EmailVerification, error)
	DeleteByUserID(ctx context.Context, userID string) error
}
