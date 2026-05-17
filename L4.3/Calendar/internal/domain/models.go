package domain

import "time"

// описывает событие в календаре
type Event struct {
	ID         string    `json:"id"`
	UserID     int       `json:"user_id"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`        // Дата самого события
	ReminderAt time.Time `json:"reminder_at"` // Время, когда нужно прислать напоминание
}
