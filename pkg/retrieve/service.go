package retrieve

type Service interface {
	RetrieveUser(int) (UserRetrieveResponse, error)
}

type Repository interface {
	GetUserById(int) (UserRetrieveResponse, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RetrieveUser(id int) (UserRetrieveResponse, error) {
	return s.r.GetUserById(id)
}
