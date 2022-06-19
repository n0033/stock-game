package authorization

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktylus/stock-game/handlers/v1/details"
)

func ApplyRoutes(router fiber.Router) {
	router.Get("/:code", details.Details)
	router.Get("/:code/data", details.GetCompanyStock)
}
