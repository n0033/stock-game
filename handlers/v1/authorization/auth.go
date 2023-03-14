package authorization

import (
	"github.com/gofiber/fiber/v2"
	user_dao "github.com/n0033/stock-game/common/dao/user"
	models_user "github.com/n0033/stock-game/common/models/mongo/user"
	models_auth "github.com/n0033/stock-game/common/models/resources/auth"
	authorization "github.com/n0033/stock-game/common/security/authorization"
	password_handler "github.com/n0033/stock-game/common/security/password_handler"
	"github.com/n0033/stock-game/config"
	"github.com/n0033/stock-game/utils"
)

func RegisterView(c *fiber.Ctx) error {
	if authorization.CookieAuthorize(c) {
		return c.Redirect("/user")
	}
	return c.Render("auth/register", fiber.Map{
		"errors":        c.Locals("errors"),
		"authenticated": authorization.CookieAuthorize(c)})

}

func LoginView(c *fiber.Ctx) error {
	if authorization.CookieAuthorize(c) {
		return c.Redirect("/user")
	}
	return c.Render("auth/login", fiber.Map{
		"errors":        c.Locals("errors"),
		"messages":      c.Locals("messages"),
		"authenticated": authorization.CookieAuthorize(c)})
}

func Register(c *fiber.Ctx) error {
	var register_request models_auth.RegisterRequest

	errors := make([]string, 0)
	messages := make([]string, 0)

	if err := c.BodyParser(&register_request); err != nil {
		return err
	}

	if register_request.Password != register_request.Password_again {
		errors = append(errors, utils.MESSAGES["ERROR_REGISTER_PASSWORD_MISMATCH"])
		c.Locals("errors", errors)
		return RegisterView(c)
	}

	dao_user := user_dao.NewDAOUser()

	success, hashed_password := password_handler.Hash(register_request.Password)

	if !success {
		errors = append(errors, utils.MESSAGES["ERROR_UNEXPECTED"])
		c.Locals("errors", errors)
		return RegisterView(c)
	}

	user_create := models_user.UserInCreate{
		UserBase: models_user.UserBase{
			Username:        register_request.Username,
			Email:           register_request.Email,
			Balance:         config.USER_DEFAULT_BALANCE,
			Hashed_password: hashed_password,
		},
	}
	_, err := dao_user.Create(user_create)

	if err != nil {
		errors = append(errors, err.Error())
		c.Locals("errors", errors)
		return RegisterView(c)
	}

	messages = append(messages, utils.MESSAGES["SUCCESS_ACCOUNT_CREATION"])
	c.Locals("messages", messages)

	return LoginView(c)
}

func Login(c *fiber.Ctx) error {
	var login_request models_auth.LoginRequest
	errors := make([]string, 0)

	if err := c.BodyParser(&login_request); err != nil {
		return err
	}

	db_user, success := authorization.Authorize(login_request)

	if !success {
		errors = append(errors, utils.MESSAGES["ERROR_LOGIN_FAILURE"])
		c.Locals("errors", errors)
		return LoginView(c)
	}

	authorization.CreateCookie(db_user, c)
	return c.Redirect("/user")
}
