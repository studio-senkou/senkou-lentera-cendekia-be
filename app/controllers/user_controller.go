package controllers

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
	gomail "github.com/studio-senkou/lentera-cendekia-be/utils/mail"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type UserController struct {
	jwtManager      *auth.JwtManager
	userRepo        *models.UserRepository
	studentRepo     *models.StudentRepository
	studentPlanRepo *models.StudentPlanRepository
	mentorRepo      *models.MentorRepository
	authRepo        *models.AuthenticationRepository
}

func NewUserController() *UserController {
	db := database.GetDB()
	authSecret := app.GetEnv("AUTH_SECRET", "")

	return &UserController{
		jwtManager:      auth.NewJwtManager(authSecret),
		userRepo:        models.NewUserRepository(db),
		studentRepo:     models.NewStudentRepository(db),
		studentPlanRepo: models.NewStudentPlanRepository(db),
		mentorRepo:      models.NewMentorRepository(db),
		authRepo:        models.NewAuthenticationRepository(db),
	}
}

func (uc *UserController) CreateNewStudent(c *fiber.Ctx) error {
	createNewStudentRequest := new(requests.CreateNewStudentRequest)

	if validationError, err := validator.ValidateRequest(c, createNewStudentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	user := &models.User{
		Name:     createNewStudentRequest.Name,
		Email:    createNewStudentRequest.Email,
		Password: "12345678",
		Role:     "user",
	}

	if err := uc.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	classID := uuid.MustParse(createNewStudentRequest.Class)
	student, err := uc.studentRepo.AddIntoClass(user.ID, classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add user into class",
		})
	}

	studentPlan := &models.StudentPlan{
		StudentID:     student.ID,
		TotalSessions: createNewStudentRequest.MinimalSessions,
	}

	if err := uc.studentPlanRepo.CreateNewStudentPlan(studentPlan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create student plan",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Student registered successfully",
		"data": fiber.Map{
			"student_class_id": classID,
			"user":             user,
		},
	})
}

func (uc *UserController) CreateNewMentor(c *fiber.Ctx) error {
	createNewMentorRequest := new(requests.CreateNewMentorRequest)

	if validationError, err := validator.ValidateRequest(c, createNewMentorRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	user := &models.User{
		Name:     createNewMentorRequest.Name,
		Email:    createNewMentorRequest.Email,
		Password: "12345678",
		Role:     "mentor",
	}

	if err := uc.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	classID := uuid.MustParse(createNewMentorRequest.Class)
	if _, err := uc.mentorRepo.AddIntoClass(user.ID, classID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add user into class",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Mentor registered successfully",
		"data": fiber.Map{
			"mentor_class_id": classID,
			"user":            user,
		},
	})
}

func (uc *UserController) ActivateUser(c *fiber.Ctx) error {
	activationRequest := new(requests.UserActivationRequest)
	if validationError, err := validator.ValidateRequest(c, activationRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	oneTimeToken, err := auth.ValidateOneTimeToken(activationRequest.ActivationToken, "account_activation")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid activation token",
			"error":   err.Error(),
		})
	}

	user, err := uc.userRepo.GetByID(oneTimeToken.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Unable to retrieve user",
			"error":   "Failed to retrieve user: " + err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot activate user",
			"error":   "User not found",
		})
	}

	if user.IsEmailVerified() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User already activated",
			"error":   "User email is already verified",
		})
	}

	user.MarkEmailAsVerified()
	user.IsActive = true

	if _, err := uc.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to activate user",
			"error":   "Failed to update user activation status: " + err.Error(),
		})
	}

	if err := uc.userRepo.UpdatePassword(user.ID, activationRequest.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to set user password",
			"error":   "Failed to set user password: " + err.Error(),
		})
	}

	// Authenticate the user after activation
	accessToken, err := uc.jwtManager.GenerateToken(auth.Payload{
		UserID: user.ID,
		Role:   user.Role,
	}, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate access token",
			"error":   "Failed to generate access token: " + err.Error(),
		})
	}

	refreshToken, err := uc.jwtManager.GenerateToken(auth.Payload{
		UserID: user.ID,
		Role:   user.Role,
	}, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate refresh token",
			"error":   "Failed to generate refresh token: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User activated successfully",
		"data": fiber.Map{
			"active_role":          user.Role,
			"access_token":         accessToken.Token,
			"access_token_expiry":  accessToken.ExpiresAt,
			"refresh_token":        refreshToken.Token,
			"refresh_token_expiry": refreshToken.ExpiresAt,
		},
	})
}

func (uc *UserController) ForceActivateUser(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := uc.userRepo.GetByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.IsEmailVerified() && user.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User already activated",
			"error":   "User email is already verified and active",
		})
	}

	user.MarkEmailAsVerified()
	user.MarkAsActive()

	if _, err := uc.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User activated successfully",
	})
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved users",
		"data": fiber.Map{
			"users": users,
		},
	})
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := uc.userRepo.GetByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved user",
		"data": fiber.Map{
			"user": user,
		},
	})
}

