package requests

type CreateMeetingSessionRequest struct {
	StudentID   uint    `json:"student_id" validate:"required"`
	MentorID    uint    `json:"mentor_id" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Time        string  `json:"time" validate:"required"`
	Duration    uint    `json:"duration" validate:"required,min=1"`
	Note        *string `json:"note" validate:"omitempty,min=3"`
	Description string  `json:"description" validate:"required,min=3"`
}

type BulkCreateMeetingSessionRequest struct {
	Sessions []CreateMeetingSessionRequest `json:"sessions" validate:"required,dive"`
}

type UpdateSessionRequest struct {
	SessionID   uint    `json:"session_id" validate:"required"`
	StudentID   uint    `json:"student_id" validate:"required"`
	MentorID    uint    `json:"mentor_id" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Time        string  `json:"time" validate:"required"`
	Duration    uint    `json:"duration" validate:"required,min=1"`
	Note        *string `json:"note" validate:"omitempty,min=3"`
	Status      string  `json:"status" validate:"required,oneof=pending confirmed completed cancelled"`
	Description string  `json:"description" validate:"min=3"`
}

type UpdateMeetingSessionRequest struct {
	Sessions    []UpdateSessionRequest `json:"sessions" validate:"required,dive"`
}

// type MentorAttendanceRequest struct {
// 	SessionFeedback *string `json:"session_feedback" validate:"omitempty,min=3"`
// }

// type MeetingSessionFilters struct {
// 	Date string `json:"date" validate:"omitempty"`
// }
