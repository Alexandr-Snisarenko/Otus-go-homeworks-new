package domain

import (
	"errors"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventIsEmpty  = errors.New("event is empty")
)
