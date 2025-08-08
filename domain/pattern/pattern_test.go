package pattern

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var Debug = false

const (
	testUserID     = "user1"
	testCategoryID = "category1"
	testPatternID  = "pattern1"
	testBoxID      = "box1"
)

func TestNewPattern(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		patternID    string
		userID       string
		patternName  string
		targetWeight string
		registeredAt time.Time
		editedAt     time.Time
		want         *Pattern
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "有効なパターン（正常系）",
			patternID:    testPatternID,
			userID:       testUserID,
			patternName:  "Standard Review",
			targetWeight: TargetWeightNormal,
			registeredAt: now,
			editedAt:     now,
			want: func() *Pattern {
				pattern, _ := ReconstructPattern(
					testPatternID,
					testUserID,
					"Standard Review",
					TargetWeightNormal,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: false,
		},
		{
			name:         "パターン名が空（異常系）",
			patternID:    "pattern2",
			userID:       testUserID,
			patternName:  "",
			targetWeight: TargetWeightNormal,
			registeredAt: now,
			editedAt:     now,
			want:         nil,
			wantErr:      true,
			errMsg:       "名前は必須です",
		},
		{
			name:         "重みが不正（異常系）",
			patternID:    "pattern3",
			userID:       testUserID,
			patternName:  "Test Pattern",
			targetWeight: "invalid",
			registeredAt: now,
			editedAt:     now,
			want:         nil,
			wantErr:      true,
			errMsg:       "重みの値が不正です",
		},
		{
			name:         "重みがHeavy（正常系）",
			patternID:    "pattern4",
			userID:       testUserID,
			patternName:  "Heavy Pattern",
			targetWeight: TargetWeightHeavy,
			registeredAt: now,
			editedAt:     now,
			want: func() *Pattern {
				pattern, _ := ReconstructPattern(
					"pattern4",
					testUserID,
					"Heavy Pattern",
					TargetWeightHeavy,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: false,
		},
		{
			name:         "重みがLight（正常系）",
			patternID:    "pattern5",
			userID:       testUserID,
			patternName:  "Light Pattern",
			targetWeight: TargetWeightLight,
			registeredAt: now,
			editedAt:     now,
			want: func() *Pattern {
				pattern, _ := ReconstructPattern(
					"pattern5",
					testUserID,
					"Light Pattern",
					TargetWeightLight,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: false,
		},
		{
			name:         "重みがUnset（正常系）",
			patternID:    "pattern6",
			userID:       testUserID,
			patternName:  "Unset Pattern",
			targetWeight: TargetWeightUnset,
			registeredAt: now,
			editedAt:     now,
			want: func() *Pattern {
				pattern, _ := ReconstructPattern(
					"pattern6",
					testUserID,
					"Unset Pattern",
					TargetWeightUnset,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pattern, err := NewPattern(tc.patternID, tc.userID, tc.patternName, tc.targetWeight, tc.registeredAt, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("unexpected error message: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// フィールドごとの比較
			if pattern.PatternID() != tc.want.PatternID() {
				t.Errorf("PatternID: got %v, want %v", pattern.PatternID(), tc.want.PatternID())
			}
			if pattern.UserID() != tc.want.UserID() {
				t.Errorf("UserID: got %v, want %v", pattern.UserID(), tc.want.UserID())
			}
			if pattern.Name() != tc.want.Name() {
				t.Errorf("Name: got %v, want %v", pattern.Name(), tc.want.Name())
			}
			if pattern.TargetWeight() != tc.want.TargetWeight() {
				t.Errorf("TargetWeight: got %v, want %v", pattern.TargetWeight(), tc.want.TargetWeight())
			}
			if !pattern.RegisteredAt().Equal(tc.want.RegisteredAt()) {
				t.Errorf("RegisteredAt: got %v, want %v", pattern.RegisteredAt(), tc.want.RegisteredAt())
			}
			if !pattern.EditedAt().Equal(tc.want.EditedAt()) {
				t.Errorf("EditedAt: got %v, want %v", pattern.EditedAt(), tc.want.EditedAt())
			}
		})
	}
}

func TestPattern_UpdatePattern(t *testing.T) {
	now := time.Now()
	pattern, err := NewPattern(testPatternID, testUserID, "Original", TargetWeightNormal, now, now)
	if err != nil {
		t.Fatalf("failed to create pattern: %v", err)
	}

	newTime := now.Add(time.Hour)

	tests := []struct {
		name         string
		newName      string
		targetWeight string
		editedAt     time.Time
		wantPattern  *Pattern
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "全項目を更新（正常系）",
			newName:      "Updated Pattern",
			targetWeight: TargetWeightHeavy,
			editedAt:     newTime,
			wantPattern: func() *Pattern {
				pattern, _ := ReconstructPattern(
					testPatternID,
					testUserID,
					"Updated Pattern",
					TargetWeightHeavy,
					now,
					newTime,
				)
				return pattern
			}(),
			wantErr: false,
		},
		{
			name:         "パターン名が空（異常系）",
			newName:      "",
			targetWeight: TargetWeightNormal,
			editedAt:     newTime,
			wantPattern: func() *Pattern {
				pattern, _ := ReconstructPattern(
					testPatternID,
					testUserID,
					"Original",
					TargetWeightNormal,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: true,
			errMsg:  "名前は必須です",
		},
		{
			name:         "重みが不正（異常系）",
			newName:      "Valid Name",
			targetWeight: "invalid",
			editedAt:     newTime,
			wantPattern: func() *Pattern {
				pattern, _ := ReconstructPattern(
					testPatternID,
					testUserID,
					"Original",
					TargetWeightNormal,
					now,
					now,
				)
				return pattern
			}(),
			wantErr: true,
			errMsg:  "重みの値が不正です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// パターンをコピー
			testPattern := *pattern

			err := testPattern.UpdatePattern(tc.newName, tc.targetWeight, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("unexpected error message: got %q, want %q", err.Error(), tc.errMsg)
				}
				if diff := cmp.Diff(tc.wantPattern, &testPattern, cmp.AllowUnexported(Pattern{})); diff != "" {
					t.Errorf("Pattern mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantPattern, &testPattern, cmp.AllowUnexported(Pattern{})); diff != "" {
				t.Errorf("Pattern mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewPatternStep(t *testing.T) {
	tests := []struct {
		name          string
		patternStepID string
		userID        string
		patternID     string
		stepNumber    int
		intervalDays  int
		want          *PatternStep
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "有効なパターンステップ（正常系）",
			patternStepID: "step1",
			userID:        testUserID,
			patternID:     testPatternID,
			stepNumber:    1,
			intervalDays:  1,
			want: func() *PatternStep {
				step, _ := ReconstructPatternStep(
					"step1",
					testUserID,
					testPatternID,
					1,
					1,
				)
				return step
			}(),
			wantErr: false,
		},
		{
			name:          "順序番号が0（異常系）",
			patternStepID: "step2",
			userID:        testUserID,
			patternID:     testPatternID,
			stepNumber:    0,
			intervalDays:  1,
			want:          nil,
			wantErr:       true,
			errMsg:        "順序番号は必須です",
		},
		{
			name:          "復習日間隔数が0（異常系）",
			patternStepID: "step3",
			userID:        testUserID,
			patternID:     testPatternID,
			stepNumber:    1,
			intervalDays:  0,
			want:          nil,
			wantErr:       true,
			errMsg:        "復習日間隔数は必須です",
		},
		{
			name:          "順序番号が負数（異常系）",
			patternStepID: "step4",
			userID:        testUserID,
			patternID:     testPatternID,
			stepNumber:    -1,
			intervalDays:  1,
			want:          nil,
			wantErr:       true,
			errMsg:        "順序番号の値が不正です",
		},
		{
			name:          "復習日間隔数が負数（異常系）",
			patternStepID: "step5",
			userID:        testUserID,
			patternID:     testPatternID,
			stepNumber:    1,
			intervalDays:  -1,
			want:          nil,
			wantErr:       true,
			errMsg:        "復習日間隔数は1以上で指定してください",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			step, err := NewPatternStep(tc.patternStepID, tc.userID, tc.patternID, tc.stepNumber, tc.intervalDays)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("unexpected error message: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, step, cmp.AllowUnexported(PatternStep{})); diff != "" {
				t.Errorf("PatternStep mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateSteps(t *testing.T) {
	if Debug {
		t.Skip("スキップさせる")
	}

	s1, _ := ReconstructPatternStep("1", "u1", "p1", 1, 1)
	s2, _ := ReconstructPatternStep("2", "u1", "p1", 2, 2)
	s3, _ := ReconstructPatternStep("3", "u1", "p1", 1, 3) // StepNumber 重複
	s4, _ := ReconstructPatternStep("4", "u1", "p1", 3, 2) // IntervalDays 重複

	tests := []struct {
		name    string
		args    []*PatternStep
		wantErr bool
		errMsg  string
	}{
		// 異常系
		{
			name:    "復習日間隔数が0（異常系）",
			args:    []*PatternStep{},
			wantErr: true,
			errMsg:  "復習日間隔数は1つ以上指定してください",
		},
		// 正常系
		{
			name:    "復習日間隔数が1つ（正常系）",
			args:    []*PatternStep{s1},
			wantErr: false,
		},
		// 正常系
		{
			name:    "復習日間隔数が2つ以上かつ昇順（正常系）",
			args:    []*PatternStep{s1, s2},
			wantErr: false,
		},
		// 異常系
		{
			name:    "順序番号が昇順でない（異常系）",
			args:    []*PatternStep{s4, s3},
			wantErr: true,
			errMsg:  "順序番号は昇順で指定してください",
		},
		// 異常系
		{
			name:    "復習日間隔数が昇順でない（異常系）",
			args:    []*PatternStep{s3, s2},
			wantErr: true,
			errMsg:  "復習日間隔数は昇順で指定してください",
		},
		// 異常系
		{
			name:    "順序番号が前の値から+1でない（公差1の等差数列でない）（異常系）",
			args:    []*PatternStep{s1, s4},
			wantErr: true,
			errMsg:  "順序番号は1ずつ増加して指定してください",
		},
		// 異常系
		{
			name:    "順序番号が重複している（異常系）",
			args:    []*PatternStep{s1, s3},
			wantErr: true,
			errMsg:  "順序番号は重複してはいけません",
		},
		// 異常系
		{
			name:    "復習日間隔数が重複している（異常系）",
			args:    []*PatternStep{s1, s4, s2},
			wantErr: true,
			errMsg:  "順序番号は1ずつ増加して指定してください",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSteps(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("err==nilのため、%qのテストが失敗しました", tt.name)
				}
				if err.Error() != tt.errMsg {
					t.Errorf("予期しないエラー:実際の結果 %q, 期待 %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("予期しないエラー: %v", err)
				}
			}
		})
	}
}
