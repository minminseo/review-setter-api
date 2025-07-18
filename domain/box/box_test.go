package box

import (
	"testing"
	"time"
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
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "有効なボックスの場合（正常系）",
			id:           "box1",
			userID:       "user1",
			categoryID:   "category1",
			patternID:    "pattern1",
			boxName:      "英単語",
			registeredAt: now,
			editedAt:     now,
			wantErr:      false,
		},
		{
			name:         "名前が空の場合（異常系）",
			id:           "box2",
			userID:       "user1",
			categoryID:   "category1",
			patternID:    "pattern1",
			boxName:      "",
			registeredAt: now,
			editedAt:     now,
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
					t.Fatalf("エラーが発生するはずですが、発生しませんでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if box.ID != tc.id {
				t.Errorf("ID: got %q, want %q", box.ID, tc.id)
			}
			if box.UserID != tc.userID {
				t.Errorf("ユーザーID: got %q, want %q", box.UserID, tc.userID)
			}
			if box.CategoryID != tc.categoryID {
				t.Errorf("カテゴリID: got %q, want %q", box.CategoryID, tc.categoryID)
			}
			if box.PatternID != tc.patternID {
				t.Errorf("パターンID: got %q, want %q", box.PatternID, tc.patternID)
			}
			if box.Name != tc.boxName {
				t.Errorf("ボックス名: got %q, want %q", box.Name, tc.boxName)
			}
			if !box.RegisteredAt.Equal(tc.registeredAt) {
				t.Errorf("登録日時: got %v, want %v", box.RegisteredAt, tc.registeredAt)
			}
			if !box.EditedAt.Equal(tc.editedAt) {
				t.Errorf("編集日時: got %v, want %v", box.EditedAt, tc.editedAt)
			}
		})
	}
}

func TestBox_Set(t *testing.T) {
	now := time.Now()
	box, err := NewBox("box1", "user1", "category1", "pattern1", "Original Box", now, now)
	if err != nil {
		t.Fatalf("failed to create box: %v", err)
	}

	newTime := now.Add(time.Hour)

	tests := []struct {
		name            string
		newPatternID    string
		newName         string
		editedAt        time.Time
		wantErr         bool
		errMsg          string
		wantSamePattern bool
	}{
		{
			name:            "同じパターンで有効な更新",
			newPatternID:    "pattern1",
			newName:         "Updated Box Name",
			editedAt:        newTime,
			wantErr:         false,
			wantSamePattern: true,
		},
		{
			name:            "異なるパターンで有効な更新",
			newPatternID:    "pattern2",
			newName:         "Updated Box Name",
			editedAt:        newTime,
			wantErr:         false,
			wantSamePattern: false,
		},
		{
			name:            "名前が空の場合",
			newPatternID:    "pattern1",
			newName:         "",
			editedAt:        newTime,
			wantErr:         true,
			errMsg:          "カテゴリー名は必須です",
			wantSamePattern: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// ボックスをコピー
			testBox := *box

			isSamePattern, err := testBox.Set(tc.newPatternID, tc.newName, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生するはずですが、発生しませんでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				if isSamePattern != tc.wantSamePattern {
					t.Errorf("isSamePattern: got %v, want %v", isSamePattern, tc.wantSamePattern)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if isSamePattern != tc.wantSamePattern {
				t.Errorf("isSamePattern: got %v, want %v", isSamePattern, tc.wantSamePattern)
			}
			if testBox.PatternID != tc.newPatternID {
				t.Errorf("パターンID: got %q, want %q", testBox.PatternID, tc.newPatternID)
			}
			if testBox.Name != tc.newName {
				t.Errorf("ボックス名: got %q, want %q", testBox.Name, tc.newName)
			}
			if !testBox.EditedAt.Equal(tc.editedAt) {
				t.Errorf("編集日時: got %v, want %v", testBox.EditedAt, tc.editedAt)
			}
		})
	}
}
