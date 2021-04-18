package register

type UserRegisterRequest struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`

	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type UserRegisterResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
