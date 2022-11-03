package controllers

import (
	"net/http"
	"time"

	"github.com/ImBIOS/go-micho-twitter/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func AddTweet(c echo.Context) (err error) {
	const timeout time.Duration = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	userID := userIDFromToken(c)

	var tweet models.Tweet

	defer cancel()

	// Validate the request body type
	if err = c.Bind(&tweet); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{Status: "error", Message: "Bad request", Data: err.Error()},
		)
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&tweet); validationErr != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.Response{
				Status:  "error",
				Message: "Validation error",
				Data:    validationErr.Error(),
			},
		)
	}

	newTweet := models.Tweet{
		ID:          primitive.NewObjectID(),
		From:        userID,
		FullText:    tweet.FullText,
		CreatedTime: primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	// Save tweet in database
	_, err = tweetCollection.InsertOne(ctx, newTweet)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusCreated,
		models.Response{
			Status:  "success",
			Message: "Tweet created successfully",
			Data:    newTweet,
		},
	)
}

func GetFeed(c echo.Context) (err error) {
	const timeout time.Duration = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	userID := userIDFromToken(c)

	defer cancel()

	// Get feed from database
	cursor, err := tweetCollection.Find(ctx, bson.M{"from": userID})
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    err.Error(),
			},
		)
	}

	var tweets []models.Tweet
	if err = cursor.All(ctx, &tweets); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.Response{
			Status:  "success",
			Message: "Tweets retrieved successfully",
			Data:    tweets,
		},
	)
}

// 	// Save post in database
// 	if err = db.DB("twitter").C("posts").Insert(p); err != nil {
// 		return
// 	}
// 	return c.JSON(http.StatusCreated, p)
// }

// func (h *Handler) FetchPost(c echo.Context) (err error) {
// 	userID := userIDFromToken(c)
// 	page, _ := strconv.Atoi(c.QueryParam("page"))
// 	limit, _ := strconv.Atoi(c.QueryParam("limit"))

// 	// Defaults
// 	if page == 0 {
// 		page = 1
// 	}
// 	if limit == 0 {
// 		limit = 100
// 	}

// 	// Retrieve posts from database
// 	posts := []*models.Post{}
// 	db := h.DB.Clone()
// 	if err = db.DB("twitter").C("posts").
// 		Find(bson.M{"to": userID}).
// 		Skip((page - 1) * limit).
// 		Limit(limit).
// 		All(&posts); err != nil {
// 		return
// 	}
// 	defer db.Close()

// 	return c.JSON(http.StatusOK, posts)
// }
