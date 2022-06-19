package details

import (
	"encoding/json"
	"math"

	"github.com/gofiber/fiber/v2"

	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
	models_av "github.com/ktylus/stock-game/common/models/resources/alpha_vantage"
	models_details "github.com/ktylus/stock-game/common/models/resources/details"

	company_dao "github.com/ktylus/stock-game/common/dao/company"

	"github.com/ktylus/stock-game/common/security/authorization"
	security_user "github.com/ktylus/stock-game/common/security/user"
	user_handlers "github.com/ktylus/stock-game/handlers/v1/user"
	adapter "github.com/ktylus/stock-game/services/adapter/alpha_vantage"
	alpha_vantage_connector "github.com/ktylus/stock-game/services/api_connector/connectors/alpha_vantage"
	asset_manager "github.com/ktylus/stock-game/services/asset/manager"
	asset_presenter "github.com/ktylus/stock-game/services/asset/presenter"
)

func Details(c *fiber.Ctx) error {
	db_user := security_user.GetUser(c)
	av_connector := alpha_vantage_connector.GetAlphaVantageConnector()
	dao_company := company_dao.NewDAOCompany()
	manager := asset_manager.NewAssetManager(db_user)
	presenter := asset_presenter.NewAssetPresenter(db_user)
	errors := make([]string, 0)

	code := c.Params("code")
	var user_resp models_user.UserInResponse
	current_value := manager.GetCurrentValueByCode(code)

	if current_value == nil {
		errors = append(errors, "Data for requested company or cryptocurrency is not available.")
		c.Locals("errors", errors)
		return user_handlers.UserView(c)
	}

	balance := math.Floor(db_user.Balance*100) / 100
	user_asset, err := presenter.GetAssetByCode(code)

	var max_buy float64 = db_user.Balance / current_value.GetValue()
	var max_sell float64
	if err != nil {
		max_sell = 0
	}
	if err == nil {
		max_sell = user_asset.Amount
	}

	var company_overview *models_av.CompanyOverview = nil
	_, err = dao_company.FindByCode(code)

	if err == nil {
		av_response := av_connector.GetCompanyOverview(code)
		parsed := adapter.ToCompanyOverview(av_response)
		company_overview = &parsed
		if len(av_response.Payload) == 0 {
			company_overview = nil
		}
	}

	details := models_details.AssetDetailsResponse{
		Symbol:  code,
		User:    user_resp.FromUserInDB(db_user),
		Company: company_overview,
		Price:   current_value.GetValue(),
	}

	return c.Render("views/company_view", fiber.Map{"authenticated": authorization.CookieAuthorize(c), "balance": balance, "data": details, "max_buy": max_buy, "max_sell": max_sell})
}

func GetCompanyStock(c *fiber.Ctx) error {
	db_user := security_user.GetUser(c)
	manager := asset_manager.NewAssetManager(db_user)

	code := c.Params("code")

	datapoints := manager.GetLatestValuesByCode(code)

	var resp_datapoints []models_details.CompanyStockDatapoint

	for _, datapoint := range datapoints {
		resp_datapoints = append(resp_datapoints, models_details.CompanyStockDatapoint{
			Timestamp: datapoint.Timestamp,
			Value:     datapoint.GetValue(),
		})
	}

	var response []map[string]interface{}

	response_bytes, _ := json.Marshal(resp_datapoints)
	json.Unmarshal(response_bytes, &response)
	return c.JSON(response)
}
