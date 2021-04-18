package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")
var UserUnauthorizedErr = errors.New("User is not authorized")
var AccessExpectedErr = errors.New("User is not authorized")

type tokenPayload struct {
	Email        string
	Expiration   uint64
	IsAuthorized bool
	Type         string
}

func validateAuth(r *http.Request) error {
	tp, err := extractTokenPayload(r)
	if err != nil {
		return err
	}
	if !tp.IsAuthorized {
		return UserUnauthorizedErr
	}
	if tp.Type != "access" {
		return AccessExpectedErr
	}

	return nil
}

func respondWithErrorMessage(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func extractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenStr := extractToken(r)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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
		return nil, err
	}

	return token, nil
}

func isTokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func extractTokenPayload(r *http.Request) (*tokenPayload, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return nil, err
		}

		exp, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if err != nil {
			return nil, err
		}

		isAuth, ok := claims["isAuthorized"].(bool)
		if !ok {
			return nil, err
		}

		tType, ok := claims["type"].(string)
		if !ok {
			return nil, err
		}

		tp := tokenPayload{
			Email:        email,
			Expiration:   exp,
			IsAuthorized: isAuth,
			Type:         tType,
		}
		return &tp, nil
	}

	return nil, err
}
