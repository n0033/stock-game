package company

import (
	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyBase struct {
	Name string `bson:"name"`
	Code string `bson:"code"`
}

type CompanyInDB struct {
	CompanyBase        `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

type CompanyInResponse struct {
	CompanyBase `bson:",inline"`
	ID          primitive.ObjectID `bson:",inline"`
}

type CompanyInUpdate struct {
	CompanyBase `bson:",inline"`
}

type CompanyInCreate struct {
	CompanyBase `bson:",inline"`
}
