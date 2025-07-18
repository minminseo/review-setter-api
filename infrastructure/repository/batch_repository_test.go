package repository

import (
	"testing"
)

func TestBatchRepository_ExecuteUpdateOverdueScheduledDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "バッチ更新の実行に成功する場合",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBatchRepository()

			err := repo.ExecuteUpdateOverdueScheduledDates(ctx)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
		})
	}
}
