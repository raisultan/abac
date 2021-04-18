package list

type Service interface {
	ListUsers(start, count int) ([]UserRetrieveResponse, error)
}

type Repository interface {
	GetAllUsers(start, count int) ([]UserRetrieveResponse, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ListUsers(start, count int) ([]UserRetrieveResponse, error) {
	return s.r.GetAllUsers(start, count)
}
