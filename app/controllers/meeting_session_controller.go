package controllers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/storage"
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

func (mc *MeetingSessionController) GetUserMeetingSession(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}

	sessions, err := mc.meetingSessionRepo.GetByUser(int(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve meeting sessions for user",
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

func (mc *MeetingSessionController) UserAttend(c *fiber.Ctx) error {
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}

	if err := mc.meetingSessionRepo.VerifyAttendance(meetingID, int(userID), false); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "You have already attended this session or session does not exist",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot verify attendance",
			"error":   "User has already attended this session or session does not exist",
		})
	}

	sessionProof, err := c.FormFile("session_proof")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if sessionProof == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Session proof image is required",
		})
	}

	sessionAttendanceProof, err := c.FormFile("session_attendance_proof")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if sessionAttendanceProof == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Session attendance proof image is required",
		})
	}

	if !storage.IsValidImageExtension(sessionProof.Filename) || !storage.IsValidImageExtension(sessionAttendanceProof.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid session proof image format",
		})
	}

	maxSize := int64(0.5 * 1024 * 1024)
	if sessionProof.Size > maxSize || sessionAttendanceProof.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Session proof image size exceeds 500KB",
		})
	}

	uploadedSessionProofPath, err := storage.UploadFileToStorage(sessionProof, "meeting_sessions", "MEET-SESSION", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to upload session proof image",
			"error":   err.Error(),
		})
	}

	uploadedSessionAttendanceProofPath, err := storage.UploadFileToStorage(sessionAttendanceProof, "meeting_sessions", "MEET-SESSION", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to upload session attendance proof image",
			"error":   err.Error(),
		})
	}

	if err := mc.meetingSessionRepo.UpdateProofs(
		meetingID,
		&uploadedSessionProofPath,
		&uploadedSessionAttendanceProofPath,
		nil,
		nil,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update meeting session proofs",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Session attendance proof uploaded successfully",
	})
}

func (mc *MeetingSessionController) MentorAttend(c *fiber.Ctx) error {
	meetingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid meeting session ID",
			"error":   "ID must be a number",
		})
	}

	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   "Invalid user ID",
		})
	}

	if err := mc.meetingSessionRepo.VerifyAttendance(meetingID, int(userID), true); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "You have already attended this session or session does not exist",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot verify attendance",
			"error":   "Mentor has already attended this session or session does not exist",
		})
	}

	mentorAttendanceProof, err := c.FormFile("session_attendance_proof")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if mentorAttendanceProof == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Mentor attendance proof image is required",
		})
	}

	if !storage.IsValidImageExtension(mentorAttendanceProof.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid mentor attendance proof image format",
		})
	}

	maxSize := int64(0.5 * 1024 * 1024)
	if mentorAttendanceProof.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Mentor attendance proof image size exceeds 500KB",
		})
	}

	mentorAttendanceRequest := new(requests.MentorAttendanceRequest)
	if validationError, err := validator.ValidateFormData(c, mentorAttendanceRequest); err != nil {
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


	uploadedMentorAttendanceProofPath, err := storage.UploadFileToStorage(mentorAttendanceProof, "meeting_sessions", "MEET-SESSION", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to upload mentor attendance proof image",
			"error":   err.Error(),
		})
	}

	if err := mc.meetingSessionRepo.UpdateProofs(
		meetingID,
		nil,
		nil,
		&uploadedMentorAttendanceProofPath,
		mentorAttendanceRequest.SessionFeedback,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update meeting session mentor attendance proof",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Mentor attendance proof uploaded successfully",
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
