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
			category: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					"新しい科目",
					time.Now(),
					time.Now(),
				)
				return cat
			}(),
			want: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					"",
					"550e8400-e29b-41d4-a716-446655440001",
					"新しい科目",
					time.Time{},
					time.Time{},
				)
				return cat
			}(),
			wantErr: false,
		},
		{
			name: "存在しないユーザーによる外部キー制約違反",
			category: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					uuid.New().String(),
					uuid.New().String(), // Non-existing user
					"外部キー制約違反テスト",
					time.Now(),
					time.Now(),
				)
				return cat
			}(),
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
			createdCategory, err := repo.GetByID(ctx, tc.category.ID(), tc.category.UserID())
			if err != nil {
				t.Errorf("created category retrieval failed: %v", err)
				return
			}

			if createdCategory == nil {
				t.Error("created category is nil")
				return
			}

			// 動的に生成されるフィールドを期待値に設定して新しい期待値を作成
			want, _ := categoryDomain.ReconstructCategory(
				createdCategory.ID(),
				tc.want.UserID(),
				tc.want.Name(),
				createdCategory.RegisteredAt(),
				createdCategory.EditedAt(),
			)

			// 期待値との比較
			if diff := cmp.Diff(want, createdCategory, cmp.AllowUnexported(categoryDomain.Category{})); diff != "" {
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
			want: func() []categoryDomain.Category {
				cat1, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"数学",
					time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				)
				cat2, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440002",
					"550e8400-e29b-41d4-a716-446655440001",
					"物理",
					time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				)
				return []categoryDomain.Category{*cat1, *cat2}
			}(),
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name:   "ユーザー2のカテゴリを取得（2件）",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: func() []categoryDomain.Category {
				cat1, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440003",
					"550e8400-e29b-41d4-a716-446655440002",
					"英語",
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				)
				cat2, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440004",
					"550e8400-e29b-41d4-a716-446655440002",
					"歴史",
					time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
				)
				return []categoryDomain.Category{*cat1, *cat2}
			}(),
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

			if diff := cmp.Diff(tc.want, categoriesSlice, cmp.AllowUnexported(categoryDomain.Category{})); diff != "" {
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
			category: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"更新された数学",
					time.Now().Add(-24 * time.Hour),
					time.Now(),
				)
				return cat
			}(),
			want: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"更新された数学",
					time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					time.Time{},
				)
				return cat
			}(),
			wantErr: false,
		},
		{
			name: "存在しないカテゴリを更新する場合",
			category: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					"存在しない科目",
					time.Now().Add(-24 * time.Hour),
					time.Now(),
				)
				return cat
			}(),
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
				updatedCategory, err := repo.GetByID(ctx, tc.category.ID(), tc.category.UserID())
				if err != nil {
					t.Errorf("updated category retrieval failed: %v", err)
					return
				}

				if updatedCategory == nil {
					t.Error("updated category is nil")
					return
				}

				// 動的に変更されるフィールドを期待値に設定して新しい期待値を作成
				want, _ := categoryDomain.ReconstructCategory(
					tc.want.ID(),
					tc.want.UserID(),
					tc.want.Name(),
					tc.want.RegisteredAt(),
					updatedCategory.EditedAt(),
				)

				// 期待値との比較
				if diff := cmp.Diff(want, updatedCategory, cmp.AllowUnexported(categoryDomain.Category{})); diff != "" {
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
			if diff := cmp.Diff(tc.want, deletedCategory, cmp.AllowUnexported(categoryDomain.Category{})); diff != "" {
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
			want: func() *categoryDomain.Category {
				cat, _ := categoryDomain.ReconstructCategory(
					"650e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"数学",
					time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				)
				return cat
			}(),
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
			if diff := cmp.Diff(tc.want, category, cmp.AllowUnexported(categoryDomain.Category{})); diff != "" {
				t.Errorf("GetByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
