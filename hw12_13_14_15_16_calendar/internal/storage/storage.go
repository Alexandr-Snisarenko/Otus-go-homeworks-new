package storage

import (
	"context"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/domain"
)

type EventStorage interface {
	CreateEvent(context.Context, *domain.Event) error
	UpdateEvent(context.Context, *domain.Event) error
	DeleteEvent(context.Context, int64) error
	GetEvent(context.Context, int64) (*domain.Event, error)
	GetEvents(context.Context, domain.EventFilter) ([]*domain.Event, error)
	Close() error
}

type UserStorage interface {
	CreateUser(context.Context, *domain.User) error
	UpdateUser(context.Context, *domain.User) error
	DeleteUser(context.Context, int64) error
	GetUser(context.Context, int64) (*domain.User, error)
	Close() error
}
