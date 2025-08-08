package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	boxDomain "github.com/minminseo/recall-setter/domain/box"
)

func TestBoxRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		box     *boxDomain.Box
		want    *boxDomain.Box
		wantErr bool
	}{
		{
			name: "ボックス作成に成功する場合",
			box: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440001",
					"新しい復習ボックス",
					time.Now(),
					time.Now(),
				)
				return box
			}(),
			want: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"",
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440001",
					"新しい復習ボックス",
					time.Time{},
					time.Time{},
				)
				return box
			}(),
			wantErr: false,
		},
		{
			name: "存在しないカテゴリによる外部キー制約違反",
			box: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					uuid.New().String(),
					"750e8400-e29b-41d4-a716-446655440001",
					"存在しないカテゴリのボックス",
					time.Now(),
					time.Now(),
				)
				return box
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "存在しないパターンによる外部キー制約違反",
			box: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					uuid.New().String(), // Does not exist in fixture
					"存在しないパターンのボックス",
					time.Now(),
					time.Now(),
				)
				return box
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			err := repo.Create(ctx, tc.box)

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
				// 作成されたボックスを取得して検証
				createdBox, err := repo.GetByID(ctx, tc.box.ID(), tc.box.CategoryID(), tc.box.UserID())
				if err != nil {
					t.Errorf("作成されたボックスの取得に失敗: %v", err)
					return
				}

				if createdBox == nil {
					t.Error("作成されたボックスがnilです")
					return
				}

				// 期待値を生成し動的な値を設定
				wantBox, _ := boxDomain.ReconstructBox(
					createdBox.ID(),
					tc.want.UserID(),
					tc.want.CategoryID(),
					tc.want.PatternID(),
					tc.want.Name(),
					createdBox.RegisteredAt(),
					createdBox.EditedAt(),
				)

				// 期待値との比較
				if diff := cmp.Diff(wantBox, createdBox, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBoxRepository_GetAllByCategoryID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name          string
		categoryID    string
		userID        string
		want          []boxDomain.Box
		wantErr       bool
		expectedCount int
	}{
		{
			name:       "カテゴリ1・ユーザー1のボックスを取得（3件）",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want: []boxDomain.Box{
				func() boxDomain.Box {
					box, _ := boxDomain.ReconstructBox(
						"950e8400-e29b-41d4-a716-446655440001",
						"550e8400-e29b-41d4-a716-446655440001",
						"650e8400-e29b-41d4-a716-446655440001",
						"750e8400-e29b-41d4-a716-446655440001",
						"代数学",
						time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					)
					return *box
				}(),
				func() boxDomain.Box {
					box, _ := boxDomain.ReconstructBox(
						"950e8400-e29b-41d4-a716-446655440002",
						"550e8400-e29b-41d4-a716-446655440001",
						"650e8400-e29b-41d4-a716-446655440001",
						"750e8400-e29b-41d4-a716-446655440002",
						"幾何学",
						time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
					)
					return *box
				}(),
				func() boxDomain.Box {
					box, _ := boxDomain.ReconstructBox(
						"950e8400-e29b-41d4-a716-446655440005",
						"550e8400-e29b-41d4-a716-446655440001",
						"650e8400-e29b-41d4-a716-446655440001",
						"750e8400-e29b-41d4-a716-446655440001",
						"復習物のないボックス",
						time.Date(2024, 1, 1, 11, 00, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 11, 00, 0, 0, time.UTC),
					)
					return *box
				}(),
			},
			wantErr:       false,
			expectedCount: 3,
		},
		{
			name:       "カテゴリ2・ユーザー1のボックスを取得（1件）",
			categoryID: "650e8400-e29b-41d4-a716-446655440002",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want: []boxDomain.Box{
				func() boxDomain.Box {
					box, _ := boxDomain.ReconstructBox(
						"950e8400-e29b-41d4-a716-446655440003",
						"550e8400-e29b-41d4-a716-446655440001",
						"650e8400-e29b-41d4-a716-446655440002",
						"750e8400-e29b-41d4-a716-446655440001",
						"力学",
						time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
					)
					return *box
				}(),
			},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name:       "カテゴリ3・ユーザー2のボックスを取得（1件）",
			categoryID: "650e8400-e29b-41d4-a716-446655440003",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			want: []boxDomain.Box{
				func() boxDomain.Box {
					box, _ := boxDomain.ReconstructBox(
						"950e8400-e29b-41d4-a716-446655440004",
						"550e8400-e29b-41d4-a716-446655440002",
						"650e8400-e29b-41d4-a716-446655440003",
						"750e8400-e29b-41d4-a716-446655440003",
						"リーディング",
						time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
					)
					return *box
				}(),
			},
			wantErr:       false,
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			boxes, err := repo.GetAllByCategoryID(ctx, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if boxes != nil {
					t.Error("ボックスがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if boxes == nil {
				t.Error("ボックスのスライスがnilです")
				return
			}

			if len(boxes) != tc.expectedCount {
				t.Errorf("期待されるボックス数: %d, 実際: %d", tc.expectedCount, len(boxes))
			}

			if tc.want != nil {
				// ボックスのスライスを作成してポインタを外す
				boxesSlice := make([]boxDomain.Box, len(boxes))
				for i, box := range boxes {
					boxesSlice[i] = *box
				}

				// 期待値との比較
				if diff := cmp.Diff(tc.want, boxesSlice, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("GetAllByCategoryID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBoxRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		box     *boxDomain.Box
		want    *boxDomain.Box
		wantErr bool
	}{
		{
			name: "ボックス更新に成功する場合",
			box: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440001", // Exists in fixture
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440002", // Change pattern
					"更新された代数学ボックス",
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), // Keep original registered_at
					time.Now(),
				)
				return box
			}(),
			want: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440001", // Pattern should remain unchanged in Update (not UpdateWithPatternID)
					"更新された代数学ボックス",
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					time.Time{},
				)
				return box
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			err := repo.Update(ctx, tc.box)

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
				// 更新されたボックスを取得して検証
				updatedBox, err := repo.GetByID(ctx, tc.box.ID(), tc.box.CategoryID(), tc.box.UserID())
				if err != nil {
					t.Errorf("更新されたボックスの取得に失敗: %v", err)
					return
				}

				if updatedBox == nil {
					t.Error("更新されたボックスがnilです")
					return
				}

				// 期待値を生成し動的な値を設定
				wantBox, _ := boxDomain.ReconstructBox(
					tc.want.ID(),
					tc.want.UserID(),
					tc.want.CategoryID(),
					tc.want.PatternID(),
					tc.want.Name(),
					tc.want.RegisteredAt(),
					updatedBox.EditedAt(),
				)

				// 期待値との比較
				if diff := cmp.Diff(wantBox, updatedBox, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBoxRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		boxID      string
		categoryID string
		userID     string
		wantBefore *boxDomain.Box
		wantErr    bool
	}{
		{
			name:       "ボックス削除に成功する場合",
			boxID:      "950e8400-e29b-41d4-a716-446655440004",
			categoryID: "650e8400-e29b-41d4-a716-446655440003",
			userID:     "550e8400-e29b-41d4-a716-446655440002",
			wantBefore: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440004",
					"550e8400-e29b-41d4-a716-446655440002",
					"650e8400-e29b-41d4-a716-446655440003",
					"750e8400-e29b-41d4-a716-446655440003",
					"リーディング",
					time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
				)
				return box
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			// 削除前にボックスが存在することを確認
			if tc.wantBefore != nil {
				boxBefore, err := repo.GetByID(ctx, tc.boxID, tc.categoryID, tc.userID)
				if err != nil {
					t.Errorf("削除前のボックス取得に失敗: %v", err)
					return
				}
				if boxBefore == nil {
					t.Error("削除前のボックスがnilです")
					return
				}

				// 削除前のボックスが期待値と一致することを確認
				if diff := cmp.Diff(tc.wantBefore, boxBefore, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("削除前のボックスが期待値と一致しません (-want +got):\n%s", diff)
				}
			}

			// 削除を実行
			err := repo.Delete(ctx, tc.boxID, tc.categoryID, tc.userID)

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

			// 削除後にボックスが存在しないことを確認
			boxAfter, err := repo.GetByID(ctx, tc.boxID, tc.categoryID, tc.userID)
			if err == nil {
				t.Error("削除後にボックスが見つかりました。削除されていません")
				return
			}
			if boxAfter != nil {
				t.Error("削除後にボックスがnilではありません")
			}
		})
	}
}

func TestBoxRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		boxID      string
		categoryID string
		userID     string
		want       *boxDomain.Box
		wantErr    bool
	}{
		{
			name:       "有効なIDでボックスを取得する場合",
			boxID:      "950e8400-e29b-41d4-a716-446655440001",
			categoryID: "650e8400-e29b-41d4-a716-446655440001",
			userID:     "550e8400-e29b-41d4-a716-446655440001",
			want: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440001",
					"代数学",
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				)
				return box
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			box, err := repo.GetByID(ctx, tc.boxID, tc.categoryID, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if box != nil {
					t.Error("expected nil box but got one")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if box == nil {
				t.Error("expected box but got nil")
				return
			}

			if tc.want != nil {
				if diff := cmp.Diff(tc.want, box, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("GetByID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBoxRepository_UpdateWithPatternID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name            string
		box             *boxDomain.Box
		want            *boxDomain.Box
		wantErr         bool
		expectedUpdated int64
	}{
		{
			name: "復習物がないボックスのパターンID更新に成功する場合",
			box: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440005",
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440002", // 変更するパターンID
					"代数学",
					time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
					time.Now(),
				)
				return box
			}(),
			want: func() *boxDomain.Box {
				box, _ := boxDomain.ReconstructBox(
					"950e8400-e29b-41d4-a716-446655440005",
					"550e8400-e29b-41d4-a716-446655440001",
					"650e8400-e29b-41d4-a716-446655440001",
					"750e8400-e29b-41d4-a716-446655440002",
					"代数学",
					time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
					time.Time{},
				)
				return box
			}(),
			wantErr:         false,
			expectedUpdated: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			updatedCount, err := repo.UpdateWithPatternID(ctx, tc.box)

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

			if updatedCount != tc.expectedUpdated {
				t.Errorf("期待される更新数: %d, 実際: %d", tc.expectedUpdated, updatedCount)
			}

			if tc.want != nil && updatedCount > 0 {
				// 更新されたボックスを取得して検証
				updatedBox, err := repo.GetByID(ctx, tc.box.ID(), tc.box.CategoryID(), tc.box.UserID())
				if err != nil {
					t.Errorf("更新されたボックスの取得に失敗: %v", err)
					return
				}

				if updatedBox == nil {
					t.Error("更新されたボックスがnilです")
					return
				}

				// 期待値を生成し動的な値を設定
				wantBox, _ := boxDomain.ReconstructBox(
					tc.want.ID(),
					tc.want.UserID(),
					tc.want.CategoryID(),
					tc.want.PatternID(),
					tc.want.Name(),
					tc.want.RegisteredAt(),
					updatedBox.EditedAt(),
				)

				// 期待値との比較
				if diff := cmp.Diff(wantBox, updatedBox, cmp.AllowUnexported(boxDomain.Box{})); diff != "" {
					t.Errorf("UpdateWithPatternID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBoxRepository_GetBoxNamesByBoxIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		boxIDs  []string
		want    []boxDomain.BoxName
		wantErr bool
	}{
		{
			name: "複数のボックスIDからボックス名を取得する場合",
			boxIDs: []string{
				"950e8400-e29b-41d4-a716-446655440001",
				"950e8400-e29b-41d4-a716-446655440002",
				"950e8400-e29b-41d4-a716-446655440003",
			},
			want: []boxDomain.BoxName{
				{
					BoxID:     "950e8400-e29b-41d4-a716-446655440001",
					Name:      "代数学",
					PatternID: "750e8400-e29b-41d4-a716-446655440001",
				},
				{
					BoxID:     "950e8400-e29b-41d4-a716-446655440002",
					Name:      "幾何学",
					PatternID: "750e8400-e29b-41d4-a716-446655440002",
				},
				{
					BoxID:     "950e8400-e29b-41d4-a716-446655440003",
					Name:      "力学",
					PatternID: "750e8400-e29b-41d4-a716-446655440001",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewBoxRepository()

			boxNames, err := repo.GetBoxNamesByBoxIDs(ctx, tc.boxIDs)

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

			// BoxNameのスライスを作成してポインタを外す
			boxNamesSlice := make([]boxDomain.BoxName, len(boxNames))
			for i, boxName := range boxNames {
				boxNamesSlice[i] = *boxName
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, boxNamesSlice, cmp.AllowUnexported(boxDomain.BoxName{})); diff != "" {
				t.Errorf("GetBoxNamesByBoxIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
