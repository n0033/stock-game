package transaction

import (
	"time"

	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionBase struct {
	Company_id primitive.ObjectID `bson:"company"`
	User_id    primitive.ObjectID `bson:"account"`
	Timestamp  time.Time          `bson:"timestamp"`
	Volume     float64            `bson:"volume"`
	Value      float64            `bson:"value"`
}

type TransactionInDB struct {
	TransactionBase    `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

type TransactionInResponse struct {
	TransactionBase `bson:",inline"`
	ID              primitive.ObjectID `bson:"_id"`
}

type TransactionInUpdate struct {
	TransactionBase `bson:",inline"`
}

type TransactionInCreate struct {
	TransactionBase `bson:",inline"`
}
