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

type QuizController struct {
	quizRepo *models.QuizRepository
}

func NewQuizController() *QuizController {
	return &QuizController{
		quizRepo: models.NewQuizRepository(database.GetDB()),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /quiz/:id
// Mengambil soal kuis. User yang sudah 'completed' tetap bisa melihat kuis.
// User yang 'in_progress' mendapatkan soal untuk dilanjutkan.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) GetQuiz(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	userID := uint(c.Locals("userID").(int))

	// 1. Cek attempt user terlebih dahulu untuk mendapatkan urutan soal (acak)
	attempt, err := qc.quizRepo.GetActiveAttempt(userID, uint(quizID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check attempt status",
			"error":   err.Error(),
		})
	}

	// 2. Jika belum ada attempt atau sudah di-reset → buat attempt baru (lazy creation + shuffle)
	if attempt == nil {
		newAttempt := &models.QuizAttempt{
			QuizID: uint(quizID),
			UserID: userID,
		}
		if err := qc.quizRepo.CreateAttempt(newAttempt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to start quiz attempt",
				"error":   err.Error(),
			})
		}
		attempt = newAttempt
	}

	// 3. Ambil kuis dan soal berdasarkan urutan yang tersimpan di attempt (acak per student)
	// 3. Ambil kuis dan soal berdasarkan urutan yang tersimpan di attempt (acak per student)
	quiz, questions, err := qc.quizRepo.GetActiveQuizWithQuestionsV2(uint(quizID), attempt)
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
			"message": "Quiz not found or inactive",
		})
	}

	// Jika user sudah 'completed', kembalikan info hasil saja — soal tidak perlu dikirim ulang
	if attempt.Status == "completed" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "You have already completed this quiz",
			"data": fiber.Map{
				"quiz": fiber.Map{
					"id":    quiz.ID,
					"title": quiz.Title,
				},
				"attempt": fiber.Map{
					"id":           attempt.ID,
					"status":       attempt.Status,
					"score":        attempt.Score,
					"submitted_at": attempt.SubmittedAt,
				},
				"questions": nil,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz retrieved successfully",
		"data": fiber.Map{
			"quiz": fiber.Map{
				"id":                  quiz.ID,
				"title":               quiz.Title,
				"description":         quiz.Description,
				"time_limit_minutes":  quiz.TimeLimitMinutes,
				"passing_score":       quiz.PassingScore,
			},
			"attempt": fiber.Map{
				"id":         attempt.ID,
				"status":     attempt.Status,
				"started_at": attempt.StartedAt,
			},
			"questions": buildQuestionsResponse(questions),
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /quiz/:id/submit
// Submit jawaban kuis. Hanya bisa jika status attempt = 'in_progress'.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) SubmitQuiz(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	userID := uint(c.Locals("userID").(int))

	req := new(requests.SubmitQuizRequest)
	if validationErrors, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationErrors,
		})
	}

	// Pastikan ada attempt aktif yang in_progress
	attempt, err := qc.quizRepo.GetActiveAttempt(userID, uint(quizID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check attempt",
			"error":   err.Error(),
		})
	}
	if attempt == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "No active attempt found. Please start the quiz first.",
		})
	}
	if attempt.Status == "completed" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "Quiz already submitted. You cannot re-submit unless reset by admin.",
		})
	}

	// Map request ke model
	answers := make([]models.QuizAnswer, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = models.QuizAnswer{
			AttemptID:  attempt.ID,
			QuestionID: a.QuestionID,
			OptionID:   a.OptionID,
		}
	}

	completedAttempt, err := qc.quizRepo.SubmitAnswers(attempt.ID, answers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to submit quiz",
			"error":   err.Error(),
		})
	}

	// Tentukan apakah user lulus
	quiz, _, err := qc.quizRepo.GetActiveQuizWithQuestionsV2(uint(quizID), attempt)
	if err != nil || quiz == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve quiz info",
		})
	}

	passed := completedAttempt.Score != nil && *completedAttempt.Score >= float64(quiz.PassingScore)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Quiz submitted successfully",
		"data": fiber.Map{
			"attempt_id":   completedAttempt.ID,
			"score":        completedAttempt.Score,
			"passing_score": quiz.PassingScore,
			"passed":       passed,
			"submitted_at": completedAttempt.SubmittedAt,
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /quiz/:id/status
// Cek status attempt user untuk kuis tertentu.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) GetQuizStatus(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	userID := uint(c.Locals("userID").(int))

	attempt, err := qc.quizRepo.GetActiveAttempt(userID, uint(quizID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check attempt status",
			"error":   err.Error(),
		})
	}
	if attempt == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "No attempt found",
			"data": fiber.Map{
				"has_attempt": false,
				"attempt":     nil,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Attempt status retrieved",
		"data": fiber.Map{
			"has_attempt": true,
			"attempt": fiber.Map{
				"id":           attempt.ID,
				"status":       attempt.Status,
				"score":        attempt.Score,
				"started_at":   attempt.StartedAt,
				"submitted_at": attempt.SubmittedAt,
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /quiz/:id/reset   [admin only]
// Mereset attempt user sehingga user bisa mengerjakan ulang.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) ResetAttempt(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}

	adminID := uint(c.Locals("userID").(int))

	req := new(requests.ResetQuizAttemptRequest)
	if validationErrors, err := validator.ValidateRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationErrors,
		})
	}

	err = qc.quizRepo.ResetAttempt(req.UserID, uint(quizID), adminID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "No active attempt found for this user on this quiz",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to reset attempt",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Attempt has been reset. The user can now retake the quiz.",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper: build soal tanpa field is_correct agar aman dikirim ke client
// ─────────────────────────────────────────────────────────────────────────────

func buildQuestionsResponse(questions []models.QuizQuestion) []fiber.Map {
	result := make([]fiber.Map, len(questions))
	for i, q := range questions {
		options := make([]fiber.Map, len(q.Options))
		for j, o := range q.Options {
		options[j] = fiber.Map{
			"id":          o.ID,
			"option_text": o.OptionText,
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

// ─────────────────────────────────────────────────────────────────────────────
// GET /quiz/:id/questions/current
// Mengambil soal yang sedang aktif berdasarkan current_question_index.
// Jika belum ada attempt, buat attempt baru dan return soal pertama.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) GetCurrentQuestion(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid quiz ID",
		})
	}
	userID := uint(c.Locals("userID").(int))

	attempt, err := qc.quizRepo.GetActiveAttempt(userID, uint(quizID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to check attempt",
		})
	}

	// Buat attempt baru jika belum ada
	if attempt == nil {
		attempt = &models.QuizAttempt{QuizID: uint(quizID), UserID: userID}
		if err := qc.quizRepo.CreateAttempt(attempt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error", "message": "Failed to start quiz attempt",
			})
		}
	}

	if attempt.Status == "completed" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "Quiz already completed",
		})
	}

	return qc.returnQuestionAtIndex(c, attempt, attempt.CurrentQuestionIndex)
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /quiz/:id/questions/next
// Maju ke soal berikutnya. Jika sudah di soal terakhir, kembalikan error.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) NextQuestion(c *fiber.Ctx) error {
	return qc.navigateQuestion(c, 1)
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /quiz/:id/questions/prev
// Mundur ke soal sebelumnya. Jika sudah di soal pertama, kembalikan error.
// ─────────────────────────────────────────────────────────────────────────────

