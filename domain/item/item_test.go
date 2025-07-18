package item

import (
	"testing"
	"time"
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
			wantErr:      false,
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
			wantErr:      false,
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
			wantErr:      false,
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

			if item.ItemID != tc.itemID {
				t.Errorf("ItemID: got %q, want %q", item.ItemID, tc.itemID)
			}
			if item.UserID != tc.userID {
				t.Errorf("UserID: got %q, want %q", item.UserID, tc.userID)
			}
			if tc.categoryID == nil {
				if item.CategoryID != nil {
					t.Errorf("CategoryID: got %v, want nil", item.CategoryID)
				}
			} else {
				if item.CategoryID == nil || *item.CategoryID != *tc.categoryID {
					t.Errorf("CategoryID: got %v, want %v", item.CategoryID, tc.categoryID)
				}
			}
			if tc.boxID == nil {
				if item.BoxID != nil {
					t.Errorf("BoxID: got %v, want nil", item.BoxID)
				}
			} else {
				if item.BoxID == nil || *item.BoxID != *tc.boxID {
					t.Errorf("BoxID: got %v, want %v", item.BoxID, tc.boxID)
				}
			}
			if tc.patternID == nil {
				if item.PatternID != nil {
					t.Errorf("PatternID: got %v, want nil", item.PatternID)
				}
			} else {
				if item.PatternID == nil || *item.PatternID != *tc.patternID {
					t.Errorf("PatternID: got %v, want %v", item.PatternID, tc.patternID)
				}
			}
			if item.Name != tc.itemName {
				t.Errorf("Name: got %q, want %q", item.Name, tc.itemName)
			}
			if item.Detail != tc.detail {
				t.Errorf("Detail: got %q, want %q", item.Detail, tc.detail)
			}
			if !item.LearnedDate.Equal(tc.learnedDate) {
				t.Errorf("LearnedDate: got %v, want %v", item.LearnedDate, tc.learnedDate)
			}
			if item.IsFinished != tc.isFinished {
				t.Errorf("IsFinished: got %v, want %v", item.IsFinished, tc.isFinished)
			}
			if !item.RegisteredAt.Equal(tc.registeredAt) {
				t.Errorf("RegisteredAt: got %v, want %v", item.RegisteredAt, tc.registeredAt)
			}
			if !item.EditedAt.Equal(tc.editedAt) {
				t.Errorf("EditedAt: got %v, want %v", item.EditedAt, tc.editedAt)
			}
		})
	}
}

func TestItem_Set(t *testing.T) {
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
			wantErr:        false,
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
			wantErr:        false,
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
			wantErr:        true,
			errMsg:         "復習物名は必須です",
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
			wantErr:        true,
			errMsg:         "学習日は必須です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// 復習物をコピー
			testItem := *item

			err := testItem.Set(tc.categoryID, tc.boxID, tc.patternID, tc.newName, tc.newDetail, tc.newLearnedDate, tc.editedAt)

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

			if tc.categoryID == nil {
				if testItem.CategoryID != nil {
					t.Errorf("CategoryIDがnilであることを期待しましたが、got %v", testItem.CategoryID)
				}
			} else {
				if testItem.CategoryID == nil || *testItem.CategoryID != *tc.categoryID {
					t.Errorf("CategoryIDが一致しません: got %v, want %v", testItem.CategoryID, tc.categoryID)
				}
			}
			if tc.boxID == nil {
				if testItem.BoxID != nil {
					t.Errorf("BoxIDがnilであることを期待しましたが、got %v", testItem.BoxID)
				}
			} else {
				if testItem.BoxID == nil || *testItem.BoxID != *tc.boxID {
					t.Errorf("BoxIDが一致しません: got %v, want %v", testItem.BoxID, tc.boxID)
				}
			}
			if tc.patternID == nil {
				if testItem.PatternID != nil {
					t.Errorf("PatternIDがnilであることを期待しましたが、got %v", testItem.PatternID)
				}
			} else {
				if testItem.PatternID == nil || *testItem.PatternID != *tc.patternID {
					t.Errorf("PatternIDが一致しません: got %v, want %v", testItem.PatternID, tc.patternID)
				}
			}
			if testItem.Name != tc.newName {
				t.Errorf("Nameが一致しません: got %q, want %q", testItem.Name, tc.newName)
			}
			if testItem.Detail != tc.newDetail {
				t.Errorf("Detailが一致しません: got %q, want %q", testItem.Detail, tc.newDetail)
			}
			if !testItem.LearnedDate.Equal(tc.newLearnedDate) {
				t.Errorf("LearnedDateが一致しません: got %v, want %v", testItem.LearnedDate, tc.newLearnedDate)
			}
			if !testItem.EditedAt.Equal(tc.editedAt) {
				t.Errorf("EditedAtが一致しません: got %v, want %v", testItem.EditedAt, tc.editedAt)
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
			wantErr:              false,
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
			wantErr:              false,
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
			wantErr:              true,
			errMsg:               "StepNumber: ステップ番号は必須です.",
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
			wantErr:              true,
			errMsg:               "StepNumber: ステップ番号の値が不正です.",
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
			wantErr:              true,
			errMsg:               "StepNumber: ステップは32768回以上は指定できません.",
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
			wantErr:              true,
			errMsg:               "ScheduledDate: スケジュール日は必須です.",
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

			if reviewdate.ReviewdateID != tc.reviewdateID {
				t.Errorf("ReviewdateIDが一致しません: got %q, want %q", reviewdate.ReviewdateID, tc.reviewdateID)
			}
			if reviewdate.UserID != tc.userID {
				t.Errorf("UserIDが一致しません: got %q, want %q", reviewdate.UserID, tc.userID)
			}
			if tc.categoryID == nil {
				if reviewdate.CategoryID != nil {
					t.Errorf("CategoryIDがnilであることを期待しましたが、got %v", reviewdate.CategoryID)
				}
			} else {
				if reviewdate.CategoryID == nil || *reviewdate.CategoryID != *tc.categoryID {
					t.Errorf("CategoryIDが一致しません: got %v, want %v", reviewdate.CategoryID, tc.categoryID)
				}
			}
			if tc.boxID == nil {
				if reviewdate.BoxID != nil {
					t.Errorf("BoxIDがnilであることを期待しましたが、got %v", reviewdate.BoxID)
				}
			} else {
				if reviewdate.BoxID == nil || *reviewdate.BoxID != *tc.boxID {
					t.Errorf("BoxIDが一致しません: got %v, want %v", reviewdate.BoxID, tc.boxID)
				}
			}
			if reviewdate.ItemID != tc.itemID {
				t.Errorf("ItemIDが一致しません: got %q, want %q", reviewdate.ItemID, tc.itemID)
			}
			if reviewdate.StepNumber != tc.stepNumber {
				t.Errorf("StepNumberが一致しません: got %d, want %d", reviewdate.StepNumber, tc.stepNumber)
			}
			if !reviewdate.InitialScheduledDate.Equal(tc.initialScheduledDate) {
				t.Errorf("InitialScheduledDateが一致しません: got %v, want %v", reviewdate.InitialScheduledDate, tc.initialScheduledDate)
			}
			if !reviewdate.ScheduledDate.Equal(tc.scheduledDate) {
				t.Errorf("ScheduledDateが一致しません: got %v, want %v", reviewdate.ScheduledDate, tc.scheduledDate)
			}
			if reviewdate.IsCompleted != tc.isCompleted {
				t.Errorf("IsCompletedが一致しません: got %v, want %v", reviewdate.IsCompleted, tc.isCompleted)
			}
		})
	}
}
