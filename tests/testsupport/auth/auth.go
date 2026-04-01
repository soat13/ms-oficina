package testauth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "test-user",
		"roles": []string{"manager"},
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func AddAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}
