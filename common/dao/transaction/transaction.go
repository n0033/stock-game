package transaction

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "github.com/ktylus/stock-game/common/models/mongo/transaction"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"
)

type DAOTransaction struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOTransaction() *DAOTransaction {
	var db *mongo.Database = database.GetDatabase()
	var dao_user DAOTransaction = DAOTransaction{
		db_client:       db,
		collection_name: config.COLLECTION["transactions"],
		collection:      db.Collection(config.COLLECTION["transactions"]),
	}
	return &dao_user
}

func (dao_transaction DAOTransaction) FindOne(id primitive.ObjectID) (models.TransactionInDB, error) {
	var db_transaction models.TransactionInDB

	err := dao_transaction.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_transaction)

	return db_transaction, err
}

func (dao_transaction DAOTransaction) FindMany(ids []primitive.ObjectID) []models.TransactionInDB {
	var retrieved_transactions []models.TransactionInDB
	cursor, err := dao_transaction.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var transaction models.TransactionInDB = convert.BsonDToStruct[models.TransactionInDB](result)
		retrieved_transactions = append(retrieved_transactions, transaction)
	}

	return retrieved_transactions
}

func (dao_transaction DAOTransaction) FindByUser(user_id primitive.ObjectID) []models.TransactionInDB {
	var retrieved_transactions []models.TransactionInDB

	cursor, err := dao_transaction.collection.Find(context.TODO(), bson.M{"user": user_id})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var transaction models.TransactionInDB = convert.BsonDToStruct[models.TransactionInDB](result)
		retrieved_transactions = append(retrieved_transactions, transaction)
	}

	return retrieved_transactions
}

func (dao_transaction DAOTransaction) Create(transaction models.TransactionInCreate) (models.TransactionInDB, error) {
	var db_transaction models.TransactionInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var create_doc bson.M = convert.StructToBsonM(transaction)

	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_transaction.collection.InsertOne(ctx, create_doc)

	if err != nil {
		return db_transaction, err
	}

	db_transaction = models.TransactionInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		TransactionBase: transaction.TransactionBase,
	}

	return db_transaction, err
}

func (dao_transaction DAOTransaction) Update(transaction_id primitive.ObjectID, transaction_update models.TransactionInUpdate) (models.TransactionInDB, error) {
	var updated_transaction models.TransactionInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	existing_transaction, err := dao_transaction.FindOne(transaction_id)

	if err != nil {
		return updated_transaction, err
	}

	if existing_transaction.User_id != transaction_update.User_id || existing_transaction.Company_id != transaction_update.Company_id {
		return updated_transaction, errors.New("cannot change transaction related account or company")
	}

	var update_doc bson.M = convert.StructToBsonM(transaction_update)
	update_doc["updated_at"] = time.Now().UTC()
	_, err = dao_transaction.collection.UpdateOne(ctx, bson.M{"_id": transaction_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_transaction, err
	}

	updated_transaction, err = dao_transaction.FindOne(transaction_id)

	return updated_transaction, err
}

func (dao_transaction DAOTransaction) Delete(transaction_id primitive.ObjectID) (models.TransactionInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_transaction models.TransactionInDB
	err := dao_transaction.collection.FindOneAndDelete(ctx, bson.M{"_id": transaction_id}).Decode(&deleted_transaction)
	return deleted_transaction, err
}
