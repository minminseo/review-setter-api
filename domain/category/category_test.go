package category

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testUserID     = "user1"
	testCategoryID = "category1"
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
		want         *Category
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "有効なカテゴリー（正常系）",
			id:           testCategoryID,
			userID:       testUserID,
			categoryName: "英語",
			registeredAt: now,
			editedAt:     now,
			want: func() *Category {
				category, _ := ReconstructCategory(
					testCategoryID,
					testUserID,
					"英語",
					now,
					now,
				)
				return category
			}(),
			wantErr: false,
		},
		{
			name:         "カテゴリー名が空（異常系）",
			id:           "category2",
			userID:       testUserID,
			categoryName: "",
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

			category, err := NewCategory(tc.id, tc.userID, tc.categoryName, tc.registeredAt, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if diff := cmp.Diff(tc.want, category, cmp.AllowUnexported(Category{})); diff != "" {
				t.Errorf("Category mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCategory_UpdateCategory(t *testing.T) {
	now := time.Now()
	category, err := NewCategory(testCategoryID, testUserID, "Original Name", now, now)
	if err != nil {
		t.Fatalf("カテゴリーの生成に失敗しました: %v", err)
	}

	newTime := now.Add(time.Hour)

	tests := []struct {
		name         string
		newName      string
		editedAt     time.Time
		wantCategory *Category
		wantErr      bool
		errMsg       string
	}{
		{
			name:     "カテゴリー名を更新（正常系）",
			newName:  "Updated Category Name",
			editedAt: newTime,
			wantCategory: func() *Category {
				category, _ := ReconstructCategory(
					testCategoryID,
					testUserID,
					"Updated Category Name",
					now,
					newTime,
				)
				return category
			}(),
			wantErr: false,
		},
		{
			name:     "カテゴリー名が空で更新（異常系）",
			newName:  "",
			editedAt: newTime,
			wantCategory: func() *Category {
				category, _ := ReconstructCategory(
					testCategoryID,
					testUserID,
					"Original Name",
					now,
					now,
				)
				return category
			}(),
			wantErr: true,
			errMsg:  "カテゴリー名は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// カテゴリーをコピー
			testCategory := *category

			err := testCategory.UpdateCategory(tc.newName, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				if diff := cmp.Diff(tc.wantCategory, &testCategory, cmp.AllowUnexported(Category{})); diff != "" {
					t.Errorf("Category mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if diff := cmp.Diff(tc.wantCategory, &testCategory, cmp.AllowUnexported(Category{})); diff != "" {
				t.Errorf("Category mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
