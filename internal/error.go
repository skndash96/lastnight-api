package api

import (
	"encoding/json"

	"fmt"

	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/repository"
	"github.com/skndash96/lastnight-backend/internal/service"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var (
		code int    = http.StatusInternalServerError
		msg  string = "internal server error"
	)

	if _, ok := err.(*json.SyntaxError); ok {
		code = http.StatusBadRequest
		msg = "malformed JSON"
	} else if unmarshallErr, ok := err.(*json.UnmarshalTypeError); ok {
		code = http.StatusBadRequest
		msg = fmt.Sprintf("expected %v but got %v for field %v", unmarshallErr.Type, unmarshallErr.Value, unmarshallErr.Field)
	} else if bindErr, ok := err.(*echo.BindingError); ok {
		code = http.StatusBadRequest
		msg = bindErr.Error()
	} else if vErr, ok := err.(validator.ValidationErrors); ok {
		msg = CustomValidationErrMsg(&vErr)
		code = http.StatusBadRequest
	} else if repoErr, ok := err.(*repository.RepoError); ok {
		err = repoErr.Unwrap()
		code = mapRepoErrorToApiError(repoErr)
		msg = repoErr.Message
	} else if srvErr, ok := err.(*service.SrvError); ok {
		err = srvErr.Unwrap()
		code = mapSrvErrorToApiError(srvErr)
		msg = srvErr.Message
	} else if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = httpErr.Message.(string)
	}

	c.JSON(code, &dto.ErrorResponse{
		Message: msg,
	})
}

func mapSrvErrorToApiError(err *service.SrvError) int {
	var (
		code int
	)

	switch err.Kind {
	case service.SrvErrInvalidInput:
		code = http.StatusBadRequest
	case service.SrvErrUnauthorized:
		code = http.StatusUnauthorized
	case service.SrvErrForbidden:
		code = http.StatusForbidden
	case service.SrvErrNotFound:
		code = http.StatusNotFound
	case service.SrvErrConflict:
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}

	return code

}

func mapRepoErrorToApiError(err *repository.RepoError) int {
	var (
		code int
	)

	switch err.Kind {
	case repository.RepoErrInvalidInput:
		code = http.StatusBadRequest
	case repository.RepoErrConflict:
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}

	return code
}

func CustomValidationErrMsg(error *validator.ValidationErrors) string {
	msg := "invalid input"

	for _, fieldErr := range *error {
		fieldName := fieldErr.Field()
		fieldParam := fieldErr.Param()

		switch fieldErr.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required and cannot be empty.", fieldName)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address.", fieldName)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters long.", fieldName, fieldParam)
		case "max":
			msg = fmt.Sprintf("%s cannot exceed %s characters.", fieldName, fieldParam)
		default:
			msg = fmt.Sprintf("%s failed validation on the '%s' rule.", fieldName, fieldErr.Tag())
		}
	}

	return msg
}
