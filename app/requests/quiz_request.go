package requests

type SubmitQuizRequest struct {
	Answers []SubmitAnswerItem `json:"answers" validate:"required,min=1,dive"`
}

type SubmitAnswerItem struct {
	QuestionID uint `json:"question_id" validate:"required"`
	OptionID   uint `json:"option_id"   validate:"required"`
}

type ResetQuizAttemptRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}
