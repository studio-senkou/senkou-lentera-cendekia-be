package client

type BadRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BadRequest) Errors() string {
	return e.Message
}

func NewBadRequest(message string) *BadRequest {
	return &BadRequest{
		Code:    400,
		Message: message,
	}
}
