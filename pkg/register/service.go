package register

import "errors"

var ErrDuplicate = errors.New("User already exists")

type Service interface {
	RegisterUser(UserRegisterRequest) (UserRegisterResponse, error)
}

type Repository interface {
	RegisterUser(UserRegisterRequest) (UserRegisterResponse, error)
	CheckUserExists(UserRegisterRequest) (bool, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterUser(ur UserRegisterRequest) (UserRegisterResponse, error) {
	exists, err := s.r.CheckUserExists(ur)
	if err != nil {
		return UserRegisterResponse{}, err
	}
	if exists {
		return UserRegisterResponse{}, ErrDuplicate
	}

	u, err := s.r.RegisterUser(ur)
	if err != nil {
		return UserRegisterResponse{}, err
	}

	return u, nil
}
