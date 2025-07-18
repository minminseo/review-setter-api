package pattern

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	itemDomain "github.com/minminseo/recall-setter/domain/item"
	patternDomain "github.com/minminseo/recall-setter/domain/pattern"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

func TestPatternUsecase_CreatePattern(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   CreatePatternInput
		setup   func(*patternDomain.MockIPatternRepository, *itemDomain.MockIItemRepository, *transaction.MockITransactionManager)
		want    *CreatePatternOutput
		wantErr bool
	}{
		{
			name: "正常系_単一ステップのパターン作成成功",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				Steps:        []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						CreatePattern(ctx, gomock.Any()).
						Return(nil).
						Times(1),
					patternRepo.EXPECT().
						CreatePatternSteps(ctx, gomock.Any()).
						Return(int64(1), nil).
						Times(1),
				)
			},
			want: &CreatePatternOutput{
				ID:           "",
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				RegisteredAt: fixedTime,
				EditedAt:     fixedTime,
				Steps: []CreatePatternStepOutput{
					{PatternStepID: "", UserID: "user-123", PatternID: "", StepNumber: 1, IntervalDays: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "正常系_複数ステップのパターン作成成功",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "複数ステップパターン",
				TargetWeight: "heavy",
				Steps: []CreatePatternStepInput{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
					{StepNumber: 3, IntervalDays: 7},
				},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						CreatePattern(ctx, gomock.Any()).
						Return(nil).
						Times(1),
					patternRepo.EXPECT().
						CreatePatternSteps(ctx, gomock.Any()).
						Return(int64(3), nil).
						Times(1),
				)
			},
			want: &CreatePatternOutput{
				ID:           "",
				UserID:       "user-123",
				Name:         "複数ステップパターン",
				TargetWeight: "heavy",
				RegisteredAt: fixedTime,
				EditedAt:     fixedTime,
				Steps: []CreatePatternStepOutput{
					{PatternStepID: "", UserID: "user-123", PatternID: "", StepNumber: 1, IntervalDays: 1},
					{PatternStepID: "", UserID: "user-123", PatternID: "", StepNumber: 2, IntervalDays: 3},
					{PatternStepID: "", UserID: "user-123", PatternID: "", StepNumber: 3, IntervalDays: 7},
				},
			},
			wantErr: false,
		},
		{
			name:  "異常系_TargetWeightが無効な値",
			input: CreatePatternInput{UserID: "user-123", Name: "テストパターン", TargetWeight: "invalid", Steps: []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}}},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
			},
			wantErr: true,
		},
		{
			name:  "異常系_Nameが空文字列",
			input: CreatePatternInput{UserID: "user-123", Name: "", TargetWeight: "light", Steps: []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}}},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
			},
			wantErr: true,
		},
		{
			name:  "異常系_Stepsが空",
			input: CreatePatternInput{UserID: "user-123", Name: "テストパターン", TargetWeight: "light", Steps: []CreatePatternStepInput{}},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
			},
			wantErr: true,
		},
		{
			name: "異常系_StepNumberが重複",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				Steps:        []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}, {StepNumber: 1, IntervalDays: 3}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
			},
			wantErr: true,
		},
		{
			name: "異常系_CreatePatternでエラー",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				Steps:        []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						CreatePattern(ctx, gomock.Any()).
						Return(errors.New("データベースエラー")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "異常系_CreatePatternStepsでエラー",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				Steps:        []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						CreatePattern(ctx, gomock.Any()).
						Return(nil).
						Times(1),
					patternRepo.EXPECT().
						CreatePatternSteps(ctx, gomock.Any()).
						Return(int64(0), errors.New("ステップ作成エラー")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "異常系_トランザクション全体でエラー",
			input: CreatePatternInput{
				UserID:       "user-123",
				Name:         "テストパターン",
				TargetWeight: "light",
				Steps:        []CreatePatternStepInput{{StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						Return(errors.New("トランザクションエラー")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			patternRepo := patternDomain.NewMockIPatternRepository(ctrl)
			itemRepo := itemDomain.NewMockIItemRepository(ctrl)
			txManager := transaction.NewMockITransactionManager(ctrl)

			tt.setup(patternRepo, itemRepo, txManager)

			uc := NewPatternUsecase(patternRepo, itemRepo, txManager)
			got, err := uc.CreatePattern(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreatePattern() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("CreatePattern() unexpected error = %v", err)
				return
			}
			if got == nil {
				t.Error("CreatePattern() got = nil, want not nil")
				return
			}

			got.ID = ""
			for i := range got.Steps {
				got.Steps[i].PatternStepID = ""
				got.Steps[i].PatternID = ""
			}
			got.RegisteredAt = fixedTime
			got.EditedAt = fixedTime

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CreatePattern() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternUsecase_GetPatternsByUserID(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		userID  string
		setup   func(*patternDomain.MockIPatternRepository, *itemDomain.MockIItemRepository, *transaction.MockITransactionManager)
		want    []*GetPatternOutput
		wantErr bool
	}{
		{
			name:   "正常系_パターンとステップが存在する場合",
			userID: "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				patterns := []*patternDomain.Pattern{{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "パターン1",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}}
				steps := []*patternDomain.PatternStep{
					{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1},
					{PatternStepID: "step-2", PatternID: "pattern-1", StepNumber: 2, IntervalDays: 3},
				}
				gomock.InOrder(
					patternRepo.EXPECT().
						GetAllPatternsByUserID(ctx, "user-123").
						Return(patterns, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByUserID(ctx, "user-123").
						Return(steps, nil).
						Times(1),
				)
			},
			want: []*GetPatternOutput{{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "パターン1",
				TargetWeight: "light",
				RegisteredAt: fixedTime,
				EditedAt:     fixedTime,
				Steps: []GetPatternStepOutput{
					{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1},
					{PatternStepID: "step-2", PatternID: "pattern-1", StepNumber: 2, IntervalDays: 3},
				},
			}},
		},
		{
			name:   "正常系_パターンは存在するがステップが存在しない",
			userID: "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				patterns := []*patternDomain.Pattern{{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "パターン1",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}}
				steps := []*patternDomain.PatternStep{}
				gomock.InOrder(
					patternRepo.EXPECT().
						GetAllPatternsByUserID(ctx, "user-123").
						Return(patterns, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByUserID(ctx, "user-123").
						Return(steps, nil).
						Times(1),
				)
			},
			want: []*GetPatternOutput{{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "パターン1",
				TargetWeight: "light",
				RegisteredAt: fixedTime,
				EditedAt:     fixedTime,
				Steps:        nil,
			}},
		},
		{
			name:   "正常系_パターンもステップも存在しない",
			userID: "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					patternRepo.EXPECT().
						GetAllPatternsByUserID(ctx, "user-123").
						Return([]*patternDomain.Pattern{}, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByUserID(ctx, "user-123").
						Return([]*patternDomain.PatternStep{}, nil).
						Times(1),
				)
			},
			want: []*GetPatternOutput{},
		},
		{
			name:   "異常系_GetAllPatternsByUserIDでエラー",
			userID: "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				patternRepo.EXPECT().
					GetAllPatternsByUserID(ctx, "user-123").
					Return(nil, errors.New("データベースエラー")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:   "異常系_GetAllPatternStepsByUserIDでエラー",
			userID: "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				patterns := []*patternDomain.Pattern{{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "パターン1",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}}
				gomock.InOrder(
					patternRepo.EXPECT().
						GetAllPatternsByUserID(ctx, "user-123").
						Return(patterns, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByUserID(ctx, "user-123").
						Return(nil, errors.New("ステップ取得エラー")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			patternRepo := patternDomain.NewMockIPatternRepository(ctrl)
			itemRepo := itemDomain.NewMockIItemRepository(ctrl)
			txManager := transaction.NewMockITransactionManager(ctrl)

			tt.setup(patternRepo, itemRepo, txManager)

			uc := NewPatternUsecase(patternRepo, itemRepo, txManager)
			got, err := uc.GetPatternsByUserID(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPatternsByUserID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("GetPatternsByUserID() unexpected error = %v", err)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetPatternsByUserID() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternUsecase_UpdatePattern(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	editedTime := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   UpdatePatternInput
		setup   func(*patternDomain.MockIPatternRepository, *itemDomain.MockIItemRepository, *transaction.MockITransactionManager)
		want    *UpdatePatternOutput
		wantErr bool
	}{
		{
			name: "正常系_パターンのみ更新成功",
			input: UpdatePatternInput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "更新されたパターン",
				TargetWeight: "heavy",
				Steps:        []UpdatePatternStepInput{{StepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				pattern := &patternDomain.Pattern{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "元のパターン",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}
				steps := []*patternDomain.PatternStep{{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}}
				gomock.InOrder(
					patternRepo.EXPECT().
						FindPatternByPatternID(ctx, "pattern-1", "user-123").
						Return(pattern, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByPatternID(ctx, "pattern-1", "user-123").
						Return(steps, nil).
						Times(1),
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						UpdatePattern(ctx, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdatePatternOutput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "更新されたパターン",
				TargetWeight: "heavy",
				RegisteredAt: fixedTime,
				EditedAt:     editedTime,
				Steps:        []UpdatePatternStepOutput{},
			},
		},
		{
			name: "正常系_ステップのみ更新成功",
			input: UpdatePatternInput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "元のパターン",
				TargetWeight: "light",
				Steps:        []UpdatePatternStepInput{{StepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 2}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				pattern := &patternDomain.Pattern{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "元のパターン",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}
				steps := []*patternDomain.PatternStep{{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}}
				gomock.InOrder(
					patternRepo.EXPECT().
						FindPatternByPatternID(ctx, "pattern-1", "user-123").
						Return(pattern, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByPatternID(ctx, "pattern-1", "user-123").
						Return(steps, nil).
						Times(1),
					itemRepo.EXPECT().
						IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
						Return(false, nil).
						Times(1),
					txManager.EXPECT().
						RunInTransaction(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					patternRepo.EXPECT().
						DeletePatternSteps(ctx, "pattern-1", "user-123").
						Return(nil).
						Times(1),
					patternRepo.EXPECT().
						CreatePatternSteps(ctx, gomock.Any()).
						Return(int64(1), nil).
						Times(1),
				)
			},
			want: &UpdatePatternOutput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "元のパターン",
				TargetWeight: "light",
				RegisteredAt: fixedTime,
				EditedAt:     fixedTime,
				Steps: []UpdatePatternStepOutput{
					{PatternStepID: "", UserID: "user-123", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 2},
				},
			},
			wantErr: false,
		},
		{
			name:  "異常系_パターンが存在しない",
			input: UpdatePatternInput{PatternID: "pattern-1", UserID: "user-123", Name: "更新されたパターン", TargetWeight: "heavy", Steps: []UpdatePatternStepInput{}},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				patternRepo.EXPECT().
					FindPatternByPatternID(ctx, "pattern-1", "user-123").
					Return(nil, patternDomain.ErrPatternNotFound).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "異常系_変更がない場合",
			input: UpdatePatternInput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "元のパターン",
				TargetWeight: "light",
				Steps:        []UpdatePatternStepInput{{StepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				pattern := &patternDomain.Pattern{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "元のパターン",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}
				steps := []*patternDomain.PatternStep{{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}}
				gomock.InOrder(
					patternRepo.EXPECT().
						FindPatternByPatternID(ctx, "pattern-1", "user-123").
						Return(pattern, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByPatternID(ctx, "pattern-1", "user-123").
						Return(steps, nil).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "異常系_ステップ変更時に復習物関連がある",
			input: UpdatePatternInput{
				PatternID:    "pattern-1",
				UserID:       "user-123",
				Name:         "元のパターン",
				TargetWeight: "light",
				Steps:        []UpdatePatternStepInput{{StepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 2}},
			},
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				pattern := &patternDomain.Pattern{
					PatternID:    "pattern-1",
					UserID:       "user-123",
					Name:         "元のパターン",
					TargetWeight: "light",
					RegisteredAt: fixedTime,
					EditedAt:     fixedTime,
				}
				steps := []*patternDomain.PatternStep{{PatternStepID: "step-1", PatternID: "pattern-1", StepNumber: 1, IntervalDays: 1}}
				gomock.InOrder(
					patternRepo.EXPECT().
						FindPatternByPatternID(ctx, "pattern-1", "user-123").
						Return(pattern, nil).
						Times(1),
					patternRepo.EXPECT().
						GetAllPatternStepsByPatternID(ctx, "pattern-1", "user-123").
						Return(steps, nil).
						Times(1),
					itemRepo.EXPECT().
						IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
						Return(true, nil).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			patternRepo := patternDomain.NewMockIPatternRepository(ctrl)
			itemRepo := itemDomain.NewMockIItemRepository(ctrl)
			txManager := transaction.NewMockITransactionManager(ctrl)

			tt.setup(patternRepo, itemRepo, txManager)

			uc := NewPatternUsecase(patternRepo, itemRepo, txManager)
			got, err := uc.UpdatePattern(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePattern() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("UpdatePattern() unexpected error = %v", err)
				return
			}

			if got != nil {
				for i := range got.Steps {
					got.Steps[i].PatternStepID = ""
				}
				if got.EditedAt.After(tt.want.EditedAt) {
					got.EditedAt = editedTime
				}
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("UpdatePattern() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternUsecase_DeletePattern(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		patternID string
		userID    string
		setup     func(*patternDomain.MockIPatternRepository, *itemDomain.MockIItemRepository, *transaction.MockITransactionManager)
		wantErr   bool
	}{
		{
			name:      "正常系_復習物関連なしでの削除成功",
			patternID: "pattern-1",
			userID:    "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					itemRepo.EXPECT().
						IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
						Return(false, nil).
						Times(1),
					patternRepo.EXPECT().
						DeletePattern(ctx, "pattern-1", "user-123").
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name:      "異常系_復習物関連がある場合",
			patternID: "pattern-1",
			userID:    "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				itemRepo.EXPECT().
					IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
					Return(true, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:      "異常系_IsPatternRelatedToItemByPatternIDでエラー",
			patternID: "pattern-1",
			userID:    "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				itemRepo.EXPECT().
					IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
					Return(false, errors.New("データベースエラー")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:      "異常系_DeletePatternでエラー",
			patternID: "pattern-1",
			userID:    "user-123",
			setup: func(patternRepo *patternDomain.MockIPatternRepository, itemRepo *itemDomain.MockIItemRepository, txManager *transaction.MockITransactionManager) {
				gomock.InOrder(
					itemRepo.EXPECT().
						IsPatternRelatedToItemByPatternID(ctx, "pattern-1", "user-123").
						Return(false, nil).
						Times(1),
					patternRepo.EXPECT().
						DeletePattern(ctx, "pattern-1", "user-123").
						Return(errors.New("削除エラー")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			patternRepo := patternDomain.NewMockIPatternRepository(ctrl)
			itemRepo := itemDomain.NewMockIItemRepository(ctrl)
			txManager := transaction.NewMockITransactionManager(ctrl)

			tt.setup(patternRepo, itemRepo, txManager)

			uc := NewPatternUsecase(patternRepo, itemRepo, txManager)
			err := uc.DeletePattern(ctx, tt.patternID, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeletePattern() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("DeletePattern() unexpected error = %v", err)
			}
		})
	}
}

func TestNewPatternUsecase(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patternRepo := patternDomain.NewMockIPatternRepository(ctrl)
	itemRepo := itemDomain.NewMockIItemRepository(ctrl)
	txManager := transaction.NewMockITransactionManager(ctrl)

	uc := NewPatternUsecase(patternRepo, itemRepo, txManager)
	if uc == nil {
		t.Error("NewPatternUsecase() returned nil")
	}
	if _, ok := uc.(*patternUsecase); !ok {
		t.Error("NewPatternUsecase() did not return *patternUsecase")
	}
}
