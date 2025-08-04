package config

import (
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
	SENDGRID_API_KEY string
	SMTP_SENDER      string
	SMTP_NAME        string
)

func init() {
	var err error

	// Load .env only if not in production
	if ENV = os.Getenv("ENV"); ENV != "production" {
		ENV = "development"
		err = godotenv.Load(".env")
		if err != nil {
			logger.Error("Could not load .env file, continuing with system environment variables")
			return
		}
		logger.Success("📦 Running in " + ENV + " environment")
	}

	PORT = os.Getenv("PORT")
	
	MONGO_URI = os.Getenv("MONGO_URI")
	TOKEN_KEY = []byte(os.Getenv("TOKEN_KEY"))

	TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("TOKEN_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ TOKEN_EXPIRY is not a valid integer")
	}

	REFRESH_TOKEN_KEY = []byte(os.Getenv("REFRESH_TOKEN_KEY"))
	REFRESH_TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ REFRESH_TOKEN_EXPIRY is not a valid integer")
	}

	OTP_EXPIRY, err = strconv.Atoi(os.Getenv("OTP_EXPIRY"))
	if err != nil {
		logger.Fatal("❌ OTP_EXPIRY is not a valid integer")
	}

	SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
	SMTP_SENDER = os.Getenv("SMTP_SENDER")
	SMTP_NAME = os.Getenv("SMTP_NAME")

	if SENDGRID_API_KEY == "" || SMTP_SENDER == "" {
		logger.Fatal("❌ SENDGRID_API_KEY or SMTP_SENDER is missing in environment")
	}
}
