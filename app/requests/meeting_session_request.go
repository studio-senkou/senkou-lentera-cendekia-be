package requests

type CreateMeetingSessionRequest struct {
	StudentID   int    `json:"student_id" validate:"required"`
	MentorID    int    `json:"mentor_id" validate:"required"`
	Date        string `json:"date" validate:"required"`
	Time        string `json:"time" validate:"required"`
	Topic       string `json:"topic" validate:"required,min=3,max=255"`
	Duration    int    `json:"duration" validate:"required,min=1"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description" validate:"omitempty,min=3"`
}

type UpdateMeetingSessionRequest struct {
	Date        string  `json:"date" validate:"omitempty"`
	Time        string  `json:"time" validate:"omitempty"`
	Topic       string  `json:"topic" validate:"omitempty,min=3,max=255"`
	Duration    int     `json:"duration" validate:"omitempty,min=1"`
	Type        string  `json:"type" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty,min=3"`
}
