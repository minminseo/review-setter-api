package item

import (
	"reflect"
	"testing"
	"time"

	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

func TestFormatWithOverdueMarkedCompleted(t *testing.T) {
	scheduler := NewScheduler()
	tests := []struct {
		name               string
		targetPatternSteps []*PatternDomain.PatternStep
		userID             string
		categoryID         *string
		boxID              *string
		itemID             string
		parsedLearnedDate  time.Time
		parsedToday        time.Time
		wantIsFinished     bool
		wantReviewdatesLen int
		wantError          bool
	}{
		{
			name: "正常なパターンステップで復習日が今日より前（完了扱い）",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     true,
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "正常なパターンステップで復習日が今日以降（未完了扱い）",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 5},
				{StepNumber: 2, IntervalDays: 10},
			},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     false,
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "単一ステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
			},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     true,
			wantReviewdatesLen: 1,
			wantError:          false,
		},
		{
			name:               "空のパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     false,
			wantReviewdatesLen: 0,
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReviewdates, gotIsFinished, err := scheduler.FormatWithOverdueMarkedCompleted(
				tt.targetPatternSteps,
				tt.userID,
				tt.categoryID,
				tt.boxID,
				tt.itemID,
				tt.parsedLearnedDate,
				tt.parsedToday,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("FormatWithOverdueMarkedCompleted() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err != nil {
				return
			}

			if gotIsFinished != tt.wantIsFinished {
				t.Errorf("FormatWithOverdueMarkedCompleted() gotIsFinished = %v, want %v", gotIsFinished, tt.wantIsFinished)
			}

			if len(gotReviewdates) != tt.wantReviewdatesLen {
				t.Errorf("FormatWithOverdueMarkedCompleted() reviewdates length = %v, want %v", len(gotReviewdates), tt.wantReviewdatesLen)
			}

			for i, rd := range gotReviewdates {
				if rd.UserID != tt.userID {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].UserID = %v, want %v", i, rd.UserID, tt.userID)
				}
				if rd.ItemID != tt.itemID {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].ItemID = %v, want %v", i, rd.ItemID, tt.itemID)
				}
				if !reflect.DeepEqual(rd.CategoryID, tt.categoryID) {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].CategoryID = %v, want %v", i, rd.CategoryID, tt.categoryID)
				}
				if !reflect.DeepEqual(rd.BoxID, tt.boxID) {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].BoxID = %v, want %v", i, rd.BoxID, tt.boxID)
				}

				expectedScheduledDate := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[i].IntervalDays)
				if !rd.ScheduledDate.Equal(expectedScheduledDate) {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].ScheduledDate = %v, want %v", i, rd.ScheduledDate, expectedScheduledDate)
				}

				expectedIsCompleted := expectedScheduledDate.Before(tt.parsedToday)
				if rd.IsCompleted != expectedIsCompleted {
					t.Errorf("FormatWithOverdueMarkedCompleted() reviewdate[%d].IsCompleted = %v, want %v", i, rd.IsCompleted, expectedIsCompleted)
				}
			}
		})
	}
}

func TestFormatWithOverdueMarkedInCompleted(t *testing.T) {
	scheduler := NewScheduler()

	tests := []struct {
		name               string
		targetPatternSteps []*PatternDomain.PatternStep
		userID             string
		categoryID         *string
		boxID              *string
		itemID             string
		parsedLearnedDate  time.Time
		parsedToday        time.Time
		wantReviewdatesLen int
		wantError          bool
	}{
		{
			name: "最初のステップが今日より前の場合",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "最初のステップが今日以降の場合",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 10},
				{StepNumber: 2, IntervalDays: 15},
			},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "単一ステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
			},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 1,
			wantError:          false,
		},
		{
			name:               "空のパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 0,
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReviewdates, err := scheduler.FormatWithOverdueMarkedInCompleted(
				tt.targetPatternSteps,
				tt.userID,
				tt.categoryID,
				tt.boxID,
				tt.itemID,
				tt.parsedLearnedDate,
				tt.parsedToday,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("FormatWithOverdueMarkedInCompleted() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err != nil {
				return
			}

			if len(gotReviewdates) != tt.wantReviewdatesLen {
				t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdates length = %v, want %v", len(gotReviewdates), tt.wantReviewdatesLen)
			}

			var expectedAddDuration int
			if len(tt.targetPatternSteps) > 0 {
				firstScheduled := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[0].IntervalDays)
				if firstScheduled.Before(tt.parsedToday) {
					expectedAddDuration = int(tt.parsedToday.Sub(firstScheduled).Hours() / 24)
				} else {
					expectedAddDuration = 0
				}
			}

			for i, rd := range gotReviewdates {
				if rd.UserID != tt.userID {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].UserID = %v, want %v", i, rd.UserID, tt.userID)
				}
				if rd.ItemID != tt.itemID {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].ItemID = %v, want %v", i, rd.ItemID, tt.itemID)
				}
				if !reflect.DeepEqual(rd.CategoryID, tt.categoryID) {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].CategoryID = %v, want %v", i, rd.CategoryID, tt.categoryID)
				}
				if !reflect.DeepEqual(rd.BoxID, tt.boxID) {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].BoxID = %v, want %v", i, rd.BoxID, tt.boxID)
				}

				expectedScheduledDate := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[i].IntervalDays+expectedAddDuration)
				if !rd.ScheduledDate.Equal(expectedScheduledDate) {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].ScheduledDate = %v, want %v", i, rd.ScheduledDate, expectedScheduledDate)
				}

				if rd.IsCompleted != false {
					t.Errorf("FormatWithOverdueMarkedInCompleted() reviewdate[%d].IsCompleted = %v, want %v", i, rd.IsCompleted, false)
				}
			}
		})
	}
}

