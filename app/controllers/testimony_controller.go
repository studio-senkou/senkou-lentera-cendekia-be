package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/storage"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type TestimonyController struct {
	testimonyRepository *models.TestimonyRepository
}

func NewTestimonyController() *TestimonyController {
	db := database.GetDB()
	testimonyRepository := models.NewTestimonyRepository(db)

	return &TestimonyController{testimonyRepository: testimonyRepository}
}

func (tc *TestimonyController) CreateTestimony(c *fiber.Ctx) error {
	createTestimonyRequest := new(requests.CreateTestimonyRequest)
	if validatorErrors, err := validator.ValidateFormData(c, createTestimonyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validatorErrors,
		})
	} else if len(validatorErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validatorErrors,
		})
	}

	testimonerPhoto, err := c.FormFile("testimoner_photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to upload asset",
			"error":   err.Error(),
		})
	}

	var uploadedTestimonerPhoto *string
	if testimonerPhoto != nil {
		maxSize := int64(1 * 1024 * 1024)
		if testimonerPhoto.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Photo size exceeds the limit of 1MB",
			})
		}

		testimonerPhotoPath, err := storage.UploadFileToStorage(testimonerPhoto, "testimoners", "TESTIMONER", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Failed to upload testimoner photo",
				"error":   err.Error(),
			})
		}

		uploadedTestimonerPhoto = &testimonerPhotoPath
	}

	newTestimony := &models.Testimony{
		TestimonerName:             createTestimonyRequest.TestimonerName,
		TestimonerCurrentPosition:  createTestimonyRequest.TestimonerCurrentPosition,
		TestimonerPreviousPosition: createTestimonyRequest.TestimonerPreviousPosition,
		TestimonerPhoto:            uploadedTestimonerPhoto,
		TestimonyText:              createTestimonyRequest.TestimonyText,
	}

	if err := tc.testimonyRepository.Create(newTestimony); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create testimony",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Testimony created successfully",
	})
}

func (tc *TestimonyController) GetAllTestimonials(c *fiber.Ctx) error {
	testimonies, err := tc.testimonyRepository.GetAllTestimonials()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve testimonies",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Testimonials retrieved successfully",
		"data": fiber.Map{
			"testimonials": testimonies,
		},
	})
}

func (tc *TestimonyController) GetTestimonyByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid testimony ID",
			"error":   err.Error(),
		})
	}

	testimony, err := tc.testimonyRepository.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve testimony",
			"error":   err.Error(),
		})
	}

	if testimony == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Testimony not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Testimony retrieved successfully",
		"data": fiber.Map{
			"testimony": testimony,
		},
	})
}

func (tc *TestimonyController) UpdateTestimony(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid testimony ID",
			"error":   err.Error(),
		})
	}

	updateTestimonyRequest := new(requests.UpdateTestimonyRequest)
	if validatorErrors, err := validator.ValidateFormData(c, updateTestimonyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validatorErrors,
		})
	} else if len(validatorErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validatorErrors,
		})
	}

	testimony, err := tc.testimonyRepository.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve testimony",
			"error":   err.Error(),
		})
	}

	if testimony == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Testimony not found",
		})
	}

	testimonerPhoto, err := c.FormFile("testimoner_photo")
	if err != nil && testimonerPhoto != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to upload asset",
		})
	}

	var uploadedTestimonerPhoto *string
	if testimonerPhoto != nil {
		maxSize := int64(1 * 1024 * 1024)
		if testimonerPhoto.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Photo size exceeds the limit of 1MB",
			})
		}
		
		if testimony.TestimonerPhoto != nil {
			if err := storage.RemoveFileFromStorage(*testimony.TestimonerPhoto); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "fail",
					"message": "Failed to remove old testimoner photo",
					"error":   err.Error(),
				})
			}
		}

		testimonerPhotoPath, err := storage.UploadFileToStorage(testimonerPhoto, "testimoners", "TESTIMONER", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Failed to upload testimoner photo",
				"error":   err.Error(),
			})
		}

		uploadedTestimonerPhoto = &testimonerPhotoPath
	}

	testimony.TestimonerName = updateTestimonyRequest.TestimonerName
	testimony.TestimonerCurrentPosition = updateTestimonyRequest.TestimonerCurrentPosition
	testimony.TestimonerPreviousPosition = updateTestimonyRequest.TestimonerPreviousPosition
	testimony.TestimonyText = updateTestimonyRequest.TestimonyText

	if uploadedTestimonerPhoto != nil {
		testimony.TestimonerPhoto = uploadedTestimonerPhoto
	}

	if err := tc.testimonyRepository.Update(testimony); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update testimony",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Testimony updated successfully",
	})
}

func (tc *TestimonyController) DeleteTestimony(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid testimony ID",
			"error":   err.Error(),
		})
	}

	testimony, err := tc.testimonyRepository.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve testimony",
			"error":   err.Error(),
		})
	}

	if testimony == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Testimony not found",
		})
	}

	if err := tc.testimonyRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete testimony",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Testimony deleted successfully",
	})
}
