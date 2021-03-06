package jsonapi

import (
	"net/http"
	"strconv"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/errors"
	"github.com/goadesign/goa"
	pkgerrors "github.com/pkg/errors"
)

const (
	ErrorCodeNotFound          = "not_found"
	ErrorCodeBadParameter      = "bad_parameter"
	ErrorCodeVersionConflict   = "version_conflict"
	ErrorCodeUnknownError      = "unknown_error"
	ErrorCodeConversionError   = "conversion_error"
	ErrorCodeInternalError     = "internal_error"
	ErrorCodeUnauthorizedError = "unauthorized_error"
	ErrorCodeJWTSecurityError  = "jwt_security_error"
)

// ErrorToJSONAPIError returns the JSONAPI representation
// of an error and the HTTP status code that will be associated with it.
// This function knows about the models package and the errors from there
// as well as goa error classes.
func ErrorToJSONAPIError(err error) (app.JSONAPIError, int) {
	detail := err.Error()
	var title, code string
	var statusCode int
	var id *string
	switch err.(type) {
	case errors.NotFoundError:
		code = ErrorCodeNotFound
		title = "Not found error"
		statusCode = http.StatusNotFound
	case errors.ConversionError:
		code = ErrorCodeConversionError
		title = "Conversion error"
		statusCode = http.StatusBadRequest
	case errors.BadParameterError:
		code = ErrorCodeBadParameter
		title = "Bad parameter error"
		statusCode = http.StatusBadRequest
	case errors.VersionConflictError:
		code = ErrorCodeVersionConflict
		title = "Version conflict error"
		statusCode = http.StatusBadRequest
	case errors.InternalError:
		code = ErrorCodeInternalError
		title = "Internal error"
		statusCode = http.StatusInternalServerError
	default:
		code = ErrorCodeUnknownError
		title = "Unknown error"
		statusCode = http.StatusInternalServerError

		cause := pkgerrors.Cause(err)
		if err, ok := cause.(goa.ServiceError); ok {
			statusCode = err.ResponseStatus()
			idStr := err.Token()
			id = &idStr
			title = http.StatusText(statusCode)
		}
		if errResp, ok := cause.(*goa.ErrorResponse); ok {
			code = errResp.Code
			detail = errResp.Detail
		}
	}
	statusCodeStr := strconv.Itoa(statusCode)
	jerr := app.JSONAPIError{
		ID:     id,
		Code:   &code,
		Status: &statusCodeStr,
		Title:  &title,
		Detail: detail,
	}
	return jerr, statusCode
}

// ErrorToJSONAPIErrors is a convenience function if you
// just want to return one error from the models package as a JSONAPI errors
// array.
func ErrorToJSONAPIErrors(err error) (*app.JSONAPIErrors, int) {
	jerr, httpStatusCode := ErrorToJSONAPIError(err)
	jerrors := app.JSONAPIErrors{}
	jerrors.Errors = append(jerrors.Errors, &jerr)
	return &jerrors, httpStatusCode
}
