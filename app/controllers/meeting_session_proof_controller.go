package controllers

import "github.com/gofiber/fiber/v2"

type MeetingSessionProofController struct {
}

func NewMeetingSessionProofController() *MeetingSessionProofController {
	return &MeetingSessionProofController{}
}

func (mc *MeetingSessionProofController) UpdateSessionProof(c *fiber.Ctx) error {
	return nil
}
