package authorization

import (
	"time"

	"github.com/gofiber/fiber/v2"
	authorization_dao "github.com/n0033/stock-game/common/dao/authorization"
	user_dao "github.com/n0033/stock-game/common/dao/user"
	user_models "github.com/n0033/stock-game/common/models/mongo/user"
	auth_models "github.com/n0033/stock-game/common/models/resources/auth"
	passwords "github.com/n0033/stock-game/common/security/password_handler"
)

func Authorize(request auth_models.LoginRequest) (user_models.UserInDB, bool) {
	var user user_models.UserInDB
	dao_user := user_dao.NewDAOUser()
	user, err := dao_user.FindByUsername(request.Username)

	if err != nil {
		return user, false
	}

	if !passwords.Compare(request.Password, user.Hashed_password) {
		return user, false
	}

	return user, true
}

func CookieAuthorize(c *fiber.Ctx) bool {
	dao_authorization := authorization_dao.NewDAOAuthorization()
	key := c.Cookies("identity")
	if key == "" {
		return false
	}

	authorization, err := dao_authorization.FindByKey(key)

	if err != nil {
		return false
	}
	if authorization.Expires.Before(time.Now()) {
		return false
	}

	dao_user := user_dao.NewDAOUser()
	user, err := dao_user.FindOne(authorization.User_id)

	if err != nil {
		return false
	}

	return passwords.Compare(user.Hashed_password, key)
}
