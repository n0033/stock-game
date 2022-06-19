package authorization

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktylus/stock-game/handlers/v1/authorization"
)

func ApplyRoutes(router fiber.Router) {
	router.Get("/login", authorization.LoginView)
	router.Post("/login", authorization.Login)

	router.Get("/register", authorization.RegisterView)
	router.Post("/register", authorization.Register)
}
