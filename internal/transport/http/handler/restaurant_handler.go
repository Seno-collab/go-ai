package handler

import (
	restaurantapp "go-ai/internal/application/restaurant"
	"go-ai/internal/domain/restaurant"
	"go-ai/internal/transport/http/response"
	"go-ai/internal/transport/http/status"
	"go-ai/pkg/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type RestaurantHandler struct {
	CreateUC *restaurantapp.CreateRestaurantUseCase
	Logger   zerolog.Logger
}

func NewRestaurantHandler(createUC *restaurantapp.CreateRestaurantUseCase) *RestaurantHandler {
	return &RestaurantHandler{
		CreateUC: createUC,
		Logger:   logger.NewLogger().With().Str("component", "Restaurant handler").Logger(),
	}
}

// CreateRestaurant godoc
// @Summary Create restaurant
// @Description Create a new restaurant with name, email, phone, logo_url, banner_url,...
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param request body restaurantapp.CreateRestaurantRequest true "Restaurant create payload"
// @Success 200 {object} app.CreateRestaurantSuccessResponseDoc "Create restaurant successfully"
// @Failure default {object} app.ErrorResponseDoc "Errors"
// @Router /api/restaurant [post]
func (h *RestaurantHandler) Create(c echo.Context) error {
	var in restaurantapp.CreateRestaurantRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	userId := c.Get("user_id")
	if userId == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	_, err := h.CreateUC.Execute(in, userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed create restaurant")
		details := response.ErrorDetail{}
		switch err {
		case status.ErrInvalidEmail:
			details = response.ErrorDetail{
				Field:   "email",
				Message: "Email is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case status.ErrInvalidName:
			details = response.ErrorDetail{
				Field:   "name",
				Message: "Name is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case restaurant.ErrInvalidAddress:
			details = response.ErrorDetail{
				Field:   "address",
				Message: "Address is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case restaurant.ErrInvalidBanner:
			details = response.ErrorDetail{
				Field:   "banner_url",
				Message: "Banner url is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case restaurant.ErrInvalidLogo:
			details = response.ErrorDetail{
				Field:   "logo_url",
				Message: "Logo url is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case restaurant.ErrInvalidPhoneNumber:
			details = response.ErrorDetail{
				Field:   "phone_numer",
				Message: "Phone number is a required field",
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		default:
			return response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	return response.Success[any](c, nil, "Create restaurant successfully")
}
