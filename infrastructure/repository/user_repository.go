package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type userRepository struct{}

func NewUserRepository() userDomain.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, u *userDomain.User) error {
	q := db.GetQuery(ctx)

	parsed, err := uuid.Parse(u.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	params := dbgen.CreateUserParams{
		ID:             pgID,
		EmailSearchKey: u.EmailSearchKey,
		Email:          u.EncryptedEmail,
		Password:       u.EncryptedPassword,
		Timezone:       u.Timezone,
		ThemeColor:     dbgen.ThemeColorEnum(u.ThemeColor),
		Language:       u.Language,
	}

	return q.CreateUser(ctx, params)
}

func (r *userRepository) FindByEmailSearchKey(ctx context.Context, searchKey string) (*userDomain.User, error) {
	q := db.GetQuery(ctx)

	row, err := q.FindUserByEmailSearchKey(ctx, searchKey)
	if err != nil {
		return nil, err
	}

	if !row.ID.Valid {
		return nil, errors.New("invalid UUID from DB")

	}
	id := uuid.UUID(row.ID.Bytes).String()

	var verifiedAt *time.Time
	if row.VerifiedAt.Valid {
		verifiedAt = &row.VerifiedAt.Time
	}

	return &userDomain.User{
		ID:                id,
		EncryptedEmail:    row.Email,
		EncryptedPassword: row.Password,
		ThemeColor:        string(row.ThemeColor),
		Language:          row.Language,
		VerifiedAt:        verifiedAt,
	}, nil
}

func (r *userRepository) GetSettingByID(ctx context.Context, userID string) (*userDomain.User, error) {
	q := db.GetQuery(ctx)

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	row, err := q.GetUserSettingByID(ctx, pgID)
	if err != nil {
		return nil, err
	}

	ud, err := userDomain.ReconstructUser(
		userID,
		row.Email,
		row.Timezone,
		string(row.ThemeColor),
		row.Language,
		nil, // VerifiedAt使わない
	)
	if err != nil {
		return nil, err
	}
	return ud, nil
}

func (r *userRepository) Update(ctx context.Context, u *userDomain.User) error {
	q := db.GetQuery(ctx)

	parsed, err := uuid.Parse(u.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	params := dbgen.UpdateUserParams{
		EmailSearchKey: u.EmailSearchKey,
		Email:          u.EncryptedEmail,
		Timezone:       u.Timezone,
		ThemeColor:     dbgen.ThemeColorEnum(u.ThemeColor),
		Language:       u.Language,
		ID:             pgID,
	}

	return q.UpdateUser(ctx, params)
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID, password string) error {
	q := db.GetQuery(ctx)

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	params := dbgen.UpdateUserPasswordParams{
		Password: password,
		ID:       pgID,
	}

	return q.UpdateUserPassword(ctx, params)
}

func (r *userRepository) UpdateVerifiedAt(ctx context.Context, verifiedAt *time.Time, userID string) error {
	q := db.GetQuery(ctx)

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	var verifiedAtPg pgtype.Timestamptz
	if verifiedAt != nil {
		verifiedAtPg = pgtype.Timestamptz{Time: *verifiedAt, Valid: true}
	} else {
		verifiedAtPg = pgtype.Timestamptz{Valid: false}
	}

	params := dbgen.UpdateVerifiedAtParams{
		VerifiedAt: verifiedAtPg,
		ID:         pgID,
	}
	return q.UpdateVerifiedAt(ctx, params)
}
