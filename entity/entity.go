package entity

import (
	"errors"
	"net/http"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (u User) IsValid() error {
	if u.Name == "" {
		return RegErr(errors.New("empty name"))
	} else if u.Balance < 300 {
		return RegErr(errors.New("low balance"))
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
	ErrInvData      = "user's data is not valid"
	ErrDB           = "database error"
	ErrID           = "wrong input id"
	ErrUserNotFound = "user not found"
	ErrUnknown      = "unknown error"
	ErrDecode       = "decoding data error"
)

func RegErr(err error) Error {
	return Error{
		Type:    ErrInvData,
		Cause:   err,
		Code:    http.StatusBadRequest,
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

func InvIDErr(err error) Error {
	return Error{
		Type:    ErrID,
		Cause:   err,
		Code:    http.StatusBadRequest,
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

func DecodeErr(err error) Error {
	return Error{
		Type:    ErrDecode,
		Cause:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}
}

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
