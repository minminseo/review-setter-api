package item

import (
	"context"
	"time"

	"github.com/google/uuid"
	BoxDomain "github.com/minminseo/recall-setter/domain/box"
	CategoryDomain "github.com/minminseo/recall-setter/domain/category"
	ItemDomain "github.com/minminseo/recall-setter/domain/item"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

type ItemUsecase struct {
	categoryRepo       CategoryDomain.ICategoryRepository
	boxRepo            BoxDomain.IBoxRepository
	itemRepo           ItemDomain.IItemRepository
	patternRepo        PatternDomain.IPatternRepository
	transactionManager transaction.ITransactionManager
}

func NewItemUsecase(
	categoryRepo CategoryDomain.ICategoryRepository,
	boxRepo BoxDomain.IBoxRepository,
	itemRepo ItemDomain.IItemRepository,
	patternRepo PatternDomain.IPatternRepository,
	transactionManager transaction.ITransactionManager,
) *ItemUsecase {
	return &ItemUsecase{
		categoryRepo:       categoryRepo,
		boxRepo:            boxRepo,
		itemRepo:           itemRepo,
		patternRepo:        patternRepo,
		transactionManager: transactionManager,
	}
}

// 復習物作成
func (iu *ItemUsecase) CreateItem(ctx context.Context, in CreateItemInput) (*CreateItemOutput, error) {
	ItemID := uuid.NewString()
	parsedLearnedDate, err := time.Parse("2006-01-02", in.LearnedDate)
	if err != nil {
		return nil, err
	}
	registeredAt := time.Now().UTC()
	editedAt := registeredAt

	newItem, err := ItemDomain.NewItem(
		ItemID,
		in.UserID,
		in.CategoryID,
		in.BoxID,
		in.PatternID,
		in.Name,
		in.Detail,
		parsedLearnedDate,
		false, // 初期状態は未完了
		registeredAt,
		editedAt,
	)
	if err != nil {
		return nil, err
	}

	// 永続化
	// patternIDがnilの場合は、復習物のみ永続化してreturn
	if in.PatternID == nil {
		err = iu.itemRepo.CreateItem(ctx, newItem)
		if err != nil {
			return nil, err
		}
	}

	var newReviewdates []*ItemDomain.Reviewdate
	if in.PatternID != nil {
		targetPatternSteps, err := iu.patternRepo.GetAllPatternStepsByPatternID(ctx, *in.PatternID, in.UserID)
		if err != nil {
			return nil, err
		}
		parsedToday, err := time.Parse("2006-01-02", in.Today)
		if err != nil {
			return nil, err
		}

		if in.IsMarkOverdueAsCompleted {
			var isFinished bool
			newReviewdates, isFinished, err = FormatWithOverdueMarkedCompleted(
				targetPatternSteps,
				in.UserID,
				in.CategoryID,
				in.BoxID,
				ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
			// もし最後のステップが今日より前なら（復習物作成の時点で全復習日完了扱いなら）、newItem.isFinishedをtrueにする
			if isFinished {
				newItem.IsFinished = true
			}
		} else {
			newReviewdates, err = FormatWithOverdueMarkedInCompleted(
				targetPatternSteps,
				in.UserID,
				in.CategoryID,
				in.BoxID,
				ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// 永続化
	// ItemとReviewDatesは別テーブルなので同一トランザクションで永続化
	err = iu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {

		err = iu.itemRepo.CreateItem(ctx, newItem)
		if err != nil {
			return err
		}

		// "_"←はCopyfromの返り値の、「挿入された行数」。使わないのでブランク識別子にする。
		_, err = iu.itemRepo.CreateReviewdates(ctx, newReviewdates)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := &CreateItemOutput{
		ItemID:       newItem.ItemID,
		UserID:       newItem.UserID,
		CategoryID:   newItem.CategoryID,
		BoxID:        newItem.BoxID,
		PatternID:    newItem.PatternID,
		Name:         newItem.Name,
		Detail:       newItem.Detail,
		LearnedDate:  (newItem.LearnedDate).Format("2006-01-02"),
		IsCompleted:  newItem.IsFinished,
		RegisteredAt: newItem.RegisteredAt,
		EditedAt:     newItem.EditedAt,
	}
	out.Reviewdates = make([]CreateReviewdateOutput, len(newReviewdates))
	for i, rs := range newReviewdates {
		out.Reviewdates[i] = CreateReviewdateOutput{
			DateID:               rs.ReviewdateID,
			UserID:               rs.UserID,
			ItemID:               rs.ItemID,
			StepNumber:           rs.StepNumber,
			InitialScheduledDate: rs.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        rs.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          rs.IsCompleted,
		}
	}
	return out, nil
}

// 復習物の更新（編集）

/*
<前提条件>
復習物は必ず学習日を持つ。

学習日が変更されているかどうかは問わない
・未設定→設定なら復習日s作成（リクエストに含めた学習日で算出）：完了済み復習日の有無の判定不要
・設定済み→未設定なら復習日s削除：完了済み復習日の有無の判定必要
・復習パターンの長さが不一致なら復習日s一括削除→一括挿入（リクエストに含めた学習日で算出）：完了済み復習日の有無の判定必要
・復習パターンの長さは一致するが、InterVal_daysが不一致なら復習日s更新（リクエストに含めた学習日で算出）：完了済み復習日の有無の判定必要
・未設定→未設定なら復習物のみ更新：完了済み復習日の有無の判定不要

学習日が変更されているかどうかを問う
・同一復習パターンステップ、かつisLearnedDateChanged=trueなら復習日s更新（リクエストに含めた学習日で算出）：完了済み復習日の有無の判定必要
・同一復習パターンステップ、かつisLearnedDateChanged=falseなら復習物のみ更新：完了済み復習日の有無の判定不要

<その他補足>
復習パターン未設定→必ず「完了済み復習日の有無の判定不要」になる
復習パターンは設定済みだが、復習日の変更の必要がない場合→「完了済み復習日の有無の判定不要」になる
*/
func (iu *ItemUsecase) UpdateItem(ctx context.Context, input UpdateItemInput) (*UpdateItemOutput, error) {
	// フラグたち作成のために対象の復習物情報を取得
	targetItem, err := iu.itemRepo.GetItemByID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	targetPatternSteps, err := iu.patternRepo.GetAllPatternStepsByPatternID(ctx, *input.PatternID, input.UserID)
	if err != nil {
		return nil, err
	}

	// 単なるreview_itemsテーブルだけの操作にとどまるかどうか判別するためのフラグ
	isNameAndDetailChanged := targetItem.Name != input.Name || targetItem.Detail != input.Detail

	// 完了済み復習物の有無の判定が必要かどうか判定するための土台づくり
	// 元の復習パターンが未設定だったかどうかのフラグ
	isOriginalPatternNil := targetItem.PatternID == nil

	// 新たに適用したい復習パターンがnilかどうかのフラグ
	isInputPatternNil := input.PatternID == nil

	/* 学習日の変更があるかどうかのフラグ*/
	parsedLearnedDate, err := time.Parse("2006-01-02", input.LearnedDate)
	if err != nil {
		return nil, err
	}
	isLearnedDateChanged := targetItem.LearnedDate != parsedLearnedDate

	// 即エラーリターン用フラグ
	isCategoryChanged := targetItem.CategoryID != input.CategoryID
	isBoxChanged := targetItem.BoxID != input.BoxID
	isPatternChangedWithNillPossible := targetItem.PatternID != input.PatternID
	var isPatternChanged bool
	if !isOriginalPatternNil && !isInputPatternNil {
		isPatternChanged = targetItem.PatternID != input.PatternID
	}
	// 変更点なしの場合、エラーを返す
	if !isCategoryChanged && !isBoxChanged && !isPatternChangedWithNillPossible && !isNameAndDetailChanged && !isLearnedDateChanged {
		return nil, ItemDomain.ErrNoDiff
	}

	/* 「(復習パターンステップが完全一致 || 同一の復習パターンID) && 学習日の変更なし」なら「完了済み復習日の有無の判定不要」になる */
	// ここで、復習日テーブルの操作する時に、「削除→挿入」か「更新」のどちらをするかの判定をするためのフラグ「isTargetPatternStepsLengthChanged」と「isTargetPatternStepsChangedWithoutLengthChange」を作成

	// 具体的な検証内容：shouldCheckRecordExistence
	// 元のパターンIDがnil（新たなパターンIDの値は問わない）：false
	// 同じIDかつ学習日の変更なし：false
	// 同じIDだが学習日の変更あり：true
	// 異なるIDだがパターンステップが完全一致、かつ学習日の変更なし：false
	// 異なるIDだがパターンステップが完全一致、かつ学習日の変更あり：true
	// 異なるIDかつパターンステップが不一致（学習日変更の有無は問わない）：true

	shouldCheckRecordExistence := true
	var isCurrentPatternStepsLengthChanged bool              // これがtrueなら復習日一括削除→一括挿入
	var isCurrentPatternStepsChangedWithoutLengthChange bool // これがtrueでisCurrentPatternStepsLengthChangedがfalseなら復習日一括更新（削除はしない）

	// 元の復習物のパターンIDがNilの場合、復習日は存在しないことが確定するので、完了済み復習日の有無の判定は不要
	if isOriginalPatternNil {
		shouldCheckRecordExistence = false

		// 元の復習物のパターンIDがNil出ない場合、新しいパターンIDフィールド値との関係性を検証
	} else {
		if !isInputPatternNil {
			/* 復習パターンIDが一致、かつ学習日の変更なし */
			if targetItem.PatternID == input.PatternID { // 復習パターンIDが一致
				if !isLearnedDateChanged { // 学習日の変更なし

					shouldCheckRecordExistence = false
				}

				/* 異なる復習パターンIDだがパターンステップが完全一致、かつ学習日の変更なし*/
			} else {
				// パターンステップを取得
				var CurrentPatternSteps []*PatternDomain.PatternStep
				CurrentPatternSteps, err = iu.patternRepo.GetAllPatternStepsByPatternID(ctx, *input.PatternID, input.UserID)
				if err != nil {
					return nil, err
				}

				// パターンステップの長さが異なるかどうか検証
				if len(CurrentPatternSteps) != len(targetPatternSteps) { // targetPatternStepsは新たに適用したい復習パターンのステップ。
					isCurrentPatternStepsLengthChanged = true
				}

				// パターンステップの長さは一致するが、IntervalDaysが異なるかどうか検証
				if !isCurrentPatternStepsLengthChanged {
					for i, step := range CurrentPatternSteps {
						if step.IntervalDays != targetPatternSteps[i].IntervalDays {
							isCurrentPatternStepsChangedWithoutLengthChange = true
							break
						}
					}
				}

				// 「異なる復習パターンIDだがパターンステップが完全一致、かつ学習日の変更なし」が真かどうか検証
				if !isCurrentPatternStepsLengthChanged && !isCurrentPatternStepsChangedWithoutLengthChange { // パターンは異なるが、パターンステップが完全一致
					if !isLearnedDateChanged { // 学習日の変更なし
						shouldCheckRecordExistence = false
					}
				}
			}
		} /*else { // これは不要なコード：elseも記述した場合↓↓↓↓。ネスト深いから書いた
			// 設定済み→Nil：復習日削除のため、完了済み復習日があるかどうかの判定が必要
			shouldCheckRecordExistence = true
		}*/
	}

	editedAt := time.Now().UTC()
	// 更新用のItem完成
	err = targetItem.Set(input.CategoryID, input.BoxID, input.PatternID, input.Name, input.Detail, parsedLearnedDate, editedAt)
	if err != nil {
		return nil, err
	}

	//isNameAndDetailChanged・shouldCheckRecordExistence・isCurrentPatternStepsLengthChanged・isCurrentPatternStepsChangedWithoutLengthChange
	if shouldCheckRecordExistence {
		HasCompleted, err := iu.itemRepo.HasCompletedReviewDateByItemID(ctx, input.ItemID, input.UserID)
		if err != nil {
			return nil, err
		}
		if HasCompleted {
			return nil, ItemDomain.ErrHasCompletedReviewDate
		}
	}

	/* 復習日のオブジェクト生成 */
	parsedToday, err := time.Parse("2006-01-02", input.Today)
	if err != nil {
		return nil, err
	}

	var reviewdatesResult []*ItemDomain.Reviewdate

	if (isCurrentPatternStepsLengthChanged) || (isOriginalPatternNil && !isInputPatternNil) {

		if input.IsMarkOverdueAsCompleted {
			var isFinished bool
			reviewdatesResult, isFinished, err = FormatWithOverdueMarkedCompleted(
				targetPatternSteps,
				input.UserID,
				input.CategoryID,
				input.BoxID,
				input.ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
			// もし最後のステップが今日より前なら（復習物作成の時点で全復習日完了扱いなら）、newItem.isFinishedをtrueにする
			if isFinished {
				targetItem.IsFinished = true
			}
		} else {
			reviewdatesResult, err = FormatWithOverdueMarkedInCompleted(
				targetPatternSteps,
				input.UserID,
				input.CategoryID,
				input.BoxID,
				input.ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if (isCurrentPatternStepsChangedWithoutLengthChange) ||
		(!isCurrentPatternStepsLengthChanged && !isCurrentPatternStepsChangedWithoutLengthChange && isLearnedDateChanged) ||
		(!isPatternChanged && isLearnedDateChanged) {
		reviewDateIDs, err := iu.itemRepo.GetReviewDateIDsByItemID(ctx, input.ItemID, input.UserID)
		if err != nil {
			return nil, err
		}
		var isFinished bool
		if input.IsMarkOverdueAsCompleted {
			reviewdatesResult, isFinished, err = FormatWithOverdueMarkedCompletedWithIDs(
				targetPatternSteps,
				reviewDateIDs,
				input.UserID,
				input.CategoryID,
				input.BoxID,
				input.ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
			// もし最後のステップが今日より前なら（復習物作成の時点で全復習日完了扱いなら）、newItem.isFinishedをtrueにする
			if isFinished {
				targetItem.IsFinished = true
			}
		} else {
			reviewdatesResult, err = FormatWithOverdueMarkedInCompletedWithIDs(
				targetPatternSteps,
				reviewDateIDs,
				input.UserID,
				input.CategoryID,
				input.BoxID,
				input.ItemID,
				parsedLearnedDate,
				parsedToday,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	err = iu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		err = iu.itemRepo.UpdateItem(ctx, targetItem)
		if err != nil {
			return err
		}
		if (!isOriginalPatternNil && isInputPatternNil) || (isCurrentPatternStepsLengthChanged) {
			err = iu.itemRepo.DeleteReviewDates(ctx, input.ItemID, input.UserID)
			if err != nil {
				return err
			}
		}
		if (isCurrentPatternStepsLengthChanged) || (isOriginalPatternNil && !isInputPatternNil) {
			_, err = iu.itemRepo.CreateReviewdates(ctx, reviewdatesResult)
			if err != nil {
				return err
			}
		}
		if (isCurrentPatternStepsChangedWithoutLengthChange) ||
			(!isCurrentPatternStepsLengthChanged && !isCurrentPatternStepsChangedWithoutLengthChange && isLearnedDateChanged) ||
			(!isPatternChanged && isLearnedDateChanged) {
			err = iu.itemRepo.UpdateReviewDates(ctx, reviewdatesResult, input.UserID)
			if err != nil {
				return err
			}
		}
		return nil
	})

	resItem := &UpdateItemOutput{
		ItemID:      targetItem.ItemID,
		UserID:      targetItem.UserID,
		CategoryID:  targetItem.CategoryID,
		BoxID:       targetItem.BoxID,
		PatternID:   targetItem.PatternID,
		Name:        targetItem.Name,
		Detail:      targetItem.Detail,
		LearnedDate: (targetItem.LearnedDate).Format("2006-01-02"),
		IsFinished:  targetItem.IsFinished,
		EditedAt:    targetItem.EditedAt,
	}
	resItem.ReviewDates = make([]UpdateReviewDateOutput, len(reviewdatesResult))
	for i, rs := range reviewdatesResult {
		resItem.ReviewDates[i] = UpdateReviewDateOutput{
			ReviewDateID:         rs.ReviewdateID,
			UserID:               rs.UserID,
			CategoryID:           rs.CategoryID,
			BoxID:                rs.BoxID,
			ItemID:               rs.ItemID,
			StepNumber:           rs.StepNumber,
			InitialScheduledDate: rs.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        rs.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          rs.IsCompleted,
		}
	}

	return resItem, nil
}

// 　復習日の更新（編集）
func (iu *ItemUsecase) UpdateReviewDates(ctx context.Context, input UpdateBackReviewDateInput) (*UpdateBackReviewDateOutput, error) {

	parsedInitialScheduledDate, err := time.Parse("2006-01-02", input.InitialScheduledDate)
	if err != nil {
		return nil, err
	}
	parsedNewScheduledDate, err := time.Parse("2006-01-02", input.RequestScheduledDate)
	if err != nil {
		return nil, err
	}
	if parsedNewScheduledDate.Before(parsedInitialScheduledDate) {
		return nil, ItemDomain.ErrNewScheduledDateBeforeInitialScheduledDate
	}

	// 復習日再生成・IsCompletedのbool値判別ロジック
	calculatedDuration := int(parsedNewScheduledDate.Sub(parsedInitialScheduledDate).Hours() / 24)
	parsedLearnedDate, err := time.Parse("2006-01-02", input.LearnedDate)
	if err != nil {
		return nil, err
	}
	FakeLearnedDate := parsedLearnedDate.AddDate(0, 0, calculatedDuration) // これでFormat〇〇系の関数を使い回せる

	parsedToday, err := time.Parse("2006-01-02", input.Today)
	if err != nil {
		return nil, err
	}

	// input.PatternStepsInReviewDateをFormatWithOverdueMarkedCompletedに渡せるように型変換
	var targetPatternSteps []*PatternDomain.PatternStep
	for _, step := range input.PatternStepsInReviewDate {
		patternStep := &PatternDomain.PatternStep{
			PatternStepID: step.PatternStepID,
			UserID:        step.UserID,
			PatternID:     step.PatternID,
			StepNumber:    step.StepNumber,
			IntervalDays:  step.IntervalDays,
		}
		targetPatternSteps = append(targetPatternSteps, patternStep)
	}

	var newReviewdates []*ItemDomain.Reviewdate

	reviewDateIDs, err := iu.itemRepo.GetReviewDateIDsByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	var isFinished bool
	if input.IsMarkOverdueAsCompleted {
		newReviewdates, isFinished, err = FormatWithOverdueMarkedCompletedWithIDs(
			targetPatternSteps,
			reviewDateIDs,
			input.UserID,
			input.CategoryID,
			input.BoxID,
			input.ItemID,
			FakeLearnedDate,
			parsedToday,
		)
		if err != nil {
			return nil, err
		}
		// isFinishedがtrueの場合、UpdateItemAsFinishedを実行
	} else {
		newReviewdates, err = FormatWithOverdueMarkedInCompletedWithIDs(
			targetPatternSteps,
			reviewDateIDs,
			input.UserID,
			input.CategoryID,
			input.BoxID,
			input.ItemID,
			FakeLearnedDate,
			parsedToday,
		)
		if err != nil {
			return nil, err
		}
	}

	// 操作対象の復習日以降の復習日のみ抽出
	filteredReviewdates := make([]*ItemDomain.Reviewdate, 0, len(newReviewdates)) // 操作対象が1個目の可能性もあるので容量はlen(newReviewdates)で初期化（最大値）
	for _, Reviewdate := range newReviewdates {
		if Reviewdate.StepNumber >= input.StepNumber {
			filteredReviewdates = append(filteredReviewdates, Reviewdate)
		}
	}

	// isFinishedがtrueの場合、review_itemテーブルも操作するため、isFinishedで分岐し適宜トランザクションをはる。
	targetEditedAt, err := iu.itemRepo.GetEditedAtByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	resultEditedAt := targetEditedAt
	if isFinished {
		err = iu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {

			resultEditedAt = time.Now().UTC()
			err = iu.itemRepo.UpdateItemAsFinished(ctx, input.ItemID, input.UserID, resultEditedAt)
			if err != nil {
				return err
			}
			err = iu.itemRepo.UpdateReviewDates(ctx, filteredReviewdates, input.UserID)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		err = iu.itemRepo.UpdateReviewDates(ctx, filteredReviewdates, input.UserID)
		if err != nil {
			return nil, err
		}
	}

	// 最新の復習日たちをDBから取得（クライアントで復習日のうち何回目以降を上書きすべきか考慮せずに済むため）
	latestReviewdates, err := iu.itemRepo.GetReviewDatesByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}

	res := &UpdateBackReviewDateOutput{
		ItemID:     input.ItemID,
		UserID:     input.UserID,
		IsFinished: isFinished,
		EditedAt:   resultEditedAt,
	}
	res.ReviewDates = make([]UpdateReviewDateOutput, len(latestReviewdates))
	for i, rs := range latestReviewdates {
		res.ReviewDates[i] = UpdateReviewDateOutput{
			ReviewDateID:         rs.ReviewdateID,
			UserID:               rs.UserID,
			CategoryID:           rs.CategoryID,
			BoxID:                rs.BoxID,
			ItemID:               rs.ItemID,
			StepNumber:           rs.StepNumber,
			InitialScheduledDate: rs.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        rs.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          rs.IsCompleted,
		}
	}

	return res, nil
}

// 復習物を手動で途中完了に更新
func (iu *ItemUsecase) UpdateItemAsFinishedForce(ctx context.Context, input UpdateItemAsFinishedForceInput) (*UpdateItemAsFinishedForceOutput, error) {

	targetItem, err := iu.itemRepo.GetItemByID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	editedAt := time.Now().UTC()
	err = iu.itemRepo.UpdateItemAsFinished(ctx, input.ItemID, input.UserID, editedAt)
	if err != nil {
		return nil, err
	}

	resItem := &UpdateItemAsFinishedForceOutput{
		ItemID:     targetItem.ItemID,
		UserID:     targetItem.UserID,
		IsFinished: true,
		EditedAt:   editedAt,
	}

	return resItem, nil
}

// 復習物の復習日を完了済みに更新
func (iu *ItemUsecase) UpdateReviewDateAsCompleted(ctx context.Context, input UpdateReviewDateAsCompletedInput) (*UpdateReviewDateAsCompletedOutput, error) {
	targetReviewdates, err := iu.itemRepo.GetReviewDatesByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}

	// targetReviewdatesの最後の復習日のStepNumberがinput.StepNumberと一致するかどうかを判定するフラグを作成
	var isLastStepNumberMatch bool
	// ここで検証
	lastReviewdate := targetReviewdates[len(targetReviewdates)-1]
	if lastReviewdate.StepNumber == input.StepNumber {
		isLastStepNumberMatch = true
	} else {
		isLastStepNumberMatch = false
	}

	targetEditedAt, err := iu.itemRepo.GetEditedAtByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	resultEditedAt := targetEditedAt
	// 最後の復習日が完了した場合、復習物を完了済みに更新
	if isLastStepNumberMatch {
		err = iu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			err = iu.itemRepo.UpdateReviewDateAsCompleted(ctx, input.ReviewDateID, input.UserID)
			if err != nil {
				return nil
			}

			resultEditedAt = time.Now().UTC()
			err = iu.itemRepo.UpdateItemAsFinished(ctx, input.ItemID, input.UserID, resultEditedAt)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		// 復習物そのものは未完了のまま
	} else {
		err = iu.itemRepo.UpdateReviewDateAsCompleted(ctx, input.ReviewDateID, input.UserID)
		if err != nil {
			return nil, err
		}
	}

	resReviewdate := &UpdateReviewDateAsCompletedOutput{
		ReviewDateID: input.ReviewDateID,
		UserID:       input.UserID,
		IsCompleted:  true,
		IsFinished:   isLastStepNumberMatch,
		EditedAt:     resultEditedAt,
	}

	return resReviewdate, nil
}

// 復習物の復習日を未完了に更新
func (iu *ItemUsecase) UpdateReviewDateAsInCompleted(ctx context.Context, input UpdateReviewDateAsInCompletedInput) (*UpdateReviewDateAsInCompletedOutput, error) {
	targetItem, err := iu.itemRepo.GetItemByID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}

	var isItemFinished bool
	if targetItem.IsFinished {
		isItemFinished = true
	}

	targetEditedAt, err := iu.itemRepo.GetEditedAtByItemID(ctx, input.ItemID, input.UserID)
	if err != nil {
		return nil, err
	}
	resultEditedAt := targetEditedAt
	if isItemFinished {
		err = iu.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
			err = iu.itemRepo.UpdateReviewDateAsInCompleted(ctx, input.ReviewDateID, input.UserID)
			if err != nil {
				return err
			}

			resultEditedAt = time.Now().UTC()
			err = iu.itemRepo.UpdateItemAsUnFinished(ctx, input.ItemID, input.UserID, resultEditedAt)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		// 復習物が未完了の場合
	} else {
		err = iu.itemRepo.UpdateReviewDateAsInCompleted(ctx, input.ReviewDateID, input.UserID)
		if err != nil {
			return nil, err
		}
	}

	resReviewdate := &UpdateReviewDateAsInCompletedOutput{
		ReviewDateID: input.ReviewDateID,
		UserID:       input.UserID,
		IsCompleted:  false,
		IsFinished:   isItemFinished,
		EditedAt:     resultEditedAt,
	}

	return resReviewdate, nil
}

// 途中完了復習物再開リクエスト
func (iu *ItemUsecase) UpdateItemAsUnFinishedForce(ctx context.Context, input UpdateItemAsUnFinishedForceInput) (*UpdateItemAsUnFinishedForceOutput, error) {
	EditedAt := time.Now().UTC()
	err := iu.itemRepo.UpdateItemAsUnFinished(ctx, input.ItemID, input.UserID, EditedAt)
	if err != nil {
		return nil, err
	}

	resItem := &UpdateItemAsUnFinishedForceOutput{
		ItemID:     input.ItemID,
		UserID:     input.UserID,
		IsFinished: false,
		EditedAt:   EditedAt,
	}

	return resItem, nil
}

// 物理削除
// TODO: 論理削除に変更する（影響範囲を確認してから）
func (iu *ItemUsecase) DeleteItem(ctx context.Context, itemID string, userID string) error {
	err := iu.itemRepo.DeleteItem(ctx, itemID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (iu *ItemUsecase) GetAllUnFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*GetItemOutput, error) {
	items, err := iu.itemRepo.GetAllUnFinishedItemsByBoxID(ctx, boxID, userID)
	if err != nil {
		return nil, err
	}

	// アイテムIDをキーに未完了アイテムをマップ化
	unfinishedItemMap := make(map[string]struct{}, len(items))
	for _, item := range items {
		unfinishedItemMap[item.ItemID] = struct{}{}
	}

	// 全復習日を取得
	reviewdates, err := iu.itemRepo.GetAllReviewDatesByBoxID(ctx, boxID, userID)
	if err != nil {
		return nil, err
	}

	// 未完了アイテムを親に持つ復習日のみ抽出
	reviewdatesByItem := make(map[string][]GetReviewDateOutput, len(items))

	for _, reviewdate := range reviewdates {
		if _, ok := unfinishedItemMap[reviewdate.ItemID]; !ok {
			continue
		}
		reviewdatesByItem[reviewdate.ItemID] = append(reviewdatesByItem[reviewdate.ItemID], GetReviewDateOutput{
			ReviewDateID:         reviewdate.ReviewdateID,
			UserID:               reviewdate.UserID,
			CategoryID:           reviewdate.CategoryID,
			BoxID:                reviewdate.BoxID,
			ItemID:               reviewdate.ItemID,
			StepNumber:           reviewdate.StepNumber,
			InitialScheduledDate: reviewdate.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        reviewdate.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          reviewdate.IsCompleted,
		})
	}

	result := make([]*GetItemOutput, 0, len(items))
	for _, item := range items {
		result = append(result, &GetItemOutput{
			ItemID:       item.ItemID,
			UserID:       item.UserID,
			CategoryID:   item.CategoryID,
			BoxID:        item.BoxID,
			PatternID:    item.PatternID,
			Name:         item.Name,
			Detail:       item.Detail,
			LearnedDate:  item.LearnedDate.Format("2006-01-02"),
			IsFinished:   item.IsFinished,
			RegisteredAt: item.RegisteredAt,
			EditedAt:     item.EditedAt,
			ReviewDates:  reviewdatesByItem[item.ItemID],
		})
	}
	return result, nil
}

func (iu *ItemUsecase) GetAllUnFinishedUnclassifiedItemsByUserID(ctx context.Context, userID string) ([]*GetItemOutput, error) {
	items, err := iu.itemRepo.GetAllUnFinishedUnclassifiedItemsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// アイテムIDをキーに未分類未完了アイテムをマップ化
	unfinishedItemMap := make(map[string]struct{}, len(items))
	for _, item := range items {
		unfinishedItemMap[item.ItemID] = struct{}{}
	}

	// 全復習日を取得
	reviewdates, err := iu.itemRepo.GetAllUnclassifiedReviewDatesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 未分類未完了アイテムを親に持つ復習日のみ抽出
	reviewdatesByItem := make(map[string][]GetReviewDateOutput, len(items))

	for _, reviewdate := range reviewdates {
		if _, ok := unfinishedItemMap[reviewdate.ItemID]; !ok {
			continue
		}
		reviewdatesByItem[reviewdate.ItemID] = append(reviewdatesByItem[reviewdate.ItemID], GetReviewDateOutput{
			ReviewDateID:         reviewdate.ReviewdateID,
			UserID:               reviewdate.UserID,
			CategoryID:           reviewdate.CategoryID,
			BoxID:                reviewdate.BoxID,
			ItemID:               reviewdate.ItemID,
			StepNumber:           reviewdate.StepNumber,
			InitialScheduledDate: reviewdate.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        reviewdate.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          reviewdate.IsCompleted,
		})
	}

	result := make([]*GetItemOutput, 0, len(items))
	for _, item := range items {
		result = append(result, &GetItemOutput{
			ItemID:       item.ItemID,
			UserID:       item.UserID,
			CategoryID:   item.CategoryID,
			BoxID:        item.BoxID,
			PatternID:    item.PatternID,
			Name:         item.Name,
			Detail:       item.Detail,
			LearnedDate:  item.LearnedDate.Format("2006-01-02"),
			IsFinished:   item.IsFinished,
			RegisteredAt: item.RegisteredAt,
			EditedAt:     item.EditedAt,
			ReviewDates:  reviewdatesByItem[item.ItemID],
		})
	}
	return result, nil
}

func (iu *ItemUsecase) GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx context.Context, userID string, categoryID string) ([]*GetItemOutput, error) {
	items, err := iu.itemRepo.GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx, userID, categoryID)
	if err != nil {
		return nil, err
	}

	// アイテムIDをキーに未分類未完了アイテムをマップ化
	unfinishedItemMap := make(map[string]struct{}, len(items))
	for _, item := range items {
		unfinishedItemMap[item.ItemID] = struct{}{}
	}

	// 全復習日を取得
	reviewdates, err := iu.itemRepo.GetAllUnclassifiedReviewDatesByCategoryID(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}

	// 未分類未完了アイテムを親に持つ復習日のみ抽出
	reviewdatesByItem := make(map[string][]GetReviewDateOutput, len(items))

	for _, reviewdate := range reviewdates {
		if _, ok := unfinishedItemMap[reviewdate.ItemID]; !ok {
			continue
		}
		reviewdatesByItem[reviewdate.ItemID] = append(reviewdatesByItem[reviewdate.ItemID], GetReviewDateOutput{
			ReviewDateID:         reviewdate.ReviewdateID,
			UserID:               reviewdate.UserID,
			CategoryID:           reviewdate.CategoryID,
			BoxID:                reviewdate.BoxID,
			ItemID:               reviewdate.ItemID,
			StepNumber:           reviewdate.StepNumber,
			InitialScheduledDate: reviewdate.InitialScheduledDate.Format("2006-01-02"),
			ScheduledDate:        reviewdate.ScheduledDate.Format("2006-01-02"),
			IsCompleted:          reviewdate.IsCompleted,
		})
	}

	result := make([]*GetItemOutput, 0, len(items))
	for _, item := range items {
		result = append(result, &GetItemOutput{
			ItemID:       item.ItemID,
			UserID:       item.UserID,
			CategoryID:   item.CategoryID,
			BoxID:        item.BoxID,
			PatternID:    item.PatternID,
			Name:         item.Name,
			Detail:       item.Detail,
			LearnedDate:  item.LearnedDate.Format("2006-01-02"),
			IsFinished:   item.IsFinished,
			RegisteredAt: item.RegisteredAt,
			EditedAt:     item.EditedAt,
			ReviewDates:  reviewdatesByItem[item.ItemID],
		})
	}
	return result, nil
}

func (iu *ItemUsecase) CountItemsGroupedByBoxByUserID(ctx context.Context, userID string) ([]*ItemCountGroupedByBoxOutput, error) {
	counts, err := iu.itemRepo.CountItemsGroupedByBoxByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*ItemCountGroupedByBoxOutput, len(counts))
	for i, count := range counts {
		result[i] = &ItemCountGroupedByBoxOutput{
			CategoryID: count.CategoryID,
			BoxID:      count.BoxID,
			Count:      count.Count,
		}
	}
	return result, nil
}

func (iu *ItemUsecase) CountUnclassifiedItemsGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*UnclassifiedItemCountGroupedByCategoryOutput, error) {
	counts, err := iu.itemRepo.CountUnclassifiedItemsGroupedByCategoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*UnclassifiedItemCountGroupedByCategoryOutput, len(counts))
	for i, count := range counts {
		result[i] = &UnclassifiedItemCountGroupedByCategoryOutput{
			CategoryID: count.CategoryID,
			Count:      count.Count,
		}
	}
	return result, nil
}

func (iu *ItemUsecase) CountUnclassifiedItemsByUserID(ctx context.Context, userID string) (int, error) {
	count, err := iu.itemRepo.CountUnclassifiedItemsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (iu *ItemUsecase) CountDailyDatesGroupedByBoxByUserID(ctx context.Context, userID string, today string) ([]*DailyCountGroupedByBoxOutput, error) {
	parsedToday, err := time.Parse("2006-01-02", today)
	if err != nil {
		return nil, err
	}

	counts, err := iu.itemRepo.CountDailyDatesGroupedByBoxByUserID(ctx, userID, parsedToday)
	if err != nil {
		return nil, err
	}

	result := make([]*DailyCountGroupedByBoxOutput, len(counts))
	for i, count := range counts {
		result[i] = &DailyCountGroupedByBoxOutput{
			CategoryID: count.CategoryID,
			BoxID:      count.BoxID,
			Count:      count.Count,
		}
	}
	return result, nil
}

func (iu *ItemUsecase) CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string, today string) ([]*UnclassifiedDailyDatesCountGroupedByCategoryOutput, error) {
	parsedToday, err := time.Parse("2006-01-02", today)
	if err != nil {
		return nil, err
	}

	counts, err := iu.itemRepo.CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx, userID, parsedToday)
	if err != nil {
		return nil, err
	}

	result := make([]*UnclassifiedDailyDatesCountGroupedByCategoryOutput, len(counts))
	for i, count := range counts {
		result[i] = &UnclassifiedDailyDatesCountGroupedByCategoryOutput{
			CategoryID: count.CategoryID,
			Count:      count.Count,
		}
	}
	return result, nil
}

func (iu *ItemUsecase) CountDailyDatesUnclassifiedByUserID(ctx context.Context, userID string, today string) (int, error) {
	parsedToday, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0, err
	}

	count, err := iu.itemRepo.CountDailyDatesUnclassifiedByUserID(ctx, userID, parsedToday)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 今日の全復習日数を取得
func (iu *ItemUsecase) CountAllDailyReviewDates(ctx context.Context, userID string, today string) (int, error) {
	parsedToday, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0, err
	}

	count, err := iu.itemRepo.CountAllDailyReviewDates(ctx, userID, parsedToday)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// TODO: ボックスレベルの完了済みの過去日の復習日を今日に変更するユースケース実装
// TODO: 完了した復習物（is_finishedがtrue）を取得するユースケース実装

func (iu *ItemUsecase) GetAllDailyReviewDates(ctx context.Context, userID string, today string) (*GetDailyReviewDatesOutput, error) {
	parsedToday, err := time.Parse("2006-01-02", today)
	if err != nil {
		return nil, err
	}

	// ユーザー直下（NULL／NULL）の未分類ボックス今日の復習日、カテゴリー毎（非NULL／NULL）の未分類ボックスの今日の復習日、ボックス毎（非NULL／非NULL）の復習日をまとめて取得。
	dailyDates, err := iu.itemRepo.GetAllDailyReviewDates(ctx, userID, parsedToday)
	if err != nil {
		return nil, err
	}

	// 一意なIDを保持するためのセットを作成
	categorySet := make(map[string]struct{})
	boxSet := make(map[string]struct{})

	// category_idとbox_idをセットに追加
	for _, d := range dailyDates {
		if d.CategoryID != nil {
			categorySet[*d.CategoryID] = struct{}{}
		}
		if d.BoxID != nil {
			boxSet[*d.BoxID] = struct{}{}
		}
	}

	// カテゴリー名を一括取得
	categoryIDs := make([]string, 0, len(categorySet))
	for id := range categorySet {
		categoryIDs = append(categoryIDs, id)
	}
	categories, err := iu.categoryRepo.GetCategoryNamesByCategoryIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}
	// ID→Nameのマップを作成
	categoryMap := make(map[string]string, len(categories))
	for _, c := range categories {
		categoryMap[c.ID] = c.Name
	}

	// ボックス名とpattern_idを一括取得
	boxIDs := make([]string, 0, len(boxSet))
	for id := range boxSet {
		boxIDs = append(boxIDs, id)
	}
	boxes, err := iu.boxRepo.GetBoxNamesByBoxIDs(ctx, boxIDs)
	if err != nil {
		return nil, err
	}

	// ID→NameとID→PatternIDのマップを作成
	boxNameMap := make(map[string]string, len(boxes))
	boxPatternMap := make(map[string]string, len(boxes))

	// 一意なパターンIDを保持するためのセット
	patternSet := make(map[string]struct{})
	for _, b := range boxes {
		boxNameMap[b.BoxID] = b.Name
		boxPatternMap[b.BoxID] = b.PatternID

		// 一意なパターンIDを保持するためのセットに追加
		patternSet[b.PatternID] = struct{}{}
	}

	// 6. パターンIDsでtarget_weightを一括取得
	patternIDs := make([]string, 0, len(patternSet))
	for id := range patternSet {
		patternIDs = append(patternIDs, id)
	}
	patterns, err := iu.patternRepo.GetPatternTargetWeightsByPatternIDs(ctx, patternIDs)
	if err != nil {
		return nil, err
	}
	// ID→TargetWeightのマップを作成
	patternMap := make(map[string]string, len(patterns))
	for _, p := range patterns {
		patternMap[p.PatternID] = p.TargetWeight
	}

	// 結果を組み立てていく。
	out := &GetDailyReviewDatesOutput{
		Categories:                    []DailyReviewDatesGroupedByCategoryOutput{},
		DailyReviewDatesGroupedByUser: []UnclassifiedDailyReviewDatesGroupedByUserOutput{},
	}
	categoryIndex := make(map[string]int)
	boxIndex := make(map[string]int)

	// 8. 一つずつマッピングとグルーピング
	for _, d := range dailyDates {
		var prev, next *string
		// ScheduledDateのstep_numberが1の時、PrevScheduledDateは存在しないという例を考慮してnil確認
		if d.PrevScheduledDate != nil {
			s := d.PrevScheduledDate.Format("2006-01-02")
			prev = &s
		}
		sched := d.ScheduledDate.Format("2006-01-02")
		// ScheduledDateのstep_numberが最後の時、NextScheduledDateは存在しないという例を考慮してnil確認
		if d.NextScheduledDate != nil {
			s := d.NextScheduledDate.Format("2006-01-02")
			next = &s
		}

		// 未分類 (category=nil && box=nil)の場合、ユーザー直下グループに追加
		if d.CategoryID == nil && d.BoxID == nil {
			out.DailyReviewDatesGroupedByUser = append(out.DailyReviewDatesGroupedByUser,
				UnclassifiedDailyReviewDatesGroupedByUserOutput{
					ReviewDateID:      d.ReviewdateID,
					StepNumber:        d.StepNumber,
					PrevScheduledDate: prev,
					ScheduledDate:     sched,
					NextScheduledDate: next,
					IsCompleted:       d.IsCompleted,
					ItemName:          d.Name,
					Detail:            d.Detail,
					RegisteredAt:      d.RegisteredAt,
					EditedAt:          d.EditedAt,
				},
			)
			// 未分類の分岐を通った場合、その後のカテゴリー・ボックス振り分け処理は不要なのでcontinue
			continue
		}

		var categoryID, boxID string
		if d.CategoryID != nil {
			categoryID = *d.CategoryID
		}
		if d.BoxID != nil {
			boxID = *d.BoxID
		}

		// カテゴリーグループ初期化
		ci, ok := categoryIndex[categoryID]
		if !ok {
			out.Categories = append(out.Categories, DailyReviewDatesGroupedByCategoryOutput{
				CategoryID:                             categoryID,
				CategoryName:                           categoryMap[categoryID],
				Boxes:                                  []DailyReviewDatesGroupedByBoxOutput{},
				UnclassifiedDailyReviewDatesByCategory: []UnclassifiedDailyReviewDatesGroupedByCategoryOutput{},
			})
			ci = len(out.Categories) - 1
			categoryIndex[categoryID] = ci
		}
		categoryGroup := &out.Categories[ci]

		// ボックス未分類 (box=nil)の場合、カテゴリー毎の未分類に追加
		if d.BoxID == nil {
			categoryGroup.UnclassifiedDailyReviewDatesByCategory = append(categoryGroup.UnclassifiedDailyReviewDatesByCategory,
				UnclassifiedDailyReviewDatesGroupedByCategoryOutput{
					ReviewDateID:      d.ReviewdateID,
					CategoryID:        categoryID,
					StepNumber:        d.StepNumber,
					PrevScheduledDate: prev,
					ScheduledDate:     sched,
					NextScheduledDate: next,
					IsCompleted:       d.IsCompleted,
					ItemName:          d.Name,
					Detail:            d.Detail,
					RegisteredAt:      d.RegisteredAt,
					EditedAt:          d.EditedAt,
				},
			)
			// ボックス未分類の分岐を通った場合、その後のボックスグループ振り分け処理は不要なのでcontinue
			continue
		}

		// ボックスグループ初期化
		key := categoryID + "|" + boxID
		bi, ok := boxIndex[key]
		if !ok {
			boxName := boxNameMap[boxID]
			patternID := boxPatternMap[boxID]
			targetWeight := patternMap[patternID]

			out.Categories[ci].Boxes = append(categoryGroup.Boxes,
				DailyReviewDatesGroupedByBoxOutput{
					BoxID:        boxID,
					CategoryID:   categoryID,
					BoxName:      boxName,
					TargetWeight: targetWeight,
					ReviewDates:  []DailyReviewDatesByBoxOutput{},
				},
			)
			bi = len(categoryGroup.Boxes) - 1
			boxIndex[key] = bi
		}
		boxGroup := &categoryGroup.Boxes[bi]

		// 今日の復習日データをボックスグループに追加
		boxGroup.ReviewDates = append(boxGroup.ReviewDates,
			DailyReviewDatesByBoxOutput{
				ReviewDateID:      d.ReviewdateID,
				CategoryID:        categoryID,
				BoxID:             boxID,
				StepNumber:        d.StepNumber,
				PrevScheduledDate: prev,
				ScheduledDate:     sched,
				NextScheduledDate: next,
				IsCompleted:       d.IsCompleted,
				ItemName:          d.Name,
				Detail:            d.Detail,
				RegisteredAt:      d.RegisteredAt,
				EditedAt:          d.EditedAt,
			},
		)
	}
	return out, nil
}
