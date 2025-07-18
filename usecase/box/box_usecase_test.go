package box

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	boxDomain "github.com/minminseo/recall-setter/domain/box"
)

func TestCreateBox(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     CreateBoxInput
		setupMock func(*boxDomain.MockIBoxRepository)
		want      *CreateBoxOutput
		wantErr   bool
	}{
		{
			name: "正常系_有効な入力での作成成功",
			input: CreateBoxInput{
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "英語学習ボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						Create(ctx, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &CreateBoxOutput{
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "英語学習ボックス",
			},
			wantErr: false,
		},
		{
			name: "異常系_空文字名前での作成失敗",
			input: CreateBoxInput{
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				// Create() should not be called
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_リポジトリ保存エラー",
			input: CreateBoxInput{
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "テストボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						Create(ctx, gomock.Any()).
						Return(errors.New("database error")).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := boxDomain.NewMockIBoxRepository(ctrl)
			tt.setupMock(mockRepo)

			usecase := NewBoxUsecase(mockRepo)
			got, err := usecase.CreateBox(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.ID == "" {
				t.Error("CreateBox() ID should not be empty")
			}
			if got.RegisteredAt.IsZero() {
				t.Error("CreateBox() RegisteredAt should not be zero")
			}
			if got.EditedAt.IsZero() {
				t.Error("CreateBox() EditedAt should not be zero")
			}

			// compare other fields
			want := tt.want
			want.ID = got.ID
			want.RegisteredAt = got.RegisteredAt
			want.EditedAt = got.EditedAt

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("CreateBox() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetBoxesByCategoryID(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		setupMock  func(*boxDomain.MockIBoxRepository)
		want       []*GetBoxOutput
		wantErr    bool
	}{
		{
			name:       "正常系_複数ボックスの取得成功",
			categoryID: "test-category-id",
			userID:     "test-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				boxes := []*boxDomain.Box{
					{ID: "box1", UserID: "test-user-id", CategoryID: "test-category-id", PatternID: "pattern1", Name: "ボックス1", RegisteredAt: mockTime, EditedAt: mockTime},
					{ID: "box2", UserID: "test-user-id", CategoryID: "test-category-id", PatternID: "pattern2", Name: "ボックス2", RegisteredAt: mockTime, EditedAt: mockTime},
				}
				gomock.InOrder(
					m.EXPECT().
						GetAllByCategoryID(ctx, "test-category-id", "test-user-id").
						Return(boxes, nil).
						Times(1),
				)
			},
			want: []*GetBoxOutput{
				{ID: "box1", UserID: "test-user-id", CategoryID: "test-category-id", PatternID: "pattern1", Name: "ボックス1", RegisteredAt: mockTime, EditedAt: mockTime},
				{ID: "box2", UserID: "test-user-id", CategoryID: "test-category-id", PatternID: "pattern2", Name: "ボックス2", RegisteredAt: mockTime, EditedAt: mockTime},
			},
			wantErr: false,
		},
		{
			name:       "正常系_空のボックスリストの取得",
			categoryID: "empty-category-id",
			userID:     "test-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						GetAllByCategoryID(ctx, "empty-category-id", "test-user-id").
						Return([]*boxDomain.Box{}, nil).
						Times(1),
				)
			},
			want:    []*GetBoxOutput{},
			wantErr: false,
		},
		{
			name:       "異常系_リポジトリアクセスエラー",
			categoryID: "error-category-id",
			userID:     "test-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						GetAllByCategoryID(ctx, "error-category-id", "test-user-id").
						Return(nil, errors.New("database error")).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := boxDomain.NewMockIBoxRepository(ctrl)
			tt.setupMock(mockRepo)

			usecase := NewBoxUsecase(mockRepo)
			got, err := usecase.GetBoxesByCategoryID(ctx, tt.categoryID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBoxesByCategoryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetBoxesByCategoryID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateBox(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     UpdateBoxInput
		setupMock func(*boxDomain.MockIBoxRepository)
		want      *UpdateBoxOutput
		wantErr   bool
	}{
		{
			name: "正常系_同じパターンIDでの名前変更",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "更新された英語学習ボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existing := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existing, nil).
						Times(1),
					m.EXPECT().
						Update(ctx, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateBoxOutput{
				ID:         "44444444-4444-4444-4444-444444444444",
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "更新された英語学習ボックス",
			},
			wantErr: false,
		},
		{
			name: "正常系_異なるパターンIDでの更新成功",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "55555555-5555-5555-5555-555555555555",
				Name:       "パターン変更されたボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existing := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existing, nil).
						Times(1),
					m.EXPECT().
						UpdateWithPatternID(ctx, gomock.Any()).
						Return(int64(1), nil).
						Times(1),
				)
			},
			want: &UpdateBoxOutput{
				ID:         "44444444-4444-4444-4444-444444444444",
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "55555555-5555-5555-5555-555555555555",
				Name:       "パターン変更されたボックス",
			},
			wantErr: false,
		},
		{
			name: "異常系_存在しないボックスでの更新失敗",
			input: UpdateBoxInput{
				ID:         "nonexistent-box-id",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				UserID:     "11111111-1111-1111-1111-111111111111",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "更新テスト",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "nonexistent-box-id", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(nil, errors.New("box not found")).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_空文字名前での更新失敗",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				UserID:     "11111111-1111-1111-1111-111111111111",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existingBox := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(gomock.Any(), "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existingBox, nil).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_同じパターンでのリポジトリ更新エラー",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				UserID:     "11111111-1111-1111-1111-111111111111",
				PatternID:  "33333333-3333-3333-3333-333333333333",
				Name:       "更新された英語学習ボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existing := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existing, nil).
						Times(1),
					m.EXPECT().
						Update(ctx, gomock.Any()).
						Return(errors.New("database error")).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_異なるパターンでのリポジトリ更新エラー",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				UserID:     "11111111-1111-1111-1111-111111111111",
				PatternID:  "55555555-5555-5555-5555-555555555555",
				Name:       "パターン変更されたボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existing := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existing, nil).
						Times(1),
					m.EXPECT().
						UpdateWithPatternID(ctx, gomock.Any()).
						Return(int64(0), errors.New("database error")).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_パターン変更時の競合エラー",
			input: UpdateBoxInput{
				ID:         "44444444-4444-4444-4444-444444444444",
				CategoryID: "22222222-2222-2222-2222-222222222222",
				UserID:     "11111111-1111-1111-1111-111111111111",
				PatternID:  "55555555-5555-5555-5555-555555555555",
				Name:       "パターン変更されたボックス",
			},
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				existing := &boxDomain.Box{
					ID:           "44444444-4444-4444-4444-444444444444",
					UserID:       "11111111-1111-1111-1111-111111111111",
					CategoryID:   "22222222-2222-2222-2222-222222222222",
					PatternID:    "33333333-3333-3333-3333-333333333333",
					Name:         "英語学習ボックス",
					RegisteredAt: mockTime,
					EditedAt:     mockTime,
				}
				gomock.InOrder(
					m.EXPECT().
						GetByID(ctx, "44444444-4444-4444-4444-444444444444", "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111").
						Return(existing, nil).
						Times(1),
					m.EXPECT().
						UpdateWithPatternID(ctx, gomock.Any()).
						Return(int64(0), nil).
						Times(1),
				)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := boxDomain.NewMockIBoxRepository(ctrl)
			tt.setupMock(mockRepo)

			usecase := NewBoxUsecase(mockRepo)
			got, err := usecase.UpdateBox(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// EditedAt is dynamic
			want := tt.want
			want.EditedAt = got.EditedAt

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("UpdateBox() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteBox(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		boxID      string
		categoryID string
		userID     string
		setupMock  func(*boxDomain.MockIBoxRepository)
		wantErr    bool
	}{
		{
			name:       "正常系_存在するボックスの削除成功",
			boxID:      "existing-box-id",
			categoryID: "valid-category-id",
			userID:     "valid-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						Delete(ctx, "existing-box-id", "valid-category-id", "valid-user-id").
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
		{
			name:       "異常系_存在しないボックスでの削除失敗",
			boxID:      "nonexistent-box-id",
			categoryID: "valid-category-id",
			userID:     "valid-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						Delete(ctx, "nonexistent-box-id", "valid-category-id", "valid-user-id").
						Return(errors.New("box not found")).
						Times(1),
				)
			},
			wantErr: true,
		},
		{
			name:       "異常系_リポジトリアクセスエラー",
			boxID:      "error-box-id",
			categoryID: "valid-category-id",
			userID:     "valid-user-id",
			setupMock: func(m *boxDomain.MockIBoxRepository) {
				gomock.InOrder(
					m.EXPECT().
						Delete(ctx, "error-box-id", "valid-category-id", "valid-user-id").
						Return(errors.New("database error")).
						Times(1),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := boxDomain.NewMockIBoxRepository(ctrl)
			tt.setupMock(mockRepo)

			usecase := NewBoxUsecase(mockRepo)
			err := usecase.DeleteBox(ctx, tt.boxID, tt.categoryID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
