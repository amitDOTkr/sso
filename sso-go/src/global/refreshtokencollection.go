package global

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var RefreshCollection mongo.Collection

func Refresh() {
	RefreshCollection = *DB.Collection("refresh")
	RefreshCollection.Indexes().CreateOne(Ctx, mongo.IndexModel{Keys: bson.D{{Key: "expireAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)})
}
