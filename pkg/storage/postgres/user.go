package postgres

type User struct {
	ID int `json:"id"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	CreatedAt  string `json:"createdAt"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}
