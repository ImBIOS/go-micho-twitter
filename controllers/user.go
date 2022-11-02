package controllers

import (
	"net/http"
	"time"

	"github.com/ImBIOS/go-micho-twitter/configs"
	"github.com/ImBIOS/go-micho-twitter/helpers"
	"github.com/ImBIOS/go-micho-twitter/models"
	"github.com/ImBIOS/go-micho-twitter/models/responses"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Signup(c echo.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// Validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error{Status: "error", Message: "Bad request", Data: err.Error()},
		)
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error{
				Status:  "error",
				Message: "Validation error",
				Data:    validationErr.Error(),
			},
		)
	}

	// Hash the password
	hashed := helpers.HashAndSalt(user.Password)
	user.Password = string(hashed)

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Email:    user.Email,
		Password: user.Password,
	}

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		// Handle for email duplicate error
		merr := err.(mongo.WriteException)
		log.Errorf("Number of errors: %d", len(merr.WriteErrors))
		errCode := merr.WriteErrors[0].Code
		if errCode == 11000 {
			return c.JSON(
				http.StatusBadRequest,
				responses.Error{
					Status:  "error",
					Message: "Email already exists",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			responses.Error{
				Status:  "error",
				Message: "Internal server error",
				Data:    err.Error(),
			},
		)
	}

	newUser.Password = "" // Remove password from response
	return c.JSON(
		http.StatusCreated,
		responses.Success{
			Status:  "success",
			Message: "User created successfully",
			Data:    newUser,
		},
	)
}

func Signin(c echo.Context) (err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	var user models.User

	// Validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error{Status: "error", Message: "Bad request", Data: err.Error()},
		)
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error{
				Status:  "error",
				Message: "Validation error",
				Data:    validationErr.Error(),
			},
		)
	}

	var result bson.M
	err = userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).
		Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.JSON(
				http.StatusInternalServerError,
				responses.Error{
					Status:  "error",
					Message: "Email or password is incorrect",
				},
			)
		}
	}

	// Compare the password
	isPasswordCorrect := helpers.ComparePasswords(result["password"].(string), user.Password)
	if !isPasswordCorrect {
		return c.JSON(
			http.StatusUnauthorized,
			responses.Error{
				Status:  "error",
				Message: "Email or password is incorrect",
			},
		)
	}

	user = models.User{
		ID:    result["_id"].(primitive.ObjectID),
		Email: result["email"].(string),
		// Don't send password
	}

	//-----
	// JWT
	//-----

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = result["id"]
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	// // Update the token
	// filter := bson.M{"_id": result["id"]}
	// update := bson.M{"$set": bson.M{"token": token}}
	// _, err = userCollection.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	return c.JSON(
	// 		http.StatusInternalServerError,
	// 		responses.Error{
	// 			Status:  "error",
	// 			Message: "Internal server error",
	// 			Data:    err.Error(),
	// 		},
	// 	)
	// }

	return c.JSON(
		http.StatusOK,
		responses.Success{
			Status:  "success",
			Message: "User logged in successfully",
			Data:    user,
		},
	)
}

// func (h *Handler) Login(c echo.Context) (err error) {
// 	// Bind
// 	u := new(models.User)
// 	if err = c.Bind(u); err != nil {
// 		return
// 	}

// 	// Find user
// 	db := h.DB.Clone()
// 	defer db.Close()
// 	if err = db.DB("twitter").C("users").
// 		Find(bson.M{"email": u.Email, "password": u.Password}).One(u); err != nil {
// 		if err == mgo.ErrNotFound {
// 			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
// 		}
// 		return
// 	}

// 	//-----
// 	// JWT
// 	//-----

// 	// Create token
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	// Set claims
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["id"] = u.ID
// 	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

// 	// Generate encoded token and send it as response
// 	u.Token, err = token.SignedString([]byte(Key))
// 	if err != nil {
// 		return err
// 	}

// 	u.Password = "" // Don't send password
// 	return c.JSON(http.StatusOK, u)
// }

// func (h *Handler) Follow(c echo.Context) (err error) {
// 	userID := userIDFromToken(c)
// 	id := c.Param("id")

// 	// Add a follower to user
// 	db := h.DB.Clone()
// 	defer db.Close()
// 	if err = db.DB("twitter").C("users").
// 		UpdateId(bson.ObjectIdHex(id), bson.M{"$addToSet": bson.M{"followers": userID}}); err != nil {
// 		if err == mgo.ErrNotFound {
// 			return echo.ErrNotFound
// 		}
// 	}

// 	return
// }

// func userIDFromToken(c echo.Context) string {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	return claims["id"].(string)
// }
