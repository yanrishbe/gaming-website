// Package entities collects data structures used both by a server and a database
package entity

import (
	"errors"
	"net/http"
)

// User struct represents a struct necessary for storing and changing user's data
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (u User) CanRegister() error {
	if u.Name == "" {
		return DBRegisterErr(errors.New("user's name is empty"))
	} else if u.Balance < 300 {
		return DBRegisterErr(errors.New("user has got not enough points"))
	}
	return nil
}

type Error struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Cause   error  `json:"-"`
	Message string `json:"message"`
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
	ErrDecode       = "decoding data error"
	ErrID           = "wrong input id"
)

func HandlerErr(err error) Error {
	e, ok := err.(Error)
	if ok {
		return e
	}
	return Error{
		Type:    ErrUnknown,
		Cause:   err,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

func DecodeErr(err error) Error {
	return Error{
		Type:    ErrDecode,
		Cause:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}
}

func DBErr(err error) Error {
	return Error{
		Type:    ErrDB,
		Cause:   err,
		Code:    http.StatusServiceUnavailable,
		Message: err.Error(),
	}
}

// you should not have deadcode in your application.
// all previosly commited code is saved in Git, so if you don't need something - just delete it, so it doesn't confuse others
// or if you want to explicitly put code here for some important reason - comment it out.
func FewBalErr(err error) Error {
	return Error{
		Type:    ErrFewBalance,
		Cause:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}
}

func InvIDErr(err error) Error {
	return Error{
		Type:    ErrID,
		Cause:   err,
		Code:    http.StatusBadRequest,
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

func UserNotFoundErr(err error) Error {
	return Error{
		Type:    ErrUserNotFound,
		Cause:   err,
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
}
