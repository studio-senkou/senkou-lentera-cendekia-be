package requests

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	// Password string `json:"password" validate:"required,min=6,max=30"`
}

type UserActivationRequest struct {
	ActivationToken string `json:"activation_token" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
