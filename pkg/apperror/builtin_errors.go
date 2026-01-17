package apperror

import "net/http"

func InvalidRequest(message string, err error) *AppError {
	return Wrap(err, CodeInvalidRequest, message, http.StatusBadRequest)
}

func NotFound(message string, err error) *AppError {
	return Wrap(err, CodeNotFound, message, http.StatusNotFound)
}

func InternalError(message string, err error) *AppError {
	return Wrap(err, CodeInternalError, message, http.StatusInternalServerError)
}

func Unauthorized(message string, err error) *AppError {
	return Wrap(err, CodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string, err error) *AppError {
	return Wrap(err, CodeForbidden, message, http.StatusForbidden)
}
