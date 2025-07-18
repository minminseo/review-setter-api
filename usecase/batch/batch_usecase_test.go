package batch

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockBatchRepository struct {
	mock.Mock
}

func (m *MockBatchRepository) ExecuteUpdateOverdueScheduledDates(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewBatchUsecase(t *testing.T) {
	tests := []struct {
		name string
		repo interface{}
		want bool // インスタンスが正常に生成されるかどうか
	}{
		{
			name: "正常なbatchRepositoryが渡された場合",
			repo: &MockBatchRepository{},
			want: true,
		},
		{
			name: "nilのbatchRepositoryが渡された場合",
			repo: nil,
			want: true, // 現在の実装ではnilでもインスタンスは生成される
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var usecase IBatchUsecase
			if tt.repo == nil {
				usecase = NewBatchUsecase(nil)
			} else {
				usecase = NewBatchUsecase(tt.repo.(*MockBatchRepository))
			}

			if tt.want {
				require.NotNil(t, usecase)
			} else {
				require.Nil(t, usecase)
			}
		})
	}
}

func TestBatchUsecase_ExecuteUpdateOverdueScheduledDates(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBatchRepository, context.Context)
		setupCtx    func() context.Context
		wantErr     bool
		wantErrType error
	}{
		{
			name: "リポジトリが正常に実行される場合",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(nil)
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
		{
			name: "リポジトリでエラーが発生する場合",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(errors.New("database connection failed"))
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "contextがキャンセルされた場合",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(context.Canceled)
			},
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr:     true,
			wantErrType: context.Canceled,
		},
		{
			name: "contextにタイムアウトが設定されている場合",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(context.DeadlineExceeded)
			},
			setupCtx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(2 * time.Millisecond)
				return ctx
			},
			wantErr:     true,
			wantErrType: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockBatchRepository{}
			usecase := NewBatchUsecase(mockRepo)
			ctx := tt.setupCtx()

			tt.setupMock(mockRepo, ctx)

			err := usecase.ExecuteUpdateOverdueScheduledDates(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrType != nil {
					require.Equal(t, tt.wantErrType, err)
				}
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBatchUsecase_ExecuteUpdateOverdueScheduledDates_LogOutput(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockBatchRepository, context.Context)
		expectedLogs   []string
		unexpectedLogs []string
	}{
		{
			name: "成功時のログ出力確認",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(nil)
			},
			expectedLogs: []string{
				"期限切れ復習日の更新処理を開始します",
				"未完了復習日の更新処理が正常に完了しました",
			},
			unexpectedLogs: []string{
				"未完了復習日の更新に失敗しました",
			},
		},
		{
			name: "エラー時のログ出力確認",
			setupMock: func(m *MockBatchRepository, ctx context.Context) {
				m.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(errors.New("update failed"))
			},
			expectedLogs: []string{
				"期限切れ復習日の更新処理を開始します",
				"未完了復習日の更新に失敗しました",
				"update failed",
			},
			unexpectedLogs: []string{
				"未完了復習日の更新処理が正常に完了しました",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			slog.SetDefault(logger)

			mockRepo := &MockBatchRepository{}
			usecase := NewBatchUsecase(mockRepo)
			ctx := context.Background()

			tt.setupMock(mockRepo, ctx)

			_ = usecase.ExecuteUpdateOverdueScheduledDates(ctx)

			logOutput := buf.String()

			for _, expectedLog := range tt.expectedLogs {
				require.Contains(t, logOutput, expectedLog)
			}

			for _, unexpectedLog := range tt.unexpectedLogs {
				require.NotContains(t, logOutput, unexpectedLog)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBatchUsecase_ExecuteUpdateOverdueScheduledDates_Idempotency(t *testing.T) {
	mockRepo := &MockBatchRepository{}
	usecase := NewBatchUsecase(mockRepo)
	ctx := context.Background()

	mockRepo.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(nil).Times(3)

	for i := 0; i < 3; i++ {
		err := usecase.ExecuteUpdateOverdueScheduledDates(ctx)
		require.NoError(t, err)
	}

	mockRepo.AssertExpectations(t)
}

func TestBatchUsecase_ExecuteUpdateOverdueScheduledDates_ErrorTypes(t *testing.T) {
	errorTypes := []struct {
		name string
		err  error
	}{
		{"DatabaseConnectionError", errors.New("database connection failed")},
		{"SQLExecutionError", errors.New("SQL execution failed")},
		{"TimeoutError", errors.New("operation timeout")},
		{"ContextCanceled", context.Canceled},
		{"ContextDeadlineExceeded", context.DeadlineExceeded},
	}

	for _, errorType := range errorTypes {
		t.Run(errorType.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockBatchRepository{}
			usecase := NewBatchUsecase(mockRepo)
			ctx := context.Background()

			mockRepo.On("ExecuteUpdateOverdueScheduledDates", ctx).Return(errorType.err)

			err := usecase.ExecuteUpdateOverdueScheduledDates(ctx)

			require.Error(t, err)
			if diff := cmp.Diff(errorType.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
