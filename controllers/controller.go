package controllers

import (
	"github.com/ImBIOS/go-micho-twitter/configs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var Key = configs.EnvSecret()
var validate = validator.New()

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var tweetCollection *mongo.Collection = configs.GetCollection(configs.DB, "tweets")
