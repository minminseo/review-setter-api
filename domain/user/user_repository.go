package user

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	GetSettingByID(ctx context.Context, userID string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID, password string) error
	UpdateVerifiedAt(ctx context.Context, verifiedAt *time.Time, userID string) error
}
