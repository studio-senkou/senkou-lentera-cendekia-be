package requests

type CreateClassRequest struct {
	ClassName string `json:"classname" validate:"required"`
}
