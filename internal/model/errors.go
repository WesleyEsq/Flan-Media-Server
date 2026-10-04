package model

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrDuplicate           = errors.New("resource already exists")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrLockedOut           = errors.New("account locked out")
	ErrInvalidInput        = errors.New("invalid input")
	ErrStorageOffline      = errors.New("storage source is offline")
	ErrMissingMarker       = errors.New("missing .flan-keep marker")
	ErrInsufficientStorage = errors.New("insufficient storage space")
	ErrBcryptBusy          = errors.New("authentication server busy")
)
