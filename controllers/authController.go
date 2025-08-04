package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/helpers"
	"github.com/itsanindyak/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		var existingUser models.User
		var newUser models.User
		var otp models.OTP

		// Parse and bind JSON input
		if err := c.BindJSON(&signupdata); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(&signupdata); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid: "+err.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"Error": strings.Join(errors, ", ")})
			return
		}

		err := User.FindOne(ctx, bson.M{"email": signupdata.Email}).Decode(existingUser)

		if err == nil {
			if existingUser.Verified {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered and verified."})
				return
			} else {
				userOtp, _ := helpers.GenerateOTP(6)
				otp = models.OTP{
					UserID: newUser.ID,
					Otp:    userOtp,
					Used:   false,
					TTL:    time.Now().Add(time.Duration(config.OTP_EXPIRY) * time.Second),
				}

				_, err := OTP.InsertOne(ctx, otp)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP"})
					return
				}

				// TODO: Send OTP to email (use SendGrid, SMTP, etc.)

				c.JSON(http.StatusOK, gin.H{
					"message": "Email already registered but not verified. OTP sent to your email.",
				})
				return
			}
		} else if err != mongo.ErrNoDocuments {
			// Real DB error
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
			return
		}

		hashPassword, err := hashPassword(&signupdata.Password)

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		newUser = models.User{
			ID:           primitive.NewObjectID(),
			FirstName:    signupdata.FirstName,
			LastName:     signupdata.LastName,
			Email:        signupdata.Email,
			Password:     hashPassword,
			UserType:     signupdata.UserType,
			Verified:     false,
			RefreshToken: "", // set below
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		token := helpers.GenerateToken(helpers.TokenInput{
			UID:      newUser.ID.Hex(),
			UserType: newUser.UserType,
		})
		newUser.RefreshToken = token.RefreshToken

		_, err = User.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
			return
		}
		
		//NOTE: otp model 
		userOtp, _ := helpers.GenerateOTP(6)
		otp = models.OTP{
			UserID: newUser.ID,
			Otp:    userOtp,
			Used:   false,
			TTL:    time.Now().Add(time.Duration(config.OTP_EXPIRY) * time.Second),
		}

		_, err = OTP.InsertOne(ctx, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP"})
			return
		}

		// c.SetCookie(
		// 	"Refreshtoken",     // Name
		// 	token.RefreshToken, // Value
		// 	3600,
		// 	"/",   // Path
		// 	"",    // Domain (empty for current domain)
		// 	false, // Secure (true for HTTPS, false for local HTTP)
		// 	true,  // HttpOnly
		// )
		// c.SetCookie(
		// 	"token",           // Name
		// 	token.SignedToken, // Value
		// 	3600,
		// 	"/",   // Path
		// 	"",    // Domain (empty for current domain)
		// 	false, // Secure (true for HTTPS, false for local HTTP)
		// 	true,  // HttpOnly
		// )

		c.JSON(http.StatusOK, gin.H{
			"Message": "OTP send to email.Please verify.",
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
