package box

type IBoxUsecase interface {
	CreateBox(box CreateBoxInput) (*CreateBoxOutput, error)
	GetBoxesByCategoryID(categoryID string, userID string) ([]*GetBoxOutput, error)
	UpdateBox(box UpdateBoxInput) (*UpdateBoxOutput, error)
	DeleteBox(boxID string, categoryID string, userID string) error
}
