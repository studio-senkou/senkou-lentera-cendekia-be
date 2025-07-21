package controllers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type UserController struct {
	userRepo *models.UserRepository
}

func NewUserController() *UserController {
	db := database.GetDB().GetDB()
	return &UserController{
		userRepo: models.NewUserRepository(db),
	}
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
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
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User created successfully",
		"data":    user,
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
