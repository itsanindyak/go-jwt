package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/models"
	"github.com/itsanindyak/go-jwt/pkg/helpers"
	"github.com/itsanindyak/go-jwt/pkg/logger"
	"github.com/itsanindyak/go-jwt/pkg/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		err := User.FindOne(ctx, bson.M{"email": signupdata.Email}).Decode(&existingUser)

		if err == nil {
			if existingUser.Verified {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered and verified."})
				return
			} else {
				userOtp, _ := helpers.GenerateOTP(6)
				otp = models.OTP{
					UserID: existingUser.ID,
					Otp:    userOtp,
					Used:   false,
					TTL:    time.Now().Add(time.Duration(config.OTP_EXPIRY) * time.Second),
				}

				_, err := OTP.InsertOne(ctx, otp)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP"})
					return
				}

				err = mail.Send(mail.MailData{Name: existingUser.FirstName + " " + existingUser.LastName, OTP: otp.Otp}, existingUser.Email)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send authentication mail"})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "Email already registered but not verified. OTP sent to your email.",
					"data": gin.H{
						"user_id": existingUser.ID.Hex(),
					},
				})
				return
			}
		} else if err != mongo.ErrNoDocuments {
			// Real DB error
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
			return
		}

		hashPassword, err := hashPassword(&signupdata.Password)

		if err != nil {
			logger.Error(err.Error())
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

		result, err := User.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
			return
		}

		insertedID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID returned from insert"})
			return
		}
		newUser.ID = insertedID

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

		err = mail.Send(mail.MailData{Name: newUser.FirstName + " " + newUser.LastName, OTP: otp.Otp}, newUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send authentication mail"})
			return
		}


		c.JSON(http.StatusOK, gin.H{
			"Message": "OTP send to email.Please verify.",
			"data": gin.H{
				"user_id": newUser.ID.Hex(),
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
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// Check verified
		if !foundUser.Verified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email before logging in."})
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

		update := bson.M{"$set": bson.M{"refresh_token": token.RefreshToken, "updated_at": time.Now()}}
		_, err = User.UpdateOne(ctx, bson.M{"_id": foundUser.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
			return
		}

		secure := os.Getenv("ENV") == "production"

		c.SetCookie(
			"Refreshtoken",     // Name
			token.RefreshToken, // Value
			3600,
			"/",    // Path
			"",     // Domain (empty for current domain)
			secure, // Secure (true for HTTPS, false for local HTTP)
			true,   // HttpOnly
		)

		c.SetCookie(
			"token",           // Name
			token.SignedToken, // Value
			3600,
			"/",    // Path
			"",     // Domain (empty for current domain)
			secure, // Secure (true for HTTPS, false for local HTTP)
			true,   // HttpOnly
		)

		c.JSON(http.StatusOK, gin.H{
			// "token":   token.SignedToken,
			"Message": "Login succesfully.",
			"data": gin.H{
				"id":        foundUser.ID,
				"name":      foundUser.FirstName + foundUser.LastName,
				"email":     foundUser.Email,
				"user_type": foundUser.UserType,
			},
		})

	}
}

func VerifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		reqID := c.Param("id")
		if reqID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId not find in url."})
			return
		}

		var reqOtp helpers.OTPReq
		var otp models.OTP

		if err := c.BindJSON(&reqOtp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(&reqOtp); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid: "+err.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"Error": strings.Join(errors, ", ")})
			return
		}

		userID, err := primitive.ObjectIDFromHex(reqID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{
			"user_id": userID,
			"used":    false,
			"expireAt": bson.M{
				"$gt": time.Now(),
			},
		}

		opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}) // latest one

		err = OTP.FindOne(ctx, filter, opts).Decode(&otp)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not found or expired"})
			return
		}

		// Check OTP value
		if otp.Otp != reqOtp.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		// Mark OTP as used
		_, err = OTP.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"used": true}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP status"})
			return
		}

		// Update user as verified
		_, err = User.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"verified": true, "updated_at": time.Now()}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully.please login again."})
	}
}
