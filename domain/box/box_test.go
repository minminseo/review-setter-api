package box

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testUserID     = "user1"
	testCategoryID = "category1"
	testPatternID  = "pattern1"
	testBoxID      = "box1"
)

func TestNewBox(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		id           string
		userID       string
		categoryID   string
		patternID    string
		boxName      string
		registeredAt time.Time
		editedAt     time.Time
		want         *Box
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "有効なボックスの場合（正常系）",
			id:           testBoxID,
			userID:       testUserID,
			categoryID:   testCategoryID,
			patternID:    testPatternID,
			boxName:      "英単語",
			registeredAt: now,
			editedAt:     now,
			want: func() *Box {
				box, _ := ReconstructBox(
					testBoxID,
					testUserID,
					testCategoryID,
					testPatternID,
					"英単語",
					now,
					now,
				)
				return box
			}(),
			wantErr: false,
		},
		{
			name:         "名前が空の場合（異常系）",
			id:           "box2",
			userID:       testUserID,
			categoryID:   testCategoryID,
			patternID:    testPatternID,
			boxName:      "",
			registeredAt: now,
			editedAt:     now,
			want:         nil,
			wantErr:      true,
			errMsg:       "カテゴリー名は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			box, err := NewBox(tc.id, tc.userID, tc.categoryID, tc.patternID, tc.boxName, tc.registeredAt, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生するはずですが、発生しませんでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if diff := cmp.Diff(tc.want, box, cmp.AllowUnexported(Box{})); diff != "" {
				t.Errorf("Box mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBox_UpdateBox(t *testing.T) {
	now := time.Now()
	box, err := NewBox(testBoxID, testUserID, testCategoryID, testPatternID, "Original Box", now, now)
	if err != nil {
		t.Fatalf("failed to create box: %v", err)
	}

	newTime := now.Add(time.Hour)

	tests := []struct {
		name            string
		newPatternID    string
		newName         string
		editedAt        time.Time
		wantBox         *Box
		wantSamePattern bool
		wantErr         bool
		errMsg          string
	}{
		{
			name:         "同じパターンで有効な更新",
			newPatternID: testPatternID,
			newName:      "Updated Box Name",
			editedAt:     newTime,
			wantBox: func() *Box {
				box, _ := ReconstructBox(
					testBoxID,
					testUserID,
					testCategoryID,
					testPatternID,
					"Updated Box Name",
					now,
					newTime,
				)
				return box
			}(),
			wantSamePattern: true,
			wantErr:         false,
		},
		{
			name:         "異なるパターンで有効な更新",
			newPatternID: "pattern2",
			newName:      "Updated Box Name",
			editedAt:     newTime,
			wantBox: func() *Box {
				box, _ := ReconstructBox(
					testBoxID,
					testUserID,
					testCategoryID,
					"pattern2",
					"Updated Box Name",
					now,
					newTime,
				)
				return box
			}(),
			wantSamePattern: false,
			wantErr:         false,
		},
		{
			name:         "名前が空の場合",
			newPatternID: testPatternID,
			newName:      "",
			editedAt:     newTime,
			wantBox: func() *Box {
				box, _ := ReconstructBox(
					testBoxID,
					testUserID,
					testCategoryID,
					testPatternID,
					"Original Box",
					now,
					now,
				)
				return box
			}(),
			wantSamePattern: true,
			wantErr:         true,
			errMsg:          "カテゴリー名は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// ボックスをコピー
			testBox := *box

			isSamePattern, err := testBox.UpdateBox(tc.newPatternID, tc.newName, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生するはずですが、発生しませんでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				if isSamePattern != tc.wantSamePattern {
					t.Errorf("isSamePattern: got %v, want %v", isSamePattern, tc.wantSamePattern)
				}
				if diff := cmp.Diff(tc.wantBox, &testBox, cmp.AllowUnexported(Box{})); diff != "" {
					t.Errorf("Box mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if isSamePattern != tc.wantSamePattern {
				t.Errorf("isSamePattern: got %v, want %v", isSamePattern, tc.wantSamePattern)
			}

			if diff := cmp.Diff(tc.wantBox, &testBox, cmp.AllowUnexported(Box{})); diff != "" {
				t.Errorf("Box mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
