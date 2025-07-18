package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	categoryDomain "github.com/minminseo/recall-setter/domain/category"
)

func TestCategoryRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name     string
		category *categoryDomain.Category
		want     *categoryDomain.Category
		wantErr  bool
	}{
		{
			name: "カテゴリ作成に成功する場合",
			category: &categoryDomain.Category{
				ID:           uuid.New().String(),
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "新しい科目",
				RegisteredAt: time.Now(),
				EditedAt:     time.Now(),
			},
			want: &categoryDomain.Category{
				UserID: "550e8400-e29b-41d4-a716-446655440001",
				Name:   "新しい科目",
			},
			wantErr: false,
		},
		{
			name: "存在しないユーザーによる外部キー制約違反",
			category: &categoryDomain.Category{
				ID:           uuid.New().String(),
				UserID:       uuid.New().String(), // Non-existing user
				Name:         "外部キー制約違反テスト",
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
			repo := NewCategoryRepository()

			err := repo.Create(ctx, tc.category)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// 作成されたカテゴリを取得して検証
			createdCategory, err := repo.GetByID(ctx, tc.category.ID, tc.category.UserID)
			if err != nil {
				t.Errorf("created category retrieval failed: %v", err)
				return
			}

			if createdCategory == nil {
				t.Error("created category is nil")
				return
			}

			// 動的に生成されるフィールドを期待値に設定
			tc.want.ID = createdCategory.ID
			tc.want.RegisteredAt = createdCategory.RegisteredAt
			tc.want.EditedAt = createdCategory.EditedAt

			// 期待値との比較
			if diff := cmp.Diff(tc.want, createdCategory); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCategoryRepository_GetAllByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		userID        string
		want          []categoryDomain.Category
		wantErr       bool
		expectedCount int
	}{
		{
			name:   "ユーザー1のカテゴリを取得（2件）",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []categoryDomain.Category{
				{
					ID:           "650e8400-e29b-41d4-a716-446655440001",
					UserID:       "550e8400-e29b-41d4-a716-446655440001",
					Name:         "数学",
					RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				},
				{
					ID:           "650e8400-e29b-41d4-a716-446655440002",
					UserID:       "550e8400-e29b-41d4-a716-446655440001",
					Name:         "物理",
					RegisteredAt: time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
			},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name:   "ユーザー2のカテゴリを取得（2件）",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: []categoryDomain.Category{
				{
					ID:           "650e8400-e29b-41d4-a716-446655440003",
					UserID:       "550e8400-e29b-41d4-a716-446655440002",
					Name:         "英語",
					RegisteredAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:           "650e8400-e29b-41d4-a716-446655440004",
					UserID:       "550e8400-e29b-41d4-a716-446655440002",
					Name:         "歴史",
					RegisteredAt: time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
				},
			},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name:          "存在しないユーザーのカテゴリを取得",
			userID:        uuid.New().String(),
			want:          []categoryDomain.Category{},
			wantErr:       false,
			expectedCount: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewCategoryRepository()

			categories, err := repo.GetAllByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if categories != nil {
					t.Error("カテゴリがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if categories == nil {
				t.Error("カテゴリのスライスがnilです")
				return
			}

			if len(categories) != tc.expectedCount {
				t.Errorf("期待されるカテゴリ数: %d, 実際: %d", tc.expectedCount, len(categories))
			}

			// 期待値との比較
			categoriesSlice := make([]categoryDomain.Category, len(categories))
			for i, category := range categories {
				categoriesSlice[i] = *category
			}

			if diff := cmp.Diff(tc.want, categoriesSlice); diff != "" {
				t.Errorf("GetAllByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name     string
		category *categoryDomain.Category
		want     *categoryDomain.Category
		wantErr  bool
	}{
		{
			name: "カテゴリ更新に成功する場合",
			category: &categoryDomain.Category{
				ID:           "650e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "更新された数学",
				RegisteredAt: time.Now().Add(-24 * time.Hour),
				EditedAt:     time.Now(),
			},
			want: &categoryDomain.Category{
				ID:           "650e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "更新された数学",
				RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "存在しないカテゴリを更新する場合",
			category: &categoryDomain.Category{
				ID:           uuid.New().String(),
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "存在しない科目",
				RegisteredAt: time.Now().Add(-24 * time.Hour),
				EditedAt:     time.Now(),
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewCategoryRepository()

			err := repo.Update(ctx, tc.category)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tc.want != nil {
				// 更新されたカテゴリを取得して検証
				updatedCategory, err := repo.GetByID(ctx, tc.category.ID, tc.category.UserID)
				if err != nil {
					t.Errorf("updated category retrieval failed: %v", err)
					return
				}

				if updatedCategory == nil {
					t.Error("updated category is nil")
					return
				}

				// 動的に変更されるフィールドを期待値に設定
				tc.want.EditedAt = updatedCategory.EditedAt

				// 期待値との比較
				if diff := cmp.Diff(tc.want, updatedCategory); diff != "" {
					t.Errorf("Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		want       *categoryDomain.Category
		wantErr    bool
	}{
		{
			name:       "カテゴリ削除に成功する場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440004",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			want:       nil, // 削除されたので取得できない
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewCategoryRepository()

			err := repo.Delete(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// 削除後に本当に削除されたかを確認
			deletedCategory, _ := repo.GetByID(ctx, tc.categoryID, tc.userID)
			// 削除された場合はエラーが発生するか、nilが返される

			// 期待値との比較
			if diff := cmp.Diff(tc.want, deletedCategory); diff != "" {
				t.Errorf("Delete() verification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCategoryRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		want       *categoryDomain.Category
		wantErr    bool
		expectName string
	}{
		{
			name:       "有効なIDでカテゴリを取得する場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want: &categoryDomain.Category{
				ID:           "650e8400-e29b-41d4-a716-446655440001",
				UserID:       "550e8400-e29b-41d4-a716-446655440001",
				Name:         "数学",
				RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EditedAt:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			wantErr:    false,
			expectName: "数学",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewCategoryRepository()

			category, err := repo.GetByID(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if category != nil {
					t.Error("expected nil category but got one")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if category == nil {
				t.Error("expected category but got nil")
				return
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, category); diff != "" {
				t.Errorf("GetByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
