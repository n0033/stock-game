package asset

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "github.com/ktylus/stock-game/common/models/mongo/asset"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"
)

type DAOAsset struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOAsset() *DAOAsset {
	var db *mongo.Database = database.GetDatabase()
	var dao_user DAOAsset = DAOAsset{
		db_client:       db,
		collection_name: config.COLLECTION["assets"],
		collection:      db.Collection(config.COLLECTION["assets"]),
	}
	return &dao_user
}

func (dao_asset DAOAsset) FindOne(id primitive.ObjectID) (models.AssetInDB, error) {
	var db_asset models.AssetInDB

	err := dao_asset.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_asset)

	return db_asset, err
}

func (dao_asset DAOAsset) FindMany(ids []primitive.ObjectID) []models.AssetInDB {
	var retrieved_assets []models.AssetInDB
	cursor, err := dao_asset.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var asset models.AssetInDB = convert.BsonDToStruct[models.AssetInDB](result)
		retrieved_assets = append(retrieved_assets, asset)
	}

	return retrieved_assets
}

func (dao_asset DAOAsset) FindByUser(user_id primitive.ObjectID) []models.AssetInDB {
	var retrieved_assets []models.AssetInDB

	cursor, err := dao_asset.collection.Find(context.TODO(), bson.M{"user": user_id})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var asset models.AssetInDB = convert.BsonDToStruct[models.AssetInDB](result)
		retrieved_assets = append(retrieved_assets, asset)
	}

	return retrieved_assets
}

func (dao_asset DAOAsset) FindByUserAndCode(user_id primitive.ObjectID, code string) (models.AssetInDB, error) {
	var db_asset models.AssetInDB

	err := dao_asset.collection.FindOne(context.TODO(), bson.M{"user": user_id, "code": code}).Decode(&db_asset)

	return db_asset, err
}

func (dao_asset DAOAsset) Create(asset models.AssetInCreate) (models.AssetInDB, error) {
	var db_asset models.AssetInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := dao_asset.FindByUserAndCode(asset.User_id, asset.Code)

	if err == nil {
		return db_asset, errors.New("user with given asset already has this company's assets")
	}

	var create_doc bson.M = convert.StructToBsonM(asset)

	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_asset.collection.InsertOne(ctx, create_doc)

	if err != nil {
		return db_asset, err
	}

	db_asset = models.AssetInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		AssetBase: asset.AssetBase,
	}

	return db_asset, err
}

func (dao_asset DAOAsset) Update(asset_id primitive.ObjectID, asset_update models.AssetInUpdate) (models.AssetInDB, error) {
	var updated_asset models.AssetInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	existing_asset, err := dao_asset.FindOne(asset_id)

	if err != nil {
		return updated_asset, errors.New("asset with given id does not exist")
	}

	if existing_asset.User_id != asset_update.User_id || existing_asset.Code != asset_update.Code {
		return updated_asset, errors.New("cannot change asset's related account or company")
	}

	var update_doc bson.M = convert.StructToBsonM(asset_update)
	update_doc["updated_at"] = time.Now().UTC()
	_, err = dao_asset.collection.UpdateOne(ctx, bson.M{"_id": asset_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_asset, err
	}

	updated_asset, err = dao_asset.FindOne(asset_id)

	if err != nil {
		return updated_asset, err
	}

	return updated_asset, err
}

func (dao_asset DAOAsset) Delete(asset_id primitive.ObjectID) (models.AssetInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_asset models.AssetInDB
	err := dao_asset.collection.FindOneAndDelete(ctx, bson.M{"_id": asset_id}).Decode(&deleted_asset)

	return deleted_asset, err
}
