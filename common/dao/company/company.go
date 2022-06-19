package company

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "github.com/ktylus/stock-game/common/models/mongo/company"
	models_utils "github.com/ktylus/stock-game/common/models/utils"
	"github.com/ktylus/stock-game/config"
	convert "github.com/ktylus/stock-game/services/bson_converter"
	"github.com/ktylus/stock-game/services/database"
)

type DAOCompany struct {
	db_client       *mongo.Database
	collection_name string
	collection      *mongo.Collection
}

func NewDAOCompany() *DAOCompany {
	var db *mongo.Database = database.GetDatabase()
	var dao_company DAOCompany = DAOCompany{
		db_client:       db,
		collection_name: config.COLLECTION["companies"],
		collection:      db.Collection(config.COLLECTION["companies"]),
	}
	return &dao_company
}

func (dao_company DAOCompany) FindOne(id primitive.ObjectID) (models.CompanyInDB, error) {
	var db_company models.CompanyInDB

	err := dao_company.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&db_company)

	return db_company, err
}

func (dao_company DAOCompany) FindMany(ids []primitive.ObjectID) []models.CompanyInDB {
	var retrieved_companies []models.CompanyInDB
	cursor, err := dao_company.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var company models.CompanyInDB = convert.BsonDToStruct[models.CompanyInDB](result)
		retrieved_companies = append(retrieved_companies, company)
	}
	return retrieved_companies
}

func (dao_company DAOCompany) FindByCode(code string) (models.CompanyInDB, error) {
	var db_company models.CompanyInDB

	err := dao_company.collection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&db_company)

	return db_company, err
}

func (dao_company DAOCompany) FindManyByCodes(codes []string) []models.CompanyInDB {
	var db_companies []models.CompanyInDB

	cursor, err := dao_company.collection.Find(context.TODO(), bson.M{"code": bson.M{"$in": codes}})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var company models.CompanyInDB = convert.BsonDToStruct[models.CompanyInDB](result)
		db_companies = append(db_companies, company)
	}
	return db_companies
}

func (dao_company DAOCompany) Search(term string) []models.CompanyInDB {
	var retrieved_companies []models.CompanyInDB
	filter := bson.D{{"$or", []bson.D{
		{{"name", bson.D{{"$regex", primitive.Regex{Pattern: term, Options: "i"}}}}},
		{{Key: "code", Value: bson.E{Key: "$regex", Value: primitive.Regex{Pattern: term, Options: "i"}}}},
	}}}

	cursor, err := dao_company.collection.Find(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		var company models.CompanyInDB = convert.BsonDToStruct[models.CompanyInDB](result)
		retrieved_companies = append(retrieved_companies, company)
	}
	return retrieved_companies
}

func (dao_company DAOCompany) Create(company models.CompanyInCreate) (models.CompanyInDB, error) {
	var db_company models.CompanyInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := dao_company.FindByCode(company.Code)

	if err == nil {
		return db_company, errors.New("company with given code already exists")
	}

	var create_doc bson.M = convert.StructToBsonM(company)
	var time_now time.Time = time.Now().UTC()
	create_doc["created_at"] = time_now
	create_doc["updated_at"] = time_now

	result, err := dao_company.collection.InsertOne(ctx, create_doc)

	if err != nil {
		log.Fatal(err)
	}

	db_company = models.CompanyInDB{
		DBModelMixin: models_utils.DBModelMixin{
			ID: result.InsertedID.(primitive.ObjectID),
			DateTime: models_utils.DateTime{
				Created_at: create_doc["created_at"].(time.Time),
				Updated_at: create_doc["updated_at"].(time.Time),
			},
		},
		CompanyBase: company.CompanyBase,
	}
	return db_company, err
}

func (dao_company DAOCompany) Update(company_id primitive.ObjectID, company_update models.CompanyInUpdate) (models.CompanyInDB, error) {
	var updated_company models.CompanyInDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	existing_company, err := dao_company.FindByCode(company_update.Code)

	if err == nil {
		if existing_company.ID != company_id {
			return updated_company, errors.New("company with given code already exists")
		}
	}

	var update_doc bson.M = convert.StructToBsonM(company_update)

	update_doc["updated_at"] = time.Now().UTC()
	_, err = dao_company.collection.UpdateOne(ctx, bson.M{"_id": company_id}, bson.M{"$set": update_doc})
	if err != nil {
		return updated_company, err
	}

	updated_company, err = dao_company.FindOne(company_id)

	if err != nil {
		return updated_company, err
	}

	return updated_company, err
}

func (dao_company DAOCompany) Delete(company_id primitive.ObjectID) (models.CompanyInDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var deleted_company models.CompanyInDB
	err := dao_company.collection.FindOneAndDelete(ctx, bson.M{"_id": company_id}).Decode(&deleted_company)

	if err != nil {
		return deleted_company, err
	}
	return deleted_company, err
}
