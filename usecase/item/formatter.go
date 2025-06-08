package item

import (
	"time"

	"github.com/google/uuid"
	ItemDomain "github.com/minminseo/recall-setter/domain/item"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

// 作成
func FormatWithOverdueMarkedCompleted(
	targetPatternSteps []*PatternDomain.PatternStep,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*ItemDomain.Reviewdate, bool, error) {

	result := make([]*ItemDomain.Reviewdate, len(targetPatternSteps))

	for i, step := range targetPatternSteps {
		reviewDateID := uuid.NewString()
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays)

		reviewdate, err := ItemDomain.NewReviewdate(
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
	lastScheduled := result[len(targetPatternSteps)-1].ScheduledDate
	isFinished := lastScheduled.Before(parsedToday)
	return result, isFinished, nil
}

func FormatWithOverdueMarkedInCompleted(
	targetPatternSteps []*PatternDomain.PatternStep,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*ItemDomain.Reviewdate, error) {
	result := make([]*ItemDomain.Reviewdate, len(targetPatternSteps))
	firstScheduled := parsedLearnedDate.AddDate(0, 0, targetPatternSteps[0].IntervalDays)
	var addDuration int
	if firstScheduled.Before(parsedToday) {
		addDuration = int(parsedToday.Sub(firstScheduled).Hours() / 24)
	} else {
		addDuration = 0
	}

	for i, step := range targetPatternSteps {
		reviewDateID := uuid.NewString()
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays+addDuration)

		reviewdate, err := ItemDomain.NewReviewdate(
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
func FormatWithOverdueMarkedCompletedWithIDs(
	targetPatternSteps []*PatternDomain.PatternStep,
	reviewDateIDs []string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*ItemDomain.Reviewdate, bool, error) {

	if len(reviewDateIDs) != len(targetPatternSteps) {
		return nil, false, ItemDomain.ErrNewScheduledDateBeforeInitialScheduledDate
	}

	result := make([]*ItemDomain.Reviewdate, len(targetPatternSteps))

	for i, step := range targetPatternSteps {
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays)

		reviewdate, err := ItemDomain.NewReviewdate(
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

	lastScheduled := result[len(targetPatternSteps)-1].ScheduledDate
	isFinished := lastScheduled.Before(parsedToday)
	return result, isFinished, nil
}

func FormatWithOverdueMarkedInCompletedWithIDs(
	targetPatternSteps []*PatternDomain.PatternStep,
	reviewDateIDs []string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	parsedLearnedDate time.Time,
	parsedToday time.Time,
) ([]*ItemDomain.Reviewdate, error) {
	if len(reviewDateIDs) != len(targetPatternSteps) {
		return nil, ItemDomain.ErrNewScheduledDateBeforeInitialScheduledDate
	}

	result := make([]*ItemDomain.Reviewdate, len(targetPatternSteps))
	firstScheduled := parsedLearnedDate.AddDate(0, 0, targetPatternSteps[0].IntervalDays)
	var addDuration int
	if firstScheduled.Before(parsedToday) {
		addDuration = int(parsedToday.Sub(firstScheduled).Hours() / 24)
	} else {
		addDuration = 0
	}

	for i, step := range targetPatternSteps {
		calculatedScheduledDate := parsedLearnedDate.AddDate(0, 0, step.IntervalDays+addDuration)

		reviewdate, err := ItemDomain.NewReviewdate(
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
