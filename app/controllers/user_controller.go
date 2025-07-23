package controllers

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
	gomail "github.com/studio-senkou/lentera-cendekia-be/utils/mail"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type UserController struct {
	userRepo *models.UserRepository
}

func NewUserController() *UserController {
	db := database.GetDB()
	return &UserController{
		userRepo: models.NewUserRepository(db),
	}
}

func (uc *UserController) CreateMentor(c *fiber.Ctx) error {
	return uc.CreateUser(c, "mentor")
}

func (uc *UserController) CreateStudent(c *fiber.Ctx) error {
	return uc.CreateUser(c, "student")
}

func (uc *UserController) CreateUser(c *fiber.Ctx, role string) error {
	createUserRequest := new(requests.CreateUserRequest)

	if validationError, err := validator.ValidateRequest(c, createUserRequest); err != nil {
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
		Name:  createUserRequest.Name,
		Email: createUserRequest.Email,
		// Password: createUserRequest.Password,
		Password: "12345678",
		Role:     role,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	activationToken, err := auth.GenerateOneTimeToken(user.ID, "account_activation", 24*time.Hour)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate activation token",
		})
	}

	email, err := gomail.NewMailFromTemplate(
		// user.Email,
		// "ajhmdni02@gmail.com",
		"studio.senkou@gmail.com",
		"Welcome aboard to Lentera Cendekia",
		"templates/emails/welcome.html",
		fiber.Map{
			"Name":           user.Name,
			"ActivationLink": "https://portal.lenteracendekia.id/activate?token=" + activationToken.Token,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate email from template",
		})
	}

	email.Send()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Email verification sent successfully, please check your user email inbox",
		"data": fiber.Map{
			"user": user,
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

	// Get user
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

	if err := uc.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to activate user",
			"error":   "Failed to update user activation status: " + err.Error(),
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := uc.userRepo.GetByID(id)
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

func (uc *UserController) GetUserAsDropdown(c *fiber.Ctx) error {
	users, err := uc.userRepo.GetUserDropdown()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	dropdownUsers := make([]fiber.Map, len(users))
	for i, user := range users {
		dropdownUsers[i] = fiber.Map{
			"id":   user.ID,
			"name": user.Name,
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
	for i, user := range users {
		dropdownUsers[i] = fiber.Map{
			"id":   user.ID,
			"name": user.Name,
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

	user := new(models.User)
	user.ID = id
	user.Name = updateUserRequest.Name
	user.Email = updateUserRequest.Email

	if err := uc.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"data": fiber.Map{
			"user": user,
		},
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
