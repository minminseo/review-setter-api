package pattern

import (
	"context"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/minminseo/recall-setter/domain/item"
	patternDomain "github.com/minminseo/recall-setter/domain/pattern"
	"github.com/minminseo/recall-setter/usecase/transaction"
)

type patternUsecase struct {
	patternRepo patternDomain.IPatternRepository
	itemRepo    itemDomain.IItemRepository
	// ここでtransactionManagerを使うのは、patternとpatternStepを同一トランザクションで永続化するため。
	transactionManeger transaction.ITransactionManager
}

func NewPatternUsecase(
	patternRepo patternDomain.IPatternRepository,
	itemRepo itemDomain.IItemRepository,
	transactionManeger transaction.ITransactionManager,
) IPatternUsecase {
	return &patternUsecase{
		patternRepo:        patternRepo,
		itemRepo:           itemRepo,
		transactionManeger: transactionManeger,
	}
}

func (pu *patternUsecase) CreatePattern(ctx context.Context, in CreatePatternInput) (*CreatePatternOutput, error) {
	patternID := uuid.NewString()

	registeredAt := time.Now().UTC()
	editedAt := registeredAt

	newPattern, err := patternDomain.NewPattern(
		patternID,
		in.UserID,
		in.Name,
		in.TargetWeight,
		registeredAt,
		editedAt,
	)
	if err != nil {
		return nil, err
	}

	newSteps := make([]*patternDomain.PatternStep, len(in.Steps))
	for i, s := range in.Steps {
		stepID := uuid.NewString()
		patternStep, err := patternDomain.NewPatternStep(
			stepID,
			in.UserID,
			patternID,
			s.StepNumber,
			s.IntervalDays,
		)
		if err != nil {
			return nil, err
		}
		newSteps[i] = patternStep
	}

	err = patternDomain.ValidateSteps(newSteps)
	if err != nil {
		return nil, err
	}

	// patternとstepは別テーブルなので同一トランザクションで永続化
	err = pu.transactionManeger.RunInTransaction(ctx, func(ctx context.Context) error {

		err = pu.patternRepo.CreatePattern(ctx, newPattern)
		if err != nil {
			return err
		}

		// "_"←はCopyfromの返り値の、「挿入された行数」。使わないのでブランク識別子にする。
		_, err = pu.patternRepo.CreatePatternSteps(ctx, newSteps)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := &CreatePatternOutput{
		ID:           newPattern.PatternID,
		UserID:       newPattern.UserID,
		Name:         newPattern.Name,
		TargetWeight: string(newPattern.TargetWeight),
		RegisteredAt: newPattern.RegisteredAt,
		EditedAt:     newPattern.EditedAt,
	}
	out.Steps = make([]CreatePatternStepOutput, len(newSteps))
	for i, ps := range newSteps {
		out.Steps[i] = CreatePatternStepOutput{
			PatternStepID: ps.PatternStepID,
			UserID:        ps.UserID,
			PatternID:     ps.PatternID,
			StepNumber:    ps.StepNumber,
			IntervalDays:  ps.IntervalDays,
		}
	}
	return out, nil
}

func (pu *patternUsecase) GetPatternsByUserID(ctx context.Context, userID string) ([]*GetPatternOutput, error) {
	allPatterns, err := pu.patternRepo.GetAllPatternsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	allPatternSteps, err := pu.patternRepo.GetAllPatternStepsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// pattern_idをキーにして、GetPatternStepOutputのスライスを値にして復習パターン毎にグルーピング
	stepsByPattern := make(map[string][]GetPatternStepOutput, len(allPatternSteps))
	for _, domainStep := range allPatternSteps {

		stepOutput := GetPatternStepOutput{
			PatternStepID: domainStep.PatternStepID,
			PatternID:     domainStep.PatternID,
			StepNumber:    domainStep.StepNumber,
			IntervalDays:  domainStep.IntervalDays,
		}
		stepsByPattern[domainStep.PatternID] = append(stepsByPattern[domainStep.PatternID], stepOutput)
	}

	// 4) stepsByPatternを使ってGetPatternOutputのスライスを生成
	var result []*GetPatternOutput
	result = make([]*GetPatternOutput, 0, len(allPatterns))
	for _, domainPattern := range allPatterns {
		patternOutput := &GetPatternOutput{
			PatternID:    domainPattern.PatternID,
			UserID:       domainPattern.UserID,
			Name:         domainPattern.Name,
			TargetWeight: domainPattern.TargetWeight,
			RegisteredAt: domainPattern.RegisteredAt,
			EditedAt:     domainPattern.EditedAt,
			Steps:        stepsByPattern[domainPattern.PatternID],
		}
		result = append(result, patternOutput)
	}

	return result, nil
}

func (pu *patternUsecase) UpdatePattern(ctx context.Context, input UpdatePatternInput) (*UpdatePatternOutput, error) {
	targetPattern, err := pu.patternRepo.FindPatternByPatternID(ctx, input.PatternID, input.UserID)
	if err != nil {
		return nil, err
	}
	targetPatternSteps, err := pu.patternRepo.GetAllPatternStepsByPatternID(ctx, input.PatternID, input.UserID)
	if err != nil {
		return nil, err
	}

	// 変更部分の判定
	// pattern
	isPatternChanged := targetPattern.Name != input.Name || targetPattern.TargetWeight != input.TargetWeight

	// steps
	isStepsChanged := false
	if len(targetPatternSteps) != len(input.Steps) {
		isStepsChanged = true
	}
	for i := range targetPatternSteps {
		if targetPatternSteps[i].IntervalDays != input.Steps[i].IntervalDays {
			isStepsChanged = true
			break
		}
	}

	if !isPatternChanged && !isStepsChanged {
		return nil, patternDomain.ErrNoDiff
	}

	if isPatternChanged {
		editedAt := time.Now().UTC()
		err = targetPattern.Set(input.Name, input.TargetWeight, editedAt)
		if err != nil {
			return nil, err
		}
	}

	var newSteps []*patternDomain.PatternStep
	if isStepsChanged {
		newSteps = make([]*patternDomain.PatternStep, len(input.Steps))
		for i, s := range input.Steps {
			stepID := uuid.NewString()
			patternStep, err := patternDomain.NewPatternStep(
				stepID,
				input.UserID,
				input.PatternID,
				s.StepNumber,
				s.IntervalDays,
			)
			if err != nil {
				return nil, err
			}
			newSteps[i] = patternStep
		}
	}

	err = patternDomain.ValidateSteps(newSteps)
	if err != nil {
		return nil, err
	}

	// patternとstepは別テーブルなので同一トランザクションで永続化
	err = pu.transactionManeger.RunInTransaction(ctx, func(ctx context.Context) error {
		// パターンに変更がある場合、パターンを更新
		if isPatternChanged {
			err = pu.patternRepo.UpdatePattern(ctx, targetPattern)
			if err != nil {
				return err
			}
		}

		// ステップに変更がある場合、古いステップを一括削除→新しいステップを一括挿入
		if isStepsChanged {
			err = pu.patternRepo.DeletePatternSteps(ctx, input.PatternID, input.UserID)
			if err != nil {
				return err
			}

			// "_"←はCopyfromの返り値の、「挿入された行数」。使わないのでブランク識別子にする。
			if _, err := pu.patternRepo.CreatePatternSteps(ctx, newSteps); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resPattern := &UpdatePatternOutput{
		PatternID:    targetPattern.PatternID,
		UserID:       targetPattern.UserID,
		Name:         targetPattern.Name,
		TargetWeight: targetPattern.TargetWeight,
		RegisteredAt: targetPattern.RegisteredAt,
		EditedAt:     targetPattern.EditedAt,
	}
	resPattern.Steps = make([]UpdatePatternStepOutput, len(newSteps))
	for i, s := range newSteps {
		resPattern.Steps[i] = UpdatePatternStepOutput{
			PatternStepID: s.PatternStepID,
			UserID:        s.UserID,
			PatternID:     s.PatternID,
			StepNumber:    s.StepNumber,
			IntervalDays:  s.IntervalDays,
		}
	}

	// 変更点がpattern側だけの場合、レスポンスするのはpatternだけ(stepsは空)
	return resPattern, nil
}

func (pu *patternUsecase) DeletePattern(ctx context.Context, patternID, userID string) error {
	// パターンに紐づく復習物が存在するか確認。
	// TODO: どの復習物に紐づいているのかを返す機能もあってもいいかも
	var isItemRelated bool
	isItemRelated, err := pu.itemRepo.IsPatternRelatedToItemByPatternID(ctx, patternID, userID)
	if err != nil {
		return err
	}
	if isItemRelated {
		return patternDomain.ErrPatternRelatedToItem
	}
	// パターンを削除
	err = pu.patternRepo.DeletePattern(ctx, patternID, userID)
	if err != nil {
		return err
	}
	return nil
}
