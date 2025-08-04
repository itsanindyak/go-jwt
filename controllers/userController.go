package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/database"
	"github.com/itsanindyak/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var User *mongo.Collection = database.GetCollection(context.Background(), "user")

var OTP *mongo.Collection = database.GetCollection(context.Background(),"otp")

func GetUser() gin.HandlerFunc {

	return func(c *gin.Context) {
		userId := c.Param("user_id")
		objId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User
		err = User.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":        user.ID,
				"name":      user.FirstName + user.LastName,
				"email":     user.Email,
				"user_type": user.UserType,
			},
		})

	}

}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("rp"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("p"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex, err := strconv.Atoi(c.Query("si"))
		if err != nil || startIndex < 0 {
			startIndex = (page - 1) * recordPerPage
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectBeforeGroup := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "password", Value: 0},
				{Key: "refreshtoken", Value: 0},
			}},
		}

		groupStage := bson.D{
			{
				Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
					{Key: "data", Value: bson.D{
						{Key: "$push", Value: "$$ROOT"},
					}},
				},
			},
		}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "user_items", Value: bson.D{
						{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
					}},
				},
			},
		}

		result, err := User.Aggregate(ctx, mongo.Pipeline{matchStage, projectBeforeGroup, groupStage, projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
			return
		}

		var allUser []bson.M

		if err = result.All(ctx, &allUser); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUser)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in params"})
			return
		}

		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
