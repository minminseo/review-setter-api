package item

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	BoxDomain "github.com/minminseo/recall-setter/domain/box"
	CategoryDomain "github.com/minminseo/recall-setter/domain/category"
	ItemDomain "github.com/minminseo/recall-setter/domain/item"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

func TestItemUsecase_CreateItem(t *testing.T) {
	ctx := context.Background()

	// テストデータの準備
	userID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	patternID := uuid.NewString()
	itemID := uuid.NewString()

	parsedLearnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	registeredAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	// PatternStepのテストデータ
	testPatternSteps := []*PatternDomain.PatternStep{
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    1,
			IntervalDays:  1,
		},
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    2,
			IntervalDays:  3,
		},
	}

	// Reviewdateのテストデータ
	testReviewdates1 := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: parsedToday,
			ScheduledDate:        parsedToday,
			IsCompleted:          false,
		},
	}

	testReviewdates2 := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: parsedLearnedDate.AddDate(0, 0, 1),
			ScheduledDate:        parsedLearnedDate.AddDate(0, 0, 1),
			IsCompleted:          true,
		},
	}

	tests := []struct {
		name      string
		input     CreateItemInput
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      *CreateItemOutput
		wantErr   bool
	}{
		{
			name: "PatternIDがnilの場合の正常系",
			input: CreateItemInput{
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                nil,
				Name:                     "Test Item",
				Detail:                   "Test Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						CreateItem(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &CreateItemOutput{
				ItemID:       itemID,
				UserID:       userID,
				CategoryID:   &categoryID,
				BoxID:        &boxID,
				PatternID:    nil,
				Name:         "Test Item",
				Detail:       "Test Detail",
				LearnedDate:  "2024-01-01",
				IsCompleted:  false,
				RegisteredAt: registeredAt,
				EditedAt:     registeredAt,
				Reviewdates:  nil,
			},
			wantErr: false,
		},
		{
			name: "PatternIDが設定されている場合の正常系（IsMarkOverdueAsCompleted=true。作成時点で全ての復習日が完了の時）",
			input: CreateItemInput{
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Test Item",
				Detail:                   "Test Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: true,
				Today:                    "2024-01-10",
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().
						GetAllPatternStepsByPatternID(gomock.Any(), patternID, userID).
						Return(testPatternSteps, nil).
						Times(1),
					mockScheduler.EXPECT().
						FormatWithOverdueMarkedCompleted(
							testPatternSteps,
							userID,
							&categoryID,
							&boxID,
							gomock.Any(),
							parsedLearnedDate,
							parsedToday,
						).
						Return(testReviewdates2, true, nil).
						Times(1),
					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					mockItemRepo.EXPECT().
						CreateItem(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),
					mockItemRepo.EXPECT().
						CreateReviewdates(gomock.Any(), testReviewdates2).
						Return(int64(1), nil).
						Times(1),
				)
			},
			want: &CreateItemOutput{
				ItemID:       itemID,
				UserID:       userID,
				CategoryID:   &categoryID,
				BoxID:        &boxID,
				PatternID:    &patternID,
				Name:         "Test Item",
				Detail:       "Test Detail",
				LearnedDate:  "2024-01-01",
				IsCompleted:  true,
				RegisteredAt: registeredAt,
				EditedAt:     registeredAt,
				Reviewdates: []CreateReviewdateOutput{
					{
						DateID:               testReviewdates2[0].ReviewdateID,
						UserID:               userID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: "2024-01-02",
						ScheduledDate:        "2024-01-02",
						IsCompleted:          true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PatternIDが設定されている場合の正常系（IsMarkOverdueAsCompleted=false）",
			input: CreateItemInput{
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Test Item",
				Detail:                   "Test Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().
						GetAllPatternStepsByPatternID(gomock.Any(), patternID, userID).
						Return(testPatternSteps, nil).
						Times(1),
					mockScheduler.EXPECT().
						FormatWithOverdueMarkedInCompleted(
							testPatternSteps,
							userID,
							&categoryID,
							&boxID,
							gomock.Any(),
							parsedLearnedDate,
							parsedToday,
						).
						Return(testReviewdates1, nil).
						Times(1),
					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					mockItemRepo.EXPECT().
						CreateItem(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1),
					mockItemRepo.EXPECT().
						CreateReviewdates(gomock.Any(), testReviewdates1).
						Return(int64(1), nil).
						Times(1),
				)
			},
			want: &CreateItemOutput{
				ItemID:       itemID,
				UserID:       userID,
				CategoryID:   &categoryID,
				BoxID:        &boxID,
				PatternID:    &patternID,
				Name:         "Test Item",
				Detail:       "Test Detail",
				LearnedDate:  "2024-01-01",
				IsCompleted:  false,
				RegisteredAt: registeredAt,
				EditedAt:     registeredAt,
				Reviewdates: []CreateReviewdateOutput{
					{
						DateID:               testReviewdates1[0].ReviewdateID,
						UserID:               userID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: "2024-01-10",
						ScheduledDate:        "2024-01-10",
						IsCompleted:          false,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.CreateItem(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != nil {
				// 時刻系フィールドは動的に生成されるため、テストでは除外
				got.RegisteredAt = tc.want.RegisteredAt
				got.EditedAt = tc.want.EditedAt

				// UUIDも動的に生成されるため、テストでは除外
				got.ItemID = tc.want.ItemID

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("CreateItem() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemUsecase_DeleteItem(t *testing.T) {
	ctx := context.Background()

	itemID := uuid.NewString()
	userID := uuid.NewString()

	tests := []struct {
		name      string
		itemID    string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name:   "正常系",
			itemID: itemID,
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						DeleteItem(gomock.Any(), itemID, userID).
						Return(nil).
						Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			err := usecase.DeleteItem(ctx, tc.itemID, tc.userID)

			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItemAsFinishedForce(t *testing.T) {
	ctx := context.Background()

	itemID := uuid.NewString()
	userID := uuid.NewString()
	editedAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testItem := &ItemDomain.Item{
		ItemID:       itemID,
		UserID:       userID,
		CategoryID:   nil,
		BoxID:        nil,
		PatternID:    nil,
		Name:         "Test Item",
		Detail:       "Test Detail",
		LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		IsFinished:   false,
		RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		EditedAt:     editedAt,
	}

	tests := []struct {
		name      string
		input     UpdateItemAsFinishedForceInput
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      *UpdateItemAsFinishedForceOutput
		wantErr   bool
	}{
		{
			name: "正常系",
			input: UpdateItemAsFinishedForceInput{
				ItemID: itemID,
				UserID: userID,
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetItemByID(gomock.Any(), itemID, userID).
						Return(testItem, nil).
						Times(1),
					mockItemRepo.EXPECT().
						UpdateItemAsFinished(gomock.Any(), itemID, userID, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateItemAsFinishedForceOutput{
				ItemID:     testItem.ItemID,
				UserID:     testItem.UserID,
				IsFinished: true,
				EditedAt:   editedAt,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.UpdateItemAsFinishedForce(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateItemAsFinishedForce() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != nil {
				// 時刻系フィールドは動的に生成されるため、テストでは除外
				got.EditedAt = tc.want.EditedAt

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("UpdateItemAsFinishedForce() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemUsecase_UpdateReviewDateAsCompleted(t *testing.T) {
	ctx := context.Background()
	reviewDateID := uuid.NewString()
	userID := uuid.NewString()
	itemID := uuid.NewString()
	editedAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           nil,
			BoxID:                nil,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           nil,
			BoxID:                nil,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		input     UpdateReviewDateAsCompletedInput
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      *UpdateReviewDateAsCompletedOutput
		wantErr   bool
	}{
		{
			name: "最後のステップの復習日完了（復習物も完了）",
			input: UpdateReviewDateAsCompletedInput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				ItemID:       itemID,
				StepNumber:   2,
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetReviewDatesByItemID(gomock.Any(), itemID, userID).
						Return(testReviewdates, nil).
						Times(1),

					mockItemRepo.EXPECT().
						GetEditedAtByItemID(gomock.Any(), itemID, userID).
						Return(editedAt, nil).
						Times(1),

					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),

					mockItemRepo.EXPECT().
						UpdateReviewDateAsCompleted(gomock.Any(), reviewDateID, userID).
						Return(nil).
						Times(1),

					mockItemRepo.EXPECT().
						UpdateItemAsFinished(gomock.Any(), itemID, userID, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateReviewDateAsCompletedOutput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				IsCompleted:  true,
				IsFinished:   true,
				EditedAt:     editedAt,
			},
			wantErr: false,
		},
		{
			name: "最後のステップではない復習日完了",
			input: UpdateReviewDateAsCompletedInput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				ItemID:       itemID,
				StepNumber:   1,
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetReviewDatesByItemID(gomock.Any(), itemID, userID).
						Return(testReviewdates, nil).
						Times(1),

					mockItemRepo.EXPECT().
						GetEditedAtByItemID(gomock.Any(), itemID, userID).
						Return(editedAt, nil).
						Times(1),

					mockItemRepo.EXPECT().
						UpdateReviewDateAsCompleted(gomock.Any(), reviewDateID, userID).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateReviewDateAsCompletedOutput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				IsCompleted:  true,
				IsFinished:   false,
				EditedAt:     editedAt,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.UpdateReviewDateAsCompleted(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateReviewDateAsCompleted() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != nil {
				// 時刻系フィールドは動的に生成されるため、テストでは除外（最後のステップの場合）
				if tc.want.IsFinished {
					got.EditedAt = tc.want.EditedAt
				}

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("UpdateReviewDateAsCompleted() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemUsecase_UpdateReviewDateAsInCompleted(t *testing.T) {
	reviewDateID := uuid.NewString()
	userID := uuid.NewString()
	itemID := uuid.NewString()
	editedAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testFinishedItem := &ItemDomain.Item{
		ItemID:       itemID,
		UserID:       userID,
		CategoryID:   nil,
		BoxID:        nil,
		PatternID:    nil,
		Name:         "Test Item",
		Detail:       "Test Detail",
		LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		IsFinished:   true,
		RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		EditedAt:     editedAt,
	}

	testUnFinishedItem := &ItemDomain.Item{
		ItemID:       itemID,
		UserID:       userID,
		CategoryID:   nil,
		BoxID:        nil,
		PatternID:    nil,
		Name:         "Test Item",
		Detail:       "Test Detail",
		LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		IsFinished:   false,
		RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		EditedAt:     editedAt,
	}

	tests := []struct {
		name      string
		input     UpdateReviewDateAsInCompletedInput
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      *UpdateReviewDateAsInCompletedOutput
		wantErr   bool
	}{
		{
			name: "復習物が完了済みの場合の復習日未完了化",
			input: UpdateReviewDateAsInCompletedInput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				ItemID:       itemID,
				StepNumber:   1,
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetItemByID(gomock.Any(), itemID, userID).
						Return(testFinishedItem, nil).
						Times(1),
					mockItemRepo.EXPECT().
						GetEditedAtByItemID(gomock.Any(), itemID, userID).
						Return(editedAt, nil).
						Times(1),
					mockTransactionManager.EXPECT().
						RunInTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						}).
						Times(1),
					mockItemRepo.EXPECT().
						UpdateReviewDateAsInCompleted(gomock.Any(), reviewDateID, userID).
						Return(nil).
						Times(1),
					mockItemRepo.EXPECT().
						UpdateItemAsUnFinished(gomock.Any(), itemID, userID, gomock.Any()).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateReviewDateAsInCompletedOutput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				IsCompleted:  false,
				IsFinished:   true,
				EditedAt:     editedAt,
			},
			wantErr: false,
		},
		{
			name: "復習物が未完了の場合の復習日未完了化",
			input: UpdateReviewDateAsInCompletedInput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				ItemID:       itemID,
				StepNumber:   1,
			},
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetItemByID(gomock.Any(), itemID, userID).
						Return(testUnFinishedItem, nil).
						Times(1),
					mockItemRepo.EXPECT().
						GetEditedAtByItemID(gomock.Any(), itemID, userID).
						Return(editedAt, nil).
						Times(1),
					mockItemRepo.EXPECT().
						UpdateReviewDateAsInCompleted(gomock.Any(), reviewDateID, userID).
						Return(nil).
						Times(1),
				)
			},
			want: &UpdateReviewDateAsInCompletedOutput{
				ReviewDateID: reviewDateID,
				UserID:       userID,
				IsCompleted:  false,
				IsFinished:   false,
				EditedAt:     editedAt,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.UpdateReviewDateAsInCompleted(context.Background(), tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateReviewDateAsInCompleted() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != nil {
				// 時刻系フィールドは動的に生成されるため、テストでは除外（復習物が完了済みの場合）
				if tc.want.IsFinished {
					got.EditedAt = tc.want.EditedAt
				}

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("UpdateReviewDateAsInCompleted() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemUsecase_UpdateItemAsUnFinishedForce(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	itemID := uuid.NewString()
	patternID := uuid.NewString()
	today := "2024-01-10"
	learnedDate := "2024-01-01"

	testReviewDates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	testPatternSteps := []*PatternDomain.PatternStep{
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    1,
			IntervalDays:  1,
		},
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    2,
			IntervalDays:  3,
		},
	}

	testNewReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         testReviewDates[0].ReviewdateID,
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
		{
			ReviewdateID:         testReviewDates[1].ReviewdateID,
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		input     UpdateItemAsUnFinishedForceInput
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name: "正常系_スケジュール更新が必要な場合",
			input: UpdateItemAsUnFinishedForceInput{
				ItemID:      itemID,
				UserID:      userID,
				CategoryID:  &categoryID,
				BoxID:       &boxID,
				PatternID:   patternID,
				LearnedDate: learnedDate,
				Today:       today,
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(testReviewDates, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompletedWithIDs(
						testPatternSteps,
						[]string{testReviewDates[0].ReviewdateID, testReviewDates[1].ReviewdateID},
						userID,
						&categoryID,
						&boxID,
						itemID,
						gomock.Any(),
						gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItemAsUnFinished(ctx, itemID, userID, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, gomock.Any(), userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(testNewReviewdates, nil).Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.UpdateItemAsUnFinishedForce(ctx, tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateItemAsUnFinishedForce() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got == nil {
				t.Error("UpdateItemAsUnFinishedForce() got = nil, want not nil")
			}
		})
	}
}

// isPatternNotNilToNil = false の場合のテスト
func TestItemUsecase_UpdateItem_PatternNotNilToNil(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
	mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
	mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
	mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
	mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
	mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

	usecase := NewItemUsecase(
		mockCategoryRepo,
		mockBoxRepo,
		mockItemRepo,
		mockPatternRepo,
		mockTransactionManager,
		mockScheduler,
	)

	userID := uuid.NewString()
	itemID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	patternID := uuid.NewString()

	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	currentItem := &ItemDomain.Item{
		ItemID:      itemID,
		UserID:      userID,
		CategoryID:  &categoryID,
		BoxID:       &boxID,
		PatternID:   &patternID,
		Name:        "Original Item",
		Detail:      "Original Detail",
		LearnedDate: learnedDate,
		IsFinished:  false,
	}

	// モック設定
	gomock.InOrder(
		mockItemRepo.EXPECT().GetItemByID(gomock.Any(), itemID, userID).Return(currentItem, nil).Times(1),
		mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(gomock.Any(), itemID, userID).Return(false, nil).Times(1),
		mockTransactionManager.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		).Times(1),
		mockItemRepo.EXPECT().UpdateItem(gomock.Any(), gomock.Any()).Return(nil).Times(1),
		mockItemRepo.EXPECT().DeleteReviewDates(gomock.Any(), itemID, userID).Return(nil).Times(1),
	)

	input := UpdateItemInput{
		ItemID:                   itemID,
		UserID:                   userID,
		CategoryID:               &categoryID,
		BoxID:                    &boxID,
		PatternID:                nil,
		Name:                     "Item",
		Detail:                   "Detail",
		LearnedDate:              "2024-01-01",
		IsMarkOverdueAsCompleted: false,
		Today:                    "2024-01-10",
	}

	got, err := usecase.UpdateItem(context.Background(), input)
	if err != nil {
		t.Errorf("UpdateItem() error = %v", err)
		return
	}

	if got.ItemID != itemID || got.PatternID != nil {
		t.Errorf("UpdateItem() result mismatch")
	}
}

func TestItemUsecase_UpdateItem_PatternNilToNotNil(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	itemID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	patternID := uuid.NewString()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	currentItem := &ItemDomain.Item{
		ItemID:      itemID,
		UserID:      userID,
		CategoryID:  &categoryID,
		BoxID:       &boxID,
		PatternID:   nil,
		Name:        "Original Item",
		Detail:      "Original Detail",
		LearnedDate: learnedDate,
		IsFinished:  false,
	}

	testPatternSteps := []*PatternDomain.PatternStep{
		{StepNumber: 1, IntervalDays: 1},
		{StepNumber: 2, IntervalDays: 3},
	}

	testNewReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		input     UpdateItemInput
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name: "PatternNilToNotNil_未完了で上書きマーク無し",
			input: UpdateItemInput{
				ItemID:                   itemID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Updated Item",
				Detail:                   "Updated Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompleted(
						testPatternSteps, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().CreateReviewdates(ctx, testNewReviewdates).Return(int64(0), nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "PatternNilToNotNil_未完了で上書きマーク有り",
			input: UpdateItemInput{
				ItemID:                   itemID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Updated Item",
				Detail:                   "Updated Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: true,
				Today:                    "2024-01-10",
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompleted(
						testPatternSteps, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, false, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().CreateReviewdates(ctx, testNewReviewdates).Return(int64(0), nil).Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItem_PatternNotNilToNotNil_LengthDiff(t *testing.T) {
	ctx := context.Background()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler) (UpdateItemInput, bool)
	}{
		{
			name: "PatternStepsLength異なる場合（未完了で上書きマーク無し）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
					{StepNumber: 3, IntervalDays: 7},
				}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           3,
						InitialScheduledDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompleted(
						newPatternSteps, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().DeleteReviewDates(ctx, itemID, userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().CreateReviewdates(ctx, testNewReviewdates).Return(int64(0), nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-01",
					IsMarkOverdueAsCompleted: false,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
		{
			name: "PatternStepsLength異なる場合（未完了で上書きマーク有り）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
					{StepNumber: 3, IntervalDays: 7},
				}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         uuid.NewString(),
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           3,
						InitialScheduledDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompleted(
						newPatternSteps, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, false, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
				)
				mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1)
				mockItemRepo.EXPECT().DeleteReviewDates(ctx, itemID, userID).Return(nil).Times(1)
				mockItemRepo.EXPECT().CreateReviewdates(ctx, testNewReviewdates).Return(int64(0), nil).Times(1)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-01",
					IsMarkOverdueAsCompleted: true,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
		{
			name: "PatternStepsLength異なる場合_HasCompletedReviewDateエラー",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
					{StepNumber: 3, IntervalDays: 7},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(true, nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-01",
					IsMarkOverdueAsCompleted: false,
					Today:                    "2024-01-10",
				}
				return input, true
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			input, wantErr := tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, input)

			if (err != nil) != wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItem_PatternNotNilToNotNil_IntervalDaysDiff(t *testing.T) {
	ctx := context.Background()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler) (UpdateItemInput, bool)
	}{
		{
			name: "PatternStepsIntervalDays異なる場合（未完了で上書きマーク無し）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 2},
					{StepNumber: 2, IntervalDays: 5},
				}

				reviewDateIDs := []string{uuid.NewString(), uuid.NewString()}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         reviewDateIDs[0],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         reviewDateIDs[1],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(reviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompletedWithIDs(
						newPatternSteps, reviewDateIDs, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, testNewReviewdates, userID).Return(nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-01",
					IsMarkOverdueAsCompleted: false,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
		{
			name: "PatternStepsIntervalDays異なる場合（未完了で上書きマーク有り）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 2},
					{StepNumber: 2, IntervalDays: 5},
				}

				reviewDateIDs := []string{uuid.NewString(), uuid.NewString()}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         reviewDateIDs[0],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         reviewDateIDs[1],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(reviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompletedWithIDs(
						newPatternSteps, reviewDateIDs, userID, &categoryID, &boxID, itemID, learnedDate, gomock.Any(),
					).Return(testNewReviewdates, false, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, testNewReviewdates, userID).Return(nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-01",
					IsMarkOverdueAsCompleted: true,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			input, wantErr := tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, input)

			if (err != nil) != wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItem_LearnedDateChanged(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	itemID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	patternID := uuid.NewString()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	currentItem := &ItemDomain.Item{
		ItemID:      itemID,
		UserID:      userID,
		CategoryID:  &categoryID,
		BoxID:       &boxID,
		PatternID:   &patternID,
		Name:        "Original Item",
		Detail:      "Original Detail",
		LearnedDate: learnedDate,
		IsFinished:  false,
	}

	patternSteps := []*PatternDomain.PatternStep{
		{StepNumber: 1, IntervalDays: 1},
		{StepNumber: 2, IntervalDays: 3},
	}

	reviewDateIDs := []string{uuid.NewString(), uuid.NewString()}

	testNewReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         reviewDateIDs[0],
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
		{
			ReviewdateID:         reviewDateIDs[1],
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		input     UpdateItemInput
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name: "LearnedDate変更_SamePatternID",
			input: UpdateItemInput{
				ItemID:                   itemID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Updated Item",
				Detail:                   "Updated Detail",
				LearnedDate:              "2024-01-02",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(reviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompletedWithIDs(
						patternSteps, reviewDateIDs, userID, &categoryID, &boxID, itemID, gomock.Any(), gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, testNewReviewdates, userID).Return(nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "LearnedDate変更_HasCompletedReviewDateエラー",
			input: UpdateItemInput{
				ItemID:                   itemID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				PatternID:                &patternID,
				Name:                     "Updated Item",
				Detail:                   "Updated Detail",
				LearnedDate:              "2024-01-02",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(true, nil).Times(1),
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

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItem_SamePatternStepsStructure_LearnedDateChanged(t *testing.T) {
	ctx := context.Background()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler) (UpdateItemInput, bool)
	}{
		{
			name: "SamePatternStepsStructure_LearnedDateChanged（未完了で上書きマーク無し）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				// 同じpatternStepsStructure（stepNumberとintervalDaysが同じ）
				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				reviewDateIDs := []string{uuid.NewString(), uuid.NewString()}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         reviewDateIDs[0],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         reviewDateIDs[1],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(reviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompletedWithIDs(
						newPatternSteps, reviewDateIDs, userID, &categoryID, &boxID, itemID, gomock.Any(), gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, testNewReviewdates, userID).Return(nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-02", // LearnedDateが変更されている
					IsMarkOverdueAsCompleted: false,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
		{
			name: "SamePatternStepsStructure_LearnedDateChanged（未完了で上書きマーク有り）",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				// 同じpatternStepsStructure（stepNumberとintervalDaysが同じ）
				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				reviewDateIDs := []string{uuid.NewString(), uuid.NewString()}

				testNewReviewdates := []*ItemDomain.Reviewdate{
					{
						ReviewdateID:         reviewDateIDs[0],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           1,
						InitialScheduledDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
					{
						ReviewdateID:         reviewDateIDs[1],
						UserID:               userID,
						CategoryID:           &categoryID,
						BoxID:                &boxID,
						ItemID:               itemID,
						StepNumber:           2,
						InitialScheduledDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						IsCompleted:          false,
					},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(false, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(reviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompletedWithIDs(
						newPatternSteps, reviewDateIDs, userID, &categoryID, &boxID, itemID, gomock.Any(), gomock.Any(),
					).Return(testNewReviewdates, false, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, testNewReviewdates, userID).Return(nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-02", // LearnedDateが変更されている
					IsMarkOverdueAsCompleted: true,
					Today:                    "2024-01-10",
				}
				return input, false
			},
		},
		{
			name: "SamePatternStepsStructure_LearnedDateChanged_HasCompletedReviewDateエラー",
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) (UpdateItemInput, bool) {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				categoryID := uuid.NewString()
				boxID := uuid.NewString()
				currentPatternID := uuid.NewString()
				newPatternID := uuid.NewString()

				currentItem := &ItemDomain.Item{
					ItemID:      itemID,
					UserID:      userID,
					CategoryID:  &categoryID,
					BoxID:       &boxID,
					PatternID:   &currentPatternID,
					Name:        "Original Item",
					Detail:      "Original Detail",
					LearnedDate: learnedDate,
					IsFinished:  false,
				}

				// 同じpatternStepsStructure（stepNumberとintervalDaysが同じ）
				currentPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				newPatternSteps := []*PatternDomain.PatternStep{
					{StepNumber: 1, IntervalDays: 1},
					{StepNumber: 2, IntervalDays: 3},
				}

				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, currentPatternID, userID).Return(currentPatternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, newPatternID, userID).Return(newPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().HasCompletedReviewDateByItemID(ctx, itemID, userID).Return(true, nil).Times(1),
				)

				input := UpdateItemInput{
					ItemID:                   itemID,
					UserID:                   userID,
					CategoryID:               &categoryID,
					BoxID:                    &boxID,
					PatternID:                &newPatternID,
					Name:                     "Updated Item",
					Detail:                   "Updated Detail",
					LearnedDate:              "2024-01-02", // LearnedDateが変更されている
					IsMarkOverdueAsCompleted: false,
					Today:                    "2024-01-10",
				}
				return input, true
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			input, wantErr := tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, input)

			if (err != nil) != wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateItem_CategoryIDBoxIDUpdate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	itemID := uuid.NewString()
	categoryID := uuid.NewString()
	newCategoryID := uuid.NewString()
	boxID := uuid.NewString()
	newBoxID := uuid.NewString()
	patternID := uuid.NewString()
	learnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	currentItem := &ItemDomain.Item{
		ItemID:      itemID,
		UserID:      userID,
		CategoryID:  &categoryID,
		BoxID:       &boxID,
		PatternID:   &patternID,
		Name:        "Original Item",
		Detail:      "Original Detail",
		LearnedDate: learnedDate,
		IsFinished:  false,
	}

	patternSteps := []*PatternDomain.PatternStep{
		{StepNumber: 1, IntervalDays: 1},
		{StepNumber: 2, IntervalDays: 3},
	}

	currentReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		input     UpdateItemInput
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name: "CategoryIDとBoxIDのみ更新",
			input: UpdateItemInput{
				ItemID:                   itemID,
				UserID:                   userID,
				CategoryID:               &newCategoryID,
				BoxID:                    &newBoxID,
				PatternID:                &patternID,
				Name:                     "Original Item",
				Detail:                   "Original Detail",
				LearnedDate:              "2024-01-01",
				IsMarkOverdueAsCompleted: false,
				Today:                    "2024-01-10",
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetItemByID(ctx, itemID, userID).Return(currentItem, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(patternSteps, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(currentReviewdates, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItem(ctx, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDates(ctx, gomock.Any(), userID).Return(nil).Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			_, err := usecase.UpdateItem(ctx, tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestItemUsecase_UpdateReviewDates(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	itemID := uuid.NewString()
	patternID := uuid.NewString()
	reviewDateID := uuid.NewString()
	today := "2024-01-10"
	learnedDate := "2024-01-01"
	initialScheduledDate := "2024-01-02"
	requestScheduledDate := "2024-01-05"

	testPatternSteps := []*PatternDomain.PatternStep{
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    1,
			IntervalDays:  1,
		},
		{
			PatternStepID: uuid.NewString(),
			UserID:        userID,
			PatternID:     patternID,
			StepNumber:    2,
			IntervalDays:  3,
		},
	}

	testReviewDateIDs := []string{reviewDateID, uuid.NewString()}

	testNewReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         reviewDateID,
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
		{
			ReviewdateID:         testReviewDateIDs[1],
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	editedAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	// 最終ステップ用のテストデータ
	testFinalReviewdate := &ItemDomain.Reviewdate{
		ReviewdateID:         reviewDateID,
		UserID:               userID,
		CategoryID:           &categoryID,
		BoxID:                &boxID,
		ItemID:               itemID,
		StepNumber:           2,
		InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		IsCompleted:          true,
	}

	// Overdue完了マーク用のテストデータ
	testOverdueCompletedReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         reviewDateID,
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
		{
			ReviewdateID:         testReviewDateIDs[1],
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           2,
			InitialScheduledDate: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
	}

	tests := []struct {
		name      string
		input     UpdateBackReviewDateInput
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name: "正常系_非最終ステップでOverdueマーク無し",
			input: UpdateBackReviewDateInput{
				ReviewDateID:             reviewDateID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				ItemID:                   itemID,
				StepNumber:               1,
				InitialScheduledDate:     initialScheduledDate,
				RequestScheduledDate:     requestScheduledDate,
				IsMarkOverdueAsCompleted: false,
				Today:                    today,
				LearnedDate:              learnedDate,
				PatternID:                patternID,
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(testReviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates(
						testPatternSteps,
						testReviewDateIDs,
						userID,
						&categoryID,
						&boxID,
						itemID,
						gomock.Any(),
						gomock.Any(),
					).Return(testNewReviewdates, nil).Times(1),
					mockItemRepo.EXPECT().GetEditedAtByItemID(ctx, itemID, userID).Return(editedAt, nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDatesBack(ctx, gomock.Any(), userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(testNewReviewdates, nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "正常系_非最終ステップでOverdueマーク有り_isFinished=false",
			input: UpdateBackReviewDateInput{
				ReviewDateID:             reviewDateID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				ItemID:                   itemID,
				StepNumber:               1,
				InitialScheduledDate:     initialScheduledDate,
				RequestScheduledDate:     requestScheduledDate,
				IsMarkOverdueAsCompleted: true,
				Today:                    today,
				LearnedDate:              learnedDate,
				PatternID:                patternID,
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(testReviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompletedWithIDs(
						testPatternSteps,
						testReviewDateIDs,
						userID,
						&categoryID,
						&boxID,
						itemID,
						gomock.Any(),
						gomock.Any(),
					).Return(testNewReviewdates, false, nil).Times(1),
					mockItemRepo.EXPECT().GetEditedAtByItemID(ctx, itemID, userID).Return(editedAt, nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDatesBack(ctx, gomock.Any(), userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(testNewReviewdates, nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "正常系_非最終ステップでOverdueマーク有り_isFinished=true",
			input: UpdateBackReviewDateInput{
				ReviewDateID:             reviewDateID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				ItemID:                   itemID,
				StepNumber:               1,
				InitialScheduledDate:     initialScheduledDate,
				RequestScheduledDate:     requestScheduledDate,
				IsMarkOverdueAsCompleted: true,
				Today:                    today,
				LearnedDate:              learnedDate,
				PatternID:                patternID,
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDateIDsByItemID(ctx, itemID, userID).Return(testReviewDateIDs, nil).Times(1),
					mockScheduler.EXPECT().FormatWithOverdueMarkedCompletedWithIDs(
						testPatternSteps,
						testReviewDateIDs,
						userID,
						&categoryID,
						&boxID,
						itemID,
						gomock.Any(),
						gomock.Any(),
					).Return(testOverdueCompletedReviewdates, true, nil).Times(1),
					mockItemRepo.EXPECT().GetEditedAtByItemID(ctx, itemID, userID).Return(editedAt, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItemAsFinished(ctx, itemID, userID, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDatesBack(ctx, gomock.Any(), userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return(testOverdueCompletedReviewdates, nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "正常系_最終ステップ",
			input: UpdateBackReviewDateInput{
				ReviewDateID:             reviewDateID,
				UserID:                   userID,
				CategoryID:               &categoryID,
				BoxID:                    &boxID,
				ItemID:                   itemID,
				StepNumber:               2, // 最終ステップ
				InitialScheduledDate:     initialScheduledDate,
				RequestScheduledDate:     requestScheduledDate,
				IsMarkOverdueAsCompleted: false,
				Today:                    today,
				LearnedDate:              learnedDate,
				PatternID:                patternID,
			},
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockPatternRepo.EXPECT().GetAllPatternStepsByPatternID(ctx, patternID, userID).Return(testPatternSteps, nil).Times(1),
					mockItemRepo.EXPECT().GetEditedAtByItemID(ctx, itemID, userID).Return(editedAt, nil).Times(1),
					mockTransactionManager.EXPECT().RunInTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					).Times(1),
					mockItemRepo.EXPECT().UpdateItemAsFinished(ctx, itemID, userID, gomock.Any()).Return(nil).Times(1),
					mockItemRepo.EXPECT().UpdateReviewDatesBack(ctx, gomock.Any(), userID).Return(nil).Times(1),
					mockItemRepo.EXPECT().GetReviewDatesByItemID(ctx, itemID, userID).Return([]*ItemDomain.Reviewdate{testFinalReviewdate}, nil).Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.UpdateReviewDates(ctx, tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateReviewDates() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got == nil {
				t.Error("UpdateReviewDates() got = nil, want not nil")
			}
		})
	}
}

// 取得系メソッドのテスト
func TestItemUsecase_GetAllUnFinishedItemsByBoxID(t *testing.T) {
	userID := uuid.NewString()
	boxID := uuid.NewString()

	testItems := []*ItemDomain.Item{
		{
			ItemID:       uuid.NewString(),
			UserID:       userID,
			CategoryID:   nil,
			BoxID:        &boxID,
			PatternID:    nil,
			Name:         "Test Item",
			Detail:       "Test Detail",
			LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			IsFinished:   false,
			RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			EditedAt:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		},
	}

	testReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           nil,
			BoxID:                &boxID,
			ItemID:               testItems[0].ItemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		boxID     string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      []*GetItemOutput
		wantErr   bool
	}{
		{
			name:   "正常系",
			boxID:  boxID,
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().
						GetAllUnFinishedItemsByBoxID(gomock.Any(), boxID, userID).
						Return(testItems, nil).
						Times(1),
					mockItemRepo.EXPECT().
						GetAllReviewDatesByBoxID(gomock.Any(), boxID, userID).
						Return(testReviewdates, nil).
						Times(1),
				)
			},
			want: []*GetItemOutput{
				{
					ItemID:       testItems[0].ItemID,
					UserID:       userID,
					CategoryID:   nil,
					BoxID:        &boxID,
					PatternID:    nil,
					Name:         "Test Item",
					Detail:       "Test Detail",
					LearnedDate:  "2024-01-01",
					IsFinished:   false,
					RegisteredAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					EditedAt:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					ReviewDates: []GetReviewDateOutput{
						{
							ReviewDateID:         testReviewdates[0].ReviewdateID,
							UserID:               userID,
							CategoryID:           nil,
							BoxID:                &boxID,
							ItemID:               testItems[0].ItemID,
							StepNumber:           1,
							InitialScheduledDate: "2024-01-02",
							ScheduledDate:        "2024-01-02",
							IsCompleted:          false,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.GetAllUnFinishedItemsByBoxID(context.Background(), tc.boxID, tc.userID)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetAllUnFinishedItemsByBoxID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("GetAllUnFinishedItemsByBoxID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemUsecase_CountAllDailyReviewDates(t *testing.T) {
	userID := uuid.NewString()
	today := "2024-01-10"
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		userID    string
		today     string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			today:  today,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().
					CountAllDailyReviewDates(gomock.Any(), userID, parsedToday).
					Return(15, nil).
					Times(1)
			},
			want:    15,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)

			got, err := usecase.CountAllDailyReviewDates(context.Background(), tc.userID, tc.today)

			if (err != nil) != tc.wantErr {
				t.Errorf("CountAllDailyReviewDates() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got != tc.want {
				t.Errorf("CountAllDailyReviewDates() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestItemUsecase_GetAllUnFinishedUnclassifiedItemsByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	itemID := uuid.NewString()
	parsedLearnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	registeredAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testItems := []*ItemDomain.Item{
		{
			ItemID:       itemID,
			UserID:       userID,
			CategoryID:   nil,
			BoxID:        nil,
			PatternID:    nil,
			Name:         "Test Item",
			Detail:       "Test Detail",
			LearnedDate:  parsedLearnedDate,
			IsFinished:   false,
			RegisteredAt: registeredAt,
			EditedAt:     registeredAt,
		},
	}

	testReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           nil,
			BoxID:                nil,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetAllUnFinishedUnclassifiedItemsByUserID(gomock.Any(), userID).Return(testItems, nil).Times(1),
					mockItemRepo.EXPECT().GetAllUnclassifiedReviewDatesByUserID(gomock.Any(), userID).Return(testReviewdates, nil).Times(1),
				)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetAllUnFinishedUnclassifiedItemsByUserID(context.Background(), tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetAllUnFinishedUnclassifiedItemsByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("GetAllUnFinishedUnclassifiedItemsByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_GetAllUnFinishedUnclassifiedItemsByCategoryID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	categoryID := uuid.NewString()
	itemID := uuid.NewString()
	parsedLearnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	registeredAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testItems := []*ItemDomain.Item{
		{
			ItemID:       itemID,
			UserID:       userID,
			CategoryID:   &categoryID,
			BoxID:        nil,
			PatternID:    nil,
			Name:         "Test Item",
			Detail:       "Test Detail",
			LearnedDate:  parsedLearnedDate,
			IsFinished:   false,
			RegisteredAt: registeredAt,
			EditedAt:     registeredAt,
		},
	}

	testReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                nil,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			IsCompleted:          false,
		},
	}

	tests := []struct {
		name       string
		userID     string
		categoryID string
		mockSetup  func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "正常系",
			userID:     userID,
			categoryID: categoryID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetAllUnFinishedUnclassifiedItemsByCategoryID(gomock.Any(), categoryID, userID).Return(testItems, nil).Times(1),
					mockItemRepo.EXPECT().GetAllUnclassifiedReviewDatesByCategoryID(gomock.Any(), categoryID, userID).Return(testReviewdates, nil).Times(1),
				)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetAllUnFinishedUnclassifiedItemsByCategoryID(context.Background(), tc.userID, tc.categoryID)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetAllUnFinishedUnclassifiedItemsByCategoryID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("GetAllUnFinishedUnclassifiedItemsByCategoryID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_CountItemsGroupedByBoxByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()

	testCounts := []*ItemDomain.ItemCountGroupedByBox{
		{
			CategoryID: categoryID,
			BoxID:      boxID,
			Count:      5,
		},
	}

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountItemsGroupedByBoxByUserID(gomock.Any(), userID).Return(testCounts, nil).Times(1)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountItemsGroupedByBoxByUserID(context.Background(), tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountItemsGroupedByBoxByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("CountItemsGroupedByBoxByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_CountUnclassifiedItemsGroupedByCategoryByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	categoryID := uuid.NewString()

	testCounts := []*ItemDomain.UnclassifiedItemCountGroupedByCategory{
		{
			CategoryID: categoryID,
			Count:      3,
		},
	}

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountUnclassifiedItemsGroupedByCategoryByUserID(gomock.Any(), userID).Return(testCounts, nil).Times(1)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountUnclassifiedItemsGroupedByCategoryByUserID(context.Background(), tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountUnclassifiedItemsGroupedByCategoryByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("CountUnclassifiedItemsGroupedByCategoryByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_CountUnclassifiedItemsByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountUnclassifiedItemsByUserID(gomock.Any(), userID).Return(10, nil).Times(1)
			},
			want:    10,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountUnclassifiedItemsByUserID(context.Background(), tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountUnclassifiedItemsByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("CountUnclassifiedItemsByUserID() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestItemUsecase_CountDailyDatesGroupedByBoxByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	today := "2024-01-10"
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	categoryID := uuid.NewString()
	boxID := uuid.NewString()

	testCounts := []*ItemDomain.DailyCountGroupedByBox{
		{
			CategoryID: categoryID,
			BoxID:      boxID,
			Count:      3,
		},
	}

	tests := []struct {
		name      string
		userID    string
		today     string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			today:  today,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountDailyDatesGroupedByBoxByUserID(gomock.Any(), userID, parsedToday).Return(testCounts, nil).Times(1)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountDailyDatesGroupedByBoxByUserID(context.Background(), tc.userID, tc.today)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountDailyDatesGroupedByBoxByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("CountDailyDatesGroupedByBoxByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_CountDailyDatesUnclassifiedGroupedByCategoryByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	today := "2024-01-10"
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	categoryID := uuid.NewString()

	testCounts := []*ItemDomain.UnclassifiedDailyDatesCountGroupedByCategory{
		{
			CategoryID: categoryID,
			Count:      2,
		},
	}

	tests := []struct {
		name      string
		userID    string
		today     string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			today:  today,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountDailyDatesUnclassifiedGroupedByCategoryByUserID(gomock.Any(), userID, parsedToday).Return(testCounts, nil).Times(1)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountDailyDatesUnclassifiedGroupedByCategoryByUserID(context.Background(), tc.userID, tc.today)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountDailyDatesUnclassifiedGroupedByCategoryByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("CountDailyDatesUnclassifiedGroupedByCategoryByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_CountDailyDatesUnclassifiedByUserID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	today := "2024-01-10"
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		userID    string
		today     string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		want      int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			today:  today,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				mockItemRepo.EXPECT().CountDailyDatesUnclassifiedByUserID(gomock.Any(), userID, parsedToday).Return(5, nil).Times(1)
			},
			want:    5,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.CountDailyDatesUnclassifiedByUserID(context.Background(), tc.userID, tc.today)
			if (err != nil) != tc.wantErr {
				t.Errorf("CountDailyDatesUnclassifiedByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("CountDailyDatesUnclassifiedByUserID() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestItemUsecase_GetAllDailyReviewDates(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	categoryID := uuid.NewString()
	boxID := uuid.NewString()
	patternID := uuid.NewString()
	today := "2024-01-10"
	parsedToday := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	testDailyReviewDates := []*ItemDomain.DailyReviewDate{
		{
			ReviewdateID:         uuid.NewString(),
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               uuid.NewString(),
			StepNumber:           1,
			InitialScheduledDate: parsedToday,
			PrevScheduledDate:    &parsedToday,
			ScheduledDate:        parsedToday,
			NextScheduledDate:    &parsedToday,
			IsCompleted:          false,
			Name:                 "Test Item",
			Detail:               "Test Detail",
			LearnedDate:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			RegisteredAt:         time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
			EditedAt:             time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		},
	}

	testCategoryNames := []*CategoryDomain.CategoryName{
		{
			ID:   categoryID,
			Name: "Test Category",
		},
	}

	testBoxNames := []*BoxDomain.BoxName{
		{
			BoxID:     boxID,
			Name:      "Test Box",
			PatternID: patternID,
		},
	}

	testTargetWeights := []*PatternDomain.TargetWeight{
		{
			PatternID:    patternID,
			TargetWeight: "Easy",
		},
	}

	tests := []struct {
		name      string
		userID    string
		today     string
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			today:  today,
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetAllDailyReviewDates(ctx, userID, parsedToday).Return(testDailyReviewDates, nil).Times(1),
					mockCategoryRepo.EXPECT().GetCategoryNamesByCategoryIDs(ctx, []string{categoryID}).Return(testCategoryNames, nil).Times(1),
					mockBoxRepo.EXPECT().GetBoxNamesByBoxIDs(ctx, []string{boxID}).Return(testBoxNames, nil).Times(1),
					mockPatternRepo.EXPECT().GetPatternTargetWeightsByPatternIDs(ctx, []string{patternID}).Return(testTargetWeights, nil).Times(1),
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetAllDailyReviewDates(ctx, tc.userID, tc.today)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetAllDailyReviewDates() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got == nil {
				t.Error("GetAllDailyReviewDates() got = nil, want not nil")
			}
		})
	}
}

func TestItemUsecase_GetFinishedItemsByBoxID(t *testing.T) {
	// テストデータの準備
	userID := uuid.NewString()
	boxID := uuid.NewString()
	categoryID := uuid.NewString()
	itemID := uuid.NewString()
	patternID := uuid.NewString()
	parsedLearnedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	registeredAt := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)

	testItems := []*ItemDomain.Item{
		{
			ItemID:       itemID,
			UserID:       userID,
			CategoryID:   &categoryID,
			BoxID:        &boxID,
			PatternID:    &patternID,
			Name:         "Finished Item",
			Detail:       "Finished Detail",
			LearnedDate:  parsedLearnedDate,
			IsFinished:   true,
			RegisteredAt: registeredAt,
			EditedAt:     registeredAt,
		},
	}

	testReviewdates := []*ItemDomain.Reviewdate{
		{
			ReviewdateID:         uuid.NewString(),
			UserID:               userID,
			CategoryID:           &categoryID,
			BoxID:                &boxID,
			ItemID:               itemID,
			StepNumber:           1,
			InitialScheduledDate: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			ScheduledDate:        time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			IsCompleted:          true,
		},
	}

	tests := []struct {
		name      string
		boxID     string
		userID    string
		mockSetup func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			boxID:  boxID,
			userID: userID,
			mockSetup: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetFinishedItemsByBoxID(gomock.Any(), boxID, userID).Return(testItems, nil).Times(1),
					mockItemRepo.EXPECT().GetAllReviewDatesByBoxID(gomock.Any(), boxID, userID).Return(testReviewdates, nil).Times(1),
				)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.mockSetup(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetFinishedItemsByBoxID(context.Background(), tc.boxID, tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetFinishedItemsByBoxID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("GetFinishedItemsByBoxID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_GetUnclassfiedFinishedItemsByCategoryID(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	categoryID := uuid.NewString()

	testItems := []*ItemDomain.Item{
		{
			ItemID:       uuid.NewString(),
			UserID:       userID,
			CategoryID:   &categoryID,
			BoxID:        nil,
			PatternID:    nil,
			Name:         "Test Item 1",
			Detail:       "Test Detail 1",
			LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			IsFinished:   true,
			RegisteredAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
			EditedAt:     time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name       string
		categoryID string
		userID     string
		setupMock  func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "正常系",
			categoryID: categoryID,
			userID:     userID,
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetUnclassfiedFinishedItemsByCategoryID(ctx, categoryID, userID).Return(testItems, nil).Times(1),
					mockItemRepo.EXPECT().GetAllUnclassifiedReviewDatesByCategoryID(ctx, categoryID, userID).Return([]*ItemDomain.Reviewdate{}, nil).Times(1),
				)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetUnclassfiedFinishedItemsByCategoryID(ctx, tc.userID, tc.categoryID)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetUnclassfiedFinishedItemsByCategoryID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("GetUnclassfiedFinishedItemsByCategoryID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}

func TestItemUsecase_GetUnclassfiedFinishedItemsByUserID(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	categoryID := uuid.NewString()

	testItems := []*ItemDomain.Item{
		{
			ItemID:       uuid.NewString(),
			UserID:       userID,
			CategoryID:   &categoryID,
			BoxID:        nil,
			PatternID:    nil,
			Name:         "Test Item 1",
			Detail:       "Test Detail 1",
			LearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			IsFinished:   true,
			RegisteredAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
			EditedAt:     time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name      string
		userID    string
		setupMock func(*CategoryDomain.MockICategoryRepository, *BoxDomain.MockIBoxRepository, *ItemDomain.MockIItemRepository, *PatternDomain.MockIPatternRepository, *transaction.MockITransactionManager, *ItemDomain.MockIScheduler)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "正常系",
			userID: userID,
			setupMock: func(mockCategoryRepo *CategoryDomain.MockICategoryRepository, mockBoxRepo *BoxDomain.MockIBoxRepository, mockItemRepo *ItemDomain.MockIItemRepository, mockPatternRepo *PatternDomain.MockIPatternRepository, mockTransactionManager *transaction.MockITransactionManager, mockScheduler *ItemDomain.MockIScheduler) {
				gomock.InOrder(
					mockItemRepo.EXPECT().GetUnclassfiedFinishedItemsByUserID(ctx, userID).Return(testItems, nil).Times(1),
					mockItemRepo.EXPECT().GetAllUnclassifiedReviewDatesByUserID(ctx, userID).Return([]*ItemDomain.Reviewdate{}, nil).Times(1),
				)
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCategoryRepo := CategoryDomain.NewMockICategoryRepository(ctrl)
			mockBoxRepo := BoxDomain.NewMockIBoxRepository(ctrl)
			mockItemRepo := ItemDomain.NewMockIItemRepository(ctrl)
			mockPatternRepo := PatternDomain.NewMockIPatternRepository(ctrl)
			mockTransactionManager := transaction.NewMockITransactionManager(ctrl)
			mockScheduler := ItemDomain.NewMockIScheduler(ctrl)

			usecase := NewItemUsecase(
				mockCategoryRepo,
				mockBoxRepo,
				mockItemRepo,
				mockPatternRepo,
				mockTransactionManager,
				mockScheduler,
			)

			tc.setupMock(mockCategoryRepo, mockBoxRepo, mockItemRepo, mockPatternRepo, mockTransactionManager, mockScheduler)
			got, err := usecase.GetUnclassfiedFinishedItemsByUserID(ctx, tc.userID)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetUnclassfiedFinishedItemsByUserID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != tc.wantLen {
				t.Errorf("GetUnclassfiedFinishedItemsByUserID() got len = %v, want %v", len(got), tc.wantLen)
			}
		})
	}
}
