package stock

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "github.com/n0033/stock-game/common/models/mongo/stock_datapoint"
	models_utils "github.com/n0033/stock-game/common/models/utils"
	"github.com/n0033/stock-game/config"
	convert "github.com/n0033/stock-game/services/bson_converter"
	"github.com/n0033/stock-game/services/database"
)

type DAOStock struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOStock() *DAOStock {
	var db *mongo.Database = database.GetDatabase()
	var dao_user DAOStock = DAOStock{
		db_client:       db,
		collection_name: config.COLLECTION["stock"],
		collection:      db.Collection(config.COLLECTION["stock"]),
	}
	return &dao_user
}

func (dao_stock DAOStock) FindOne(id primitive.ObjectID) (models.StockDatapointInDB, error) {
	var db_stock models.StockDatapointInDB

	err := dao_stock.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_stock)

	return db_stock, err
}

func (dao_stock DAOStock) FindByCodeAndTimestamp(code string, timestamp time.Time) (models.StockDatapointInDB, error) {
	var db_stock models.StockDatapointInDB

	err := dao_stock.collection.FindOne(context.TODO(), bson.M{"code": code, "timestamp": timestamp}).Decode(&db_stock)

	return db_stock, err
}

func (dao_stock DAOStock) FindMany(ids []primitive.ObjectID) []models.StockDatapointInDB {
	var retrieved_datapoints []models.StockDatapointInDB
	cursor, err := dao_stock.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var datapoint models.StockDatapointInDB = convert.BsonDToStruct[models.StockDatapointInDB](result)
		retrieved_datapoints = append(retrieved_datapoints, datapoint)
	}

	return retrieved_datapoints
}

func (dao_stock DAOStock) FindByCodeAndInterval(code string, date_from time.Time, date_to time.Time) []models.StockDatapointInDB {
	var datapoints_sorted []models.StockDatapointInDB

	options := options.Find()
	options.SetSort(bson.M{"timestamp": -1})
	sort_cursor, err := dao_stock.collection.Find(context.TODO(), bson.M{
		"code": code,
		"timestamp": bson.M{
			"$gte": date_from,
			"$lte": date_to,
		},
	}, options)

	if err != nil {
		log.Fatal(err)
	}

	if err = sort_cursor.All(context.TODO(), &datapoints_sorted); err != nil {
		log.Fatal(err)
	}
	return datapoints_sorted
}

func (dao_stock DAOStock) Find300LatestByCode(code string) []models.StockDatapointInDB {
	var datapoints_sorted []models.StockDatapointInDB

	options := options.Find()
	options.SetSort(bson.M{"timestamp": -1})
	options.SetLimit(300)
	sort_cursor, err := dao_stock.collection.Find(context.TODO(), bson.M{
		"code": code,
	}, options)

	if err != nil {
		log.Fatal(err)
	}

	if err = sort_cursor.All(context.TODO(), &datapoints_sorted); err != nil {
		log.Fatal(err)
	}
	var ids []primitive.ObjectID
	for _, datapoint := range datapoints_sorted {
		ids = append(ids, datapoint.ID)
	}
	dao_stock.collection.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"last_used": time.Now().UTC()}})

	return datapoints_sorted
}

func (dao_stock DAOStock) FindLatestByCode(code string) (models.StockDatapointInDB, error) {
	var db_stock models.StockDatapointInDB

	options := options.Find()
	options.SetSort(bson.M{"timestamp": -1})
	sort_cursor, err := dao_stock.collection.Find(context.TODO(), bson.M{"code": code}, options)

	if err != nil {
		log.Fatal(err)
	}

	var datapoints_sorted []bson.D

	if err = sort_cursor.All(context.TODO(), &datapoints_sorted); err != nil {
		log.Fatal(err)
	}
	if len(datapoints_sorted) == 0 {
		err = errors.New("datapoints for requested company do not exist")
	}

	if len(datapoints_sorted) > 0 {
		datapoint := datapoints_sorted[0]
		db_stock = convert.BsonDToStruct[models.StockDatapointInDB](datapoint)
	}
	dao_stock.collection.UpdateOne(context.TODO(), bson.M{"_id": db_stock.ID}, bson.M{"$set": bson.M{"last_used": time.Now().UTC()}})
	return db_stock, err
}

func (dao_stock DAOStock) Create(datapoint models.StockDatapointInCreate) (models.StockDatapointInDB, error) {
	var db_stock models.StockDatapointInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := dao_stock.FindByCodeAndTimestamp(datapoint.Code, datapoint.Timestamp)

	if err == nil {
		return db_stock, errors.New("such stock datapoint already exists")
	}

	var create_doc bson.M = convert.StructToBsonM(datapoint)

	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_stock.collection.InsertOne(ctx, create_doc)

	if err != nil {
		return db_stock, err
	}

	db_stock = models.StockDatapointInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		StockDatapointBase: datapoint.StockDatapointBase,
	}

	return db_stock, nil
}

func (dao_stock DAOStock) Update(datapoint_id primitive.ObjectID, datapoint_update models.StockDatapointInUpdate) (models.StockDatapointInDB, error) {
	var updated_datapoint models.StockDatapointInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var update_doc bson.M = convert.StructToBsonM(datapoint_update)
	update_doc["updated_at"] = time.Now().UTC()
	_, err := dao_stock.collection.UpdateOne(ctx, bson.M{"_id": datapoint_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_datapoint, err
	}
	updated_datapoint, err = dao_stock.FindOne(datapoint_id)

	return updated_datapoint, err
}

func (dao_stock DAOStock) Delete(datapoint_id primitive.ObjectID) (models.StockDatapointInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_datapoint models.StockDatapointInDB
	err := dao_stock.collection.FindOneAndDelete(ctx, bson.M{"_id": datapoint_id}).Decode(&deleted_datapoint)

	return deleted_datapoint, err
}
