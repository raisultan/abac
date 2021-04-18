package update

type Service interface {
	UpdateUser(UserUpdateRequest) (UserRetrieveResponse, error)
}

type Repository interface {
	UpdateUser(UserUpdateRequest) (UserRetrieveResponse, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) UpdateUser(r UserUpdateRequest) (UserRetrieveResponse, error) {
	return s.r.UpdateUser(r)
}
