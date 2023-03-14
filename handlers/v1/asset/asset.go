package asset

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	models_asset "github.com/n0033/stock-game/common/models/mongo/asset"
	models_asset_view "github.com/n0033/stock-game/common/models/resources/asset"
	"github.com/n0033/stock-game/services/asset/manager"

	security_user "github.com/n0033/stock-game/common/security/user"
)

func Buy(c *fiber.Ctx) error {
	var request models_asset_view.AssetTransactionRequest
	var response models_asset_view.AssetTransactionResponse

	messages := make([]string, 0)

	if err := c.BodyParser(&request); err != nil {
		return err
	}
	db_user := security_user.GetUser(c)
	asset_manager := manager.NewAssetManager(db_user)
	db_asset, err := asset_manager.BuyAsset(request.Code, request.Amount)

	if err != nil {
		messages = append(messages, err.Error())
		response = models_asset_view.AssetTransactionResponse{
			Messages: messages,
			Asset:    nil,
			Max_buy:  nil,
			Max_sell: nil,
		}
	}

	current_value := asset_manager.GetCurrentValueByCode(request.Code)
	db_user = security_user.GetUser(c)

	var max_buy float64 = db_user.Balance / current_value.GetValue()
	var max_sell float64 = db_asset.Amount

	if err == nil {
		var asset_resp models_asset.AssetInResponse
		asset_resp = asset_resp.FromAssetInDB(db_asset)
		messages = append(messages, "Asset was successfuly bought.")
		response = models_asset_view.AssetTransactionResponse{
			Messages: messages,
			Asset:    &asset_resp,
			Max_buy:  &max_buy,
			Max_sell: &max_sell,
		}
	}

	var resp_json map[string]interface{}
	resp_bytes, _ := json.Marshal(response)
	json.Unmarshal(resp_bytes, &resp_json)
	return c.JSON(resp_json)
}

func Sell(c *fiber.Ctx) error {
	var request models_asset_view.AssetTransactionRequest
	var response models_asset_view.AssetTransactionResponse

	messages := make([]string, 0)

	if err := c.BodyParser(&request); err != nil {
		return err
	}
	db_user := security_user.GetUser(c)
	asset_manager := manager.NewAssetManager(db_user)
	db_asset, err := asset_manager.SellAsset(request.Code, request.Amount)

	if err != nil {
		messages = append(messages, err.Error())
		response = models_asset_view.AssetTransactionResponse{
			Messages: messages,
			Asset:    nil,
			Max_buy:  nil,
			Max_sell: nil,
		}
	}

	db_user = security_user.GetUser(c)
	current_value := asset_manager.GetCurrentValueByCode(request.Code)

	var max_buy float64 = db_user.Balance / current_value.GetValue()
	var max_sell float64 = db_asset.Amount

	if err == nil {
		var asset_resp models_asset.AssetInResponse
		asset_resp = asset_resp.FromAssetInDB(db_asset)
		messages = append(messages, "Asset was successfuly sold.")
		response = models_asset_view.AssetTransactionResponse{
			Messages: messages,
			Asset:    &asset_resp,
			Max_buy:  &max_buy,
			Max_sell: &max_sell,
		}
	}

	var resp_json map[string]interface{}
	resp_bytes, _ := json.Marshal(response)
	json.Unmarshal(resp_bytes, &resp_json)
	return c.JSON(resp_json)
}
