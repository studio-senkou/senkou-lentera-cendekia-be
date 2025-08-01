package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupMeetingSessionRoutes(router fiber.Router) {
	meetingSessionController := controllers.NewMeetingSessionController()

	router.Post("/meeting-sessions", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), meetingSessionController.CreateMeetingSession)
	router.Post("/meeting-sessions/:id/student-attend", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("user"), meetingSessionController.UserAttend)
	router.Post("/meeting-sessions/:id/mentor-attend", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("mentor"), meetingSessionController.MentorAttend)
	router.Get("/meeting-sessions", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), meetingSessionController.GetMeetingSession)
	router.Get("/meeting-sessions/me", middlewares.AuthMiddleware(), meetingSessionController.GetUserMeetingSession)
	router.Get("/meeting-sessions/:id", middlewares.AuthMiddleware(), meetingSessionController.GetMeetingSessionByID)
	router.Put("/meeting-sessions/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), meetingSessionController.UpdateMeetingSession)
	router.Delete("/meeting-sessions/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), meetingSessionController.DeleteMeetingSession)
	router.Patch("/meeting-sessions/:id/:status", middlewares.AuthMiddleware(), meetingSessionController.UpdateMeetingSessionStatus)
}
