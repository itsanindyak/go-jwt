package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/itsanindyak/go-jwt/database"
	"github.com/itsanindyak/go-jwt/models"
	"github.com/itsanindyak/go-jwt/pkg/helpers"
	"github.com/itsanindyak/go-jwt/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var User *mongo.Collection = database.GetCollection(context.Background(), "user")

var OTP *mongo.Collection = database.GetCollection(context.Background(), "otp")

func init() {
	indexModel := mongo.IndexModel{
		Keys: bson.M{"expireAt": 1},
		Options: options.Index().
			SetExpireAfterSeconds(0),
	}

	_, err := OTP.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logger.Error("⚠️ TTL index creation failed: " + err.Error())
	} else {
		logger.Success("✅ TTL index ensured")
	}
}

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
				{Key: "refresh_token", Value: 0},
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
		reqID := c.Param("user_id")
		if reqID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in params"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(reqID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err = User.DeleteOne(ctx, bson.M{"_id": userID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "user delete Succesffully"})

	}
}

func UpdateUserData() gin.HandlerFunc {
	return func(c *gin.Context) {

		reqService := c.Param("service")

		userID, err := primitive.ObjectIDFromHex(c.Param("user_id"))
		if err != nil || userID.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing User ID"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var updateData helpers.UpdateData

		
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, helpers.Response{StatusCode: http.StatusBadRequest, MessageType: helpers.Error, Message: "Invalid request."})
			return
		}

		updateFields := bson.M{"updated_at": time.Now()}
		var validationErr error

		switch reqService {
		case "name":
			if updateData.FirstName == "" && updateData.LastName == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "At least one of first_name or last_name is required"})
				return
			}
			if validationErr = validate.StructPartial(&updateData, "FirstName", "LastName"); validationErr != nil {
				break
			} else {
				if updateData.FirstName != "" {
					updateFields["first_name"] = updateData.FirstName
				}
				if updateData.LastName != "" {
					updateFields["last_name"] = updateData.LastName
				}
			}
		case "email":
			if updateData.Email == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
				return
			}
			if validationErr = validate.StructPartial(&updateData, "Email"); validationErr != nil {
				break
			} else {
				if updateData.Email != "" {
					updateFields["email"] = updateData.Email
				}
			}
		case "type":
			if updateData.UserType == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User type is required"})
				return
			}
			if validationErr = validate.StructPartial(&updateData, "UserType"); validationErr != nil {
				break
			} else {
				if updateData.UserType != "" {
					updateFields["user_type"] = updateData.UserType
				}
			}
		default:
			
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown update service"})
			return
		}

		if validationErr != nil {
			
			var errors []string
			if verrs, ok := validationErr.(validator.ValidationErrors); ok {
				for _, err := range verrs {
					errors = append(errors, err.Field()+" is invalid: "+err.Tag())
				}
			} else {
				errors = append(errors, validationErr.Error())
			}
			
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errors, ", ")})
			return
		}

		result, err := User.UpdateByID(ctx, userID, bson.M{"$set": updateFields})

		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.Response{StatusCode: http.StatusInternalServerError, MessageType: helpers.Error, Message: "Failed to update user"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, helpers.Response{StatusCode: http.StatusNotFound, MessageType: helpers.Error, Message: "No user found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

	}
}

// func UpdateUserVerify() gin.HandlerFunc{
// 	return  func(c *gin.Context) {

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
// 		defer cancel()

// 		userID, err := primitive.ObjectIDFromHex(c.Param("user_id"))
// 		if err != nil || userID.IsZero() {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		var reqOtp helpers.OTPReq

// 		if err := c.BindJSON(&reqOtp); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		if err := validate.Struct(&reqOtp); err != nil {
// 			var errors []string
// 			for _, err := range err.(validator.ValidationErrors) {
// 				errors = append(errors, err.Field()+" is invalid: "+err.Tag())
// 			}
// 			c.JSON(http.StatusBadRequest, gin.H{"Error": strings.Join(errors, ", ")})
// 			return
// 		}

		
// 		filter := bson.M{
// 			"user_id": userID,
// 			"used":    false,
// 			"expireAt": bson.M{
// 				"$gt": time.Now(),
// 			},
// 		}

// 		opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}) // latest one

// 		err = OTP.FindOne(ctx, filter, opts).Decode(&otp)

// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not found or expired"})
// 			return
// 		}

// 		// Check OTP value
// 		if otp.Otp != reqOtp.OTP {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
// 			return
// 		}

// 		// Mark OTP as used
// 		_, err = OTP.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"used": true}})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP status"})
// 			return
// 		}

// 		// Update user as verified
// 		_, err = User.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"verified": true, "updated_at": time.Now()}})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully.please login again."})


// 	}
// }
