package login

type Service interface {
	LoginUser(userLoginRequest) (userLoginJWTResponse, error)
}

type Repository interface {
	LoginUser(userLoginRequest) (userLoginJWTResponse, error)
	GetUserByEmail(userLoginRequest) (userLoginRequest, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) LoginUser(ulr userLoginRequest) (userLoginJWTResponse, error) {
	u, err := s.r.GetUserByEmail(ulr)
	if err != nil {
		return userLoginJWTResponse{}, err
	}

	uJWT, err := s.r.LoginUser(u)
	if err != nil {
		return userLoginJWTResponse{}, err
	}

	return uJWT, nil
}
