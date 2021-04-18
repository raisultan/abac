package jwt_refresh

type UserJWTRefreshRequest struct {
	Refresh string `json:"refresh" validate:"required"`
}

type UserJWTRefreshResponse struct {
	Access string `json:"access"`
}
