package client

type NotFound struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *NotFound) Errors() string {
	return e.Message
}

func NewNotFound(message string) *NotFound {
	return &NotFound{
		Code:    404,
		Message: message,
	}
}
