package memory

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/domain"
	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[int64]*domain.Event
	mu     sync.RWMutex
}

func (s *Storage) CreateEvent(_ context.Context, event *domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	evID := NewID()
	event.ID = evID
	s.events[evID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, eventID)
	return nil
}

func (s *Storage) GetEvent(_ context.Context, eventID int64) (*domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[eventID]
	if !ok {
		return nil, domain.ErrEventNotFound
	}
	return event, nil
}

// Отбираем все события по фильтру.
func (s *Storage) GetEvents(_ context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*domain.Event, 0, len(s.events))

	for _, event := range s.events {
		// Фильтр по времени
		if filter.DateFrom != nil && event.StartTime.Before(*filter.DateFrom) {
			continue
		}
		if filter.DateTo != nil && event.StartTime.After(*filter.DateTo) {
			continue
		}

		// Фильтр по пользователю
		if filter.UserID != nil && event.UserID != *filter.UserID {
			continue
		}

		events = append(events, event)
	}

	// Сортируем по дате старта
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events, nil
}

func (s *Storage) Close() error {
	return nil
}

func NewID() int64 {
	now := time.Now().UnixNano()       // наносекунды (int64)
	randPart := rand.Int63n(1_000_000) // //nolint:gosec // not used for security
	return now*1_000_000 + randPart
}

func New() storage.EventStorage {
	return &Storage{events: make(map[int64]*domain.Event)}
}
