package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	userDomain "github.com/minminseo/recall-setter/domain/user"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type emailVerificationRepository struct{}

func NewEmailVerificationRepository() userDomain.EmailVerificationRepository {
	return &emailVerificationRepository{}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *userDomain.EmailVerification) error {
	q := db.GetQuery(ctx)

	pgVerificationID, err := toUUID(ev.ID)
	if err != nil {
		return err
	}

	pgUserID, err := toUUID(ev.UserID)
	if err != nil {
		return err
	}

	params := dbgen.CreateEmailVerificationParams{
		ID:        pgVerificationID,
		UserID:    pgUserID,
		CodeHash:  ev.CodeHash,
		ExpiresAt: pgtype.Timestamptz{Time: ev.ExpiresAt, Valid: true},
	}

	return q.CreateEmailVerification(ctx, params)
}

func (r *emailVerificationRepository) FindByUserID(ctx context.Context, userID string) (*userDomain.EmailVerification, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}

	row, err := q.FindEmailVerificationByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	return &userDomain.EmailVerification{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		CodeHash:  row.CodeHash,
		ExpiresAt: row.ExpiresAt.Time,
	}, nil
}

func (r *emailVerificationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	return q.DeleteEmailVerificationByUserID(ctx, pgUserID)
}