func TestFormatWithOverdueMarkedCompletedWithIDs(t *testing.T) {
	scheduler := NewScheduler()
	tests := []struct {
		name               string
		targetPatternSteps []*PatternDomain.PatternStep
		reviewDateIDs      []string
		userID             string
		categoryID         *string
		boxID              *string
		itemID             string
		parsedLearnedDate  time.Time
		parsedToday        time.Time
		wantIsFinished     bool
		wantReviewdatesLen int
		wantError          bool
		wantErrorType      error
	}{
		{
			name: "正常なID付きパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1", "rd2"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     true,
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "IDとステップ数の不一致エラー",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     false,
			wantReviewdatesLen: 0,
			wantError:          true,
			wantErrorType:      ErrNewScheduledDateBeforeInitialScheduledDate,
		},
		{
			name: "今日以降の最終ステップ（未完了扱い）",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 10},
			},
			reviewDateIDs:      []string{"rd1", "rd2"},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantIsFinished:     false,
			wantReviewdatesLen: 2,
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReviewdates, gotIsFinished, err := scheduler.FormatWithOverdueMarkedCompletedWithIDs(
				tt.targetPatternSteps,
				tt.reviewDateIDs,
				tt.userID,
				tt.categoryID,
				tt.boxID,
				tt.itemID,
				tt.parsedLearnedDate,
				tt.parsedToday,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && tt.wantErrorType != nil {
				if err != tt.wantErrorType {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() error = %v, want %v", err, tt.wantErrorType)
				}
				return
			}

			if err != nil {
				return
			}

			if gotIsFinished != tt.wantIsFinished {
				t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() gotIsFinished = %v, want %v", gotIsFinished, tt.wantIsFinished)
			}

			if len(gotReviewdates) != tt.wantReviewdatesLen {
				t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdates length = %v, want %v", len(gotReviewdates), tt.wantReviewdatesLen)
			}

			for i, rd := range gotReviewdates {
				if rd.ReviewdateID != tt.reviewDateIDs[i] {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdate[%d].ReviewdateID = %v, want %v", i, rd.ReviewdateID, tt.reviewDateIDs[i])
				}
				if rd.UserID != tt.userID {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdate[%d].UserID = %v, want %v", i, rd.UserID, tt.userID)
				}
				if rd.ItemID != tt.itemID {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdate[%d].ItemID = %v, want %v", i, rd.ItemID, tt.itemID)
				}

				expectedScheduledDate := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[i].IntervalDays)
				if !rd.ScheduledDate.Equal(expectedScheduledDate) {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdate[%d].ScheduledDate = %v, want %v", i, rd.ScheduledDate, expectedScheduledDate)
				}

				expectedIsCompleted := expectedScheduledDate.Before(tt.parsedToday)
				if rd.IsCompleted != expectedIsCompleted {
					t.Errorf("FormatWithOverdueMarkedCompletedWithIDs() reviewdate[%d].IsCompleted = %v, want %v", i, rd.IsCompleted, expectedIsCompleted)
				}
			}
		})
	}
}

