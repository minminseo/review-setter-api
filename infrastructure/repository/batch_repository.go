package repository

import (
	"context"

	"github.com/minminseo/recall-setter/infrastructure/db"
)

type IBatchRepository interface {
	ExecuteUpdateOverdueScheduledDates(ctx context.Context) error
}

type batchRepository struct{}

func NewBatchRepository() IBatchRepository {
	return &batchRepository{}
}

func (r *batchRepository) ExecuteUpdateOverdueScheduledDates(ctx context.Context) error {
	q := db.GetQuery(ctx)
	return q.UpdateOverdueScheduledDatesAndSlideFutureDates(ctx)
}
