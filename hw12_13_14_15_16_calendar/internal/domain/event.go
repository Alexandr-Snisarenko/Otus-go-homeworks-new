package domain

import "time"

// Event - структура для объекта "Событие"
// ID - уникальный идентификатор события (BIGINT);
// Заголовок - короткий текст;
// Дата и время события;
// Длительность события (или дата и время окончания);
// Описание события - длинный текст, опционально;
// ID пользователя, владельца события;
// За сколько времени высылать уведомление, опционально.
type Event struct {
	ID           int64         `db:"id"`
	Title        string        `db:"title"`
	Description  string        `db:"description"`
	StartTime    time.Time     `db:"start_time"`
	EndTime      time.Time     `db:"end_time"`
	UserID       int64         `db:"user_id"`
	NotifyPeriod time.Duration `db:"notify_period"`
}

type EventFilter struct {
	UserID   *int64
	DateFrom *time.Time
	DateTo   *time.Time
}
