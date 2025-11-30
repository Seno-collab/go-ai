package authapp

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpresIn     int    `json:"expires_in"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
type RegisterSuccess struct {
}

type GetProfileRequest struct {
}

type GetProfileResponse struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}
