package home

import "github.com/gofiber/fiber/v2"

func Homepage(c *fiber.Ctx) error {
	return c.Render("index/index", fiber.Map{})
}
