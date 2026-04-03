package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupQuizRoutes(router fiber.Router) {
	quizController := controllers.NewQuizController()

	router.Get(
		"/quiz/histories",
		middlewares.AuthMiddleware(),
		quizController.GetStudentQuizHistories,
	)

	router.Get(
		"/quiz/code/:code",
		middlewares.AuthMiddleware(),
		quizController.GetQuizByCode,
	)

	router.Get(
		"/quiz/:id",
		middlewares.AuthMiddleware(),
		quizController.GetQuiz,
	)

	router.Post(
		"/quiz/:id/submit",
		middlewares.AuthMiddleware(),
		quizController.SubmitQuiz,
	)

	router.Get(
		"/quiz/:id/status",
		middlewares.AuthMiddleware(),
		quizController.GetQuizStatus,
	)

	router.Get(
		"/quiz/:id/questions/current",
		middlewares.AuthMiddleware(),
		quizController.GetCurrentQuestion,
	)

	router.Post(
		"/quiz/:id/questions/next",
		middlewares.AuthMiddleware(),
		quizController.NextQuestion,
	)

	router.Post(
		"/quiz/:id/questions/prev",
		middlewares.AuthMiddleware(),
		quizController.PrevQuestion,
	)

	router.Post(
		"/quiz/:id/reset",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("admin"),
		quizController.ResetAttempt,
	)
}
