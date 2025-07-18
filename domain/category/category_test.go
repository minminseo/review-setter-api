package category

import (
	"testing"
	"time"
)

func TestNewCategory(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		id           string
		userID       string
		categoryName string
		registeredAt time.Time
		editedAt     time.Time
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "有効なカテゴリー（正常系）",
			id:           "category1",
			userID:       "user1",
			categoryName: "英語",
			registeredAt: now,
			editedAt:     now,
			wantErr:      false,
		},
		{
			name:         "カテゴリー名が空（異常系）",
			id:           "category2",
			userID:       "user1",
			categoryName: "",
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

			category, err := NewCategory(tc.id, tc.userID, tc.categoryName, tc.registeredAt, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if category.ID != tc.id {
				t.Errorf("IDが一致しません: got %q, want %q", category.ID, tc.id)
			}
			if category.UserID != tc.userID {
				t.Errorf("ユーザーIDが一致しません: got %q, want %q", category.UserID, tc.userID)
			}
			if category.Name != tc.categoryName {
				t.Errorf("カテゴリー名が一致しません: got %q, want %q", category.Name, tc.categoryName)
			}
			if !category.RegisteredAt.Equal(tc.registeredAt) {
				t.Errorf("登録日時が一致しません: got %v, want %v", category.RegisteredAt, tc.registeredAt)
			}
			if !category.EditedAt.Equal(tc.editedAt) {
				t.Errorf("編集日時が一致しません: got %v, want %v", category.EditedAt, tc.editedAt)
			}
		})
	}
}

func TestCategory_Set(t *testing.T) {
	now := time.Now()
	category, err := NewCategory("category1", "user1", "Original Name", now, now)
	if err != nil {
		t.Fatalf("カテゴリーの生成に失敗しました: %v", err)
	}

	newTime := now.Add(time.Hour)

	tests := []struct {
		name     string
		newName  string
		editedAt time.Time
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "カテゴリー名を更新（正常系）",
			newName:  "Updated Category Name",
			editedAt: newTime,
			wantErr:  false,
		},
		{
			name:     "カテゴリー名が空で更新（異常系）",
			newName:  "",
			editedAt: newTime,
			wantErr:  true,
			errMsg:   "カテゴリー名は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// カテゴリーをコピー
			testCategory := *category

			err := testCategory.Set(tc.newName, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if testCategory.Name != tc.newName {
				t.Errorf("カテゴリー名: got %q, want %q", testCategory.Name, tc.newName)
			}
			if !testCategory.EditedAt.Equal(tc.editedAt) {
				t.Errorf("編集日時: got %v, want %v", testCategory.EditedAt, tc.editedAt)
			}
		})
	}
}
