package response

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ResponseDTO[T any] struct {
	Message      string    `json:"message"`
	ResponseCode string    `json:"response_code,omitempty"`
	Data         *T        `json:"data,omitempty"`
	Error        *ErrorObj `json:"error,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorObj struct {
	Details []ErrorDetail `json:"details,omitempty"`
}

func Success[T any](ctx echo.Context, data *T, message string) error {
	setJSON(ctx)
	resp := &ResponseDTO[T]{
		Data:         data,
		ResponseCode: strconv.Itoa(http.StatusOK),
		Message:      message,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func Error(ctx echo.Context, code int, msg string, details ...ErrorDetail) error {
	setJSON(ctx)
	resp := &ResponseDTO[any]{
		Message:      msg,
		ResponseCode: strconv.Itoa(code),
		Error:        &ErrorObj{Details: details},
	}
	return ctx.JSON(code, resp)
}

func setJSON(ctx echo.Context) {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}

// func MapValidationErrors(err error) []ErrorDetail {
// 	var verrs validator.ValidationErrors
// 	if !errors.As(err, &verrs) {
// 		return []ErrorDetail{{Message: err.Error()}}
// 	}
// 	res := make([]ErrorDetail, 0, len(verrs))
// 	for _, fe := range verrs {
// 		d := ErrorDetail{
// 			Field:   fe.Field(),
// 			Message: fe.Message(),
// 		}
// 		res = append(res, d)
// 	}
// 	return res
// }
