package user

import (
	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBase struct {
	Username        string  `bson:"username,omitempty"`
	Email           string  `bson:"email,omitempty"`
	Balance         float64 `bson:"balance,omitempty"`
	Hashed_password string  `bson:"hashed_password,omitempty"`
}

type UserInDB struct {
	utils.DBModelMixin `bson:",inline"`
	UserBase           `bson:",inline"`
}

type UserInResponse struct {
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Balance  float64            `json:"balance"`
	ID       primitive.ObjectID `json:"id"`
}

func (*UserInResponse) FromUserInDB(db_user UserInDB) UserInResponse {
	return UserInResponse{
		Username: db_user.Username,
		Email:    db_user.Email,
		Balance:  db_user.Balance,
		ID:       db_user.ID,
	}
}

type UserInUpdate struct {
	UserBase `bson:",inline"`
}

func (user_update *UserInUpdate) FromUserInDB(db_user UserInDB) UserInUpdate {
	return UserInUpdate{
		UserBase{
			Username:        db_user.Username,
			Email:           db_user.Email,
			Balance:         db_user.Balance,
			Hashed_password: db_user.Hashed_password,
		},
	}
}

type UserInCreate struct {
	UserBase `bson:",inline"`
}

type UserCreateRequest struct {
	Username       string
	Email          string
	Password       string
	Password_again string
}
