package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	v1_router "github.com/ktylus/stock-game/routes/v1"
)

func initiateApp() *fiber.App {
	tpl_engine := html.New("./templates", ".html")
	var app *fiber.App = fiber.New(fiber.Config{
		Views: tpl_engine,
	})
	app.Static("/static", "./static")
	return app
}

func main() {
	var app *fiber.App = initiateApp()
	v1_router.ApplyRoutes(app)
	app.Listen(":3000")
}
