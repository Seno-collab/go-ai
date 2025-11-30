package app

import (
	authapp "go-ai/internal/application/auth"
	restaurantapp "go-ai/internal/application/restaurant"
	uploadapp "go-ai/internal/application/upload"
	"go-ai/internal/transport/http/response"
)

type ErrorResponseDoc struct {
	ResponseCode string                `json:"response_code,omitempty"`
	Message      string                `json:"message"`
	Error        *response.ErrorDetail `json:"error,omitempty"`
}

type RegisterSuccessResponseDoc struct {
	Message      string                   `json:"message"`
	ResponseCode string                   `json:"response_code,omitempty"`
	Data         *authapp.RegisterSuccess `json:"data,omitempty"`
}

type GetProfileSuccessResponseDoc struct {
	Message      string                      `json:"message"`
	ResponseCode string                      `json:"response_code,omitempty"`
	Data         *authapp.GetProfileResponse `json:"data,omitempty"`
}

type RefreshTokenSuccessResponseDoc struct {
	Message      string                        `json:"message"`
	ResponseCode string                        `json:"response_code,omitempty"`
	Data         *authapp.RefreshTokenResponse `json:"data,omitempty"`
}

type LoginSuccessResponseDoc struct {
	Message      string                 `json:"message"`
	ResponseCode string                 `json:"response_code,omitempty"`
	Data         *authapp.LoginResponse `json:"data,omitempty"`
}

type UploadLogoSuccessResponseDoc struct {
	Message      string                        `json:"message"`
	ResponseCode string                        `json:"response_code,omitempty"`
	Data         *uploadapp.UploadLogoResponse `json:"data,omitempty"`
}

type CreateRestaurantSuccessResponseDoc struct {
	Message      string                                 `json:"message"`
	ResponseCode string                                 `json:"response_code,omitempty"`
	Data         *restaurantapp.CreateRestaurantRequest `json:"data,omitempty"`
}
