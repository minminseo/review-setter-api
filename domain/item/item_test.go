package item

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewItem(t *testing.T) {
	now := time.Now()
	learnedDate := now.Add(-24 * time.Hour)
	categoryID := "category1"
	boxID := "box1"
	patternID := "pattern1"

	tests := []struct {
		name         string
		itemID       string
		userID       string
		categoryID   *string
		boxID        *string
		patternID    *string
		itemName     string
		detail       string
		learnedDate  time.Time
		isFinished   bool
		registeredAt time.Time
		editedAt     time.Time
		want         *Item
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "全項目が有効な復習物（正常系）",
			itemID:       "item1",
			userID:       "user1",
			categoryID:   &categoryID,
			boxID:        &boxID,
			patternID:    &patternID,
			itemName:     "Apple",
			detail:       "Apple - りんご",
			learnedDate:  learnedDate,
			isFinished:   false,
			registeredAt: now,
			editedAt:     now,
			want: func() *Item {
				item, _ := ReconstructItem(
					"item1",
					"user1",
					&categoryID,
					&boxID,
					&patternID,
					"Apple",
					"Apple - りんご",
					learnedDate,
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:         "nil項目を含む有効な復習物（正常系）",
			itemID:       "item2",
			userID:       "user1",
			categoryID:   nil,
			boxID:        nil,
			patternID:    nil,
			itemName:     "Apple",
			detail:       "Apple - りんご",
			learnedDate:  learnedDate,
			isFinished:   true,
			registeredAt: now,
			editedAt:     now,
			want: func() *Item {
				item, _ := ReconstructItem(
					"item2",
					"user1",
					nil,
					nil,
					nil,
					"Apple",
					"Apple - りんご",
					learnedDate,
					true,
					now,
					now,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:         "復習物名が空文字（異常系）",
			itemID:       "item3",
			userID:       "user1",
			categoryID:   &categoryID,
			boxID:        &boxID,
			patternID:    &patternID,
			itemName:     "",
			detail:       "復習物名が空の詳細",
			learnedDate:  learnedDate,
			isFinished:   false,
			registeredAt: now,
			editedAt:     now,
			want:         nil,
			wantErr:      true,
			errMsg:       "復習物名は必須です",
		},
		{
			name:         "学習日がゼロ値",
			itemID:       "item4",
			userID:       "user1",
			categoryID:   &categoryID,
			boxID:        &boxID,
			patternID:    &patternID,
			itemName:     "Test Item",
			detail:       "Test detail",
			learnedDate:  time.Time{},
			isFinished:   false,
			registeredAt: now,
			editedAt:     now,
			want:         nil,
			wantErr:      true,
			errMsg:       "学習日は必須です",
		},
		{
			name:         "詳細が空文字でも有効（正常系）",
			itemID:       "item4",
			userID:       "user1",
			categoryID:   &categoryID,
			boxID:        &boxID,
			patternID:    &patternID,
			itemName:     "Apple",
			detail:       "",
			learnedDate:  learnedDate,
			isFinished:   false,
			registeredAt: now,
			editedAt:     now,
			want: func() *Item {
				item, _ := ReconstructItem(
					"item4",
					"user1",
					&categoryID,
					&boxID,
					&patternID,
					"Apple",
					"",
					learnedDate,
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item, err := NewItem(
				tc.itemID,
				tc.userID,
				tc.categoryID,
				tc.boxID,
				tc.patternID,
				tc.itemName,
				tc.detail,
				tc.learnedDate,
				tc.isFinished,
				tc.registeredAt,
				tc.editedAt,
			)

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

			if diff := cmp.Diff(tc.want, item, cmp.AllowUnexported(Item{})); diff != "" {
				t.Errorf("Item mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestItem_UpdateItem(t *testing.T) {
	now := time.Now()
	learnedDate := now.Add(-24 * time.Hour)
	categoryID := "category1"
	boxID := "box1"
	patternID := "pattern1"

	item, err := NewItem("item1", "user1", &categoryID, &boxID, &patternID, "Original Item", "Original detail", learnedDate, false, now, now)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	newTime := now.Add(time.Hour)
	newLearnedDate := now.Add(-12 * time.Hour)
	newCategoryID := "category2"
	newBoxID := "box2"
	newPatternID := "pattern2"

	tests := []struct {
		name           string
		categoryID     *string
		boxID          *string
		patternID      *string
		newName        string
		newDetail      string
		newLearnedDate time.Time
		editedAt       time.Time
		wantItem       *Item
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "全項目を更新（正常系）",
			categoryID:     &newCategoryID,
			boxID:          &newBoxID,
			patternID:      &newPatternID,
			newName:        "更新後復習物",
			newDetail:      "更新後詳細",
			newLearnedDate: newLearnedDate,
			editedAt:       newTime,
			wantItem: func() *Item {
				item, _ := ReconstructItem(
					"item1",
					"user1",
					&newCategoryID,
					&newBoxID,
					&newPatternID,
					"更新後復習物",
					"更新後詳細",
					newLearnedDate,
					false,
					now,
					newTime,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:           "nil項目を含めて更新（正常系）",
			categoryID:     nil,
			boxID:          nil,
			patternID:      nil,
			newName:        "nil項目で更新",
			newDetail:      "更新後詳細",
			newLearnedDate: newLearnedDate,
			editedAt:       newTime,
			wantItem: func() *Item {
				item, _ := ReconstructItem(
					"item1",
					"user1",
					nil,
					nil,
					nil,
					"nil項目で更新",
					"更新後詳細",
					newLearnedDate,
					false,
					now,
					newTime,
				)
				return item
			}(),
			wantErr: false,
		},
		{
			name:           "復習物名が空文字（異常系）",
			categoryID:     &categoryID,
			boxID:          &boxID,
			patternID:      &patternID,
			newName:        "",
			newDetail:      "詳細",
			newLearnedDate: newLearnedDate,
			editedAt:       newTime,
			wantItem: func() *Item {
				item, _ := ReconstructItem(
					"item1",
					"user1",
					&categoryID,
					&boxID,
					&patternID,
					"Original Item",
					"Original detail",
					learnedDate,
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: true,
			errMsg:  "復習物名は必須です",
		},
		{
			name:           "学習日がゼロ値（異常系）",
			categoryID:     &categoryID,
			boxID:          &boxID,
			patternID:      &patternID,
			newName:        "有効な名前",
			newDetail:      "有効な詳細",
			newLearnedDate: time.Time{},
			editedAt:       newTime,
			wantItem: func() *Item {
				item, _ := ReconstructItem(
					"item1",
					"user1",
					&categoryID,
					&boxID,
					&patternID,
					"Original Item",
					"Original detail",
					learnedDate,
					false,
					now,
					now,
				)
				return item
			}(),
			wantErr: true,
			errMsg:  "学習日は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// 復習物をコピー
			testItem := *item

			err := testItem.UpdateItem(tc.categoryID, tc.boxID, tc.patternID, tc.newName, tc.newDetail, tc.newLearnedDate, tc.editedAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				if diff := cmp.Diff(tc.wantItem, &testItem, cmp.AllowUnexported(Item{})); diff != "" {
					t.Errorf("Item mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if diff := cmp.Diff(tc.wantItem, &testItem, cmp.AllowUnexported(Item{})); diff != "" {
				t.Errorf("Item mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewReviewdate(t *testing.T) {
	now := time.Now()
	initialDate := now.Add(-24 * time.Hour)
	scheduledDate := now.Add(24 * time.Hour)
	categoryID := "category1"
	boxID := "box1"

	tests := []struct {
		name                 string
		reviewdateID         string
		userID               string
		categoryID           *string
		boxID                *string
		itemID               string
		stepNumber           int
		initialScheduledDate time.Time
		scheduledDate        time.Time
		isCompleted          bool
		want                 *Reviewdate
		wantErr              bool
		errMsg               string
	}{
		{
			name:                 "全項目が有効な復習日（正常系）",
			reviewdateID:         "review1",
			userID:               "user1",
			categoryID:           &categoryID,
			boxID:                &boxID,
			itemID:               "item1",
			stepNumber:           1,
			initialScheduledDate: initialDate,
			scheduledDate:        scheduledDate,
			isCompleted:          false,
			want: func() *Reviewdate {
				reviewdate, _ := ReconstructReviewdate(
					"review1",
					"user1",
					&categoryID,
					&boxID,
					"item1",
					1,
					initialDate,
					scheduledDate,
					false,
				)
				return reviewdate
			}(),
			wantErr: false,
		},
		{
			name:                 "nil項目を含む復習日（正常系）",
			reviewdateID:         "review2",
			userID:               "user1",
			categoryID:           nil,
			boxID:                nil,
			itemID:               "item1",
			stepNumber:           1,
			initialScheduledDate: initialDate,
			scheduledDate:        scheduledDate,
			isCompleted:          true,
			want: func() *Reviewdate {
				reviewdate, _ := ReconstructReviewdate(
					"review2",
					"user1",
					nil,
					nil,
					"item1",
					1,
					initialDate,
					scheduledDate,
					true,
				)
				return reviewdate
			}(),
			wantErr: false,
		},
		{
			name:                 "ステップ番号が0（異常系）",
			reviewdateID:         "review3",
			userID:               "user1",
			categoryID:           &categoryID,
			boxID:                &boxID,
			itemID:               "item1",
			stepNumber:           0,
			initialScheduledDate: initialDate,
			scheduledDate:        scheduledDate,
			isCompleted:          false,
			want:                 nil,
			wantErr:              true,
			errMsg:               "stepNumber: ステップ番号は必須です.",
		},
		{
			name:                 "ステップ番号が負数（異常系）",
			reviewdateID:         "review4",
			userID:               "user1",
			categoryID:           &categoryID,
			boxID:                &boxID,
			itemID:               "item1",
			stepNumber:           -1,
			initialScheduledDate: initialDate,
			scheduledDate:        scheduledDate,
			isCompleted:          false,
			want:                 nil,
			wantErr:              true,
			errMsg:               "stepNumber: ステップ番号の値が不正です.",
		},
		{
			name:                 "ステップ番号が32768以上（異常系）",
			reviewdateID:         "review5",
			userID:               "user1",
			categoryID:           &categoryID,
			boxID:                &boxID,
			itemID:               "item1",
			stepNumber:           32768,
			initialScheduledDate: initialDate,
			scheduledDate:        scheduledDate,
			isCompleted:          false,
			want:                 nil,
			wantErr:              true,
			errMsg:               "stepNumber: ステップは32768回以上は指定できません.",
		},
		{
			name:                 "スケジュール日がゼロ値（異常系）",
			reviewdateID:         "review6",
			userID:               "user1",
			categoryID:           &categoryID,
			boxID:                &boxID,
			itemID:               "item1",
			stepNumber:           1,
			initialScheduledDate: initialDate,
			scheduledDate:        time.Time{},
			isCompleted:          false,
			want:                 nil,
			wantErr:              true,
			errMsg:               "scheduledDate: スケジュール日は必須です.",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reviewdate, err := NewReviewdate(
				tc.reviewdateID,
				tc.userID,
				tc.categoryID,
				tc.boxID,
				tc.itemID,
				tc.stepNumber,
				tc.initialScheduledDate,
				tc.scheduledDate,
				tc.isCompleted,
			)

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

			if diff := cmp.Diff(tc.want, reviewdate, cmp.AllowUnexported(Reviewdate{})); diff != "" {
				t.Errorf("Reviewdate mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
