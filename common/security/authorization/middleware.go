package authorization

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	authorization_dao "github.com/ktylus/stock-game/common/dao/authorization"
	user_dao "github.com/ktylus/stock-game/common/dao/user"
	models_user "github.com/ktylus/stock-game/common/models/mongo/user"
	passwords "github.com/ktylus/stock-game/common/security/password_handler"
	"github.com/ktylus/stock-game/config"
)

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Authorize    func(c *fiber.Ctx) bool
}

var ConfigDefault = Config{
	Filter:       nil,
	Unauthorized: nil,
	Authorize:    nil,
}

func defaultFilter(c *fiber.Ctx) bool {
	key := c.Cookies("identity")
	if key == "" {
		return false
	}
	return passwords.Compare(config.AUTH_SECRET, key)
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = defaultFilter
	}

	if cfg.Authorize == nil {
		cfg.Authorize = func(c *fiber.Ctx) bool {
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
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.Redirect("/auth/login")
		}
	}

	return cfg
}

func CreateCookie(user models_user.UserInDB, c *fiber.Ctx) {
	dao_authorization := authorization_dao.NewDAOAuthorization()
	db_authorization, err := dao_authorization.Create(user)

	if err != nil {
		log.Fatal(err.Error())
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "identity"
	cookie.Value = db_authorization.Key
	cookie.Expires = db_authorization.Expires
	c.Cookie(cookie)
}

func New(config Config) fiber.Handler {
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		authorized := cfg.Authorize(c)

		if authorized {
			c.Locals("is_authenticated", authorized)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}
