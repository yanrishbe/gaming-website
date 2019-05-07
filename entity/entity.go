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
	}
	return nil
}

type UserTourn struct {
	ID   int    `json:"userId"`
	Name string `json:"name"`
}

func (u UserTourn) IsValid() error {
	if u.Name == "" {
		return RegErr(errors.New("empty name"))
	}
	return nil
}

type Winner struct {
	ID     int    `json:"userId"`
	Name   string `json:"name"`
	Winner bool   `json:"winner"`
}

type TournFinished struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Winner int      `json:"winner"`
	Prize  int      `json:"prize"`
	Users  []Winner `json:"users"`
	Status Status   `json:"status"`
}

type Tournament struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Deposit int         `json:"deposit"`
	Prize   int         `json:"prize"`
	Users   []UserTourn `json:"users"`
	Status  Status      `json:"status"`
}
type Status string

const (
	Active   Status = "active"
	Finished Status = "finished"
)

func (t Tournament) IsValid() error {
	if t.Name == "" {
		return RegErr(errors.New("empty name"))
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
	ErrPoints       = "wrong input points"
	ErrInvReq       = "wrong request"
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

func PointsErr(err error) Error {
	return Error{
		Type:    ErrPoints,
		Cause:   err,
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

func FinishErr(err error) Error {
	return Error{
		Type:    ErrInvReq,
		Cause:   err,
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}
