package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
