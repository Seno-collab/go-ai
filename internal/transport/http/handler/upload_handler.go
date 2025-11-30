package handler

import (
	uploadapp "go-ai/internal/application/upload"
	"go-ai/internal/infra/storage"
	"go-ai/internal/transport/http/response"
	"go-ai/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UpLoadHandler struct {
	MC     *storage.MinioClient
	Logger zerolog.Logger
}

func NewUploadHandler(mc *storage.MinioClient) *UpLoadHandler {
	return &UpLoadHandler{
		MC:     mc,
		Logger: logger.NewLogger().With().Str("component", "Upload logo handler").Logger(),
	}
}

// UploadLogoHandler godoc
// @Summary Upload logo file
// @Description Upload a logo image to storage and return the public URL
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param logo formData file true "Logo file (png, jpg, jpeg)"
// @Success 200 {object} app.UploadLogoSuccessResponseDoc "Upload logo success"
// @Failure default {object} app.ErrorResponseDoc "Errors"
// @Router /api/upload/logo [post]
func (h *UpLoadHandler) UploadLogoHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("logo")
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: missing file")
			return response.Error(c, http.StatusBadRequest, "logo file is required")
		}
		file, err := fileHeader.Open()
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: cannot open file")
			return response.Error(c, http.StatusBadRequest, "unable to open file")
		}
		defer file.Close()
		url, err := h.MC.UploadLogo(c.Request().Context(), file, fileHeader)
		if err != nil {
			h.Logger.Error().Err(err).Msg("Upload logo: MinIO upload error")
			return response.Error(c, http.StatusBadRequest, "upload to storage failed")
		}
		return response.Success(c, &uploadapp.UploadLogoResponse{
			Url: url,
		}, "Upload logo successfully")
	}
}
