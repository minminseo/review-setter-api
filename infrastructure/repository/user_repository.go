package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type userRepository struct {
	q *dbgen.Queries
}

func NewUserRepository(db dbgen.DBTX) userDomain.UserRepository {
	return &userRepository{
		q: dbgen.New(db),
	}
}

func (r *userRepository) Create(u *userDomain.User) error {
	ctx := context.Background()

	// string型のidをバイナリ形式のUUIDに変換
	parsed, err := uuid.Parse(u.ID)
	if err != nil {
		return err
	}

	// pgtype.UUID に変換
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	params := dbgen.CreateUserParams{
		ID:         pgID,
		Email:      u.Email,
		Password:   u.EncryptedPassword,
		Timezone:   u.Timezone,
		ThemeColor: dbgen.ThemeColorEnum(u.ThemeColor), // dbgenで定義している列挙型に変換
		Language:   u.Language,
	}

	return r.q.CreateUser(ctx, params)
}

func (r *userRepository) FindByEmail(email string) (*userDomain.User, error) {
	ctx := context.Background()

	row, err := r.q.FindUserByEmail(ctx, email)
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

func (r *userRepository) GetSettingByID(userID string) (*userDomain.User, error) {
	ctx := context.Background()

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	row, err := r.q.GetUserSettingByID(ctx, pgID)
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

func (r *userRepository) Update(u *userDomain.User) error {
	ctx := context.Background()

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

	return r.q.UpdateUser(ctx, params)
}

func (r *userRepository) UpdatePassword(userID, password string) error {
	ctx := context.Background()

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsed, Valid: true}

	params := dbgen.UpdateUserPasswordParams{
		Password: password,
		ID:       pgID,
	}

	return r.q.UpdateUserPassword(ctx, params)
}
