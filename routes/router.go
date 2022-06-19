package routes

import (
	"github.com/gofiber/fiber/v2"
	v1_routes "github.com/ktylus/stock-game/routes/v1"
)

func ApplyRoutes(app *fiber.App) {
	v1_routes.ApplyRoutes(app)
}
