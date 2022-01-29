package userservices

import (
	"github.com/amitdotkr/sso/sso-go/src/global"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection mongo.Collection

func init() {
	UserCollection = *global.DB.Collection("users")
}
