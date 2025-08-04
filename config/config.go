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
)

func init() {
	var err error

	// Load .env only if not in production
	if ENV = os.Getenv("ENV"); ENV != "production" {
		ENV = "local"
		err = godotenv.Load(".env")
		if err != nil {
			logger.Error("Could not load .env file, continuing with system environment variables")
			return
		}
		logger.Success("📦 Running in " + ENV + " environment")
	}

	// port

	PORT = os.Getenv("PORT")

	// Mongodb url

	MONGO_URI = os.Getenv("MONGO_URI")

	// Token key

	TOKEN_KEY = []byte(os.Getenv("TOKEN_KEY"))

	// Token Expire

	TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("TOKEN_EXPIRY"))

	if err != nil {
		logger.Fatal("❌ TOKEN_EXPIRY is not a valid integer")
	}

	// Token key

	REFRESH_TOKEN_KEY = []byte(os.Getenv("REFRESH_TOKEN_KEY"))

	// Token Expire

	REFRESH_TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))

	if err != nil {
		logger.Fatal("❌ REFRESH_TOKEN_EXPIRY is not a valid integer")
	}

	OTP_EXPIRY, err = strconv.Atoi(os.Getenv("OTP_EXPIRY"))

	if err != nil {
		logger.Fatal("❌ OTP_EXPIRY is not a valid integer")
	}

}
