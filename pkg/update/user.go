package update

type UserUpdateRequest struct {
	ID int `json:"-"`

	Email      string `json:"-"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	IsAdmin    bool   `json:"-"`
	IsApproved bool   `json:"-"`
}

type UserRetrieveResponse struct {
	ID int `json:"id"`

	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}
