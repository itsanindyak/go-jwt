package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/itsanindyak/go-jwt/helpers"
	"github.com/itsanindyak/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func hashPassword(password *string) (hashPassword string, err error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)

	if err != nil {

		return "Encryption failed", err
	}

	return string(byteHash), nil

}

func verifyPassword(hashPassword string, userPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(userPassword))
	if err != nil {

		return false, "Password does not match"
	}

	return true, "Password matches"

}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var signupdata helpers.Signup
		var user models.User

		// Parse and bind JSON input
		if err := c.BindJSON(&signupdata); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// // Check for required fields
		// if signupdata.Email == "" || signupdata.Password == "" || signupdata.FirstName == "" || signupdata.LastName == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email and password are required"})
		// 	return
		// }

		if err := validate.Struct(&signupdata); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid: "+err.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"Error": strings.Join(errors, ", ")})
			return
		}

		count, err := User.CountDocuments(ctx, bson.M{"email": signupdata.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return

		}

		user.Password, err = hashPassword(&signupdata.Password)

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.FirstName = signupdata.FirstName
		user.LastName = signupdata.LastName
		user.Email = signupdata.Email
		user.UserType = signupdata.UserType
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()
		token := helpers.GenerateToken(helpers.TokenInput{
			UID:      user.ID.Hex(),
			UserType: signupdata.UserType,
		})

		user.RefreshToken = token.RefreshToken
		result, err := User.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"User create error": err.Error()})
			return
		}
		c.SetCookie(
			"Refreshtoken",     // Name
			token.RefreshToken, // Value
			3600,
			"/",   // Path
			"",    // Domain (empty for current domain)
			false, // Secure (true for HTTPS, false for local HTTP)
			true,  // HttpOnly
		)
		c.SetCookie(
			"token",           // Name
			token.SignedToken, // Value
			3600,
			"/",   // Path
			"",    // Domain (empty for current domain)
			false, // Secure (true for HTTPS, false for local HTTP)
			true,  // HttpOnly
		)

		c.JSON(http.StatusOK, gin.H{
			// "token":   token.SignedToken,
			"Message": "User is added succesfully.",
			"data": gin.H{
				"id":        result.InsertedID,
				"name":      user.FirstName + user.LastName,
				"email":     user.Email,
				"user_type": user.UserType,
			},
		})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var loginData helpers.Login
		var foundUser models.User

		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(&loginData); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid: "+err.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"Error": strings.Join(errors, ", ")})
			return
		}

		err := User.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "Message": "User not found with this email"})
			return
		}
		isMatched, msg := verifyPassword(foundUser.Password, loginData.Password)

		if !isMatched {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		token := helpers.GenerateToken(helpers.TokenInput{
			UID:      foundUser.ID.Hex(),
			UserType: foundUser.UserType,
		})

		update := bson.M{"$set": bson.M{"refreshtoken": token.RefreshToken, "updatedat": time.Now()}}
		_, err = User.UpdateOne(ctx, bson.M{"_id": foundUser.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
			return
		}

		c.SetCookie(
			"Refreshtoken",     // Name
			token.RefreshToken, // Value
			3600,
			"/",   // Path
			"",    // Domain (empty for current domain)
			false, // Secure (true for HTTPS, false for local HTTP)
			true,  // HttpOnly
		)

		c.SetCookie(
			"token",           // Name
			token.SignedToken, // Value
			3600,
			"/",   // Path
			"",    // Domain (empty for current domain)
			false, // Secure (true for HTTPS, false for local HTTP)
			true,  // HttpOnly
		)

		c.JSON(http.StatusOK, gin.H{
			// "token":   token.SignedToken,
			"Message": "User is added succesfully.",
			"data": gin.H{
				"id":        foundUser.ID,
				"name":      foundUser.FirstName + foundUser.LastName,
				"email":     foundUser.Email,
				"user_type": foundUser.UserType,
			},
		})

	}
}
