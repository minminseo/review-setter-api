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
			want: &Pattern{
				PatternID:    testPatternID,
				UserID:       testUserID,
				Name:         "Standard Review",
				TargetWeight: TargetWeightNormal,
				RegisteredAt: now,
				EditedAt:     now,
			},
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
			want: &Pattern{
				PatternID:    "pattern4",
				UserID:       testUserID,
				Name:         "Heavy Pattern",
				TargetWeight: TargetWeightHeavy,
				RegisteredAt: now,
				EditedAt:     now,
			},
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
			want: &Pattern{
				PatternID:    "pattern5",
				UserID:       testUserID,
				Name:         "Light Pattern",
				TargetWeight: TargetWeightLight,
				RegisteredAt: now,
				EditedAt:     now,
			},
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
			want: &Pattern{
				PatternID:    "pattern6",
				UserID:       testUserID,
				Name:         "Unset Pattern",
				TargetWeight: TargetWeightUnset,
				RegisteredAt: now,
				EditedAt:     now,
			},
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

			if diff := cmp.Diff(tc.want, pattern); diff != "" {
				t.Errorf("Pattern mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPattern_Set(t *testing.T) {
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
			wantPattern: &Pattern{
				PatternID:    testPatternID,
				UserID:       testUserID,
				Name:         "Updated Pattern",
				TargetWeight: TargetWeightHeavy,
				RegisteredAt: now,
				EditedAt:     newTime,
			},
			wantErr: false,
		},
		{
			name:         "パターン名が空（異常系）",
			newName:      "",
			targetWeight: TargetWeightNormal,
			editedAt:     newTime,
			wantPattern: &Pattern{
				PatternID:    testPatternID,
				UserID:       testUserID,
				Name:         "Original",
				TargetWeight: TargetWeightNormal,
				RegisteredAt: now,
				EditedAt:     now,
			},
			wantErr: true,
			errMsg:  "名前は必須です",
		},
		{
			name:         "重みが不正（異常系）",
			newName:      "Valid Name",
			targetWeight: "invalid",
			editedAt:     newTime,
			wantPattern: &Pattern{
				PatternID:    testPatternID,
				UserID:       testUserID,
				Name:         "Original",
				TargetWeight: TargetWeightNormal,
				RegisteredAt: now,
				EditedAt:     now,
			},
			wantErr: true,
			errMsg:  "重みの値が不正です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// パターンをコピー
			testPattern := *pattern

			err := testPattern.Set(tc.newName, tc.targetWeight, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("unexpected error message: got %q, want %q", err.Error(), tc.errMsg)
				}
				if diff := cmp.Diff(tc.wantPattern, &testPattern); diff != "" {
					t.Errorf("Pattern mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantPattern, &testPattern); diff != "" {
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
			want: &PatternStep{
				PatternStepID: "step1",
				UserID:        testUserID,
				PatternID:     testPatternID,
				StepNumber:    1,
				IntervalDays:  1,
			},
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

			if diff := cmp.Diff(tc.want, step); diff != "" {
				t.Errorf("PatternStep mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateSteps(t *testing.T) {
	if Debug {
		t.Skip("スキップさせる")
	}

	s1 := &PatternStep{PatternStepID: "1", UserID: "u1", PatternID: "p1", StepNumber: 1, IntervalDays: 1}
	s2 := &PatternStep{PatternStepID: "2", UserID: "u1", PatternID: "p1", StepNumber: 2, IntervalDays: 2}
	s3 := &PatternStep{PatternStepID: "3", UserID: "u1", PatternID: "p1", StepNumber: 1, IntervalDays: 3} // StepNumber 重複
	s4 := &PatternStep{PatternStepID: "4", UserID: "u1", PatternID: "p1", StepNumber: 3, IntervalDays: 2} // IntervalDays 重複

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
