package authorization

import (
	"context"
	"errors"
	"time"

	models "github.com/ktylus/stock-game/common/models/mongo/authorization"
	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	password "github.com/ktylus/stock-game/common/security/password_handler"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DAOAuthorization struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOAuthorization() *DAOAuthorization {
	var db *mongo.Database = database.GetDatabase()
	var dao_company DAOAuthorization = DAOAuthorization{
		db_client:       db,
		collection_name: config.COLLECTION["authorizations"],
		collection:      db.Collection(config.COLLECTION["authorizations"]),
	}
	return &dao_company
}

func (dao_authorization DAOAuthorization) FindByKey(hash string) (models.AuthorizationInDB, error) {
	var db_authorization models.AuthorizationInDB
	err := dao_authorization.collection.FindOne(context.TODO(), bson.M{"key": hash}).Decode(&db_authorization)

	return db_authorization, err
}

func (dao_authorization DAOAuthorization) Create(user models_user.UserInDB) (models.AuthorizationInDB, error) {
	var db_authorization models.AuthorizationInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	success, hash := password.Hash(user.Hashed_password)

	if !success {
		return db_authorization, errors.New("failed to hash password")

	}

	var authorization models.AuthorizationInCreate = models.AuthorizationInCreate{
		AuthorizationBase: models.AuthorizationBase{
			User_id: user.ID,
			Key:     hash,
			Expires: time.Now().Add(time.Duration(config.AUTH_COOKIE_EXPIRY * 1e9)),
		},
	}

	var create_doc bson.M = convert.StructToBsonM(authorization)
	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_authorization.collection.InsertOne(ctx, create_doc)

	if err != nil {
		return db_authorization, err
	}

	db_authorization = models.AuthorizationInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		AuthorizationBase: authorization.AuthorizationBase,
	}

	return db_authorization, err
}
