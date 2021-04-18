package login

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	LoginUser(UserLoginRequest) (UserLoginJWTResponse, error)
}

type Repository interface {
	GetUserByEmail(UserLoginRequest) (UserLoginRequest, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) LoginUser(ulr UserLoginRequest) (UserLoginJWTResponse, error) {
	var err error
	u := UserLoginRequest{Email: ulr.Email}
	u, err = s.r.GetUserByEmail(ulr)
	if err != nil {
		return UserLoginJWTResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(ulr.Password))
	at, err := createAccessToken(ulr.Email)
	if err != nil {
		return UserLoginJWTResponse{}, err
	}

	rt, err := createRefreshToken(ulr.Email)
	if err != nil {
		return UserLoginJWTResponse{}, err
	}

	uJWT := UserLoginJWTResponse{Access: at, Refresh: rt}
	return uJWT, nil
}

// TODO: transfer to internal
func createAccessToken(email string) (string, error) {
	// TODO: make global
	var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")

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
	// TODO: make global
	var jwtKey = []byte("23a93f6d2673b09c1d3d063cf7a97fc20d0054dfce65c3737455dbec25439938")

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
