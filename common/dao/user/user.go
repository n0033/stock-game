package user

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "github.com/ktylus/stock-game/common/models/mongo/user"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"
	"github.com/ktylus/stock-game/utils"
)

type DAOUser struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOUser() *DAOUser {
	var db *mongo.Database = database.GetDatabase()
	var dao_user DAOUser = DAOUser{
		db_client:       db,
		collection_name: config.COLLECTION["users"],
		collection:      db.Collection(config.COLLECTION["users"]),
	}
	return &dao_user
}

func (dao_user DAOUser) FindOne(id primitive.ObjectID) (models.UserInDB, error) {
	var db_user models.UserInDB

	err := dao_user.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_user)

	return db_user, err
}

func (dao_user DAOUser) FindMany(ids []primitive.ObjectID) []models.UserInDB {
	var retrieved_users []models.UserInDB
	cursor, err := dao_user.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var user models.UserInDB = convert.BsonDToStruct[models.UserInDB](result)
		retrieved_users = append(retrieved_users, user)
	}

	return retrieved_users
}

func (dao_user DAOUser) FindByEmail(email string) (models.UserInDB, error) {
	var db_user models.UserInDB

	err := dao_user.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&db_user)

	return db_user, err
}

func (dao_user DAOUser) FindByUsername(username string) (models.UserInDB, error) {
	var db_user models.UserInDB

	err := dao_user.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&db_user)

	return db_user, err
}

func (dao_user DAOUser) Create(user models.UserInCreate) (models.UserInDB, error) {
	var db_user models.UserInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := dao_user.FindByEmail(user.Email)

	if err == nil {
		return db_user, errors.New(utils.MESSAGES["ERROR_REGISTER_EMAIL_ALREADY_EXISTS"])
	}

	_, err = dao_user.FindByUsername(user.Username)

	if err == nil {
		return db_user, errors.New(utils.MESSAGES["ERROR_REGISTER_USERNAME_ALREADY_EXISTS"])
	}

	var create_doc bson.M = convert.StructToBsonM(user)

	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_user.collection.InsertOne(ctx, create_doc)

	if err != nil {
		log.Fatal(err)
	}

	db_user = models.UserInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		UserBase: user.UserBase,
	}

	return db_user, err
}

func (dao_user DAOUser) Update(user_id primitive.ObjectID, user_update models.UserInUpdate) (models.UserInDB, error) {
	var updated_user models.UserInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	existing_user, err := dao_user.FindByEmail(user_update.Email)
	if err == nil {
		if existing_user.ID != user_id {
			return updated_user, errors.New("user with given email already exists")
		}
	}

	var update_doc bson.M = convert.StructToBsonM(user_update)
	update_doc["updated_at"] = time.Now().UTC()
	_, err = dao_user.collection.UpdateOne(ctx, bson.M{"_id": user_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_user, err
	}

	updated_user, err = dao_user.FindOne(user_id)

	return updated_user, err
}

func (dao_user DAOUser) Delete(user_id primitive.ObjectID) (models.UserInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_user models.UserInDB
	err := dao_user.collection.FindOneAndDelete(ctx, bson.M{"_id": user_id}).Decode(&deleted_user)

	return deleted_user, err
}
