package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyAccountRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
