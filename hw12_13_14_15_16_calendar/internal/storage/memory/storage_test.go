package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/domain"
)

var ctx = context.Background()

func TestNewAndCreateGet(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatalf("New returned nil")
	}
	// Create event
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{
		Title:        "meeting",
		Description:  "desc",
		StartTime:    base,
		EndTime:      base.Add(time.Hour),
		UserID:       42,
		NotifyPeriod: time.Minute,
	}
	if err := s.CreateEvent(ctx, ev); err != nil {
		t.Fatalf("CreateEvent error: %v", err)
	}
	if ev.ID == 0 {
		t.Fatalf("expected non-zero ID assigned")
	}

	got, err := s.GetEvent(ctx, ev.ID)
	if err != nil {
		t.Fatalf("GetEvent error: %v", err)
	}
	if got == nil || got.ID != ev.ID || got.Title != ev.Title {
		t.Fatalf("GetEvent returned wrong event: %#v", got)
	}
}

func TestGetEvent_NotFound(t *testing.T) {
	s := New()
	if _, err := s.GetEvent(ctx, 123); err == nil {
		t.Fatalf("expected ErrEventNotFound, got nil")
	}
}

func TestUpdateEvent(t *testing.T) {
	s := New()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{Title: "old", StartTime: base, EndTime: base}
	_ = s.CreateEvent(ctx, ev)

	ev.Title = "new title"
	if err := s.UpdateEvent(ctx, ev); err != nil {
		t.Fatalf("UpdateEvent error: %v", err)
	}

	got, _ := s.GetEvent(ctx, ev.ID)
	if got.Title != "new title" {
		t.Fatalf("Update didn't persist changes: got=%q", got.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := New()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{Title: "to delete", StartTime: base, EndTime: base}
	_ = s.CreateEvent(ctx, ev)

	if err := s.DeleteEvent(ctx, ev.ID); err != nil {
		t.Fatalf("DeleteEvent error: %v", err)
	}

	if _, err := s.GetEvent(ctx, ev.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}

func TestGetEventsByStartTimeStrictBounds(t *testing.T) {
	s := New()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)

	// dateFrom .. dateTo
	dtFrom := base
	dtTo := base.Add(4 * time.Hour)

	// Events at various times
	evBefore := &domain.Event{Title: "before", StartTime: base.Add(-time.Hour)}
	evAtFrom := &domain.Event{Title: "atFrom", StartTime: dtFrom}
	evInside1 := &domain.Event{Title: "inside1", StartTime: base.Add(time.Hour)}
	evInside2 := &domain.Event{Title: "inside2", StartTime: base.Add(2 * time.Hour)}
	evAtTo := &domain.Event{Title: "atTo", StartTime: dtTo}
	evAfter := &domain.Event{Title: "after", StartTime: dtTo.Add(time.Minute)}

	for _, e := range []*domain.Event{evBefore, evAtFrom, evInside1, evInside2, evAtTo, evAfter} {
		_ = s.CreateEvent(ctx, e)
	}

	filter := domain.EventFilter{DateFrom: &dtFrom, DateTo: &dtTo}
	got, err := s.GetEvents(ctx, filter)
	if err != nil {
		t.Fatalf("GetEvents error: %v", err)
	}

	// Ожидаем 4 события (включая граничные)
	if len(got) != 4 {
		t.Fatalf("expected 4 events, got %d", len(got))
	}
	titles := map[string]bool{}
	for _, e := range got {
		titles[e.Title] = true
	}
	if !titles["inside1"] || !titles["inside2"] {
		t.Fatalf("unexpected set of events: %+v", titles)
	}
}

func TestConcurrentCreateUniqueness(t *testing.T) {
	s := New()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			e := &domain.Event{Title: "e", StartTime: base.Add(time.Duration(i) * time.Minute)}
			_ = s.CreateEvent(ctx, e)
		}(i)
	}
	wg.Wait()

	// Collect events; use very wide interval to include all
	dtFrom := base.Add(-24 * time.Hour)
	dtTo := base.Add(24 * time.Hour)
	filter := domain.EventFilter{DateFrom: &dtFrom, DateTo: &dtTo}

	got, err := s.GetEvents(ctx, filter)
	if err != nil {
		t.Fatalf("GetEvents error: %v", err)
	}

	// Ensure unique IDs counted in map as well
	if len(got) != n {
		t.Fatalf("expected %d events created, got %d", n, len(got))
	}

	seen := make(map[int64]struct{}, n)
	for _, e := range got {
		if e.ID == 0 {
			t.Fatalf("found zero ID")
		}
		if _, ok := seen[e.ID]; ok {
			t.Fatalf("duplicate ID detected: %d", e.ID)
		}
		seen[e.ID] = struct{}{}
	}
}
