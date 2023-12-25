package main
import (
	"time"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var JWT_SECRET_KEY string
func init(){
		if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	if JWT_SECRET_KEY == "" {
		log.Fatal("MONGO_URI not set")
	}
}

func generateJwt(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
		Subject: username,
	})
	log.Println(JWT_SECRET_KEY)
	tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
	return tokenString, err
}