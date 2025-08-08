package category

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	categoryDomain "github.com/minminseo/recall-setter/domain/category"
)

func TestCreateCategory(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     CreateCategoryInput
		mockSetup func(*categoryDomain.MockICategoryRepository)
		wantErr   bool
		validate  func(*testing.T, *CreateCategoryOutput, error)
	}{
		{
			name:  "正常系_有効な入力でカテゴリ作成成功",
			input: CreateCategoryInput{UserID: "valid-user-id", Name: "テストカテゴリ"},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						Create(ctx, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result *CreateCategoryOutput, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil {
					t.Fatal("expected result, got nil")
				}
				if result.ID == "" {
					t.Error("expected ID to be set")
				}
				if result.UserID != "valid-user-id" {
					t.Errorf("expected UserID %s, got %s", "valid-user-id", result.UserID)
				}
				if result.Name != "テストカテゴリ" {
					t.Errorf("expected Name %s, got %s", "テストカテゴリ", result.Name)
				}
				if result.RegisteredAt.IsZero() {
					t.Error("expected RegisteredAt to be set")
				}
				if result.EditedAt.IsZero() {
					t.Error("expected EditedAt to be set")
				}
			},
		},
		{
			name:      "異常系_空のNameでの作成失敗",
			input:     CreateCategoryInput{UserID: "valid-user-id", Name: ""},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {},
			wantErr:   true,
			validate: func(t *testing.T, result *CreateCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
		{
			name:  "異常系_リポジトリ保存エラー",
			input: CreateCategoryInput{UserID: "valid-user-id", Name: "テストカテゴリ"},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						Create(ctx, gomock.Any()).
						Return(errors.New("repository error")).
						Times(1),
				)
			},
			wantErr: true,
			validate: func(t *testing.T, result *CreateCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := categoryDomain.NewMockICategoryRepository(ctrl)
			usecase := NewCategoryUsecase(mockRepo)

			tc.mockSetup(mockRepo)
			result, err := usecase.CreateCategory(ctx, tc.input)
			tc.validate(t, result, err)
		})
	}
}

func TestGetCategoriesByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	category1, _ := categoryDomain.ReconstructCategory("category-1", "user-1", "カテゴリ1", now, now)
	category2, _ := categoryDomain.ReconstructCategory("category-2", "user-1", "カテゴリ2", now, now)
	testCategories := []*categoryDomain.Category{category1, category2}

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*categoryDomain.MockICategoryRepository)
		wantErr   bool
		validate  func(*testing.T, []*GetCategoryOutput, error)
	}{
		{
			name:   "正常系_複数カテゴリの取得",
			userID: "user-1",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetAllByUserID(ctx, "user-1").
						Return(testCategories, nil).
						Times(1),
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result []*GetCategoryOutput, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(result) != 2 {
					t.Fatalf("expected 2 categories, got %d", len(result))
				}
				if result[0].Name != "カテゴリ1" {
					t.Errorf("expected first Name %s, got %s", "カテゴリ1", result[0].Name)
				}
				if result[1].Name != "カテゴリ2" {
					t.Errorf("expected second Name %s, got %s", "カテゴリ2", result[1].Name)
				}
			},
		},
		{
			name:   "正常系_カテゴリが存在しない場合の空配列返却",
			userID: "user-empty",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetAllByUserID(ctx, "user-empty").
						Return([]*categoryDomain.Category{}, nil).
						Times(1),
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result []*GetCategoryOutput, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(result) != 0 {
					t.Errorf("expected 0 categories, got %d", len(result))
				}
			},
		},
		{
			name:   "異常系_リポジトリエラー",
			userID: "user-error",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetAllByUserID(ctx, "user-error").
						Return(nil, errors.New("repository error")).
						Times(1),
				)
			},
			wantErr: true,
			validate: func(t *testing.T, result []*GetCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := categoryDomain.NewMockICategoryRepository(ctrl)
			usecase := NewCategoryUsecase(mockRepo)

			tc.mockSetup(mockRepo)
			result, err := usecase.GetCategoriesByUserID(ctx, tc.userID)
			tc.validate(t, result, err)
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	existing, _ := categoryDomain.ReconstructCategory(
		"category-1",
		"user-1",
		"古いカテゴリ名",
		now,
		now,
	)

	tests := []struct {
		name      string
		input     UpdateCategoryInput
		mockSetup func(*categoryDomain.MockICategoryRepository)
		wantErr   bool
		validate  func(*testing.T, *UpdateCategoryOutput, error)
	}{
		{
			name:  "正常系_有効なパラメータでの更新",
			input: UpdateCategoryInput{ID: "category-1", UserID: "user-1", Name: "新しいカテゴリ名"},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetByID(ctx, "category-1", "user-1").
						Return(existing, nil).
						Times(1),
					repo.EXPECT().
						Update(ctx, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result *UpdateCategoryOutput, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil {
					t.Fatal("expected result, got nil")
				}
				if result.Name != "新しいカテゴリ名" {
					t.Errorf("expected Name %s, got %s", "新しいカテゴリ名", result.Name)
				}
				if !result.EditedAt.After(now) {
					t.Error("expected EditedAt to be updated")
				}
			},
		},
		{
			name:  "異常系_存在しないカテゴリIDでの更新失敗",
			input: UpdateCategoryInput{ID: "non-existent-id", UserID: "user-1", Name: "更新名"},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetByID(ctx, "non-existent-id", "user-1").
						Return(nil, errors.New("category not found")).
						Times(1),
				)
			},
			wantErr: true,
			validate: func(t *testing.T, result *UpdateCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
		{
			name:  "異常系_空のNameでの更新失敗",
			input: UpdateCategoryInput{ID: "category-1", UserID: "user-1", Name: ""},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetByID(ctx, "category-1", "user-1").
						Return(existing, nil).
						Times(1),
				)
			},
			wantErr: true,
			validate: func(t *testing.T, result *UpdateCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
		{
			name:  "異常系_リポジトリ更新エラー",
			input: UpdateCategoryInput{ID: "category-1", UserID: "user-1", Name: "新しいカテゴリ名"},
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						GetByID(ctx, "category-1", "user-1").
						Return(existing, nil).
						Times(1),
					repo.EXPECT().
						Update(ctx, gomock.Any()).
						Return(errors.New("repository update error")).
						Times(1),
				)
			},
			wantErr: true,
			validate: func(t *testing.T, result *UpdateCategoryOutput, err error) {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := categoryDomain.NewMockICategoryRepository(ctrl)
			usecase := NewCategoryUsecase(mockRepo)

			tc.mockSetup(mockRepo)
			result, err := usecase.UpdateCategory(ctx, tc.input)
			tc.validate(t, result, err)
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		categoryID string
		userID     string
		mockSetup  func(*categoryDomain.MockICategoryRepository)
		wantErr    bool
	}{
		{
			name:       "正常系_有効なパラメータでの削除成功",
			categoryID: "category-1",
			userID:     "user-1",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						Delete(ctx, "category-1", "user-1").
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name:       "異常系_存在しないカテゴリIDでの削除失敗",
			categoryID: "non-existent-id",
			userID:     "user-1",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						Delete(ctx, "non-existent-id", "user-1").
						Return(errors.New("category not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name:       "異常系_リポジトリ削除エラー",
			categoryID: "category-1",
			userID:     "user-1",
			mockSetup: func(repo *categoryDomain.MockICategoryRepository) {
				gomock.InOrder(
					repo.EXPECT().
						Delete(ctx, "category-1", "user-1").
						Return(errors.New("repository delete error")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := categoryDomain.NewMockICategoryRepository(ctrl)
			usecase := NewCategoryUsecase(mockRepo)

			tc.mockSetup(mockRepo)
			err := usecase.DeleteCategory(ctx, tc.categoryID, tc.userID)
			if tc.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestCategoryUsecase_DTOMapping(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := categoryDomain.NewMockICategoryRepository(ctrl)
	usecase := NewCategoryUsecase(mockRepo)
	ctx := context.Background()
	now := time.Now().UTC()

	category, _ := categoryDomain.ReconstructCategory(
		"test-id",
		"test-user",
		"テストカテゴリ",
		now,
		now,
	)

	mockRepo.EXPECT().
		GetAllByUserID(ctx, "test-user").
		Return([]*categoryDomain.Category{category}, nil).
		Times(1)

	result, err := usecase.GetCategoriesByUserID(ctx, "test-user")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := &GetCategoryOutput{
		ID:           "test-id",
		UserID:       "test-user",
		Name:         "テストカテゴリ",
		RegisteredAt: now,
		EditedAt:     now,
	}

	if diff := cmp.Diff(expected, result[0]); diff != "" {
		t.Errorf("DTO mapping mismatch (-expected +actual):\n%s", diff)
	}
}