func (uc *UserController) GetUserMe(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}

	user, err := uc.userRepo.GetByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user",
			"error":   "Failed to retrieve user: " + err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved user",
		"data": fiber.Map{
			"user": user,
		},
	})
}

func (uc *UserController) GetActiveUser(c *fiber.Ctx) error {
	users, err := uc.userRepo.GetUserCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user count",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved user count",
		"data": fiber.Map{
			"users": users,
		},
	})
}

func (uc *UserController) GetUserAsDropdown(c *fiber.Ctx) error {
	users, err := uc.userRepo.GetUserDropdown()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	dropdownUsers := make([]fiber.Map, len(users))
	for i, student := range users {
		dropdownUsers[i] = fiber.Map{
			"id":   student.ID,
			"name": student.User.Name,
			"email": student.User.Email,
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved users for dropdown",
		"data": fiber.Map{
			"users": dropdownUsers,
		},
	})
}

func (uc *UserController) GetMentorDropdown(c *fiber.Ctx) error {
	users, err := uc.userRepo.GetMentorDropdown()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve mentors",
		})
	}

	dropdownUsers := make([]fiber.Map, len(users))
	for i, mentor := range users {
		dropdownUsers[i] = fiber.Map{
			"id":   mentor.ID,
			"name": mentor.User.Name,
			"email": mentor.User.Email,
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully retrieved mentors for dropdown",
		"data": fiber.Map{
			"mentors": dropdownUsers,
		},
	})
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	updateUserRequest := new(requests.UpdateUserRequest)
	if validationError, err := validator.ValidateRequest(c, updateUserRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	user := &models.User{
		ID:    uint(id),
		Name:  updateUserRequest.Name,
		Email: updateUserRequest.Email,
	}

	updatedEmail, err := uc.userRepo.Update(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	if updatedEmail != "" {
		activationToken, err := auth.GenerateOneTimeToken(user.ID, "account_activation", 24*time.Hour)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate activation token",
			})
		}

		email, err := gomail.NewMailFromTemplate(
			user.Email,
			"Email Change Verification",
			"templates/emails/email_change_verification.html",
			fiber.Map{
				"Name":           user.Name,
				"ActivationLink": fmt.Sprintf("%s/verify?token=%s", app.GetEnv("APP_FE_URL", "http://localhost:3000"), activationToken.Token),
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate email from template",
			})
		}

		email.Send()
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"data": fiber.Map{
			"user": user,
		},
	})
}

func (uc *UserController) ResetPassword(c *fiber.Ctx) error {
	resetPasswordRequest := new(requests.ResetPasswordRequest)
	if validationError, err := validator.ValidateRequest(c, resetPasswordRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	user, err := uc.userRepo.GetByEmail(resetPasswordRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user",
			"error":   "Failed to retrieve user: " + err.Error(),
		})
	} else if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   "User not found",
		})
	}

	if !user.IsEmailVerified() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Email not verified",
			"error":   "User email is not verified",
		})
	}

	resetToken, err := auth.GenerateOneTimeToken(user.ID, "password_reset", 15*time.Minute)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate reset token",
			"error":   "Failed to generate reset token: " + err.Error(),
		})
	}

	email, err := gomail.NewMailFromTemplate(
		user.Email,
		fmt.Sprintf("Reset Password for %s", user.Name),
		"templates/emails/reset_password.html",
		fiber.Map{
			"Name":       user.Name,
			"ResetLink":  fmt.Sprintf("%s/reset-password?token=%s", app.GetEnv("APP_FE_URL", "http://localhost:3000"), resetToken.Token),
			"ExpiryTime": resetToken.ExpiresAt.Format(time.RFC1123),
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate email from template",
			"error":   "Failed to generate email from template: " + err.Error(),
		})
	}
	email.Send()

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Password reset email sent successfully",
	})
}

func (uc *UserController) UpdatePasswordByToken(c *fiber.Ctx) error {
	updatePasswordRequest := new(requests.UpdatePasswordRequest)
	if validationError, err := validator.ValidateRequest(c, updatePasswordRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	oneTimeToken, err := auth.ValidateOneTimeToken(updatePasswordRequest.Token, "password_reset")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid password reset token",
			"error":   err.Error(),
		})
	}

	if err := uc.userRepo.UpdatePassword(oneTimeToken.UserID, updatePasswordRequest.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to update password",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

func (uc *UserController) UpdateUserPassword(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}

	updatePasswordRequest := new(requests.UpdateUserPasswordRequest)
	if validationError, err := validator.ValidateRequest(c, updatePasswordRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationError,
		})
	}

	if isMatch, err := uc.userRepo.VerifyOldPassword(int(userID), updatePasswordRequest.OldPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to verify old password",
			"error":   err.Error(),
		})
	} else if !isMatch {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Old password is incorrect",
			"error":   "Old password does not match",
		})
	}

	if err := uc.userRepo.UpdatePassword(uint(userID), updatePasswordRequest.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to update password",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := uc.userRepo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
