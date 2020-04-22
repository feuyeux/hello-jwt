package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	/* use base64 to generate secret string:
	 * ▶ input=jwt_test_20200422_at
	 * ▶ echo -n $input | openssl base64
	 * and0X3Rlc3RfMjAyMDA0MjJfYXQ=
	 * ▶ input=jwt_test_20200422_rt
	 * ▶ echo -n $input | openssl base64
	 * and0X3Rlc3RfMjAyMDA0MjJfcnQ=
	 */
	err := godotenv.Load("hello_jwt.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("[Config] Initialized!\nACCESS_SECRET=%s\nREFRESH_SECRET=%s\nREDIS_DSN=%s",
		os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"), os.Getenv("REDIS_DSN"))
}

func Get(key string) string {
	return os.Getenv(key)
}
