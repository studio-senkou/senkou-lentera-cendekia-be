package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupQuizAdminRoutes(router fiber.Router) {
	ac := controllers.NewQuizAdminController()

	// Semua endpoint admin quiz wajib login dan role admin
	admin := router.Group("/admin/quizzes",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("admin"),
	)

	// ── Quiz CRUD ──────────────────────────────────────────────────────────────
	admin.Get("", ac.ListQuizzes)
	admin.Post("", ac.CreateQuiz)
	admin.Get("/:id", ac.GetQuiz)
	admin.Put("/:id", ac.UpdateQuiz)
	admin.Delete("/:id", ac.DeleteQuiz)

	// ── Question CRUD ──────────────────────────────────────────────────────────
	admin.Post("/:id/questions", ac.CreateQuestion)
	admin.Put("/:id/questions/:qid", ac.UpdateQuestion)
	admin.Delete("/:id/questions/:qid", ac.DeleteQuestion)

	// ── Option CRUD ────────────────────────────────────────────────────────────
	admin.Post("/:id/questions/:qid/options", ac.CreateOption)
	admin.Put("/:id/questions/:qid/options/:oid", ac.UpdateOption)
	admin.Delete("/:id/questions/:qid/options/:oid", ac.DeleteOption)

	// ── Attempt Management ─────────────────────────────────────────────────────
	// GET  /admin/quizzes/:id/attempts       → lihat semua attempt
	// POST /quiz/:id/reset                   → sudah ada di quiz_routes.go
	admin.Get("/:id/attempts", ac.ListAttempts)
}
