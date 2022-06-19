package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktylus/stock-game/handlers/v1/search"
)

func ApplyRoutes(router fiber.Router) {
	router.Post("/", search.CompanySearch)
}
