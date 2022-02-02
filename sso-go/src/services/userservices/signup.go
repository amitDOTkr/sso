package userservices

import (
	"time"

	"github.com/amitdotkr/sso/sso-go/src/entities"
	"github.com/amitdotkr/sso/sso-go/src/global"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Signup(c *fiber.Ctx) error {

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

	isExist, err := IsEmailAlreadyExist(user.Email)
	if err != nil {
		return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
			"error": entities.Error{Type: "DataBase Error", Detail: err.Error()},
		})
	}
	if isExist {
		return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
			"error": entities.Error{Type: "Email Already Exist"},
		})
	}

	hashedpassword, err := global.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": entities.Error{Type: "Password Hashing", Detail: err.Error()},
		})
	}

	sellerData := entities.User{
		Name:       user.Name,
		Email:      user.Email,
		Password:   hashedpassword,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	sellerRes, err := UserCollection.InsertOne(global.Ctx, sellerData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}

	oid, _ := sellerRes.InsertedID.(primitive.ObjectID)

	if err := CreateTokenPairGo(c, oid.Hex(), "user"); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Token Generation Error",
				Detail: err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"seller": entities.User{
			ID:              oid,
			Name:            user.Name,
			Email:           user.Email,
			IsEmailVerified: user.IsEmailVerified,
			ProfileImage:    user.ProfileImage,
			IsActice:        user.IsActice,
		},
	})
}

func IsEmailAlreadyExist(email string) (bool, error) {
	count, err := UserCollection.CountDocuments(global.Ctx, bson.M{"email": email}, options.Count())
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
