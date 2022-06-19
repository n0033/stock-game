package cryptocurrency

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "github.com/ktylus/stock-game/common/models/mongo/cryptocurrency"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"
)

type DAOCryptocurrency struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOCryptocurrency() *DAOCryptocurrency {
	var db *mongo.Database = database.GetDatabase()
	var dao_crypto DAOCryptocurrency = DAOCryptocurrency{
		db_client:       db,
		collection_name: config.COLLECTION["cryptocurrencies"],
		collection:      db.Collection(config.COLLECTION["cryptocurrencies"]),
	}
	return &dao_crypto
}

func (dao_crypto DAOCryptocurrency) FindOne(id primitive.ObjectID) (models.CryptocurrencyInDB, error) {
	var db_crypto models.CryptocurrencyInDB

	err := dao_crypto.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_crypto)

	return db_crypto, err
}

func (dao_crypto DAOCryptocurrency) FindMany(ids []primitive.ObjectID) []models.CryptocurrencyInDB {
	var retrieved_crypto []models.CryptocurrencyInDB
	cursor, err := dao_crypto.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var crypto models.CryptocurrencyInDB = convert.BsonDToStruct[models.CryptocurrencyInDB](result)
		retrieved_crypto = append(retrieved_crypto, crypto)
	}
	return retrieved_crypto
}

func (dao_crypto DAOCryptocurrency) FindByCode(code string) (models.CryptocurrencyInDB, error) {
	var db_crypto models.CryptocurrencyInDB

	err := dao_crypto.collection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&db_crypto)

	return db_crypto, err
}

func (dao_crypto DAOCryptocurrency) FindManyByCodes(codes []string) []models.CryptocurrencyInDB {
	var db_crypto []models.CryptocurrencyInDB

	cursor, err := dao_crypto.collection.Find(context.TODO(), bson.M{"code": bson.M{"$in": codes}})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var cryptocurrency models.CryptocurrencyInDB = convert.BsonDToStruct[models.CryptocurrencyInDB](result)
		db_crypto = append(db_crypto, cryptocurrency)
	}
	return db_crypto
}

func (dao_crypto DAOCryptocurrency) Search(term string) []models.CryptocurrencyInDB {
	var retrieved_crypto []models.CryptocurrencyInDB
	filter := bson.D{{"$or", []bson.D{
		{{"name", bson.D{{"$regex", primitive.Regex{Pattern: term, Options: "i"}}}}},
		{{Key: "code", Value: bson.E{Key: "$regex", Value: primitive.Regex{Pattern: term, Options: "i"}}}},
	}}}

	cursor, err := dao_crypto.collection.Find(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var crypto models.CryptocurrencyInDB = convert.BsonDToStruct[models.CryptocurrencyInDB](result)
		retrieved_crypto = append(retrieved_crypto, crypto)
	}
	return retrieved_crypto
}

func (dao_crypto DAOCryptocurrency) Create(crypto models.CryptocurrencyInDB) (models.CryptocurrencyInDB, error) {
	var db_crypto models.CryptocurrencyInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := dao_crypto.FindByCode(crypto.Code)

	if err == nil {
		return db_crypto, errors.New("cryptocurrency with given code already exists")
	}

	var create_doc bson.M = convert.StructToBsonM(crypto)
	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_crypto.collection.InsertOne(ctx, create_doc)

	if err != nil {
		return db_crypto, err
	}

	db_crypto = models.CryptocurrencyInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		CryptocurrencyBase: crypto.CryptocurrencyBase,
	}
	return db_crypto, err
}

func (dao_crypto DAOCryptocurrency) Update(crypto_id primitive.ObjectID, crypto_update models.CryptocurrencyInUpdate) (models.CryptocurrencyInDB, error) {
	var updated_crypto models.CryptocurrencyInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	existing_crypto, err := dao_crypto.FindByCode(crypto_update.Code)

	if err == nil {
		if existing_crypto.ID != crypto_id {
			return updated_crypto, errors.New("cryptocurrency with given code already exists")
		}
	}

	var update_doc bson.M = convert.StructToBsonM(crypto_update)

	update_doc["updated_at"] = time.Now().UTC()
	_, err = dao_crypto.collection.UpdateOne(ctx, bson.M{"_id": crypto_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_crypto, err
	}

	updated_crypto, err = dao_crypto.FindOne(crypto_id)

	if err != nil {
		return updated_crypto, err
	}

	return updated_crypto, err
}

func (dao_crypto DAOCryptocurrency) Delete(crypto_id primitive.ObjectID) (models.CryptocurrencyInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_crypto models.CryptocurrencyInDB
	err := dao_crypto.collection.FindOneAndDelete(ctx, bson.M{"_id": crypto_id}).Decode(&deleted_crypto)

	return deleted_crypto, err
}
