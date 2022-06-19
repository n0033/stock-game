package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DateTime struct {
	Created_at time.Time `bson:"created_at"`
	Updated_at time.Time `bson:"updated_at"`
}

type DBModelMixin struct {
	ID       primitive.ObjectID `bson:"_id"`
	DateTime `bson:",inline"`
}
