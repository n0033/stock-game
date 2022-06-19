package asset

import (
	models_asset "github.com/ktylus/stock-game/common/models/mongo/asset"
)

type AssetTransactionRequest struct {
	Code   string  `json:"code"`
	Amount float64 `json:"amount"`
}

type AssetTransactionResponse struct {
	Messages []string                      `json:"messages"`
	Asset    *models_asset.AssetInResponse `json:"asset"`
	Max_buy  *float64                      `json:"max_buy"`
	Max_sell *float64                      `json:"max_sell"`
}

type AssetDetailedResponse struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Amount      float64 `json:"amount"`
	Total_value float64 `bson:"total_value"`
}
