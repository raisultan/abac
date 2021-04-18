package login

import (
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
	if err != nil {
		return UserLoginJWTResponse{}, err
	}

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
