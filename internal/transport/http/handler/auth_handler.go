package handler

import (
	"errors"
	auth "go-ai/internal/domain/auth"
	auth_service "go-ai/internal/service/auth"
	"go-ai/pkg/common"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	svc *auth_service.Service
}

func NewAuthHandler(svc *auth_service.Service) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

type ErrorResponseDoc struct {
	ResponseCode string              `json:"response_code,omitempty"`
	Message      string              `json:"message"`
	Error        *common.ErrorDetail `json:"error,omitempty"`
}
type RegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
type RegisterSuccess struct {
}

type RegisterSuccessResponse struct {
	Message      string           `json:"message"`
	ResponseCode string           `json:"response_code,omitempty"`
	Data         *RegisterSuccess `json:"data,omitempty"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and full name
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterReq true "User registration payload"
// @Success 200 {object} RegisterSuccessResponse "User created successfully"
// @Failure default {object} ErrorResponseDoc "Invalid input"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	logger := zerolog.Ctx(c.Request().Context())
	var in RegisterReq
	if err := c.Bind(&in); err != nil {
		return common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}
	_, err := h.svc.Register(&auth.Entity{Email: in.Email, FullName: in.FullName, Password: in.Password})
	if err != nil {
		logger.Error().Err(err).Msg("failed to register user")
		switch err {
		case common.ErrInvalidEmail, common.ErrInvalidName, common.ErrInvalidPassword, common.ErrUserAlreadyExists, common.ErrNameAlreadyExists:
			details := common.ErrorDetail{}
			if errors.Is(err, common.ErrInvalidEmail) {
				details = common.ErrorDetail{
					Field:   "email",
					Message: "Email is a required field",
				}
			}
			if errors.Is(err, common.ErrInvalidPassword) {
				details = common.ErrorDetail{
					Field:   "password",
					Message: "Password is a required field",
				}
			}
			if errors.Is(err, common.ErrInvalidName) {
				details = common.ErrorDetail{
					Field:   "name",
					Message: "Name is a required field",
				}
			}
			return common.ErrorResponse(c, http.StatusBadRequest, err.Error(), details)
		case common.ErrConflict:
			return common.ErrorResponse(c, http.StatusConflict, err.Error())
		default:
			return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	return common.SuccessResponse[any](c, nil, "Create user success")
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}
type ResponseLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpresIn     int    `json:"expires_in"`
}
type LoginSuccessResponse struct {
	Message      string         `json:"message"`
	ResponseCode string         `json:"response_code,omitempty"`
	Data         *ResponseLogin `json:"data,omitempty"`
}

// Login godoc
// @Summary User login 
// @Description Authenticate user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RequestLogin true "User login payload"
// @Success 200 {object} LoginSuccessResponse "Login successful"
// @Failure default {object} ErrorResponseDoc "Invalid credentials"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	logger := zerolog.Ctx(c.Request().Context())
	var in RequestLogin
	if err := c.Bind(&in); err != nil {
		return common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}
	accessToken, refreshToken, expiresIn, err := h.svc.Login(&auth.Entity{Email: in.Email, Password: in.Password})
	if err != nil {
		logger.Error().Err(err).Msg("failed to login user")
		switch err {
		case common.ErrInvalidEmail, common.ErrInvalidPassword, common.ErrNotFound, auth.ErrPasswordVerifyFail, common.ErrUserInactive:
			details := common.ErrorDetail{}
			if errors.Is(err, common.ErrInvalidEmail) {
				details = common.ErrorDetail{
					Field:   "email",
					Message: "Email is a required field",
				}
			}
			if errors.Is(err, common.ErrInvalidPassword) {
				details = common.ErrorDetail{
					Field:   "password",
					Message: "Password is a required field",
				}
			}
			return common.ErrorResponse(c, http.StatusBadRequest, "Invalid email or password", details)
		default:
			return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	if accessToken == "" || refreshToken == "" {
		logger.Error().Msg("Failed to login user: invalid credentials")
		return common.ErrorResponse(c, http.StatusBadRequest, "Invalid email or password")
	}

	resp := &ResponseLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpresIn:     expiresIn,
	}
	return common.SuccessResponse[ResponseLogin](c, resp, "login success")

}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenSuccessResponse struct {
	Message      string                `json:"message"`
	ResponseCode string                `json:"response_code,omitempty"`
	Data         *RefreshTokenResponse `json:"data,omitempty"`
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RefreshTokenSuccessResponse true "Refresh token payload"
// @Success 200 {object} RefreshTokenResponse "Token refreshed successfully"
// @Failure default {object} ErrorResponseDoc "Invalid refresh token"
// @Router /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {

	logger := zerolog.Ctx(c.Request().Context())
	var in RefreshTokenRequest
	if err := c.Bind(&in); err != nil {
		return common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	accessToken, refreshToken, expiresIn, err := h.svc.RefreshToken(in.RefreshToken)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to refresh token")
		switch err {
		case auth.ErrorRefreshTokenEmpty, auth.ErrTokenInvalid, auth.ErrTokenExpired, auth.ErrTokenMalformed:
			details := common.ErrorDetail{}
			if errors.Is(err, auth.ErrorRefreshTokenEmpty) {
				details = common.ErrorDetail{
					Field:   "refresh_token",
					Message: "Refresh token is a required field",
				}
			}
			return common.ErrorResponse(c, http.StatusBadRequest, "Invalid refresh token", details)
		default:
			return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	resp := &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	return common.SuccessResponse[RefreshTokenResponse](c, resp, "Token refreshed successfully")
}

type GetProfileRequest struct {
}

type GetProfileResponse struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type GetProfileSuccessResponse struct {
	Message      string              `json:"message"`
	ResponseCode string              `json:"response_code,omitempty"`
	Data         *GetProfileResponse `json:"data,omitempty"`
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the profile information of the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} GetProfileSuccessResponse "Profile retrieved successfully"
// @Failure default {object} ErrorResponseDoc "Unauthorized"
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	logger := zerolog.Ctx(c.Request().Context())
	userId := c.Get("user_id")
	if userId == nil {
		return common.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		logger.Error().Msg("failed to get profile: invalid user ID type")
		return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}
	profile, err := h.svc.GetProfile(userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get profile")
		switch err {
		case common.ErrNotFound:
			return common.ErrorResponse(c, http.StatusNotFound, "User not found")
		default:
			return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	resp := &GetProfileResponse{
		Email:    profile.Email,
		FullName: profile.FullName,
		Role:     profile.Role,
		IsActive: profile.IsActive,
	}
	return common.SuccessResponse[GetProfileResponse](c, resp, "Profile retrieved successfully")
}
