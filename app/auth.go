package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type fieldValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type validationError struct {
	Details []fieldValidationError `json:"details"`
}

type userLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type userJWTRefreshRequest struct {
	Refresh string `json:"refresh" validate:"required"`
}

type userJWTRefreshResponse struct {
	Access string `json:"access"`
}

func (r *userJWTRefreshResponse) refresh(a string) error {
	claims := &jwtClaims{}
	t, err := jwt.ParseWithClaims(a, claims, func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("Access token is invalid")
	}

	jwtExpirationMins := 5
	expirationTime := time.Now().Add(time.Duration(jwtExpirationMins) * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		return err
	}

	r.Access = tokenStr
	return nil
}

type tokenPayload struct {
	Email        string
	Expiration   uint64
	IsAuthorized bool
	Type         string
}

type userLoginJWTResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (r *userLoginJWTResponse) generate(uEmail string) error {
	jwtExpirationMins := 5
	expirationTime := time.Now().Add(time.Duration(jwtExpirationMins) * time.Minute)

	claims := jwtClaims{
		Email: uEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	jwtRefreshExpirationMins := 30
	refreshExpirationTime := time.Now().Add(time.Duration(jwtRefreshExpirationMins) * time.Minute)

	refreshClaims := jwtClaims{
		Email: uEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(jwtKey)

	if err != nil {
		return err
	}

	r.Access = tokenStr
	r.Refresh = refreshTokenStr
	return nil
}

type jwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type userRegisterRequest struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`

	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type userRegisterResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *userRegisterRequest) register(db *sql.DB) error {
	err := u.setPassword()

	if err != nil {
		return err
	}

	err = db.QueryRow(
		"INSERT INTO users(email, password, firstName, lastName) VALUES($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.Password,
		u.FirstName,
		u.LastName,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *userLoginRequest) getUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email, &u.Password)
}

func (u *userRegisterRequest) checkUserExists(db *sql.DB) (bool, error) {
	err := db.QueryRow(
		"SELECT email FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email)

	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (u *userRegisterRequest) setPassword() error {
	hashingCost := 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashingCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
