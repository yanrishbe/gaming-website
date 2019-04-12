// Package entities collects data structures used both by a server and a database
package entity

import "net/http"

// User struct represents a struct necessary for storing and changing user's data
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type Error struct {
	Type    string
	Code    int
	Cause   error
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrFewBalance   = "balance is not enough"
	ErrCannotReg    = "user's data is not valid"
	ErrUserNotFound = "user not found"
	ErrDB           = "database error"
	ErrUnknown      = "unknown error"
)

func HandlerError(err error) Error {
	e, ok := err.(Error)
	if ok {
		return e
	}
	return Error{
		Type:    ErrUnknown,
		Cause:   err,
		Code:    500,
		Message: err.Error(),
	}
}

func DBError(err error) Error {
	return Error{
		Type:    ErrDB,
		Cause:   err,
		Code:    503,
		Message: err.Error(),
	}
}

func FewBal(err error) Error {
	return Error{
		Type:    ErrCannotReg,
		Cause:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}
}

func DBRegisterErr(err error) Error {
	return Error{
		Type:    ErrCannotReg,
		Cause:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}

}

func UserNotFoundError(err error) Error {
	return Error{
		Type:    ErrUserNotFound,
		Cause:   err,
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
}
