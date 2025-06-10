package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	patternDomain "github.com/minminseo/recall-setter/domain/pattern"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type patternRepository struct {
}

func NewPatternRepository() patternDomain.IPatternRepository {
	return &patternRepository{}
}

func (r *patternRepository) CreatePattern(ctx context.Context, p *patternDomain.Pattern) error {
	q := db.GetQuery(ctx)
	pgID, _ := toUUID(p.PatternID)
	pgUserID, _ := toUUID(p.UserID)
	pgReg := pgtype.Timestamptz{Time: p.RegisteredAt, Valid: true}
	pgEdit := pgtype.Timestamptz{Time: p.EditedAt, Valid: true}

	params := dbgen.CreatePatternParams{
		ID:           pgID,
		UserID:       pgUserID,
		Name:         p.Name,
		TargetWeight: dbgen.TargetWeightEnum(p.TargetWeight),
		RegisteredAt: pgReg,
		EditedAt:     pgEdit,
	}

	return q.CreatePattern(ctx, params)
}

func (r *patternRepository) CreatePatternSteps(ctx context.Context, steps []*patternDomain.PatternStep) (int64, error) {
	q := db.GetQuery(ctx)
	colums := []string{"id", "user_id", "pattern_id", "step_number", "interval_days"}
	cps := make([]dbgen.CreatePatternStepsParams, len(steps))
	rows := make([][]any, len(steps))
	for i, s := range steps {
		pgStepID, _ := toUUID(s.PatternStepID)
		pgUserID, _ := toUUID(s.UserID)
		pgPatternID, _ := toUUID(s.PatternID)
		cps[i] = dbgen.CreatePatternStepsParams{
			ID:           pgStepID,
			UserID:       pgUserID,
			PatternID:    pgPatternID,
			StepNumber:   int16(s.StepNumber),
			IntervalDays: int16(s.IntervalDays),
		}
		rows[i] = []any{
			cps[i].ID,
			cps[i].UserID,
			cps[i].PatternID,
			cps[i].StepNumber,
			cps[i].IntervalDays,
		}
	}
	return q.CopyFrom(
		ctx,
		pgx.Identifier{"pattern_steps"},
		colums,
		pgx.CopyFromRows(rows),
	)
}

func (r *patternRepository) GetAllPatternsByUserID(ctx context.Context, userID string) ([]*patternDomain.Pattern, error) {
	q := db.GetQuery(ctx)
	pgUserID, _ := toUUID(userID)

	rows, err := q.GetAllPatternsByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*patternDomain.Pattern, len(rows))
	for i, row := range rows {
		patternID := uuid.UUID(row.ID.Bytes).String()
		userID := uuid.UUID(row.UserID.Bytes).String()
		results[i], err = patternDomain.ReconstructPattern(
			patternID,
			userID,
			row.Name,
			string(row.TargetWeight),
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (r *patternRepository) GetAllPatternStepsByUserID(ctx context.Context, userID string) ([]*patternDomain.PatternStep, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.GetAllPatternStepsByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*patternDomain.PatternStep, len(rows))
	for i, row := range rows {
		patternStepID := uuid.UUID(row.ID.Bytes).String()
		userID := uuid.UUID(row.UserID.Bytes).String()
		patternID := uuid.UUID(row.PatternID.Bytes).String()
		results[i], err = patternDomain.ReconstructPatternStep(
			patternStepID,
			userID,
			patternID,
			int(row.StepNumber),
			int(row.IntervalDays),
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *patternRepository) UpdatePattern(ctx context.Context, p *patternDomain.Pattern) error {
	q := db.GetQuery(ctx)
	pgID, _ := toUUID(p.PatternID)
	pgUserID, _ := toUUID(p.UserID)
	pgEdit := pgtype.Timestamptz{Time: p.EditedAt, Valid: true}

	params := dbgen.UpdatePatternParams{
		Name:         p.Name,
		TargetWeight: dbgen.TargetWeightEnum(p.TargetWeight),
		EditedAt:     pgEdit,
		ID:           pgID,
		UserID:       pgUserID,
	}
	return q.UpdatePattern(ctx, params)
}

func (r *patternRepository) DeletePattern(ctx context.Context, patternID string, userID string) error {
	q := db.GetQuery(ctx)
	pgID, _ := toUUID(patternID)
	pgUserID, _ := toUUID(userID)
	patternParams := dbgen.DeletePatternParams{
		ID:     pgID,
		UserID: pgUserID,
	}

	return q.DeletePattern(ctx, patternParams)
}

func (r *patternRepository) DeletePatternSteps(ctx context.Context, patternID string, userID string) error {
	q := db.GetQuery(ctx)
	pgID, _ := toUUID(patternID)
	pgUserID, _ := toUUID(userID)
	params := dbgen.DeletePatternStepsParams{
		PatternID: pgID,
		UserID:    pgUserID,
	}
	return q.DeletePatternSteps(ctx, params)
}

// パターン単体取得
// パターン更新前の取得用
func (r *patternRepository) FindPatternByPatternID(ctx context.Context, patternID string, userID string) (*patternDomain.Pattern, error) {
	q := db.GetQuery(ctx)
	pgID, err := toUUID(patternID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}

	params := dbgen.GetPatternByIDParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	row, err := q.GetPatternByID(ctx, params)
	if err != nil {
		return nil, err
	}

	pd, err := patternDomain.ReconstructPattern(
		patternID,
		userID,
		row.Name,
		string(row.TargetWeight),
		row.RegisteredAt.Time,
		row.EditedAt.Time,
	)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

// パターンXが持つstepsを取得
// パターン更新前の取得用
func (r *patternRepository) GetAllPatternStepsByPatternID(ctx context.Context, patternID string, userID string) ([]*patternDomain.PatternStep, error) {
	q := db.GetQuery(ctx)
	pgID, err := toUUID(patternID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetPatternStepsByPatternIDParams{
		PatternID: pgID,
		UserID:    pgUserID,
	}

	rows, err := q.GetPatternStepsByPatternID(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]*patternDomain.PatternStep, len(rows))
	for i, row := range rows {
		PatternStepID := uuid.UUID(row.ID.Bytes).String()
		out[i], err = patternDomain.ReconstructPatternStep(
			PatternStepID,
			userID,
			patternID,
			int(row.StepNumber),
			int(row.IntervalDays),
		)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *patternRepository) GetPatternTargetWeightsByPatternIDs(ctx context.Context, ids []string) ([]*patternDomain.TargetWeight, error) {
	q := db.GetQuery(ctx)
	pgIDs := make([]pgtype.UUID, len(ids))
	for i, sid := range ids {
		u, err := uuid.Parse(sid)
		if err != nil {
			return nil, err
		}
		pgIDs[i] = pgtype.UUID{Bytes: u, Valid: true}
	}
	rows, err := q.GetPatternTargetWeightsByPatternIDs(ctx, pgIDs)
	if err != nil {
		return nil, err
	}
	out := make([]*patternDomain.TargetWeight, len(rows))
	for i, row := range rows {
		out[i] = &patternDomain.TargetWeight{
			PatternID:    uuid.UUID(row.ID.Bytes).String(),
			TargetWeight: string(row.TargetWeight),
		}
	}
	return out, nil
}
