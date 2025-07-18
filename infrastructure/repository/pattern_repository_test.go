package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	patternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

func TestPatternRepository_CreatePattern(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		pattern *patternDomain.Pattern
		want    *patternDomain.Pattern
		wantErr bool
	}{
		{
			name: "パターン作成に成功する場合",
			pattern: &patternDomain.Pattern{
				PatternID:    uuid.New().String(),
				UserID:       "550e8400-e29b-41d4-a716-446655440001", // Exists in fixture
				Name:         "新しいパターン",
				TargetWeight: "normal",
				RegisteredAt: time.Now(),
				EditedAt:     time.Now(),
			},
			want: &patternDomain.Pattern{
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "新しいパターン",
				TargetWeight: "normal",
			},
			wantErr: false,
		},
		{
			name: "存在しないユーザーによる外部キー制約違反",
			pattern: &patternDomain.Pattern{
				PatternID:    uuid.New().String(),
				UserID:       uuid.New().String(), // Does not exist in fixture
				Name:         "存在しないユーザーパターン",
				TargetWeight: "normal",
				RegisteredAt: time.Now(),
				EditedAt:     time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "無効な重みで作成する場合",
			pattern: &patternDomain.Pattern{
				PatternID:    uuid.New().String(),
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "無効な重みパターン",
				TargetWeight: "invalid_weight", // Invalid enum value
				RegisteredAt: time.Now(),
				EditedAt:     time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			err := repo.CreatePattern(ctx, tc.pattern)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 作成されたパターンを取得して検証
			createdPattern, err := repo.FindPatternByPatternID(ctx, tc.pattern.PatternID, tc.pattern.UserID)
			if err != nil {
				t.Errorf("created pattern retrieval failed: %v", err)
				return
			}

			if createdPattern == nil {
				t.Error("created pattern is nil")
				return
			}

			// 動的に生成されるフィールドを期待値に設定
			tc.want.PatternID = createdPattern.PatternID
			tc.want.RegisteredAt = createdPattern.RegisteredAt
			tc.want.EditedAt = createdPattern.EditedAt

			// 期待値との比較
			if diff := cmp.Diff(tc.want, createdPattern); diff != "" {
				t.Errorf("CreatePattern() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_GetAllPatternsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		userID        string
		want          []patternDomain.Pattern
		wantErr       bool
		expectedCount int
	}{
		{
			name:   "ユーザー1のパターンを取得（3件）",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []patternDomain.Pattern{
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440001",
					UserID:       "550e8400-e29b-41d4-a716-446655440001",
					Name:         "フィボナッチパターン",
					TargetWeight: "normal",
					RegisteredAt: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
				},
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440002",
					UserID:       "550e8400-e29b-41d4-a716-446655440001",
					Name:         "エビングハウスパターン",
					TargetWeight: "heavy",
					RegisteredAt: time.Date(2024, 1, 1, 8, 30, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 8, 30, 0, 0, time.UTC),
				},
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440005",
					UserID:       "550e8400-e29b-41d4-a716-446655440001",
					Name:         "ステップ未作成のパターン",
					TargetWeight: "light",
					RegisteredAt: time.Date(2024, 1, 1, 9, 00, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 9, 00, 0, 0, time.UTC),
				},
			},
			wantErr:       false,
			expectedCount: 3,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			patterns, err := repo.GetAllPatternsByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if patterns != nil {
					t.Error("パターンがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if patterns == nil {
				t.Error("パターンのスライスがnilです")
				return
			}

			if len(patterns) != tc.expectedCount {
				t.Errorf("期待されるパターン数: %d, 実際: %d", tc.expectedCount, len(patterns))
			}

			// パターンのスライスを作成してポインタを外す
			patternsSlice := make([]patternDomain.Pattern, len(patterns))
			for i, pattern := range patterns {
				patternsSlice[i] = *pattern
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, patternsSlice); diff != "" {
				t.Errorf("GetAllPatternsByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_UpdatePattern(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		pattern *patternDomain.Pattern
		want    *patternDomain.Pattern
		wantErr bool
	}{
		{
			name: "パターン更新に成功する場合",
			pattern: &patternDomain.Pattern{
				PatternID:    "750e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "更新されたフィボナッチパターン",
				TargetWeight: "heavy",
				RegisteredAt: time.Now().Add(-24 * time.Hour),
				EditedAt:     time.Now(),
			},
			want: &patternDomain.Pattern{
				PatternID:    "750e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "更新されたフィボナッチパターン",
				TargetWeight: "heavy",
				RegisteredAt: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "無効な重みで更新する場合",
			pattern: &patternDomain.Pattern{
				PatternID:    "750e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "パターン",
				TargetWeight: "invalid_weight",
				RegisteredAt: time.Now().Add(-24 * time.Hour),
				EditedAt:     time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			err := repo.UpdatePattern(ctx, tc.pattern)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新されたパターンを取得して検証
				updatedPattern, err := repo.FindPatternByPatternID(ctx, tc.pattern.PatternID, tc.pattern.UserID)
				if err != nil {
					t.Errorf("更新されたパターンの取得に失敗: %v", err)
					return
				}

				if updatedPattern == nil {
					t.Error("更新されたパターンがnilです")
					return
				}

				// 動的に変更されるフィールドを期待値に設定
				tc.want.EditedAt = updatedPattern.EditedAt

				// 期待値との比較
				if diff := cmp.Diff(tc.want, updatedPattern); diff != "" {
					t.Errorf("UpdatePattern() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPatternRepository_DeletePattern(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		patternID string
		userID    string
		want      *patternDomain.Pattern
		wantErr   bool
	}{
		{
			name:      "パターン削除に成功する場合",
			patternID: "750e8400-e29b-41d4-a716-446655440004",
			userID:    "550e8400-e29b-41d4-a716-446655440002",
			want:      nil, // 削除されたので取得できない
			wantErr:   false,
		},
		{
			name:      "外部キー制約違反でパターン削除に失敗する場合",
			patternID: "750e8400-e29b-41d4-a716-446655440003",
			userID:    "550e8400-e29b-41d4-a716-446655440002",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			err := repo.DeletePattern(ctx, tc.patternID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 削除後に本当に削除されたかを確認
			deletedPattern, _ := repo.FindPatternByPatternID(ctx, tc.patternID, tc.userID)
			// 削除された場合はエラーが発生するか、nilが返される

			// 期待値との比較
			if diff := cmp.Diff(tc.want, deletedPattern); diff != "" {
				t.Errorf("DeletePattern() verification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_FindPatternByPatternID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name         string
		patternID    string
		userID       string
		want         *patternDomain.Pattern
		wantErr      bool
		expectName   string
		expectWeight string
	}{
		{
			name:      "有効なIDでパターンを取得する場合",
			patternID: "750e8400-e29b-41d4-a716-446655440001",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want: &patternDomain.Pattern{
				PatternID:    "750e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "フィボナッチパターン",
				TargetWeight: "normal",
				RegisteredAt: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
				EditedAt:     time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			},
			wantErr:      false,
			expectName:   "フィボナッチパターン",
			expectWeight: "normal",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			pattern, err := repo.FindPatternByPatternID(ctx, tc.patternID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if pattern != nil {
					t.Error("パターンがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if pattern == nil {
				t.Error("パターンがnilです")
				return
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, pattern); diff != "" {
				t.Errorf("FindPatternByPatternID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_CreatePatternSteps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		steps   []*patternDomain.PatternStep
		want    []*patternDomain.PatternStep
		wantErr bool
	}{
		{
			name: "パターンステップ作成に成功する場合",
			steps: []*patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440100",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440005",
					StepNumber:    1,
					IntervalDays:  1,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440101",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440005",
					StepNumber:    2,
					IntervalDays:  3,
				},
			},
			want: []*patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440100",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440005",
					StepNumber:    1,
					IntervalDays:  1,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440101",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440005",
					StepNumber:    2,
					IntervalDays:  3,
				},
			},
			wantErr: false,
		},
		{
			name: "存在しないパターンIDで作成失敗する場合",
			steps: []*patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440102",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440999", // 存在しないパターンID
					StepNumber:    1,
					IntervalDays:  1,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			_, err := repo.CreatePatternSteps(ctx, tc.steps)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 作成されたパターンステップを取得して検証
				createdSteps, err := repo.GetAllPatternStepsByPatternID(ctx, tc.steps[0].PatternID, tc.steps[0].UserID)
				if err != nil {
					t.Errorf("作成されたパターンステップの取得に失敗: %v", err)
					return
				}

				// パターンステップのスライスを作成してポインタを外す
				stepsSlice := make([]patternDomain.PatternStep, len(createdSteps))
				for i, step := range createdSteps {
					stepsSlice[i] = *step
				}

				// 期待値のスライスを作成
				wantSlice := make([]patternDomain.PatternStep, len(tc.want))
				for i, step := range tc.want {
					wantSlice[i] = *step
				}

				// 期待値との比較
				if diff := cmp.Diff(wantSlice, stepsSlice); diff != "" {
					t.Errorf("CreatePatternSteps() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPatternRepository_GetAllPatternStepsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		userID        string
		want          []patternDomain.PatternStep
		wantErr       bool
		expectedCount int
	}{
		{
			name:   "ユーザー1のパターンステップを取得（5件）",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440001",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    1,
					IntervalDays:  1,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440002",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    2,
					IntervalDays:  2,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440003",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    3,
					IntervalDays:  3,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440004",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440002",
					StepNumber:    1,
					IntervalDays:  1,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440005",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440002",
					StepNumber:    2,
					IntervalDays:  5,
				},
			},
			wantErr:       false,
			expectedCount: 5,
		},
		{
			name:   "ユーザー2のパターンステップを取得（2件）",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: []patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440006",
					UserID:        "550e8400-e29b-41d4-a716-446655440002",
					PatternID:     "750e8400-e29b-41d4-a716-446655440003",
					StepNumber:    1,
					IntervalDays:  7,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440007",
					UserID:        "550e8400-e29b-41d4-a716-446655440002",
					PatternID:     "750e8400-e29b-41d4-a716-446655440004",
					StepNumber:    1,
					IntervalDays:  3,
				},
			},
			wantErr:       false,
			expectedCount: 2,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			steps, err := repo.GetAllPatternStepsByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if steps != nil {
					t.Error("パターンステップがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if steps == nil {
				t.Error("パターンステップのスライスがnilです")
				return
			}

			if len(steps) != tc.expectedCount {
				t.Errorf("期待されるパターンステップ数: %d, 実際: %d", tc.expectedCount, len(steps))
			}

			// パターンステップのスライスを作成してポインタを外す
			stepsSlice := make([]patternDomain.PatternStep, len(steps))
			for i, step := range steps {
				stepsSlice[i] = *step
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, stepsSlice); diff != "" {
				t.Errorf("GetAllPatternStepsByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_DeletePatternSteps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		patternID string
		userID    string
		want      []patternDomain.PatternStep
		wantErr   bool
	}{
		{
			name:      "パターンステップ削除に成功する場合",
			patternID: "750e8400-e29b-41d4-a716-446655440001",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want:      []patternDomain.PatternStep{}, // 削除されたので空のスライス
			wantErr:   false,
		},
		{
			name:      "存在しないパターンIDで削除する場合",
			patternID: "750e8400-e29b-41d4-a716-446655440999",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want:      []patternDomain.PatternStep{},
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			err := repo.DeletePatternSteps(ctx, tc.patternID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 削除後に本当に削除されたかを確認
			steps, err := repo.GetAllPatternStepsByPatternID(ctx, tc.patternID, tc.userID)
			if err != nil {
				t.Errorf("削除後の確認でエラー: %v", err)
				return
			}

			// パターンステップのスライスを作成してポインタを外す
			stepsSlice := make([]patternDomain.PatternStep, len(steps))
			for i, step := range steps {
				stepsSlice[i] = *step
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, stepsSlice); diff != "" {
				t.Errorf("DeletePatternSteps() verification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_GetAllPatternStepsByPatternID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		patternID string
		userID    string
		want      []patternDomain.PatternStep
		wantErr   bool
	}{
		{
			name:      "パターンID1のステップを取得（3件）",
			patternID: "750e8400-e29b-41d4-a716-446655440001",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want: []patternDomain.PatternStep{
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440001",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    1,
					IntervalDays:  1,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440002",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    2,
					IntervalDays:  2,
				},
				{
					PatternStepID: "850e8400-e29b-41d4-a716-446655440003",
					UserID:        "550e8400-e29b-41d4-a716-446655440001",
					PatternID:     "750e8400-e29b-41d4-a716-446655440001",
					StepNumber:    3,
					IntervalDays:  3,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			steps, err := repo.GetAllPatternStepsByPatternID(ctx, tc.patternID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// パターンステップのスライスを作成してポインタを外す
			stepsSlice := make([]patternDomain.PatternStep, len(steps))
			for i, step := range steps {
				stepsSlice[i] = *step
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, stepsSlice); diff != "" {
				t.Errorf("GetAllPatternStepsByPatternID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatternRepository_GetPatternTargetWeightsByPatternIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		patternIDs []string
		want       []patternDomain.TargetWeight
		wantErr    bool
	}{
		{
			name: "複数のパターンIDから重みを取得する場合",
			patternIDs: []string{
				"750e8400-e29b-41d4-a716-446655440001",
				"750e8400-e29b-41d4-a716-446655440002",
				"750e8400-e29b-41d4-a716-446655440003",
			},
			want: []patternDomain.TargetWeight{
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440001",
					TargetWeight: "normal",
				},
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440002",
					TargetWeight: "heavy",
				},
				{
					PatternID:    "750e8400-e29b-41d4-a716-446655440003",
					TargetWeight: "light",
				},
			},
			wantErr: false,
		},
		{
			name:       "存在しないパターンIDで取得する場合",
			patternIDs: []string{"750e8400-e29b-41d4-a716-446655440999"},
			want:       []patternDomain.TargetWeight{},
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewPatternRepository()

			weights, err := repo.GetPatternTargetWeightsByPatternIDs(ctx, tc.patternIDs)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// TargetWeightのスライスを作成してポインタを外す
			weightsSlice := make([]patternDomain.TargetWeight, len(weights))
			for i, weight := range weights {
				weightsSlice[i] = *weight
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, weightsSlice); diff != "" {
				t.Errorf("GetPatternTargetWeightsByPatternIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
