package models

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"string"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	CreatedAt  string `json:"createdAt"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}
