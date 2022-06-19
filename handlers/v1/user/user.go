package user

import (
	"math"

	"github.com/gofiber/fiber/v2"

	models_user "github.com/ktylus/stock-game/common/models/mongo/user"

	models_portfolio "github.com/ktylus/stock-game/common/models/resources/portfolio"

	"github.com/ktylus/stock-game/common/security/authorization"

	security_user "github.com/ktylus/stock-game/common/security/user"
	asset_presenter "github.com/ktylus/stock-game/services/asset/presenter"
)

func UserView(c *fiber.Ctx) error {
	db_user := security_user.GetUser(c)
	presenter := asset_presenter.NewAssetPresenter(db_user)

	var resp_user models_user.UserInResponse

	resp_user = resp_user.FromUserInDB(db_user)
	balance := math.Floor(db_user.Balance*100) / 100
	total_value := math.Floor(presenter.GetOverallValue()*100) / 100
	assets_details := presenter.GetAssetDetails()

	resp := models_portfolio.PortfolioResponse{
		User:        resp_user,
		Assets:      assets_details,
		Total_value: total_value,
	}

	return c.Render("views/portfolio", fiber.Map{
		"authenticated": authorization.CookieAuthorize(c), "balance": balance, "user": db_user, "data": resp, "errors": c.Locals("errors")})
}
