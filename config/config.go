package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
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
	if os.Getenv("ENV") != "production" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("⚠️ Could not load .env file, continuing with system environment variables")
		}
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
		log.Fatal("❌ TOKEN_EXPIRY is not a valid integer")
	}

	// Token key

	REFRESH_TOKEN_KEY = []byte(os.Getenv("REFRESH_TOKEN_KEY"))

	// Token Expire

	REFRESH_TOKEN_EXPIRY, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))

	if err != nil {
		log.Fatal("❌ REFRESH_TOKEN_EXPIRY is not a valid integer")
	}

	OTP_EXPIRY, err = strconv.Atoi(os.Getenv("OTP_EXPIRY"))

	if err != nil {
		log.Fatal("❌ OTP_EXPIRY is not a valid integer")
	}

}
