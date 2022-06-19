package portfolio

import (
	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
	models_asset "github.com/ktylus/stock-game/common/models/resources/asset"
)

type PortfolioResponse struct {
	User        models_user.UserInResponse
	Assets      []models_asset.AssetDetailedResponse
	Total_value float64
}
