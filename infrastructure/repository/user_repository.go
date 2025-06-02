package repository

import (
	"context"
	"errors"

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
		ID:         pgID,
		Email:      u.Email,
		Password:   u.EncryptedPassword,
		Timezone:   u.Timezone,
		ThemeColor: dbgen.ThemeColorEnum(u.ThemeColor),
		Language:   u.Language,
	}

	return q.CreateUser(ctx, params)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	q := db.GetQuery(ctx)

	row, err := q.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !row.ID.Valid {
		return nil, errors.New("invalid UUID from DB")

	}
	id := uuid.UUID(row.ID.Bytes).String()

	return &userDomain.User{
		ID:                id,
		Email:             email,
		EncryptedPassword: row.Password,
		ThemeColor:        string(row.ThemeColor),
		Language:          row.Language,
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
		Email:      u.Email,
		Timezone:   u.Timezone,
		ThemeColor: dbgen.ThemeColorEnum(u.ThemeColor),
		Language:   u.Language,
		ID:         pgID,
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
