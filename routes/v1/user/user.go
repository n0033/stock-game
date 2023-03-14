package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/n0033/stock-game/handlers/v1/user"
)

func ApplyRoutes(router fiber.Router) {
	router.Get("", user.UserView)
}
