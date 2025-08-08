package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	itemDomain "github.com/minminseo/recall-setter/domain/item"
)

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestItemRepository_CreateItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	now := time.Now()
	tests := []struct {
		name    string
		item    *itemDomain.Item
		wantErr bool
	}{
		{
			name: "ボックス・パターンありで復習物作成に成功する場合",
			item: func() *itemDomain.Item {
				item, _ := itemDomain.NewItem(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					stringPtr("750e8400-e29b-41d4-a716-446655440001"),
					"新しい数学問題",
					"微分積分の応用問題",
					time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name: "ボックスなし（未分類）で復習物作成に成功する場合",
			item: func() *itemDomain.Item {
				item, _ := itemDomain.NewItem(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					nil,
					nil,
					"未分類の問題",
					"まだボックスに分類されていない問題",
					time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name: "存在しないカテゴリによる外部キー制約違反",
			item: func() *itemDomain.Item {
				item, _ := itemDomain.NewItem(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr(uuid.New().String()), // Does not exist in fixture
					nil,
					nil,
					"存在しないカテゴリの問題",
					"詳細",
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.CreateItem(ctx, tc.item)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestItemRepository_CreateReviewdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name        string
		reviewdates []*itemDomain.Reviewdate
		want        []*itemDomain.Reviewdate
		wantCount   int64
		wantErr     bool
	}{
		{
			name: "復習日の作成に成功する場合",
			reviewdates: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"c50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						3,
						time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"c50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						4,
						time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"c50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						3,
						time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"c50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						4,
						time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			count, err := repo.CreateReviewdates(ctx, tc.reviewdates)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if count != tc.wantCount {
				t.Errorf("期待される作成数: %d, 実際: %d", tc.wantCount, count)
				return
			}

			// 作成された復習日を取得して比較
			createdReviewdates, err := repo.GetReviewDatesByItemID(ctx, tc.reviewdates[0].ItemID(), tc.reviewdates[0].UserID())
			if err != nil {
				t.Errorf("作成された復習日の取得に失敗: %v", err)
				return
			}

			// 作成された復習日の中から今回作成したもののみを抽出
			var actualReviewdates []*itemDomain.Reviewdate
			for _, created := range createdReviewdates {
				for _, expected := range tc.want {
					if created.ReviewdateID() == expected.ReviewdateID() {
						actualReviewdates = append(actualReviewdates, created)
						break
					}
				}
			}

			if diff := cmp.Diff(tc.want, actualReviewdates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("CreateReviewdates() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetItemByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    *itemDomain.Item
		wantErr bool
	}{
		{
			name:   "有効なIDで復習物を取得する場合",
			itemID: "a50e8400-e29b-41d4-a716-446655440001",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: func() *itemDomain.Item {
				item, _ := itemDomain.ReconstructItem(
					"a50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					stringPtr("750e8400-e29b-41d4-a716-446655440001"),
					"二次方程式",
					"ax^2 + bx + c = 0の解の公式を覚える",
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					false,
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:    "存在しない復習物を取得する場合",
			itemID:  uuid.New().String(),
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			item, err := repo.GetItemByID(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if item != nil {
					t.Error("復習物がnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if item == nil {
				t.Error("復習物がnilです")
				return
			}

			if tc.want != nil {
				if diff := cmp.Diff(tc.want, item, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
					t.Errorf("GetItemByID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_HasCompletedReviewDateByItemID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    bool
		wantErr bool
	}{
		{
			name:    "完了した復習日を持つ復習物の場合",
			itemID:  "a50e8400-e29b-41d4-a716-446655440002",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    true,
			wantErr: false,
		},
		{
			name:    "完了した復習日を持たない復習物の場合",
			itemID:  "a50e8400-e29b-41d4-a716-446655440001",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    false,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			hasCompleted, err := repo.HasCompletedReviewDateByItemID(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if hasCompleted != tc.want {
				t.Errorf("期待される結果: %v, 実際: %v", tc.want, hasCompleted)
			}
		})
	}
}

func TestItemRepository_GetReviewDateIDsByItemID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    []string
		wantErr bool
	}{
		{
			name:   "復習日IDを取得する場合",
			itemID: "a50e8400-e29b-41d4-a716-446655440001",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []string{
				"b50e8400-e29b-41d4-a716-446655440001",
				"b50e8400-e29b-41d4-a716-446655440002",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewDateIDs, err := repo.GetReviewDateIDsByItemID(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewDateIDs); diff != "" {
				t.Errorf("GetReviewDateIDsByItemID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_UpdateItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		item    *itemDomain.Item
		want    *itemDomain.Item
		wantErr bool
	}{
		{
			name: "復習物の更新に成功する場合",
			item: func() *itemDomain.Item {
				item, _ := itemDomain.ReconstructItem(
					"a50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					stringPtr("750e8400-e29b-41d4-a716-446655440001"),
					"更新された二次方程式",
					"更新された詳細",
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					false,
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Now(),
				)
				return item
			}(),
			want: func() *itemDomain.Item {
				item, _ := itemDomain.ReconstructItem(
					"a50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					stringPtr("750e8400-e29b-41d4-a716-446655440001"),
					"更新された二次方程式",
					"更新された詳細",
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					false,
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Time{}, // EditedAtは動的に設定
				)
				return item
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateItem(ctx, tc.item)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習物を取得して検証
				updatedItem, err := repo.GetItemByID(ctx, tc.item.ItemID(), tc.item.UserID())
				if err != nil {
					t.Errorf("更新された復習物の取得に失敗: %v", err)
					return
				}

				if updatedItem == nil {
					t.Error("更新された復習物がnilです")
					return
				}

				// 期待値を生成し動的な値を設定
				want, _ := itemDomain.ReconstructItem(
					tc.want.ItemID(),
					tc.want.UserID(),
					tc.want.CategoryID(),
					tc.want.BoxID(),
					tc.want.PatternID(),
					tc.want.Name(),
					tc.want.Detail(),
					tc.want.LearnedDate(),
					tc.want.IsFinished(),
					tc.want.RegisteredAt(),
					updatedItem.EditedAt(), // 動的に変更された値を使用
				)

				// 期待値との比較
				if diff := cmp.Diff(want, updatedItem, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
					t.Errorf("UpdateItem() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_UpdateReviewDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name        string
		reviewdates []*itemDomain.Reviewdate
		userID      string
		want        []*itemDomain.Reviewdate
		wantErr     bool
	}{
		{
			name: "復習日の更新に成功する場合",
			reviewdates: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						1,
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // 更新されたスケジュール日
						false,
					)
					return reviewdate
				}(),
			},
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						1,
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateReviewDates(ctx, tc.reviewdates, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習日を取得して検証
				updatedReviewDates, err := repo.GetReviewDatesByItemID(ctx, tc.reviewdates[0].ItemID(), tc.userID)
				if err != nil {
					t.Errorf("更新された復習日の取得に失敗: %v", err)
					return
				}

				if len(updatedReviewDates) == 0 {
					t.Error("更新された復習日が見つかりません")
					return
				}

				// 更新された復習日を探して比較
				for _, updatedReviewDate := range updatedReviewDates {
					for _, wantReviewDate := range tc.want {
						if updatedReviewDate.ReviewdateID() == wantReviewDate.ReviewdateID() {
							if diff := cmp.Diff(wantReviewDate, updatedReviewDate, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
								t.Errorf("UpdateReviewDates() mismatch (-want +got):\n%s", diff)
							}
							break
						}
					}
				}
			}
		})
	}
}

func TestItemRepository_UpdateReviewDatesBack(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name        string
		reviewdates []*itemDomain.Reviewdate
		userID      string
		want        []*itemDomain.Reviewdate
		wantErr     bool
	}{
		{
			name: "復習日の巻き戻し更新に成功する場合",
			reviewdates: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						2,
						time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), // 巻き戻されたスケジュール日
						false,
					)
					return reviewdate
				}(),
			},
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						2,
						time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateReviewDatesBack(ctx, tc.reviewdates, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習日を取得して検証
				updatedReviewDates, err := repo.GetReviewDatesByItemID(ctx, tc.reviewdates[0].ItemID(), tc.userID)
				if err != nil {
					t.Errorf("更新された復習日の取得に失敗: %v", err)
					return
				}

				if len(updatedReviewDates) == 0 {
					t.Error("更新された復習日が見つかりません")
					return
				}

				// 更新された復習日を探して比較
				for _, updatedReviewDate := range updatedReviewDates {
					for _, wantReviewDate := range tc.want {
						if updatedReviewDate.ReviewdateID() == wantReviewDate.ReviewdateID() {
							if diff := cmp.Diff(wantReviewDate, updatedReviewDate, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
								t.Errorf("UpdateReviewDatesBack() mismatch (-want +got):\n%s", diff)
							}
							break
						}
					}
				}
			}
		})
	}
}

func TestItemRepository_UpdateItemAsFinished(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	finishedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		itemID     string
		userID     string
		finishedAt time.Time
		want       *itemDomain.Item
		wantErr    bool
	}{
		{
			name:       "復習物を完了状態に更新（未完了→完了）",
			itemID:     "a50e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			finishedAt: finishedAt,
			want: func() *itemDomain.Item {
				item, _ := itemDomain.ReconstructItem(
					"a50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					stringPtr("750e8400-e29b-41d4-a716-446655440001"),
					"二次方程式",
					"ax^2 + bx + c = 0の解の公式を覚える",
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					true,
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					finishedAt,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:       "存在しない復習物を更新する場合",
			itemID:     uuid.New().String(),
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			finishedAt: finishedAt,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "他ユーザーの復習物を更新する場合",
			itemID:     "a50e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			finishedAt: finishedAt,
			want:       nil,
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateItemAsFinished(ctx, tc.itemID, tc.userID, tc.finishedAt)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 更新が成功した場合、結果を確認
			if tc.want != nil {
				actualItem, err := repo.GetItemByID(ctx, tc.itemID, tc.userID)
				if err != nil {
					t.Errorf("復習物の取得に失敗: %v", err)
					return
				}

				if diff := cmp.Diff(tc.want, actualItem, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
					t.Errorf("UpdateItemAsFinished() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_UpdateItemAsUnFinished(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name     string
		itemID   string
		userID   string
		editedAt time.Time
		want     *itemDomain.Item
		wantErr  bool
	}{
		{
			name:     "復習物を未完了状態に更新する場合",
			itemID:   "a50e8400-e29b-41d4-a716-446655440002", // フィクスチャでは完了済み
			userID:   "550e8400-e29b-41d4-a716-446655440001",
			editedAt: time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
			want: func() *itemDomain.Item {
				item, _ := itemDomain.ReconstructItem(
					"a50e8400-e29b-41d4-a716-446655440002",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440002"),
					stringPtr("750e8400-e29b-41d4-a716-446655440002"),
					"円の面積",
					"π × r^2の公式を理解する",
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					false, // 未完了に更新される
					time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
					time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
				)
				return item
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateItemAsUnFinished(ctx, tc.itemID, tc.userID, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習物を取得して検証
				updatedItem, err := repo.GetItemByID(ctx, tc.itemID, tc.userID)
				if err != nil {
					t.Errorf("更新された復習物の取得に失敗: %v", err)
					return
				}

				if updatedItem == nil {
					t.Error("更新された復習物がnilです")
					return
				}

				// 期待値との比較
				if diff := cmp.Diff(tc.want, updatedItem, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
					t.Errorf("UpdateItemAsUnFinished() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_UpdateReviewDateAsCompleted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name         string
		reviewdateID string
		userID       string
		want         *itemDomain.Reviewdate
		wantErr      bool
	}{
		{
			name:         "復習日を完了状態に更新する場合",
			reviewdateID: "b50e8400-e29b-41d4-a716-446655440001", // フィクスチャでは未完了
			userID:       "550e8400-e29b-41d4-a716-446655440001",
			want: func() *itemDomain.Reviewdate {
				reviewdate, _ := itemDomain.ReconstructReviewdate(
					"b50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					"a50e8400-e29b-41d4-a716-446655440001",
					1,
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					true, // 完了状態に更新される
				)
				return reviewdate
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateReviewDateAsCompleted(ctx, tc.reviewdateID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習日を取得して検証
				updatedReviewDates, err := repo.GetReviewDatesByItemID(ctx, tc.want.ItemID(), tc.userID)
				if err != nil {
					t.Errorf("更新された復習日の取得に失敗: %v", err)
					return
				}

				var updatedReviewDate *itemDomain.Reviewdate
				for _, rd := range updatedReviewDates {
					if rd.ReviewdateID() == tc.reviewdateID {
						updatedReviewDate = rd
						break
					}
				}

				if updatedReviewDate == nil {
					t.Error("更新された復習日が見つかりません")
					return
				}

				// 期待値との比較
				if diff := cmp.Diff(tc.want, updatedReviewDate, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
					t.Errorf("UpdateReviewDateAsCompleted() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_UpdateReviewDateAsInCompleted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name         string
		reviewdateID string
		userID       string
		want         *itemDomain.Reviewdate
		wantErr      bool
	}{
		{
			name:         "復習日を未完了状態に更新する場合",
			reviewdateID: "b50e8400-e29b-41d4-a716-446655440003", // フィクスチャでは完了済み
			userID:       "550e8400-e29b-41d4-a716-446655440001",
			want: func() *itemDomain.Reviewdate {
				reviewdate, _ := itemDomain.ReconstructReviewdate(
					"b50e8400-e29b-41d4-a716-446655440003",
					"550e8400-e29b-41d4-a716-446655440001",
					stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					stringPtr("950e8400-e29b-41d4-a716-446655440002"),
					"a50e8400-e29b-41d4-a716-446655440002",
					1,
					time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					false, // 未完了状態に更新される
				)
				return reviewdate
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.UpdateReviewDateAsInCompleted(ctx, tc.reviewdateID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tc.want != nil {
				// 更新された復習日を取得して検証
				updatedReviewDates, err := repo.GetReviewDatesByItemID(ctx, tc.want.ItemID(), tc.userID)
				if err != nil {
					t.Errorf("更新された復習日の取得に失敗: %v", err)
					return
				}

				var updatedReviewDate *itemDomain.Reviewdate
				for _, rd := range updatedReviewDates {
					if rd.ReviewdateID() == tc.reviewdateID {
						updatedReviewDate = rd
						break
					}
				}

				if updatedReviewDate == nil {
					t.Error("更新された復習日が見つかりません")
					return
				}

				// 期待値との比較
				if diff := cmp.Diff(tc.want, updatedReviewDate, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
					t.Errorf("UpdateReviewDateAsInCompleted() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestItemRepository_GetReviewDatesByItemID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    []*itemDomain.Reviewdate
		wantErr bool
	}{
		{
			name:   "復習物に関連する復習日を取得する場合",
			itemID: "a50e8400-e29b-41d4-a716-446655440001",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						1,
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						2,
						time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewdates, err := repo.GetReviewDatesByItemID(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewdates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("GetReviewDatesByItemID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_DeleteItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name             string
		itemID           string
		userID           string
		want             *itemDomain.Item // 削除後の取得結果（nilなら削除成功）
		wantErr          bool
		wantGetItemError bool // GetItemByIDでエラーが発生するか
	}{
		{
			name:             "復習物削除に成功する場合",
			itemID:           "a50e8400-e29b-41d4-a716-446655440005", // 明治維新 - user2の復習物
			userID:           "550e8400-e29b-41d4-a716-446655440002",
			want:             nil,
			wantErr:          false,
			wantGetItemError: true, // 削除後は取得できない
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			err := repo.DeleteItem(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 削除後の状態を確認
			actualItem, err := repo.GetItemByID(ctx, tc.itemID, tc.userID)

			if tc.wantGetItemError {
				if err == nil {
					t.Error("GetItemByIDでエラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("GetItemByIDで予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, actualItem, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("DeleteItem() 削除後の状態 mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_DeleteReviewDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    []*itemDomain.Reviewdate // 削除後の期待結果
		wantErr bool
	}{
		{
			name:    "復習日の削除に成功する場合",
			itemID:  "a50e8400-e29b-41d4-a716-446655440001", // 二次方程式の復習物
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.Reviewdate{}, // 削除後は空のスライス
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			// 削除を実行
			err := repo.DeleteReviewDates(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 削除後の状態を確認
			actualReviewDates, err := repo.GetReviewDatesByItemID(ctx, tc.itemID, tc.userID)
			if err != nil {
				t.Errorf("削除後の復習日取得に失敗: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, actualReviewDates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("DeleteReviewDates() 削除後の状態 mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetAllUnFinishedItemsByBoxID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		boxID   string
		userID  string
		want    []*itemDomain.Item
		wantErr bool
	}{
		{
			name:   "ボックス内の未完了復習物を取得する場合",
			boxID:  "950e8400-e29b-41d4-a716-446655440001",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Item{
				func() *itemDomain.Item {
					item, _ := itemDomain.ReconstructItem(
						"a50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						stringPtr("750e8400-e29b-41d4-a716-446655440001"),
						"二次方程式",
						"ax^2 + bx + c = 0の解の公式を覚える",
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						false,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					)
					return item
				}(),
			},
			wantErr: false,
		},
		{
			name:    "完了済み復習物のみのボックスの場合（空の結果）",
			boxID:   "950e8400-e29b-41d4-a716-446655440002",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.Item{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetAllUnFinishedItemsByBoxID(ctx, tc.boxID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, items, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("GetAllUnFinishedItemsByBoxID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetAllReviewDatesByBoxID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		boxID   string
		userID  string
		want    []*itemDomain.Reviewdate
		wantErr bool
	}{
		{
			name:   "ボックス内の全復習日を取得する場合",
			boxID:  "950e8400-e29b-41d4-a716-446655440001",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						1,
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440001"),
						"a50e8400-e29b-41d4-a716-446655440001",
						2,
						time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewdates, err := repo.GetAllReviewDatesByBoxID(ctx, tc.boxID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewdates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("GetAllReviewDatesByBoxID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetAllUnFinishedUnclassifiedItemsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		userID        string
		wantErr       bool
		expectedCount int
	}{
		{
			name:          "ユーザー2の未分類復習物を取得（完全未分類なし）",
			userID:        "550e8400-e29b-41d4-a716-446655440002",
			wantErr:       false,
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetAllUnFinishedUnclassifiedItemsByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if items == nil {
				t.Error("復習物のスライスがnilです")
				return
			}

			if len(items) != tc.expectedCount {
				t.Errorf("期待される復習物数: %d, 実際: %d", tc.expectedCount, len(items))
			}

			// 返された全ての復習物が指定ユーザーのものであり、ボックスがnilであることを検証
			for _, item := range items {
				if item.UserID() != tc.userID {
					t.Errorf("期待されるユーザーID: %s, 実際: %s", tc.userID, item.UserID())
				}
				if item.BoxID() != nil {
					t.Errorf("復習物のBoxIDはnilであるべきですが、%sが設定されています", *item.BoxID())
				}
			}
		})
	}
}

func TestItemRepository_GetAllUnclassifiedReviewDatesByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		userID  string
		want    []*itemDomain.Reviewdate
		wantErr bool
	}{
		{
			name:   "未分類復習日を取得する場合",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440006",
						"550e8400-e29b-41d4-a716-446655440002",
						nil, // 真の未分類のため nil
						nil, // 未分類のため nil
						"a50e8400-e29b-41d4-a716-446655440006",
						1,
						time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
		{
			name:    "未分類復習日がないユーザーの場合",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.Reviewdate{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewdates, err := repo.GetAllUnclassifiedReviewDatesByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewdates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("GetAllUnclassifiedReviewDatesByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetAllUnFinishedUnclassifiedItemsByCategoryID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		want       []*itemDomain.Item
		wantErr    bool
	}{
		{
			name:       "カテゴリの未分類未完了復習物を取得する場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440004",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			want: []*itemDomain.Item{
				func() *itemDomain.Item {
					item, _ := itemDomain.ReconstructItem(
						"a50e8400-e29b-41d4-a716-446655440005",
						"550e8400-e29b-41d4-a716-446655440002",
						stringPtr("650e8400-e29b-41d4-a716-446655440004"),
						nil, // 未分類のため nil
						nil,
						"明治維新",
						"1868年の政治変革について",
						time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
						false,
						time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
					)
					return item
				}(),
			},
			wantErr: false,
		},
		{
			name:       "未分類未完了復習物がないカテゴリの場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want:       []*itemDomain.Item{},
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, items, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("GetAllUnFinishedUnclassifiedItemsByCategoryID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetAllUnclassifiedReviewDatesByCategoryID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		want       []*itemDomain.Reviewdate
		wantErr    bool
	}{
		{
			name:       "カテゴリの未分類復習日を取得する場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440004",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			want: []*itemDomain.Reviewdate{
				func() *itemDomain.Reviewdate {
					reviewdate, _ := itemDomain.ReconstructReviewdate(
						"b50e8400-e29b-41d4-a716-446655440007",
						"550e8400-e29b-41d4-a716-446655440002",
						stringPtr("650e8400-e29b-41d4-a716-446655440004"),
						nil, // 未分類のため nil
						"a50e8400-e29b-41d4-a716-446655440005",
						1,
						time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						false,
					)
					return reviewdate
				}(),
			},
			wantErr: false,
		},
		{
			name:       "未分類復習日がないカテゴリの場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want:       []*itemDomain.Reviewdate{},
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewdates, err := repo.GetAllUnclassifiedReviewDatesByCategoryID(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewdates, cmp.AllowUnexported(itemDomain.Reviewdate{})); diff != "" {
				t.Errorf("GetAllUnclassifiedReviewDatesByCategoryID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_CountItemsGroupedByBoxByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		userID  string
		want    []*itemDomain.ItemCountGroupedByBox
		wantErr bool
	}{
		{
			name:   "ユーザーのボックス毎の復習物数を取得する場合",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.ItemCountGroupedByBox{
				{
					CategoryID: "650e8400-e29b-41d4-a716-446655440001",
					BoxID:      "950e8400-e29b-41d4-a716-446655440001",
					Count:      1, // 二次方程式の復習物
				},
				{
					CategoryID: "650e8400-e29b-41d4-a716-446655440002",
					BoxID:      "950e8400-e29b-41d4-a716-446655440003",
					Count:      1, // ニュートンの第一法則の復習物
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			counts, err := repo.CountItemsGroupedByBoxByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, counts, cmp.AllowUnexported(itemDomain.ItemCountGroupedByBox{})); diff != "" {
				t.Errorf("CountItemsGroupedByBoxByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_CountUnclassifiedItemsGroupedByCategoryByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		userID  string
		want    []*itemDomain.UnclassifiedItemCountGroupedByCategory
		wantErr bool
	}{
		{
			name:   "ユーザーのカテゴリ毎の未分類復習物数を取得する場合",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: []*itemDomain.UnclassifiedItemCountGroupedByCategory{
				{
					CategoryID: "650e8400-e29b-41d4-a716-446655440004",
					Count:      1, // 明治維新の復習物（未分類）
				},
				{
					CategoryID: "00000000-0000-0000-0000-000000000000", // category_id=nullの場合のUUID表現
					Count:      1,                                      // Goroutineの復習物（完全未分類）
				},
			},
			wantErr: false,
		},
		{
			name:    "未分類復習物がないユーザーの場合",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.UnclassifiedItemCountGroupedByCategory{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			counts, err := repo.CountUnclassifiedItemsGroupedByCategoryByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, counts, cmp.AllowUnexported(itemDomain.UnclassifiedItemCountGroupedByCategory{})); diff != "" {
				t.Errorf("CountUnclassifiedItemsGroupedByCategoryByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_CountUnclassifiedItemsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		userID        string
		wantErr       bool
		expectedCount int64
	}{
		{
			name:          "ユーザー2の未分類復習物数をカウント（2件）",
			userID:        "550e8400-e29b-41d4-a716-446655440002",
			wantErr:       false,
			expectedCount: 2, // 明治維新とGoroutineの2件
		},
		{
			name:          "ユーザー1の未分類復習物数をカウント（0件）",
			userID:        "550e8400-e29b-41d4-a716-446655440001",
			wantErr:       false,
			expectedCount: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			count, err := repo.CountUnclassifiedItemsByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if int64(count) != tc.expectedCount {
				t.Errorf("expected count %d, got %d", tc.expectedCount, count)
			}
		})
	}
}

func TestItemRepository_CountDailyDatesGroupedByBoxByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		userID     string
		targetDate time.Time
		want       []*itemDomain.DailyCountGroupedByBox
		wantErr    bool
	}{
		{
			name:       "特定日のボックス毎の復習日数を取得する場合",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			targetDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want: []*itemDomain.DailyCountGroupedByBox{
				{
					CategoryID: "650e8400-e29b-41d4-a716-446655440001",
					BoxID:      "950e8400-e29b-41d4-a716-446655440001",
					Count:      1, // 2024-01-02にスケジュールされた復習日
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			counts, err := repo.CountDailyDatesGroupedByBoxByUserID(ctx, tc.userID, tc.targetDate)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, counts, cmp.AllowUnexported(itemDomain.DailyCountGroupedByBox{})); diff != "" {
				t.Errorf("CountDailyDatesGroupedByBoxByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_CountDailyDatesUnclassifiedGroupedByCategoryByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		userID     string
		targetDate time.Time
		want       []*itemDomain.UnclassifiedDailyDatesCountGroupedByCategory
		wantErr    bool
	}{
		{
			name:       "特定日のカテゴリ毎の未分類復習日数を取得する場合",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			targetDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			want: []*itemDomain.UnclassifiedDailyDatesCountGroupedByCategory{
				{
					CategoryID: "650e8400-e29b-41d4-a716-446655440004",
					Count:      1, // 2024-01-06にスケジュールされた未分類復習日
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			counts, err := repo.CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx, tc.userID, tc.targetDate)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, counts, cmp.AllowUnexported(itemDomain.UnclassifiedDailyDatesCountGroupedByCategory{})); diff != "" {
				t.Errorf("CountDailyDatesUnclassifiedGroupedByCategoryByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_CountDailyDatesUnclassifiedByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		userID     string
		targetDate time.Time
		want       int
		wantErr    bool
	}{
		{
			name:       "特定日のユーザーの未分類復習日数を取得する場合",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			targetDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			want:       1, // 2024-01-06にスケジュールされた未分類復習日
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			count, err := repo.CountDailyDatesUnclassifiedByUserID(ctx, tc.userID, tc.targetDate)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if count != tc.want {
				t.Errorf("CountDailyDatesUnclassifiedByUserID() = %d, want %d", count, tc.want)
			}
		})
	}
}

func TestItemRepository_GetEditedAtByItemID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		itemID  string
		userID  string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "復習物のEditedAtを取得する場合",
			itemID:  "a50e8400-e29b-41d4-a716-446655440001",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			editedAt, err := repo.GetEditedAtByItemID(ctx, tc.itemID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if !editedAt.Equal(tc.want) {
				t.Errorf("GetEditedAtByItemID() = %v, want %v", editedAt, tc.want)
			}
		})
	}
}

func TestItemRepository_CountAllDailyReviewDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name        string
		userID      string
		parsedToday time.Time
		want        int
		wantErr     bool
	}{
		{
			name:        "特定日の全復習日数を取得する場合",
			userID:      "550e8400-e29b-41d4-a716-446655440001",
			parsedToday: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want:        1, // 2024-01-02にスケジュールされた復習日
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			count, err := repo.CountAllDailyReviewDates(ctx, tc.userID, tc.parsedToday)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if count != tc.want {
				t.Errorf("CountAllDailyReviewDates() = %d, want %d", count, tc.want)
			}
		})
	}
}

func TestItemRepository_GetAllDailyReviewDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name        string
		userID      string
		parsedToday time.Time
		want        []*itemDomain.DailyReviewDate
		wantErr     bool
	}{
		{
			name:        "特定日の全復習日詳細を取得する場合",
			userID:      "550e8400-e29b-41d4-a716-446655440001",
			parsedToday: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want: []*itemDomain.DailyReviewDate{
				{
					ReviewdateID:         "b50e8400-e29b-41d4-a716-446655440001",
					CategoryID:           stringPtr("650e8400-e29b-41d4-a716-446655440001"),
					BoxID:                stringPtr("950e8400-e29b-41d4-a716-446655440001"),
					StepNumber:           1,
					InitialScheduledDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					PrevScheduledDate:    nil,
					ScheduledDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					NextScheduledDate:    timePtr(time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)),
					IsCompleted:          false,
					ItemID:               "a50e8400-e29b-41d4-a716-446655440001",
					Name:                 "二次方程式",
					Detail:               "ax^2 + bx + c = 0の解の公式を覚える",
					LearnedDate:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					RegisteredAt:         time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					EditedAt:             time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:        "復習日がない日付の場合",
			userID:      "550e8400-e29b-41d4-a716-446655440001",
			parsedToday: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			want:        []*itemDomain.DailyReviewDate{},
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			reviewDates, err := repo.GetAllDailyReviewDates(ctx, tc.userID, tc.parsedToday)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, reviewDates, cmp.AllowUnexported(itemDomain.DailyReviewDate{})); diff != "" {
				t.Errorf("GetAllDailyReviewDates() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetFinishedItemsByBoxID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		boxID   string
		userID  string
		want    []*itemDomain.Item
		wantErr bool
	}{
		{
			name:   "ボックス内の完了済み復習物を取得する場合",
			boxID:  "950e8400-e29b-41d4-a716-446655440002",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: []*itemDomain.Item{
				func() *itemDomain.Item {
					item, _ := itemDomain.ReconstructItem(
						"a50e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						stringPtr("650e8400-e29b-41d4-a716-446655440001"),
						stringPtr("950e8400-e29b-41d4-a716-446655440002"),
						stringPtr("750e8400-e29b-41d4-a716-446655440002"),
						"円の面積",
						"π × r^2の公式を理解する",
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						true,
						time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
					)
					return item
				}(),
			},
			wantErr: false,
		},
		{
			name:    "未完了復習物のみのボックスの場合（空の結果）",
			boxID:   "950e8400-e29b-41d4-a716-446655440001",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.Item{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetFinishedItemsByBoxID(ctx, tc.boxID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, items, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("GetFinishedItemsByBoxID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetUnclassfiedFinishedItemsByCategoryID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		categoryID string
		userID     string
		want       []*itemDomain.Item
		wantErr    bool
	}{
		{
			name:       "未分類完了済み復習物がないカテゴリの場合（GetUnclassfiedFinishedItemsByCategoryIDは真の未分類のみ）",
			categoryID: "650e8400-e29b-41d4-a716-446655440003",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			want:       []*itemDomain.Item{},
			wantErr:    false,
		},
		{
			name:       "未分類完了済み復習物がないカテゴリの場合",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want:       []*itemDomain.Item{},
			wantErr:    false,
		},
		{
			name:       "存在しないカテゴリの場合",
			categoryID: uuid.New().String(),
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want:       []*itemDomain.Item{},
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetUnclassfiedFinishedItemsByCategoryID(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, items, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("GetUnclassfiedFinishedItemsByCategoryID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_GetUnclassfiedFinishedItemsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		userID  string
		want    []*itemDomain.Item
		wantErr bool
	}{
		{
			name:   "ユーザーの未分類完了済み復習物を取得する場合",
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: []*itemDomain.Item{
				func() *itemDomain.Item {
					item, _ := itemDomain.ReconstructItem(
						"a50e8400-e29b-41d4-a716-446655440006",
						"550e8400-e29b-41d4-a716-446655440002",
						nil, // 真の未分類のため nil
						nil, // 未分類のため nil
						nil,
						"江戸時代",
						"江戸時代の政治と社会を学ぶ",
						time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
						true,
						time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC),
					)
					return item
				}(),
			},
			wantErr: false,
		},
		{
			name:    "未分類完了済み復習物がないユーザーの場合",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			want:    []*itemDomain.Item{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			items, err := repo.GetUnclassfiedFinishedItemsByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if diff := cmp.Diff(tc.want, items, cmp.AllowUnexported(itemDomain.Item{})); diff != "" {
				t.Errorf("GetUnclassfiedFinishedItemsByUserID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItemRepository_IsPatternRelatedToItemByPatternID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		patternID string
		userID    string
		want      bool
		wantErr   bool
	}{
		{
			name:      "パターンが復習物に関連している場合",
			patternID: "750e8400-e29b-41d4-a716-446655440001",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "パターンが復習物に関連していない場合(異なるユーザーのパターン)",
			patternID: "750e8400-e29b-41d4-a716-446655440003",
			userID:    "550e8400-e29b-41d4-a716-446655440001",
			want:      false,
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewItemRepository()

			isRelated, err := repo.IsPatternRelatedToItemByPatternID(ctx, tc.patternID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if isRelated != tc.want {
				t.Errorf("IsPatternRelatedToItemByPatternID() = %v, want %v", isRelated, tc.want)
			}
		})
	}
}
