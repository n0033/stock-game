package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktylus/stock-game/common/security/authorization"
	asset_route "github.com/ktylus/stock-game/routes/v1/asset"
	auth_route "github.com/ktylus/stock-game/routes/v1/authorization"
	details_route "github.com/ktylus/stock-game/routes/v1/details"
	home_route "github.com/ktylus/stock-game/routes/v1/home"
	search_route "github.com/ktylus/stock-game/routes/v1/search"
	user_route "github.com/ktylus/stock-game/routes/v1/user"
)

func ApplyRoutes(app *fiber.App) {
	var asset fiber.Router = app.Group("/asset")
	var auth fiber.Router = app.Group("/auth")
	var details fiber.Router = app.Group("/details", authorization.New(authorization.Config{}))
	var search fiber.Router = app.Group("/search", authorization.New(authorization.Config{}))
	var user fiber.Router = app.Group("/user", authorization.New(authorization.Config{}))
	asset_route.ApplyRoutes(asset)
	auth_route.ApplyRoutes(auth)
	details_route.ApplyRoutes(details)
	home_route.ApplyRoutes(app)
	search_route.ApplyRoutes(search)
	user_route.ApplyRoutes(user)
}
