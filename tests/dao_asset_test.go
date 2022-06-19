package tests

import (
	"testing"

	asset_dao "github.com/ktylus/stock-game/common/dao/asset"
	models_asset "github.com/ktylus/stock-game/common/models/mongo/asset"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CREATED_ASSET_ID *primitive.ObjectID
var dao_asset = asset_dao.NewDAOAsset()

func TestCreateAsset(t *testing.T) {
	asset_create := models_asset.AssetInCreate{
		AssetBase: models_asset.AssetBase{
			User_id: TEST_USER_ID,
			Code:    APPLE_CODE,
			Amount:  3.33,
		},
	}
	db_asset, err := dao_asset.Create(asset_create)
	if err != nil {
		t.Fail()
	}
	CREATED_ASSET_ID = &db_asset.ID
}

func TestCreateExistingAsset(t *testing.T) {
	asset_create := models_asset.AssetInCreate{
		AssetBase: models_asset.AssetBase{
			User_id: TEST_USER_ID,
			Code:    APPLE_CODE,
			Amount:  3.33,
		},
	}
	_, err := dao_asset.Create(asset_create)
	if err == nil {
		t.Fail()
	}
}

func TestFindOneAsset(t *testing.T) {
	_, err := dao_asset.FindOne(*CREATED_ASSET_ID)
	if err != nil {
		t.Fail()
	}
}

func TestFindManyAssets(t *testing.T) {
	asset_ids := make([]primitive.ObjectID, 0)
	asset_ids = append(asset_ids, *CREATED_ASSET_ID)
	assets := dao_asset.FindMany(asset_ids)
	if len(assets) == 0 {
		t.Fail()
	}
}

func TestDeleteAsset(t *testing.T) {
	_, err := dao_asset.Delete(*CREATED_ASSET_ID)
	if err != nil {
		t.Fail()
	}
}

func TestFindInexistentAsset(t *testing.T) {
	_, err := dao_asset.FindOne(*CREATED_ASSET_ID)
	if err == nil {
		t.Fail()
	}
}
