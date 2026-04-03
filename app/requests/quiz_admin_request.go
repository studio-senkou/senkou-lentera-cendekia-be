package requests

type CreateQuizRequest struct {
	Title            string  `json:"title"             validate:"required,min=3,max=255"`
	Description      *string `json:"description"       validate:"omitempty,min=3"`
	PassingScore     int     `json:"passing_score"     validate:"required,min=0,max=100"`
	TimeLimitMinutes *int    `json:"time_limit_minutes" validate:"omitempty,min=1"`
	IsActive         bool    `json:"is_active"`
}

type UpdateQuizRequest struct {
	Title            string  `json:"title"             validate:"required,min=3,max=255"`
	Description      *string `json:"description"       validate:"omitempty,min=3"`
	PassingScore     int     `json:"passing_score"     validate:"required,min=0,max=100"`
	TimeLimitMinutes *int    `json:"time_limit_minutes" validate:"omitempty,min=1"`
	IsActive         bool    `json:"is_active"`
}

type CreateQuestionRequest struct {
	QuestionText string `json:"question_text" validate:"required,min=3"`
}

type UpdateQuestionRequest struct {
	QuestionText string `json:"question_text" validate:"required,min=3"`
}

type CreateOptionRequest struct {
	OptionText string `json:"option_text" validate:"required,min=1"`
	IsCorrect  bool   `json:"is_correct"`
}

type UpdateOptionRequest struct {
	OptionText string `json:"option_text" validate:"required,min=1"`
	IsCorrect  bool   `json:"is_correct"`
}
