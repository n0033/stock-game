package details

import (
	"time"

	models_av "github.com/n0033/stock-game/common/models/resources/alpha_vantage"

	models_user "github.com/n0033/stock-game/common/models/mongo/user"
)

type AssetDetailsResponse struct {
	Symbol  string                     `json:"symbol"`
	User    models_user.UserInResponse `json:"user"`
	Company *models_av.CompanyOverview `json:"company"`
	Price   float64                    `json:"price"`
}

type CompanyStockDataResponse struct {
	Timestamps []time.Time `json:"timestamps"`
	Values     []float64   `json:"values"`
}

type CompanyStockDatapoint struct {
	Timestamp time.Time `json:"date"`
	Value     float64   `json:"value"`
}
