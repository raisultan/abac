package delete

type Service interface {
	DeleteUser(int) error
}

type Repository interface {
	DeleteUser(int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) DeleteUser(id int) error {
	return s.r.DeleteUser(id)
}