func TestFormatWithOverdueMarkedInCompletedWithIDs(t *testing.T) {
	scheduler := NewScheduler()
	tests := []struct {
		name               string
		targetPatternSteps []*PatternDomain.PatternStep
		reviewDateIDs      []string
		userID             string
		categoryID         *string
		boxID              *string
		itemID             string
		parsedLearnedDate  time.Time
		parsedToday        time.Time
		wantReviewdatesLen int
		wantError          bool
		wantErrorType      error
	}{
		{
			name: "正常なID付きパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1", "rd2"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "IDとステップ数の不一致エラー",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 0,
			wantError:          true,
			wantErrorType:      ErrNewScheduledDateBeforeInitialScheduledDate,
		},
		{
			name:               "空のパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{},
			reviewDateIDs:      []string{},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 0,
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReviewdates, err := scheduler.FormatWithOverdueMarkedInCompletedWithIDs(
				tt.targetPatternSteps,
				tt.reviewDateIDs,
				tt.userID,
				tt.categoryID,
				tt.boxID,
				tt.itemID,
				tt.parsedLearnedDate,
				tt.parsedToday,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && tt.wantErrorType != nil {
				if err != tt.wantErrorType {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() error = %v, want %v", err, tt.wantErrorType)
				}
				return
			}

			if err != nil {
				return
			}

			if len(gotReviewdates) != tt.wantReviewdatesLen {
				t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdates length = %v, want %v", len(gotReviewdates), tt.wantReviewdatesLen)
			}

			var expectedAddDuration int
			if len(tt.targetPatternSteps) > 0 {
				firstScheduled := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[0].IntervalDays)
				if firstScheduled.Before(tt.parsedToday) {
					expectedAddDuration = int(tt.parsedToday.Sub(firstScheduled).Hours() / 24)
				} else {
					expectedAddDuration = 0
				}
			}

			for i, rd := range gotReviewdates {
				if rd.ReviewdateID != tt.reviewDateIDs[i] {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdate[%d].ReviewdateID = %v, want %v", i, rd.ReviewdateID, tt.reviewDateIDs[i])
				}
				if rd.UserID != tt.userID {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdate[%d].UserID = %v, want %v", i, rd.UserID, tt.userID)
				}
				if rd.ItemID != tt.itemID {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdate[%d].ItemID = %v, want %v", i, rd.ItemID, tt.itemID)
				}

				expectedScheduledDate := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[i].IntervalDays+expectedAddDuration)
				if !rd.ScheduledDate.Equal(expectedScheduledDate) {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdate[%d].ScheduledDate = %v, want %v", i, rd.ScheduledDate, expectedScheduledDate)
				}

				if rd.IsCompleted != false {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDs() reviewdate[%d].IsCompleted = %v, want %v", i, rd.IsCompleted, false)
				}
			}
		})
	}
}

func TestFormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates(t *testing.T) {
	scheduler := NewScheduler()
	tests := []struct {
		name               string
		targetPatternSteps []*PatternDomain.PatternStep
		reviewDateIDs      []string
		userID             string
		categoryID         *string
		boxID              *string
		itemID             string
		parsedLearnedDate  time.Time
		parsedToday        time.Time
		wantReviewdatesLen int
		wantError          bool
		wantErrorType      error
	}{
		{
			name: "正常なバック復習日処理",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1", "rd2"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 2,
			wantError:          false,
		},
		{
			name: "IDとステップ数の不一致エラー",
			targetPatternSteps: []*PatternDomain.PatternStep{
				{StepNumber: 1, IntervalDays: 1},
				{StepNumber: 2, IntervalDays: 3},
			},
			reviewDateIDs:      []string{"rd1"},
			userID:             "user123",
			categoryID:         stringPtr("cat123"),
			boxID:              stringPtr("box123"),
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 0,
			wantError:          true,
			wantErrorType:      ErrNewScheduledDateBeforeInitialScheduledDate,
		},
		{
			name:               "空のパターンステップでの処理",
			targetPatternSteps: []*PatternDomain.PatternStep{},
			reviewDateIDs:      []string{},
			userID:             "user123",
			categoryID:         nil,
			boxID:              nil,
			itemID:             "item123",
			parsedLearnedDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			parsedToday:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			wantReviewdatesLen: 0,
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReviewdates, err := scheduler.FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates(
				tt.targetPatternSteps,
				tt.reviewDateIDs,
				tt.userID,
				tt.categoryID,
				tt.boxID,
				tt.itemID,
				tt.parsedLearnedDate,
				tt.parsedToday,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && tt.wantErrorType != nil {
				if err != tt.wantErrorType {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() error = %v, want %v", err, tt.wantErrorType)
				}
				return
			}

			if err != nil {
				return
			}

			if len(gotReviewdates) != tt.wantReviewdatesLen {
				t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdates length = %v, want %v", len(gotReviewdates), tt.wantReviewdatesLen)
			}

			for i, rd := range gotReviewdates {
				if rd.ReviewdateID != tt.reviewDateIDs[i] {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdate[%d].ReviewdateID = %v, want %v", i, rd.ReviewdateID, tt.reviewDateIDs[i])
				}
				if rd.UserID != tt.userID {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdate[%d].UserID = %v, want %v", i, rd.UserID, tt.userID)
				}
				if rd.ItemID != tt.itemID {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdate[%d].ItemID = %v, want %v", i, rd.ItemID, tt.itemID)
				}

				expectedScheduledDate := tt.parsedLearnedDate.AddDate(0, 0, tt.targetPatternSteps[i].IntervalDays)
				if !rd.ScheduledDate.Equal(expectedScheduledDate) {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdate[%d].ScheduledDate = %v, want %v", i, rd.ScheduledDate, expectedScheduledDate)
				}

				if rd.IsCompleted != false {
					t.Errorf("FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates() reviewdate[%d].IsCompleted = %v, want %v", i, rd.IsCompleted, false)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
