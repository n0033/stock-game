package manager

import (
	"fmt"
	"log"
	"sort"
	"time"

	models_asset "github.com/n0033/stock-game/common/models/mongo/asset"
	models_stock "github.com/n0033/stock-game/common/models/mongo/stock_datapoint"
	models_user "github.com/n0033/stock-game/common/models/mongo/user"
	models_av "github.com/n0033/stock-game/common/models/resources/alpha_vantage"
	"github.com/n0033/stock-game/config"
	adapter "github.com/n0033/stock-game/services/adapter/alpha_vantage"
	alpha_vantage_connector "github.com/n0033/stock-game/services/api_connector/connectors/alpha_vantage"

	asset_dao "github.com/n0033/stock-game/common/dao/asset"
	company_dao "github.com/n0033/stock-game/common/dao/company"
	crypto_dao "github.com/n0033/stock-game/common/dao/cryptocurrency"
	stock_dao "github.com/n0033/stock-game/common/dao/stock"
	user_dao "github.com/n0033/stock-game/common/dao/user"
)

type AssetManager struct {
	user models_user.UserInDB
}

func NewAssetManager(user models_user.UserInDB) *AssetManager {
	return &AssetManager{user: user}
}

func (manager *AssetManager) GetCurrentValueByCode(code string) *models_stock.StockDatapointInDB {
	dao_stock := stock_dao.NewDAOStock()

	datapoint, err := dao_stock.FindLatestByCode(code)

	if err != nil || (time.Since(datapoint.Timestamp).Hours() > 1 && time.Since(datapoint.Last_used).Hours() > 1) {
		dao_crypto := crypto_dao.NewDAOCryptocurrency()
		dao_company := company_dao.NewDAOCompany()
		connector := alpha_vantage_connector.GetAlphaVantageConnector()
		is_crypto := false
		var parsed []models_av.StockDatapoint

		_, err := dao_crypto.FindByCode(code)
		if err == nil {
			is_crypto = true
		}
		_, err = dao_company.FindByCode(code)
		if !is_crypto && err != nil {
			log.Fatal(err)
		}

		if is_crypto {
			resp := connector.GetCurrentCryptoData(code, "USD")
			parsed = adapter.CryptoToStockDatapoints(resp)
		}

		if !is_crypto {
			resp := connector.GetCurrentCompanyData(code, config.AV_RESPONSE_SIZE["default"])
			parsed = adapter.ToStockDatapoints(resp)
		}

		for _, val := range parsed {
			var datapoint_create models_stock.StockDatapointInCreate
			datapoint_create = datapoint_create.FromAVStockDatapoint(val)
			if datapoint.Timestamp.Before(datapoint_create.Timestamp) {
				dao_stock.Create(datapoint_create)
			}
		}
		datapoint, err = dao_stock.FindLatestByCode(code)
		if err != nil {
			return nil
		}
	}

	return &datapoint
}

func (manager *AssetManager) GetLatestValuesByCode(code string) []models_stock.StockDatapointInDB {
	dao_stock := stock_dao.NewDAOStock()

	datapoints := dao_stock.Find300LatestByCode(code)

	last_datapoint := datapoints[len(datapoints)-1]

	if len(datapoints) == 0 || (time.Since(last_datapoint.Timestamp).Hours() > 1 && time.Since(last_datapoint.Last_used).Hours() > 1) {
		connector := alpha_vantage_connector.GetAlphaVantageConnector()
		dao_crypto := crypto_dao.NewDAOCryptocurrency()
		dao_company := company_dao.NewDAOCompany()
		is_crypto := false
		var parsed []models_av.StockDatapoint

		_, err := dao_crypto.FindByCode(code)
		if err == nil {
			is_crypto = true
		}
		_, err = dao_company.FindByCode(code)
		if !is_crypto && err != nil {
			log.Fatal(err)
		}

		if is_crypto {
			response := connector.GetCurrentCryptoData(code, "USD")
			parsed = adapter.CryptoToStockDatapoints(response)
		}

		if !is_crypto {
			response := connector.GetCurrentCompanyData(code, config.AV_RESPONSE_SIZE["default"])
			parsed = adapter.ToStockDatapoints(response)
		}

		sort.Slice(parsed, func(i, j int) bool {
			return parsed[i].Timestamp.Before(parsed[j].Timestamp)
		})

		for _, val := range parsed {
			var datapoint_create models_stock.StockDatapointInCreate
			datapoint_create = datapoint_create.FromAVStockDatapoint(val)
			if len(datapoints) == 0 || datapoints[len(datapoints)-1].Timestamp.Before(datapoint_create.Timestamp) {
				db_stock, err := dao_stock.Create(datapoint_create)
				if err == nil {
					datapoints = append(datapoints, db_stock)
				}
			}
		}
	}
	return datapoints
}

func (manager *AssetManager) BuyAsset(code string, amount float64) (models_asset.AssetInDB, error) {
	var db_asset models_asset.AssetInDB
	dao_asset := asset_dao.NewDAOAsset()
	dao_user := user_dao.NewDAOUser()

	latest_datapoint := manager.GetCurrentValueByCode(code)
	current_value := latest_datapoint.GetValue()
	asset_value := current_value * amount

	if manager.user.Balance < asset_value {
		return db_asset, fmt.Errorf("not enough balance to purchase this amount of '%s' assets", code)
	}

	existing, err := dao_asset.FindByUserAndCode(manager.user.ID, code)

	if err == nil {
		var asset_update models_asset.AssetInUpdate
		existing.Amount += amount
		asset_update = asset_update.FromAssetInDB(existing)
		db_asset, err = dao_asset.Update(existing.ID, asset_update)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		asset_create := models_asset.AssetInCreate{
			AssetBase: models_asset.AssetBase{
				User_id: manager.user.ID,
				Code:    code,
				Amount:  amount,
			},
		}
		db_asset, err = dao_asset.Create(asset_create)
		if err != nil {
			return db_asset, err
		}
	}

	var user_update models_user.UserInUpdate
	db_user := manager.user
	db_user.Balance -= asset_value
	user_update = user_update.FromUserInDB(db_user)
	db_user, err = dao_user.Update(db_user.ID, user_update)
	if err != nil {
		return db_asset, err
	}
	manager.user = db_user

	return db_asset, err
}

func (manager *AssetManager) SellAsset(code string, amount float64) (models_asset.AssetInDB, error) {
	var db_asset models_asset.AssetInDB
	dao_asset := asset_dao.NewDAOAsset()
	dao_user := user_dao.NewDAOUser()

	db_asset, err := dao_asset.FindByUserAndCode(manager.user.ID, code)

	if err != nil || db_asset.Amount < amount {
		return db_asset, fmt.Errorf("user with id %s does not have enough assets of '%s' company", manager.user.ID, code)
	}

	latest_datapoint := manager.GetCurrentValueByCode(code)
	current_value := latest_datapoint.GetValue()
	asset_value := current_value * amount

	var asset_update models_asset.AssetInUpdate
	db_asset.Amount -= amount
	asset_update = asset_update.FromAssetInDB(db_asset)
	db_asset, err = dao_asset.Update(db_asset.ID, asset_update)

	if err != nil {
		return db_asset, err
	}

	var user_update models_user.UserInUpdate
	db_user := manager.user
	db_user.Balance += asset_value
	user_update = user_update.FromUserInDB(db_user)
	db_user, err = dao_user.Update(db_user.ID, user_update)
	if err != nil {
		return db_asset, err
	}
	manager.user = db_user

	return db_asset, err
}
