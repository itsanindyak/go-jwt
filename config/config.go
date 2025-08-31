package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/itsanindyak/go-jwt/pkg/logger"
	"github.com/joho/godotenv"
)

var (
	ENV                  string
	PORT                 string
	MONGO_URI            string
	TOKEN_KEY            []byte
	TOKEN_EXPIRY         int
	REFRESH_TOKEN_KEY    []byte
	REFRESH_TOKEN_EXPIRY int
	OTP_EXPIRY           int
	BASE_URL             string

	//-----------------
	// SendGrid Mail
	//-----------------
	SENDGRID_API_KEY string
	SMTP_SENDER      string
	SMTP_NAME        string

	//-----------------
	// Google OAuth
	//-----------------
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string

	//-----------------
	// GitHub OAuth
	//-----------------
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	GITHUB_REDIRECT_URL  string
)

func init() {
	var err error

	//-----------------
	// Environment
	//-----------------
	if ENV = os.Getenv("ENV"); ENV != "production" {
		ENV = "development"
		err = godotenv.Load(".env")
		if err != nil {
			logger.Error("Could not load .env file, continuing with system environment variables")
			return
		}
	}

	PORT = os.Getenv("PORT")

	base := os.Getenv("BASE_URL")
	if base == "" {
		logger.Fatal("❌ BASE_URL is missing in environment")
	}

	// add port only in development
	if ENV == "development" {
		BASE_URL = fmt.Sprintf("%s:%s", base, PORT)
	} else {
		BASE_URL = base
	}

	//-----------------
	// MongoDB
	//-----------------
	MONGO_URI = os.Getenv("MONGO_URI")

	//-----------------
	// Access Token
	//-----------------
	TOKEN_KEY = []byte(os.Getenv("TOKEN_KEY"))
	TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("TOKEN_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ TOKEN_EXPIRY is not a valid integer")
	}

	//-----------------
	// Refresh Token
	//-----------------
	REFRESH_TOKEN_KEY = []byte(os.Getenv("REFRESH_TOKEN_KEY"))
	REFRESH_TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ REFRESH_TOKEN_EXPIRY is not a valid integer")
	}

	//-----------------
	// OTP Expiry
	//-----------------
	OTP_EXPIRY, err = strconv.Atoi(os.Getenv("OTP_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ OTP_EXPIRY is not a valid integer")
	}

	//-----------------
	// SendGrid
	//-----------------
	SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
	SMTP_SENDER = os.Getenv("SMTP_SENDER")
	SMTP_NAME = os.Getenv("SMTP_NAME")

	if SENDGRID_API_KEY == "" || SMTP_SENDER == "" {
		logger.Fatal("❌ SENDGRID_API_KEY or SMTP_SENDER is missing in environment")
	}

	//-----------------
	// Google OAuth
	//-----------------
	GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	GOOGLE_REDIRECT_URL = BASE_URL + "/auth/google/callback"

	if GOOGLE_CLIENT_ID == "" || GOOGLE_CLIENT_SECRET == "" {
		logger.Fatal("❌ Google OAuth credentials are missing in environment")
	}

	//-----------------
	// GitHub OAuth
	//-----------------
	GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
	GITHUB_REDIRECT_URL = BASE_URL + "/auth/github/callback"

	if GITHUB_CLIENT_ID == "" || GITHUB_CLIENT_SECRET == "" {
		logger.Fatal("❌ GitHub OAuth credentials are missing in environment")
	}
}
