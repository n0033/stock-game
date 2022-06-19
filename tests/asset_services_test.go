package tests

import (
	"fmt"
	"testing"
	"time"

	user_dao "github.com/ktylus/stock-game/common/dao/user"
	models_asset "github.com/ktylus/stock-game/common/models/mongo/asset"
	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
	"github.com/ktylus/stock-game/services/asset/manager"
	"github.com/ktylus/stock-game/services/asset/presenter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CREATED_ASSET_ID_PRESENTER *primitive.ObjectID

// var dao_asset = asset_dao.NewDAOAsset()
var dao_user = user_dao.NewDAOUser()

var asset_create = models_asset.AssetInCreate{
	AssetBase: models_asset.AssetBase{
		User_id: TEST_USER_ID,
		Code:    APPLE_CODE,
		Amount:  3.33,
	},
}
var user, _ = dao_user.FindOne(TEST_USER_ID)

var asset_presenter = presenter.NewAssetPresenter(user)
var asset_manager = manager.NewAssetManager(user)

// presenter
func TestGetOne(t *testing.T) {
	db_asset, _ := dao_asset.Create(asset_create)
	CREATED_ASSET_ID_PRESENTER = &db_asset.ID
	db_asset, err := asset_presenter.Get(APPLE_CODE)
	if err != nil {
		t.Fail()
	}
	if db_asset.AssetBase != asset_create.AssetBase {
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	db_assets := asset_presenter.GetAll()
	if len(db_assets) == 0 || len(db_assets) > 1 {
		t.Fail()
	}
	if db_assets[0].AssetBase != asset_create.AssetBase {
		t.Fail()
	}
}

func TestGetAssetByCode(t *testing.T) {
	db_asset, err := asset_presenter.GetAssetByCode(APPLE_CODE)
	if err != nil {
		t.Fail()
	}
	if db_asset.AssetBase != asset_create.AssetBase {
		t.Fail()
	}
}

func TestGetOverallValue(t *testing.T) {
	returned_value := asset_presenter.GetOverallValue()
	actual_value := asset_manager.GetCurrentValueByCode(APPLE_CODE).GetValue() * asset_create.AssetBase.Amount
	if returned_value != actual_value {
		t.Fail()
	}
}

// manager
func TestGetCurrentValueByCode(t *testing.T) {
	datapoint, _ := dao_stock.FindLatestByCode(APPLE_CODE)
	returned_datapoint := *asset_manager.GetCurrentValueByCode(APPLE_CODE)
	if time.Since(datapoint.Timestamp).Hours() > 1 && time.Since(datapoint.Last_used).Hours() > 1 {
		if datapoint.ID == returned_datapoint.ID {
			t.Fail()
		}
	} else {
		if datapoint.ID != returned_datapoint.ID {
			t.Fail()
		}
	}
}

func TestGetLatestValuesByCode(t *testing.T) {
	datapoints := asset_manager.GetLatestValuesByCode(APPLE_CODE)
	latest_datapoint, _ := dao_stock.FindLatestByCode(APPLE_CODE)
	if len(datapoints) == 0 {
		t.Fail()
	}
	if latest_datapoint.ID != datapoints[0].ID {
		t.Fail()
	}
}

func TestBuyAsset(t *testing.T) {
	user.Balance = 500
	apple_price := asset_manager.GetCurrentValueByCode(APPLE_CODE).GetValue()

	// buying too much - can't afford
	_, err := asset_manager.BuyAsset(APPLE_CODE, user.Balance/apple_price*10)
	if err == nil {
		t.Fail()
	}
	asset, _ := dao_asset.FindByUserAndCode(TEST_USER_ID, APPLE_CODE)
	if asset.Amount != 3.33 || user.Balance != 500 {
		t.Fail()
	}
	// buying an amount that the user can afford
	amount := user.Balance / (apple_price * 10)
	new_balance := user.Balance - apple_price*amount
	_, err = asset_manager.BuyAsset(APPLE_CODE, amount)
	if err != nil {
		t.Fail()
	}
	user, _ = dao_user.FindOne(user.ID)
	asset, _ = dao_asset.FindByUserAndCode(TEST_USER_ID, APPLE_CODE)
	if asset.Amount != (3.33+amount) || user.Balance != new_balance {
		t.Fail()
	}

	var user_update models_user.UserInUpdate
	user.Balance = 500
	user_update = user_update.FromUserInDB(user)
	dao_user.Update(user.ID, user_update)
	dao_asset.Delete(*CREATED_ASSET_ID_PRESENTER)
}

func TestSellAsset(t *testing.T) {
	asset_manager = manager.NewAssetManager(user)
	db_asset, _ := dao_asset.Create(asset_create)
	asset, _ := dao_asset.FindByUserAndCode(TEST_USER_ID, APPLE_CODE)
	amount := asset.Amount
	apple_price := asset_manager.GetCurrentValueByCode(APPLE_CODE).GetValue()

	// selling more assets than the user has
	_, err := asset_manager.SellAsset(APPLE_CODE, 2*amount)
	if err == nil {
		t.Fail()
	}
	if asset.Amount != 3.33 || user.Balance != 500 {
		t.Fail()
	}

	// selling a valid amount of assets
	new_balance := user.Balance + (amount/2)*apple_price
	_, err = asset_manager.SellAsset(APPLE_CODE, amount/2)
	if err != nil {
		t.Fail()
	}
	user, _ = dao_user.FindOne(user.ID)
	asset, _ = dao_asset.FindByUserAndCode(TEST_USER_ID, APPLE_CODE)
	if asset.Amount != (3.33-amount/2) || user.Balance != new_balance {
		fmt.Print(user.Balance, new_balance)
		t.Fail()
	}

	CREATED_ASSET_ID_PRESENTER = &db_asset.ID
	var user_update models_user.UserInUpdate
	user.Balance = 500
	user_update = user_update.FromUserInDB(user)
	dao_user.Update(user.ID, user_update)
	dao_asset.Delete(*CREATED_ASSET_ID_PRESENTER)
}
