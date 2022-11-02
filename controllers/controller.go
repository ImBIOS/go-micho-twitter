package controllers

import (
	"github.com/ImBIOS/go-micho-twitter/configs"
	"github.com/ImBIOS/go-micho-twitter/models"
)

// type (
// 	Handler struct {
// 		DB *mongo.Client
// 	}
// )

func NoPassword(u models.User) models.User {
	u.Password = ""
	return u
}

var (
	Key = configs.EnvSecret()
)
