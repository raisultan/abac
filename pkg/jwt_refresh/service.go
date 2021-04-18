package jwt_refresh

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var RefreshExpectedErr = errors.New("Refresh token is expected")
var InvalidRefreshErr = errors.New("Invalid refresh token received")

type Service interface {
	RefreshJWT(userJWTRefreshRequest) (userJWTRefreshResponse, error)
}

type Repository interface{}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RefreshJWT(r userJWTRefreshRequest) (userJWTRefreshResponse, error) {
	// TODO: make global
	var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")

	token, err := jwt.Parse(r.Refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(
				fmt.Sprintf(
					"unexpected signing method: %v",
					token.Header["alg"],
				),
			)
		}
		return jwtKey, nil
	})

	if err != nil {
		return userJWTRefreshResponse{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return userJWTRefreshResponse{}, err
		}

		tType, ok := claims["type"].(string)
		if !ok {
			return userJWTRefreshResponse{}, err
		}
		if tType != "refresh" {
			return userJWTRefreshResponse{}, RefreshExpectedErr
		}

		atClaims := jwt.MapClaims{}
		atClaims["isAuthorized"] = true
		atClaims["email"] = email
		atClaims["type"] = "access"
		atClaims["exp"] = time.Now().Add(5 * time.Minute).Unix()

		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		atStr, err := at.SignedString(jwtKey)
		if err != nil {
			return userJWTRefreshResponse{}, err
		}

		return userJWTRefreshResponse{Access: atStr}, nil
	}

	return userJWTRefreshResponse{}, InvalidRefreshErr
}
