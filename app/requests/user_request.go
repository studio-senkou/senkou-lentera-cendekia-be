package requests

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	// Password string `json:"password" validate:"required,min=6,max=30"`
}

type CreateNewStudentRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Class           string `json:"class" validate:"required"`
	MinimalSessions uint   `json:"minimal_sessions" validate:"required,min=1"`
}

type CreateNewMentorRequest struct {
	Name    string   `json:"name" validate:"required"`
	Email   string   `json:"email" validate:"required,email"`
	Classes []string `json:"classes" validate:"required,dive"`
}

type UserActivationRequest struct {
	ActivationToken string `json:"activation_token" validate:"required"`
	Password        string `json:"password" validate:"required,min=6,max=30"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirmRequest struct {
	ResetToken      string `json:"reset_token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=30,eqfield=NewPassword"`
}

type UpdatePasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=30,eqfield=NewPassword"`
}

type UpdateUserPasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required,min=6,max=30"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=30,eqfield=NewPassword"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
