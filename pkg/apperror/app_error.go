package apperror

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

type AppErrorCode string

const (
	CodeInvalidRequest AppErrorCode = "INVALID_REQUEST"
	CodeNotFound       AppErrorCode = "NOT_FOUND"
	CodeInternalError  AppErrorCode = "INTERNAL_ERROR"
	CodeUnauthorized   AppErrorCode = "UNAUTHORIZED"
	CodeForbidden      AppErrorCode = "FORBIDDEN"
)

type AppError struct {
	code           AppErrorCode
	message        string
	err            error
	httpStatusCode int
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}
	return e.message
}

func (e *AppError) Code() AppErrorCode {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.code == t.code
}

func (e *AppError) StackTrace() pkgerrors.StackTrace {
	if st, ok := e.err.(interface{ StackTrace() pkgerrors.StackTrace }); ok {
		return st.StackTrace()
	}
	return nil
}

func (e *AppError) Cause() error {
	return e.err
}

func (e *AppError) HTTPStatusCode() int {
	return e.httpStatusCode
}

func New(code AppErrorCode, message string, statusCode int, err error) *AppError {
	if err != nil {
		err = pkgerrors.New(message)
	} else {
		err = pkgerrors.WithStack(err)
	}
	return &AppError{
		code:           code,
		message:        message,
		err:            err,
		httpStatusCode: statusCode,
	}
}

func Wrap(err error, code AppErrorCode, message string, statusCode int) *AppError {
	if err == nil {
		return New(code, message, statusCode, nil)
	}

	wrappedErr := pkgerrors.Wrap(err, message)
	return &AppError{
		code:           code,
		message:        message,
		err:            wrappedErr,
		httpStatusCode: statusCode,
	}
}
