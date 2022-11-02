package controllers

import (
	"github.com/ImBIOS/go-micho-twitter/configs"
)

// type (
// 	Handler struct {
// 		DB *mongo.Client
// 	}
// )

var (
	Key = configs.EnvSecret()
)