func (qc *QuizController) PrevQuestion(c *fiber.Ctx) error {
	return qc.navigateQuestion(c, -1)
}

// navigateQuestion adalah helper bersama untuk next (+1) dan prev (-1)
func (qc *QuizController) navigateQuestion(c *fiber.Ctx, direction int) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": "Invalid quiz ID",
		})
	}
	userID := uint(c.Locals("userID").(int))

	attempt, err := qc.quizRepo.GetActiveAttempt(userID, uint(quizID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to get attempt",
		})
	}
	if attempt == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "No active attempt. Please start the quiz first.",
		})
	}
	if attempt.Status == "completed" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "fail", "message": "Quiz already completed",
		})
	}

	newIndex := attempt.CurrentQuestionIndex + direction
	total := len(attempt.QuestionIDs)

	if newIndex < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Already at the first question",
		})
	}
	if newIndex >= total {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Already at the last question",
		})
	}

	// Persist indeks baru ke database
	if err := qc.quizRepo.UpdateAttemptIndex(attempt.ID, newIndex); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to update question index",
		})
	}
	attempt.CurrentQuestionIndex = newIndex

	return qc.returnQuestionAtIndex(c, attempt, newIndex)
}

// returnQuestionAtIndex adalah helper untuk mengambil & mengembalikan soal berdasarkan indeks
func (qc *QuizController) returnQuestionAtIndex(c *fiber.Ctx, attempt *models.QuizAttempt, index int) error {
	total := len(attempt.QuestionIDs)
	question, err := qc.quizRepo.GetQuestionByAttemptIndex(attempt, index)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to retrieve question",
		})
	}
	if question == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "fail", "message": "Question not found",
		})
	}

	options := make([]fiber.Map, len(question.Options))
	for i, o := range question.Options {
		options[i] = fiber.Map{
			"id":          o.ID,
			"option_text": o.OptionText,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Question retrieved successfully",
		"data": fiber.Map{
			"navigation": fiber.Map{
				"current_index":  index,
				"total_questions": total,
				"has_prev":        index > 0,
				"has_next":        index < total-1,
			},
			"question": fiber.Map{
				"id":            question.ID,
				"question_text": question.QuestionText,
				"options":       options,
			},
		},
	})
}
