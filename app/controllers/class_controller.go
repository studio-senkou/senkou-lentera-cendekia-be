package controllers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (cc *ClassController) GetClassDropdown(c *fiber.Ctx) error {
	classes, err := cc.classRepo.FindAllForDropdown()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve classes as dropdown",
			"error":   err.Error(),
		})
	}

	dropdowns := []fiber.Map{}
	for _, class := range classes {
		dropdowns = append(dropdowns, fiber.Map{
			"id":   class.ID.String(),
			"name": class.ClassName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully retrieve all classes as dropdown",
		"data": dropdowns,
	})
}

func (cc *ClassController) UpdateClass(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid class ID",
			"error":   err.Error(),
		})
	}

	updateClassRequest := new(requests.UpdateClassRequest)
	if validationError, err := validator.ValidateRequest(c, updateClassRequest); err != nil {
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

	class, err := cc.classRepo.Update(classID, updateClassRequest.ClassName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Class not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update class",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully update class",
		"data": fiber.Map{
			"class": class,
		},
	})
}

func (cc *ClassController) DeleteClass(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid class ID",
			"error":   err.Error(),
		})
	}

	if err := cc.classRepo.Delete(classID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Class not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete class",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully delete class",
	})
}
