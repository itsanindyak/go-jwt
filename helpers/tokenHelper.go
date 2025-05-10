package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsanindyak/go-jwt/config"
)

type TokenInput struct {
	UID      string
	UserType string
	jwt.RegisteredClaims
}

type TokenOutput struct {
	SignedToken  string
	RefreshToken string
	Err          error
}

func GenerateToken(input TokenInput) (token TokenOutput) {

	claims := &TokenInput{
		UID:      input.UID,
		UserType: input.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(config.TOKEN_EXPIRY) * time.Hour)),
		},
	}
	refreshClaim := &TokenInput{
		UID:      input.UID,
		UserType: input.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(config.REFRESH_TOKEN_EXPIRY) * time.Hour)),
		},
	}

	token.SignedToken, token.Err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.TOKEN_KEY)

	if token.Err != nil {
		return token // Return early if access token signing fails
	}

	token.RefreshToken, token.Err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim).SignedString(config.REFRESH_TOKEN_KEY)

	if token.Err != nil {
		return token // Return early if refresh token signing fails
	}

	return token
}

func ParseToken(signedToken string) (tokenData TokenInput, msg string) {
	parse, err := jwt.ParseWithClaims(
		signedToken,
		&TokenInput{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.TOKEN_KEY), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := parse.Claims.(*TokenInput)

	if !ok {
		msg = fmt.Sprintln("the token is invalid")

		return 
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		// Token has expired
		msg = "Token is expired."
		return 
	}
	msg= ""
	return *claims,msg
}


