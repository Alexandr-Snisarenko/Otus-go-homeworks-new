package postgresql

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/config"
	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/domain"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// helper to create sqlx DB from sqlmock.
func newMock() (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	return sqlx.NewDb(db, "pgx"), mock, func() { db.Close() }
}

// construct Storage with mocked db.
func newStorageWithMock() (*Storage, sqlmock.Sqlmock, func()) {
	db, mock, closeFn := newMock()
	return &Storage{db: db}, mock, closeFn
}

func TestCreateEvent(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{
		Title:        "ttl",
		Description:  "d",
		StartTime:    base,
		EndTime:      base.Add(time.Hour),
		UserID:       7,
		NotifyPeriod: time.Minute,
	}

	// Expect INSERT with RETURNING id
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO events (title, description, start_time, end_time, user_id, notify_period)
    	VALUES ($1, $2, $3, $4, $5, $6)
    	RETURNING id`)).
		WithArgs(ev.Title, ev.Description, ev.StartTime, ev.EndTime, ev.UserID, ev.NotifyPeriod).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(123)))

	if err := s.CreateEvent(ctx, ev); err != nil {
		t.Fatalf("CreateEvent error: %v", err)
	}
	if ev.ID != 123 {
		t.Fatalf("expected returned ID=123, got %d", ev.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateEvent(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{
		ID:           10,
		Title:        "ttl",
		Description:  "d",
		StartTime:    base,
		EndTime:      base.Add(time.Hour),
		UserID:       7,
		NotifyPeriod: time.Minute,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE events SET
        title = $1,
        description = $2,
        start_time = $3,
        end_time = $4,
        user_id = $5,
        notify_period = $6
    	WHERE id = $7`)).
		WithArgs(ev.Title, ev.Description, ev.StartTime, ev.EndTime, ev.UserID, ev.NotifyPeriod, ev.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := s.UpdateEvent(ctx, ev); err != nil {
		t.Fatalf("UpdateEvent error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateEvent_NotFound(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	ev := &domain.Event{
		ID:           999,
		Title:        "ttl",
		Description:  "d",
		StartTime:    base,
		EndTime:      base.Add(time.Hour),
		UserID:       7,
		NotifyPeriod: time.Minute,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE events SET
     	title = $1,
     	description = $2,
     	start_time = $3,
     	end_time = $4,
     	user_id = $5,
     	notify_period = $6
    	WHERE id = $7`)).
		WithArgs(ev.Title, ev.Description, ev.StartTime, ev.EndTime, ev.UserID, ev.NotifyPeriod, ev.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.UpdateEvent(ctx, ev)
	if err == nil || !errors.Is(err, domain.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM events WHERE id = $1`)).
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := s.DeleteEvent(ctx, 42); err != nil {
		t.Fatalf("DeleteEvent error: %v", err)
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM events WHERE id = $1`)).
		WithArgs(int64(777)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.DeleteEvent(ctx, 777)
	if err == nil || !errors.Is(err, domain.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got %v", err)
	}
}

func TestGetEvent(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time", "user_id", "notify_period"}).
		AddRow(int64(5), "ttl", "d", base, base.Add(time.Hour), int64(7), time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, title, description, start_time, end_time, user_id, notify_period
        FROM events
        WHERE id = $1`)).
		WithArgs(int64(5)).WillReturnRows(rows)

	ev, err := s.GetEvent(ctx, 5)
	if err != nil {
		t.Fatalf("GetEvent error: %v", err)
	}
	if ev.ID != 5 || ev.Title != "ttl" || ev.UserID != 7 {
		t.Fatalf("unexpected event: %#v", ev)
	}
}

func TestGetEvents_WithFiltersAndOrder(t *testing.T) {
	s, mock, closeFn := newStorageWithMock()
	defer closeFn()

	ctx := context.Background()
	base := time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	dtFrom := base
	dtTo := base.Add(2 * time.Hour)
	userID := int64(7)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time", "user_id", "notify_period"}).
		AddRow(int64(1), "a", "", base, base.Add(time.Minute), userID, time.Minute).
		AddRow(int64(2), "b", "", base.Add(time.Hour), base.Add(61*time.Minute), userID, time.Minute)

	// The Squirrel builder produces:
	// SELECT ... FROM events
	// WHERE start_time >= $1 AND start_time <= $2 AND user_id = $3
	// ORDER BY start_time
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, description, start_time, end_time, user_id, notify_period 
		FROM events 	
		WHERE start_time >= $1 
		AND start_time <= $2 
		AND user_id = $3 
		ORDER BY start_time`)).
		WithArgs(dtFrom, dtTo, userID).
		WillReturnRows(rows)

	got, err := s.GetEvents(ctx, domain.EventFilter{DateFrom: &dtFrom, DateTo: &dtTo, UserID: &userID})
	if err != nil {
		t.Fatalf("GetEvents error: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestOpenDB_DSNFromConfig(t *testing.T) {
	// Using sqlmock directly for Open isn't straightforward because sqlx.Open will dial driver by name.
	// Here we only verify error path when workmode is not postgresql
	// and DSN building when dsn is empty is not unit-testable without driver.

	cfg := config.Database{Workmode: "memory"}
	if _, err := OpenDB(cfg); err == nil {
		t.Fatalf("expected error for non-postgresql workmode")
	}
}
