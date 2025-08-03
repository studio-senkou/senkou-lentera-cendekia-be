package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type BlogController struct {
	blogRepository *models.BlogRepository
}

func NewBlogController() *BlogController {
	db := database.GetDB()
	blogRepository := models.NewBlogRepository(db)

	return &BlogController{blogRepository: blogRepository}
}

func (bc *BlogController) CreateBlog(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}
	
	createBlogRequest := new(requests.CreateBlogRequest)
	if validatorErrors, err := validator.ValidateRequest(c, createBlogRequest); err != nil {
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

	newBlog := &models.Blog{
		Title:    createBlogRequest.Title,
		Content:  createBlogRequest.Content,
		AuthorID: int(userID),
	}

	if err := bc.blogRepository.Create(newBlog); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create blog",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Blog created successfully",
	})
}

func (bc *BlogController) GetAllBlogs(c *fiber.Ctx) error {
	blogs, err := bc.blogRepository.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve blogs",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Blogs retrieved successfully",
		"data":   blogs,
	})
}

func (bc *BlogController) GetBlogByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid blog ID",
			"error":   err.Error(),
		})
	}

	blog, err := bc.blogRepository.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Blog not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   blog,
	})
}

func (bc *BlogController) UpdateBlog(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid blog ID",
			"error":   err.Error(),
		})
	}

	updateBlogRequest := new(requests.UpdateBlogRequest)
	if validatorErrors, err := validator.ValidateRequest(c, updateBlogRequest); err != nil {
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

	updateData := &models.Blog{
		ID:      id,
		Title:   updateBlogRequest.Title,
		Content: updateBlogRequest.Content,
	}

	if err := bc.blogRepository.Update(updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update blog",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Blog updated successfully",
	})
}

func (bc *BlogController) DeleteBlog(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid blog ID",
			"error":   err.Error(),
		})
	}

	if err := bc.blogRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete blog",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Blog deleted successfully",
	})
}
