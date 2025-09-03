package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/datetime"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type MeetingSessionController struct {
	meetingSessionRepo *models.MeetingSessionRepository
	studentPlanRepo    *models.StudentPlanRepository
}

func NewMeetingSessionController() *MeetingSessionController {
	db := database.GetDB()

	return &MeetingSessionController{
		meetingSessionRepo: models.NewMeetingSessionRepository(db),
		studentPlanRepo:    models.NewStudentPlanRepository(db),
	}
}

func (mc *MeetingSessionController) CreateMeetingSession(c *fiber.Ctx) error {
	createMeetingSessionRequest := new(requests.CreateMeetingSessionRequest)

	if validationError, err := validator.ValidateRequest(c, createMeetingSessionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationError,
		})
	}

	sessionDate, err := datetime.ParseDateOnly(createMeetingSessionRequest.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid date format",
			"error":   "Date format must be YYYY-MM-DD",
		})
	}

	sessionTime, err := datetime.ParseTimeOnly(createMeetingSessionRequest.Time)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid time format",
			"error":   "Time must be in format HH:MM:SS or HH:MM",
		})
	}

	meetingSession := &models.MeetingSession{
		StudentID:   createMeetingSessionRequest.StudentID,
		MentorID:    createMeetingSessionRequest.MentorID,
		Date:        sessionDate,
		Time:        sessionTime,
		Duration:    createMeetingSessionRequest.Duration,
		Note:        createMeetingSessionRequest.Note,
		Description: createMeetingSessionRequest.Description,
		Status:      "pending",
	}

	if _, err := mc.meetingSessionRepo.Create(meetingSession); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create meeting session",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session created successfully",
		"data": fiber.Map{
			"id":           meetingSession.ID,
			"session_date": meetingSession.Date,
			"session_time": meetingSession.Time,
			"duration":     meetingSession.Duration,
			"status":       meetingSession.Status,
			"note":         meetingSession.Note,
			"description":  meetingSession.Description,
		},
	})
}

func (mc *MeetingSessionController) BulkCreateMeetingSessions(c *fiber.Ctx) error {
	bulkCreateMeetingSessions := new(requests.BulkCreateMeetingSessionRequest)
	
	if validationError, err := validator.ValidateRequest(c, bulkCreateMeetingSessions); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationError,
		})
	}

	sessions := make([]*models.MeetingSession, len(bulkCreateMeetingSessions.Sessions))

	for i, session := range bulkCreateMeetingSessions.Sessions {
		sessionDate, err := datetime.ParseDateOnly(session.Date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid date format",
				"error":   "Date format must be YYYY-MM-DD",
			})
		}

		sessionTime, err := datetime.ParseTimeOnly(session.Time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid time format",
				"error":   "Time must be in format HH:MM:SS or HH:MM",
			})
		}
		
		sessions[i] = &models.MeetingSession{
			StudentID:   session.StudentID,
			MentorID:    session.MentorID,
			Date:        sessionDate,
			Time:        sessionTime,
			Duration:    session.Duration,
			Note:        session.Note,
			Description: session.Description,
			Status:      "pending",
		}
	}

	if err := mc.meetingSessionRepo.BulkCreateSessions(sessions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail",
			"message": "Failed to bulk create meeting sessions",
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"message": "Successfully created meeting sessions",
	})
}

func (mc *MeetingSessionController) GetMeetingSessions(c *fiber.Ctx) error {

	user, err := strconv.Atoi(c.Query("user", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}

	meetingSessions, err := mc.meetingSessionRepo.GetAll(uint(user))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to retrieve meeting sessions",
			"error":   err.Error(),
		})
	}

	data := make(fiber.Map)
	sessions := make([]map[string]any, 0)
	for _, session := range meetingSessions {
		record := fiber.Map{
			"id":         session.ID,
			"student_id": session.StudentID,
			"student": fiber.Map{
				"id":    session.Student.ID,
				"name":  session.Student.User.Name,
				"email": session.Student.User.Email,
			},
			"mentor": fiber.Map{
				"id":    session.Mentor.ID,
				"name":  session.Mentor.User.Name,
				"email": session.Mentor.User.Email,
			},
			"mentor_id":    session.MentorID,
			"session_date": session.Date,
			"session_time": session.Time,
			"duration":     session.Duration,
			"status":       session.Status,
			"note":         session.Note,
			"description":  session.Description,
		}

		sessions = append(sessions, record)
	}

	data["sessions"] = sessions

	studentPlan, err := mc.studentPlanRepo.GetCurrentStudentPlan(uint(user))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve student plan",
			"error":   err.Error(),
		})
	}

	if user != 0 && studentPlan != nil {
		data["total_sessions"] = studentPlan.TotalSessions
		data["student_id"] = studentPlan.StudentID
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting sessions retrieved successfully",
		"data":    data,
	})
}

func (mc *MeetingSessionController) GetMeetingSessionByID(c *fiber.Ctx) error {
	sessionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid session ID",
			"error":   err.Error(),
		})
	}

	meetingSession, err := mc.meetingSessionRepo.GetByID(uint(sessionID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Meeting session not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session retrieved successfully",
		"data": fiber.Map{
			"id":         meetingSession.ID,
			"student_id": meetingSession.StudentID,
			"mentor_id":  meetingSession.MentorID,
			"student": fiber.Map{
				"id":    meetingSession.Student.ID,
				"name":  meetingSession.Student.User.Name,
				"email": meetingSession.Student.User.Email,
			},
			"mentor": fiber.Map{
				"id":    meetingSession.Mentor.ID,
				"name":  meetingSession.Mentor.User.Name,
				"email": meetingSession.Mentor.User.Email,
			},
			"session_date": meetingSession.Date,
			"session_time": meetingSession.Time,
			"duration":     meetingSession.Duration,
			"status":       meetingSession.Status,
			"note":         meetingSession.Note,
			"description":  meetingSession.Description,
		},
	})
}

func (mc *MeetingSessionController) UpdateMeetingSession(c *fiber.Ctx) error {
	updateMeetingSessionRequest := new(requests.UpdateMeetingSessionRequest)

	if validationError, err := validator.ValidateRequest(c, updateMeetingSessionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationError,
		})
	}

	sessions := make([]*models.MeetingSession, len(updateMeetingSessionRequest.Sessions))

	for i, sessionReq := range updateMeetingSessionRequest.Sessions {
		sessionDate, err := datetime.ParseDateOnly(sessionReq.Date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid session date",
				"error":   err.Error(),
			})
		}

		sessionTime, err := datetime.ParseTimeOnly(sessionReq.Time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid session time",
				"error":   err.Error(),
			})
		}

		sessions[i] = &models.MeetingSession{
			ID:          sessionReq.SessionID,
			StudentID:   sessionReq.StudentID,
			MentorID:    sessionReq.MentorID,
			Date:        sessionDate,
			Time:        sessionTime,
			Duration:    sessionReq.Duration,
			Note:        sessionReq.Note,
			Description: sessionReq.Description,
			Status:      sessionReq.Status,
		}
	}

	if err := mc.meetingSessionRepo.BulkUpdate(sessions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update meeting session",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session updated successfully",
	})
}

func (mc *MeetingSessionController) DeleteMeetingSession(c *fiber.Ctx) error {
	sessionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid session ID",
			"error":   err.Error(),
		})
	}

	if err := mc.meetingSessionRepo.Delete(uint(sessionID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete meeting session",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session deleted successfully",
	})
}
