package list

type UserRetrieveResponse struct {
	ID int `json:"id"`

	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}
