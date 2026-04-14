package model

/*
| Case                 | Status | Code                |
| -------------------- | ------ | ------------------- |
| invalid input        | 400    | INVALID_INPUT       |
| unauthorized         | 401    | UNAUTHORIZED        |
| invalid password     | 401    | INVALID_CREDENTIALS |
| forbidden            | 403    | FORBIDDEN           |
| not found            | 404    | NOT_FOUND           |
| conflict (duplicate) 	| 409    | USER_ALREADY_EXISTS |
*/

type AppError struct {
	Code       ErrorCode
	Message    string
	Err        error
	HTTPStatus int
}

type ErrorCode string

const (
	CodeInvalidInput       ErrorCode = "INVALID_INPUT"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeForbidden          ErrorCode = "FORBIDDEN"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeUserAlreadyExists  ErrorCode = "USER_ALREADY_EXISTS"
	CodeDeckAlreadyExists  ErrorCode = "DECK_ALREADY_EXISTS"
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
)

func ErrUnAuthorized(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: 401,
	}
}

func ErrConflict(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: 409,
	}
}

func ErrNotFound(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: 404,
	}
}

func ErrBadRequest(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: 400,
	}
}

func ErrInternal(err error) *AppError {
	return &AppError{
		Code:       CodeInternal,
		Message:    "internal server error: " + err.Error(),
		Err:        err,
		HTTPStatus: 500,
	}
}
