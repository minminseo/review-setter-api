package item

import (
	"time"

	"github.com/google/uuid"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

// 復習日の計算を担うドメインサービス
type scheduler struct{}

func NewScheduler() IScheduler {
	return &scheduler{}
}

// 作成
func (s *scheduler) FormatWithOverdueMarkedCompleted(
	targetPatternSteps []*PatternDomain.PatternStep,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*Reviewdate, bool, error) {

	result := make([]*Reviewdate, len(targetPatternSteps))

	for i, step := range targetPatternSteps {
		reviewDateID := uuid.NewString()
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays)

		reviewdate, err := NewReviewdate(
			reviewDateID,
			userID,
			categoryID,
			boxID,
			itemID,
			step.StepNumber,
			calculatedScheduledDate,
			calculatedScheduledDate,
			calculatedScheduledDate.Before(parsedToday), // 今日より前なら完了/逆は未完了
		)
		if err != nil {
			return nil, false, err
		}
		result[i] = reviewdate
	}

	// もし最後のステップが今日より前なら（復習物作成の時点で全復習日完了扱いなら）、すべてを完了扱い
	isFinished := false
	if len(targetPatternSteps) > 0 {
		lastScheduled := result[len(targetPatternSteps)-1].ScheduledDate
		isFinished = lastScheduled.Before(parsedToday)
	}
	return result, isFinished, nil
}

func (s *scheduler) FormatWithOverdueMarkedInCompleted(
	targetPatternSteps []*PatternDomain.PatternStep,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*Reviewdate, error) {
	result := make([]*Reviewdate, len(targetPatternSteps))

	var addDuration int
	if len(targetPatternSteps) > 0 {
		firstScheduled := parsedLearnedDate.AddDate(0, 0, targetPatternSteps[0].IntervalDays)
		if firstScheduled.Before(parsedToday) {
			addDuration = int(parsedToday.Sub(firstScheduled).Hours() / 24)
		} else {
			addDuration = 0
		}
	}

	for i, step := range targetPatternSteps {
		reviewDateID := uuid.NewString()
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays+addDuration)

		reviewdate, err := NewReviewdate(
			reviewDateID,
			userID,
			categoryID,
			boxID,
			itemID,
			step.StepNumber,
			calculatedScheduledDate,
			calculatedScheduledDate,
			false, //全部未完了扱い
		)
		if err != nil {
			return nil, err
		}
		result[i] = reviewdate
	}

	return result, nil
}

// 更新
func (s *scheduler) FormatWithOverdueMarkedCompletedWithIDs(
	targetPatternSteps []*PatternDomain.PatternStep,
	reviewDateIDs []string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*Reviewdate, bool, error) {

	if len(reviewDateIDs) != len(targetPatternSteps) {
		return nil, false, ErrNewScheduledDateBeforeInitialScheduledDate
	}

	result := make([]*Reviewdate, len(targetPatternSteps))

	for i, step := range targetPatternSteps {
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays)

		reviewdate, err := NewReviewdate(
			reviewDateIDs[i],
			userID,
			categoryID,
			boxID,
			itemID,
			step.StepNumber,
			calculatedScheduledDate,
			calculatedScheduledDate,
			calculatedScheduledDate.Before(parsedToday),
		)
		if err != nil {
			return nil, false, err
		}
		result[i] = reviewdate
	}

	isFinished := false
	if len(targetPatternSteps) > 0 {
		lastScheduled := result[len(targetPatternSteps)-1].ScheduledDate
		isFinished = lastScheduled.Before(parsedToday)
	}
	return result, isFinished, nil
}

func (s *scheduler) FormatWithOverdueMarkedInCompletedWithIDs(
	targetPatternSteps []*PatternDomain.PatternStep,
	reviewDateIDs []string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*Reviewdate, error) {
	if len(reviewDateIDs) != len(targetPatternSteps) {
		return nil, ErrNewScheduledDateBeforeInitialScheduledDate
	}

	result := make([]*Reviewdate, len(targetPatternSteps))

	var addDuration int
	if len(targetPatternSteps) > 0 {
		firstScheduled := parsedLearnedDate.AddDate(0, 0, targetPatternSteps[0].IntervalDays)
		if firstScheduled.Before(parsedToday) {
			addDuration = int(parsedToday.Sub(firstScheduled).Hours() / 24)
		} else {
			addDuration = 0
		}
	}

	for i, step := range targetPatternSteps {
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays+addDuration)
		reviewdate, err := NewReviewdate(
			reviewDateIDs[i],
			userID,
			categoryID,
			boxID,
			itemID,
			step.StepNumber,
			calculatedScheduledDate,
			calculatedScheduledDate,
			false,
		)
		if err != nil {
			return nil, err
		}
		result[i] = reviewdate
	}

	return result, nil
}

func (s *scheduler) FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates(
	targetPatternSteps []*PatternDomain.PatternStep,
	reviewDateIDs []string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	diff time.Duration,
) ([]*Reviewdate, error) {
	if len(reviewDateIDs) != len(targetPatternSteps) {
		return nil, ErrNewScheduledDateBeforeInitialScheduledDate
	}

	result := make([]*Reviewdate, len(targetPatternSteps))

	for i, step := range targetPatternSteps {
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays)

		reviewdate, err := NewReviewdate(
			reviewDateIDs[i],
			userID,
			categoryID,
			boxID,
			itemID,
			step.StepNumber,
			calculatedScheduledDate,
			calculatedScheduledDate.Add(-diff),
			false, // 全部未完了扱い
		)
		if err != nil {
			return nil, err
		}
		result[i] = reviewdate
	}

	return result, nil
}
