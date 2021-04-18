package login

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginJWTResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
