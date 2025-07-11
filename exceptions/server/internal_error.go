package server

type InternalError struct {
	Code   int    `json:"code"`
	Message string `json:"message"`
}

func (e *InternalError) Errors() string {
	return e.Message
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		Code:    500,
		Message: message,
	}
}