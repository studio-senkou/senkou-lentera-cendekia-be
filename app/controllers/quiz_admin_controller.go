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

type QuizAdminController struct {
	adminRepo *models.QuizAdminRepository
	quizRepo  *models.QuizRepository
}

func NewQuizAdminController() *QuizAdminController {
	db := database.GetDB()
	return &QuizAdminController{
		adminRepo: models.NewQuizAdminRepository(db),
		quizRepo:  models.NewQuizRepository(db),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /admin/quizzes
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) ListQuizzes(c *fiber.Ctx) error {
	quizzes, err := ac.adminRepo.ListQuizzes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve quizzes",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quizzes retrieved successfully",
		"data":    quizzes,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /admin/quizzes/:id
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) GetQuiz(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	quiz, questions, err := ac.adminRepo.GetQuizDetail(quizID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve quiz",
			"error":   err.Error(),
		})
	}
	if quiz == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Quiz not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz retrieved successfully",
		"data": fiber.Map{
			"quiz":      quiz,
			"questions": buildAdminQuestionsResponse(questions),
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /admin/quizzes
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) CreateQuiz(c *fiber.Ctx) error {
	req := new(requests.CreateQuizRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	quiz := &models.QuizQuiz{
		Title:            req.Title,
		Description:      req.Description,
		PassingScore:     req.PassingScore,
		TimeLimitMinutes: req.TimeLimitMinutes,
		IsActive:         req.IsActive,
	}

	if err := ac.adminRepo.CreateQuiz(quiz); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create quiz",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz created successfully",
		"data":    quiz,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// PUT /admin/quizzes/:id
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) UpdateQuiz(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	req := new(requests.UpdateQuizRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	quiz := &models.QuizQuiz{
		ID:               quizID,
		Title:            req.Title,
		Description:      req.Description,
		PassingScore:     req.PassingScore,
		TimeLimitMinutes: req.TimeLimitMinutes,
		IsActive:         req.IsActive,
	}

	if err := ac.adminRepo.UpdateQuiz(quiz); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Quiz not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update quiz",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz updated successfully",
		"data":    quiz,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// DELETE /admin/quizzes/:id
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) DeleteQuiz(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	if err := ac.adminRepo.DeleteQuiz(quizID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Quiz not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete quiz",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz deleted successfully",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /admin/quizzes/:id/questions
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) CreateQuestion(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	req := new(requests.CreateQuestionRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	question := &models.QuizQuestion{
		QuizID:       quizID,
		QuestionText: req.QuestionText,
	}

	if err := ac.adminRepo.CreateQuestion(question); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create question",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Question created successfully",
		"data":    question,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// PUT /admin/quizzes/:id/questions/:qid
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) UpdateQuestion(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	questionID, err := parseID(c, "qid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid question ID",
		})
	}

	req := new(requests.UpdateQuestionRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	question := &models.QuizQuestion{
		ID:           questionID,
		QuizID:       quizID,
		QuestionText: req.QuestionText,
	}

	if err := ac.adminRepo.UpdateQuestion(question); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Question not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update question",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Question updated successfully",
		"data":    question,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// DELETE /admin/quizzes/:id/questions/:qid
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) DeleteQuestion(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	questionID, err := parseID(c, "qid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid question ID",
		})
	}

	if err := ac.adminRepo.DeleteQuestion(questionID, quizID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Question not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete question",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Question deleted successfully",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /admin/quizzes/:id/questions/:qid/options
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) CreateOption(c *fiber.Ctx) error {
	_, err := parseID(c, "id") // pastikan quiz param ada
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	questionID, err := parseID(c, "qid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid question ID",
		})
	}

	req := new(requests.CreateOptionRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	option := &models.QuizOption{
		QuestionID:  questionID,
		OptionText:  req.OptionText,
		IsCorrect:   req.IsCorrect,
	}

	if err := ac.adminRepo.CreateOption(option); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create option",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Option created successfully",
		"data": fiber.Map{
			"id":           option.ID,
			"question_id":  option.QuestionID,
			"option_text":  option.OptionText,
			"is_correct":   option.IsCorrect,
			"created_at":   option.CreatedAt,
			"updated_at":   option.UpdatedAt,
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// PUT /admin/quizzes/:id/questions/:qid/options/:oid
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) UpdateOption(c *fiber.Ctx) error {
	_, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	questionID, err := parseID(c, "qid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid question ID",
		})
	}
	optionID, err := parseID(c, "oid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid option ID",
		})
	}

	req := new(requests.UpdateOptionRequest)
	if ve, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(ve) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  ve,
		})
	}

	option := &models.QuizOption{
		ID:          optionID,
		QuestionID:  questionID,
		OptionText:  req.OptionText,
		IsCorrect:   req.IsCorrect,
	}

	if err := ac.adminRepo.UpdateOption(option); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Option not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update option",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Option updated successfully",
		"data": fiber.Map{
			"id":           option.ID,
			"question_id":  option.QuestionID,
			"option_text":  option.OptionText,
			"is_correct":   option.IsCorrect,
			"updated_at":   option.UpdatedAt,
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// DELETE /admin/quizzes/:id/questions/:qid/options/:oid
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) DeleteOption(c *fiber.Ctx) error {
	_, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	questionID, err := parseID(c, "qid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid question ID",
		})
	}
	optionID, err := parseID(c, "oid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid option ID",
		})
	}

	if err := ac.adminRepo.DeleteOption(optionID, questionID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "Option not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete option",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Option deleted successfully",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /admin/quizzes/:id/attempts
// ─────────────────────────────────────────────────────────────────────────────

func (ac *QuizAdminController) ListAttempts(c *fiber.Ctx) error {
	quizID, err := parseID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	attempts, err := ac.adminRepo.ListAttempts(quizID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve attempts",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Attempts retrieved successfully",
		"data":    attempts,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// parseID mengambil path param berdasarkan nama dan mengembalikannya sebagai uint.
func parseID(c *fiber.Ctx, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Params(param), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// buildAdminQuestionsResponse membangun response soal untuk admin (termasuk is_correct).
func buildAdminQuestionsResponse(questions []models.QuizQuestion) []fiber.Map {
	result := make([]fiber.Map, len(questions))
	for i, q := range questions {
		options := make([]fiber.Map, len(q.Options))
		for j, o := range q.Options {
			options[j] = fiber.Map{
				"id":           o.ID,
				"option_text":  o.OptionText,
			"is_correct":   o.IsCorrect,
		}
	}
	result[i] = fiber.Map{
		"id":            q.ID,
		"question_text": q.QuestionText,
		"options":       options,
	}
}
return result
}
