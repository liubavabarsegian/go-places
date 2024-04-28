package auth

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
)

var SecretKey = []byte("SecretYouShouldHide")

type GenerateJWTResponse struct {
	Token string `json:"token"`
}

func GetToken(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := generateJWT()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseParams := GenerateJWTResponse{token}
		responseOK(w, r, &responseParams)
	}
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func responseOK(w http.ResponseWriter, r *http.Request, responseParams *GenerateJWTResponse) {
	render.JSON(w, r, GenerateJWTResponse{
		Token: responseParams.Token,
	})
}
