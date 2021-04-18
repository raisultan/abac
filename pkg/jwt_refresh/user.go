package jwt_refresh

type userJWTRefreshRequest struct {
	Refresh string `json:"refresh" validate:"required"`
}

type userJWTRefreshResponse struct {
	Access string `json:"access"`
}
