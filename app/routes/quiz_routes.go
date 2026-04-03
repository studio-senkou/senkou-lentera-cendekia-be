package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupQuizRoutes(router fiber.Router) {
	quizController := controllers.NewQuizController()

	// GET /quiz/:id → ambil soal kuis (semua user terautentikasi)
	router.Get(
		"/quiz/:id",
		middlewares.AuthMiddleware(),
		quizController.GetQuiz,
	)

	// POST /quiz/:id/submit → submit jawaban kuis
	router.Post(
		"/quiz/:id/submit",
		middlewares.AuthMiddleware(),
		quizController.SubmitQuiz,
	)

	// GET /quiz/:id/status → cek status attempt user
	router.Get(
		"/quiz/:id/status",
		middlewares.AuthMiddleware(),
		quizController.GetQuizStatus,
	)

	// GET /quiz/:id/questions/current → ambil soal aktif
	router.Get(
		"/quiz/:id/questions/current",
		middlewares.AuthMiddleware(),
		quizController.GetCurrentQuestion,
	)

	// POST /quiz/:id/questions/next → maju ke soal berikutnya
	router.Post(
		"/quiz/:id/questions/next",
		middlewares.AuthMiddleware(),
		quizController.NextQuestion,
	)

	// POST /quiz/:id/questions/prev → mundur ke soal sebelumnya
	router.Post(
		"/quiz/:id/questions/prev",
		middlewares.AuthMiddleware(),
		quizController.PrevQuestion,
	)

	// POST /quiz/:id/reset → reset attempt user (admin only)
	router.Post(
		"/quiz/:id/reset",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("admin"),
		quizController.ResetAttempt,
	)
}
