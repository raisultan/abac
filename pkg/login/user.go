package login

type userLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type userLoginJWTResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
