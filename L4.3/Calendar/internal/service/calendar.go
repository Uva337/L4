package service

import (
	"fmt"
	"time"

	"calendar/internal/domain"
	"calendar/internal/storage"
)

type CalendarService struct {
	repo         *storage.MemoryStorage
	reminderChan chan domain.Event
}

func NewCalendarService(repo *storage.MemoryStorage) *CalendarService {
	s := &CalendarService{
		repo:         repo,
		reminderChan: make(chan domain.Event, 100),
	}

	go s.reminderWorker()
	go s.archiverWorker()

	return s
}

// добавляет событие и планирует напоминание
func (s *CalendarService) CreateEvent(e domain.Event) error {
	if err := s.repo.Create(e); err != nil {
		return err
	}

	if !e.ReminderAt.IsZero() && e.ReminderAt.After(time.Now()) {
		s.reminderChan <- e
	}
	return nil
}

func (s *CalendarService) UpdateEvent(e domain.Event) error {
	return s.repo.Update(e)
}

func (s *CalendarService) DeleteEvent(id string) error {
	return s.repo.Delete(id)
}

func (s *CalendarService) GetEvents(userID int, period string) ([]domain.Event, error) {
	now := time.Now()
	var start, end time.Time

	switch period {
	case "day":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.Add(24 * time.Hour).Add(-time.Nanosecond)
	case "week":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	default:
		return nil, fmt.Errorf("неизвестный период: %s", period)
	}

	return s.repo.GetEvents(userID, start, end), nil
}

// слушает канал и ставит таймеры на отправку напоминаний
func (s *CalendarService) reminderWorker() {
	for event := range s.reminderChan {
		waitDuration := time.Until(event.ReminderAt)

		// Запускаем таймер в фоне, чтобы воркер мог сразу брать следующие задачи
		go func(e domain.Event, wait time.Duration) {
			time.Sleep(wait)
			fmt.Printf("\n🔔 НАПОМИНАНИЕ: Событие '%s' начнется %s!\n", e.Title, e.Date.Format("2006-01-02 15:04"))
		}(event, waitDuration)
	}
}

// раз в минуту чистит старые события
func (s *CalendarService) archiverWorker() {
	// Для тестов поставим 1 минуту.
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		deleted := s.repo.DeleteBefore(time.Now())
		if deleted > 0 {
			fmt.Printf("🧹 АРХИВАТОР: Удалено старых событий: %d\n", deleted)
		}
	}
}
