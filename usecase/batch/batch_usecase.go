package batch

import (
	"context"
	"log/slog"

	"github.com/minminseo/recall-setter/infrastructure/repository"
)

type IBatchUsecase interface {
	ExecuteUpdateOverdueScheduledDates(ctx context.Context) error
}

type batchUsecase struct {
	batchRepo repository.IBatchRepository
}

func NewBatchUsecase(batchRepo repository.IBatchRepository) IBatchUsecase {
	return &batchUsecase{batchRepo: batchRepo}
}

func (u *batchUsecase) ExecuteUpdateOverdueScheduledDates(ctx context.Context) error {
	slog.Info("期限切れ復習日の更新処理を開始します。")

	err := u.batchRepo.ExecuteUpdateOverdueScheduledDates(ctx)
	if err != nil {
		slog.Error("未完了復習日の更新に失敗しました。", "error", err)
		return err
	}

	slog.Info("未完了復習日の更新処理が正常に完了しました。")
	return nil
}
