package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type ClassController struct {
	classRepo *models.ClassRepository
}

func NewClassController() *ClassController {
	db := database.GetDB()
	classRepository := models.NewClassRepository(db)

	return &ClassController{
		classRepo: classRepository,
	}
}

func (cc *ClassController) CreateNewClass(c *fiber.Ctx) error {
	createClassRequest := new(requests.CreateClassRequest)

	if validationError, err := validator.ValidateRequest(c, createClassRequest); err != nil {
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

	class, err := cc.classRepo.Store(createClassRequest.ClassName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create new class",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully create new class",
		"data": fiber.Map{
			"class": class,
		},
	})
}

func (cc *ClassController) GetAllClasses(c *fiber.Ctx) error {
	classes, err := cc.classRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve classes",
			"error":   "Unable to get classes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully retrieve all classes",
		"data": fiber.Map{
			"classes": classes,
		},
	})
}
