package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type MeetingSessionController struct {
	meetingSessionRepo *models.MeetingSessionRepository
}

func NewMeetingSessionController() *MeetingSessionController {
	db := database.GetDB()

	return &MeetingSessionController{
		meetingSessionRepo: models.NewMeetingSessionRepository(db),
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

	session := &models.MeetingSession{
		UserID:             createMeetingSessionRequest.StudentID,
		MentorID:           createMeetingSessionRequest.MentorID,
		SessionDate:        createMeetingSessionRequest.Date,
		SessionTime:        createMeetingSessionRequest.Time,
		SessionDuration:    createMeetingSessionRequest.Duration,
		SessionType:        createMeetingSessionRequest.Type,
		SessionTopic:       createMeetingSessionRequest.Topic,
		SessionDescription: &createMeetingSessionRequest.Description,
	}

	if err := mc.meetingSessionRepo.Create(session); err != nil {
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
			"session": session,
		},
	})
}

func (mc *MeetingSessionController) GetMeetingSession(c *fiber.Ctx) error {
	sessions, err := mc.meetingSessionRepo.GetAll()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve meeting sessions",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting sessions retrieved successfully",
		"data": fiber.Map{
			"sessions": sessions,
		},
	})
}

func (mc *MeetingSessionController) GetMeetingSessionByID(c *fiber.Ctx) error {
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

	session, err := mc.meetingSessionRepo.GetByID(meetingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve meeting session",
			"error":   err.Error(),
		})
	}

	if session == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Meeting session not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session retrieved successfully",
		"data": fiber.Map{
			"session": session,
		},
	})
}

func (mc *MeetingSessionController) UpdateMeetingSession(c *fiber.Ctx) error {
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

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

	session := &models.MeetingSession{
		ID:                 meetingID,
		SessionDate:        updateMeetingSessionRequest.Date,
		SessionTime:        updateMeetingSessionRequest.Time,
		SessionDuration:    updateMeetingSessionRequest.Duration,
		SessionTopic:       updateMeetingSessionRequest.Topic,
		SessionType:        updateMeetingSessionRequest.Type,
		SessionDescription: updateMeetingSessionRequest.Description,
	}

	if err := mc.meetingSessionRepo.Update(session); err != nil {
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
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

	if err := mc.meetingSessionRepo.Delete(meetingID); err != nil {
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

func (mc *MeetingSessionController) UpdateMeetingSessionStatus(c *fiber.Ctx) error {
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Status cannot be empty",
		})
	}

	if err := mc.meetingSessionRepo.UpdateStatus(meetingID, status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update meeting session status",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Meeting session status updated successfully",
	})
}
