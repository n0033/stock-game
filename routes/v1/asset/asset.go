package asset

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktylus/stock-game/handlers/v1/asset"
)

func ApplyRoutes(router fiber.Router) {
	router.Post("/buy", asset.Buy)
	router.Post("/sell", asset.Sell)
}
