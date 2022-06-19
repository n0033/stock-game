package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	authorization_dao "github.com/ktylus/stock-game/common/dao/authorization"
	user_dao "github.com/ktylus/stock-game/common/dao/user"
	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
)

func GetUser(c *fiber.Ctx) models_user.UserInDB {
	dao_authorization := authorization_dao.NewDAOAuthorization()
	dao_user := user_dao.NewDAOUser()

	key := c.Cookies("identity")
	if key == "" {
		log.Fatal("Identity cookie does not exist.")
	}

	db_authorization, err := dao_authorization.FindByKey(key)

	if err != nil {
		log.Fatal(err)
	}

	db_user, err := dao_user.FindOne(db_authorization.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return db_user
}
