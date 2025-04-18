package models

import "errors"

var (
	ErrPVZNotFound            = errors.New("no such PVZ")
	ErrNoActiveReception      = errors.New("no active reception")
	ErrPVZReceptionIsClosed   = errors.New("reception is closed")
	ErrNoProductsInReception  = errors.New("no products in reception")
	ErrPVZAlreadyHasReception = errors.New("PVZ already has active reception")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidCredentials     = errors.New("user login invalid credentials")
	ErrDBInsert               = errors.New("failed to insert into DB")
	ErrDBRead                 = errors.New("failed to read from DB")
	ErrDBUpdate               = errors.New("failer to update in DB")
)
