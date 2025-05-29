package box

type CreateBoxRequest struct {
	PatternID string `json:"pattern_id"`
	Name      string `json:"name"`
}

type UpdateBoxRequest struct {
	PatternID string `json:"pattern_id"`
	Name      string `json:"name"`
}
