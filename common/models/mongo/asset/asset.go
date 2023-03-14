package asset

import (
	"github.com/n0033/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetBase struct {
	User_id primitive.ObjectID `bson:"user"`
	Code    string             `bson:"code"`
	Amount  float64            `bson:"amount"`
}

type AssetInDB struct {
	AssetBase          `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

type AssetInResponse struct {
	AssetBase `bson:",inline"`
	ID        primitive.ObjectID `bson:"_id"`
}

func (*AssetInResponse) FromAssetInDB(db_asset AssetInDB) AssetInResponse {
	return AssetInResponse{
		AssetBase: AssetBase{
			User_id: db_asset.User_id,
			Code:    db_asset.Code,
			Amount:  db_asset.Amount,
		},
		ID: db_asset.ID,
	}
}

type AssetInUpdate struct {
	AssetBase `bson:",inline"`
}

func (asset_update *AssetInUpdate) FromAssetInDB(db_asset AssetInDB) AssetInUpdate {
	return AssetInUpdate{
		AssetBase: AssetBase{
			User_id: db_asset.User_id,
			Code:    db_asset.Code,
			Amount:  db_asset.Amount,
		},
	}
}

type AssetInCreate struct {
	AssetBase `bson:",inline"`
}
