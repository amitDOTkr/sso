package userservices

import (
	"github.com/amitdotkr/sso/sso-go/src/entities"
	"github.com/amitdotkr/sso/sso-go/src/global"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// func RegenerateToken(c *fiber.Ctx) error {
// 	global.RegenerateTokenUsingRefreshToken(c)

// 	return c.SendStatus(fiber.StatusOK)
// }

func Signin(c *fiber.Ctx) error {

	var user entities.User

	validate := validator.New()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": entities.Error{Type: "Validation Error", Detail: err.Error()},
		})
	}

	data := &entities.User{}
	res := UserCollection.FindOne(global.Ctx, bson.M{"email": user.Email})
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Database/Json Error",
				Detail: err.Error()},
		})
	}

	validUser := global.CheckPasswordHash(user.Password, data.Password)
	if !validUser {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Authentication Error",
				Detail: "The Credentials you provided cannot be Authenticated."},
		})
	}

	if err := CreateTokenPairGo(c, data.ID.Hex(), "user"); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Token Generation Error",
				Detail: err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"seller": entities.User{
			ID:              data.ID,
			Name:            data.Name,
			Email:           data.Email,
			IsEmailVerified: data.IsEmailVerified,
			ProfileImage:    data.ProfileImage,
			IsActice:        data.IsActice,
		},
	})
}

func CreateTokenPairGo(c *fiber.Ctx, userid, role string) error {

	if err := global.CreateAccessTokenGo(c, userid, role); err != nil {
		return err
	}

	if err := global.CreateRefreshTokenGo(c, userid, role); err != nil {
		return err
	}
	return nil
}
