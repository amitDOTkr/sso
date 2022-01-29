package routes

import (
	"github.com/amitdotkr/sso/sso-go/src/services/oauthservices"
	"github.com/amitdotkr/sso/sso-go/src/services/userservices"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"api_working": true,
		})
	})

	user := app.Group("user")
	user.Post("/signup", userservices.Signup)
	user.Post("/signin", userservices.Signin)
	user.Get("/signout", userservices.SignOut)

	oauth := app.Group("oauth")
	oauth.Post("/generatetoken", oauthservices.Generatetoken)

}
