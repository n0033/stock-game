package presenter

import (
	"math"

	asset_dao "github.com/ktylus/stock-game/common/dao/asset"
	company_dao "github.com/ktylus/stock-game/common/dao/company"
	crypto_dao "github.com/ktylus/stock-game/common/dao/cryptocurrency"

	mongo_asset "github.com/ktylus/stock-game/common/models/mongo/asset"
	mongo_user "github.com/ktylus/stock-game/common/models/mongo/user"
	models_asset "github.com/ktylus/stock-game/common/models/resources/asset"

	av_manager "github.com/ktylus/stock-game/services/asset/manager"
)

type AssetPresenter struct {
	user mongo_user.UserInDB
}

func NewAssetPresenter(user mongo_user.UserInDB) *AssetPresenter {
	return &AssetPresenter{user: user}
}

func (presenter *AssetPresenter) Get(code string) (mongo_asset.AssetInDB, error) {
	dao_asset := asset_dao.NewDAOAsset()
	db_asset, err := dao_asset.FindByUserAndCode(presenter.user.ID, code)
	return db_asset, err
}

func (presenter *AssetPresenter) GetAll() []mongo_asset.AssetInDB {
	var db_assets []mongo_asset.AssetInDB
	dao_asset := asset_dao.NewDAOAsset()

	db_assets = dao_asset.FindByUser(presenter.user.ID)
	return db_assets
}

func (presenter *AssetPresenter) GetAssetByCode(code string) (mongo_asset.AssetInDB, error) {
	var db_asset mongo_asset.AssetInDB
	dao_asset := asset_dao.NewDAOAsset()

	db_asset, err := dao_asset.FindByUserAndCode(presenter.user.ID, code)
	return db_asset, err
}

func (presenter *AssetPresenter) GetOverallValue() float64 {
	dao_asset := asset_dao.NewDAOAsset()
	db_assets := dao_asset.FindByUser(presenter.user.ID)
	manager := av_manager.NewAssetManager(presenter.user)
	var value float64 = 0
	for _, asset := range db_assets {
		asset_value := manager.GetCurrentValueByCode(asset.Code)
		value += asset_value.GetValue() * asset.Amount
	}
	return value
}

func (presenter *AssetPresenter) GetAssetDetails() []models_asset.AssetDetailedResponse {
	dao_company := company_dao.NewDAOCompany()
	dao_crypto := crypto_dao.NewDAOCryptocurrency()
	manager := av_manager.NewAssetManager(presenter.user)
	assets := presenter.GetAll()
	var result []models_asset.AssetDetailedResponse
	var codes []string = make([]string, 0)
	amounts_by_codes := make(map[string]float64)

	for _, asset := range assets {
		codes = append(codes, asset.Code)
		amounts_by_codes[asset.Code] = asset.Amount
	}

	companies := dao_company.FindManyByCodes(codes)
	crypto := dao_crypto.FindManyByCodes(codes)

	names_by_codes := make(map[string]string)
	values_by_codes := make(map[string]float64)

	for _, company := range companies {
		names_by_codes[company.Code] = company.Name
		values_by_codes[company.Code] = manager.GetCurrentValueByCode(company.Code).GetValue()
	}

	for _, crypto := range crypto {
		names_by_codes[crypto.Code] = crypto.Name
		values_by_codes[crypto.Code] = manager.GetCurrentValueByCode(crypto.Code).GetValue()
	}

	for key, value := range names_by_codes {
		result = append(result, models_asset.AssetDetailedResponse{
			Name:        value,
			Code:        key,
			Amount:      math.Floor(amounts_by_codes[key]*100) / 100,
			Total_value: math.Floor(values_by_codes[key]*amounts_by_codes[key]*100) / 100,
		})
	}

	return result
}
