package authorization

import (
	"time"

	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthorizationBase struct {
	User_id primitive.ObjectID `bson:"user"`
	Key     string             `bson:"key"`
	Expires time.Time          `bson:"expires"`
}

type AuthorizationInDB struct {
	AuthorizationBase  `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

type AuthorizationInCreate struct {
	AuthorizationBase `bson:",inline"`
}
