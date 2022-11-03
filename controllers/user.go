//nolint:errcheck // False positive
package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ImBIOS/go-micho-twitter/helpers"
	"github.com/ImBIOS/go-micho-twitter/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func AddUser(c echo.Context) (err error) {
	const timeout time.Duration = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	var user models.User

	defer cancel()

	// Validate the request body type
	if err = c.Bind(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{Status: "error", Message: "Bad request", Data: err.Error()},
		)
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{
				Status:  "error",
				Message: "Validation error",
				Data:    validationErr.Error(),
			},
		)
	}

	// Hash the password
	hashed := helpers.HashAndSalt(user.Password)
	user.Password = hashed

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

		const duplicateKeyErrorCode = 11000
		if errCode == duplicateKeyErrorCode {
			return c.JSON(
				http.StatusBadRequest,
				models.Response{
					Status:  "error",
					Message: "Email already exists",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    err.Error(),
			},
		)
	}

	newUser.Password = "" // Remove password from response

	return c.JSON(
		http.StatusCreated,
		models.Response{
			Status:  "success",
			Message: "User created successfully",
			Data:    newUser,
		},
	)
}

func Authenticate(c echo.Context) (err error) {
	var user models.User

	// Validate the request body
	if err = c.Bind(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{Status: "error", Message: "Bad request", Data: err.Error()},
		)
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{
				Status:  "error",
				Message: "Validation error",
				Data:    validationErr.Error(),
			},
		)
	}

	var result bson.M
	err = userCollection.FindOne(context.TODO(),
		bson.D{{
			Key:   "email",
			Value: user.Email,
		}},
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.JSON(
				http.StatusInternalServerError,
				models.Response{
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
			models.Response{
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
	claims["id"] = user.ID.Hex() // Convert ObjectID to string

	const duration time.Duration = 72 * time.Hour
	claims["exp"] = time.Now().Add(duration).Unix()

	// Generate encoded token and send it as response
	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusCreated,
		models.Response{
			Status:  "success",
			Message: "User logged in successfully",
			Data:    user,
		},
	)
}

func AddFollowing(c echo.Context) (err error) {
	userID := userIDFromToken(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.Response{
				Status:  "error",
				Message: "Invalid ID",
			},
		)
	}

	// Add a follower to user
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$addToSet", Value: bson.M{"followers": userID}}}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)

	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.Response{
				Status:  "error",
				Message: "ID not found",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.Response{
			Status:  "success",
			Message: "Followed successfully",
		},
	)
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["id"].(string)
}
