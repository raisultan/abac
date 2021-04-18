package login

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")

// TODO: transfer to internal
func createAccessToken(email string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["isAuthorized"] = true
	atClaims["email"] = email
	atClaims["type"] = "access"
	atClaims["exp"] = time.Now().Add(5 * time.Minute).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createRefreshToken(email string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["isAuthorized"] = true
	atClaims["email"] = email
	atClaims["type"] = "refresh"
	atClaims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
