package user

import "context"

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *EmailVerification) error
	FindByUserID(ctx context.Context, userID string) (*EmailVerification, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
