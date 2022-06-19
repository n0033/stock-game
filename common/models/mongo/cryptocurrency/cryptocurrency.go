package cryptocurrency

import (
	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CryptocurrencyBase struct {
	Name string `bson:"name"`
	Code string `bson:"code"`
}

type CryptocurrencyInDB struct {
	CryptocurrencyBase `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

type CryptocurrencyInResponse struct {
	CryptocurrencyBase `bson:",inline"`
	ID                 primitive.ObjectID `bson:",inline"`
}

type CryptocurrencyInUpdate struct {
	CryptocurrencyBase `bson:",inline"`
}

type CryptocurrencInCreate struct {
	CryptocurrencyBase `bson:",inline"`
}
