package oauthservices

import (
	"github.com/amitdotkr/sso/sso-go/src/entities"
	"github.com/amitdotkr/sso/sso-go/src/global"
	"github.com/amitdotkr/sso/sso-go/src/services/userservices"
	"github.com/gofiber/fiber/v2"
)

func Generatetoken(c *fiber.Ctx) error {
	// var query []primitive.M

	userId, err := global.ValidatingUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	if err := userservices.CreateTokenPairGo(c, userId, "user"); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Token Generation Error",
				Detail: err.Error()},
		})
	}

	// log.Printf("userId: %v", userId)
	return c.Status(fiber.StatusOK).SendString(userId)
}
