package service

import (
	"testing"
	"time"

	"calendar/internal/domain"
	"calendar/internal/storage"
)

func TestCalendarService_CreateAndGetEvents(t *testing.T) {
	// Инициализируем хранилище и сервис
	repo := storage.NewMemoryStorage()
	svc := NewCalendarService(repo)

	userID := 1
	now := time.Now()

	event := domain.Event{
		ID:     "test_1",
		UserID: userID,
		Title:  "Тестовое событие",
		Date:   now,
	}

	// добавление
	err := svc.CreateEvent(event)
	if err != nil {
		t.Fatalf("ожидалось успешное создание, получена ошибка: %v", err)
	}

	// ошибку при дублировании ID
	err = svc.CreateEvent(event)
	if err == nil {
		t.Error("ожидалась ошибка при добавлении события с существующим ID, но её нет")
	}

	// выборку "за день"
	events, err := svc.GetEvents(userID, "day")
	if err != nil {
		t.Fatalf("ошибка при получении событий: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("ожидалось 1 событие, получено %d", len(events))
	}

	if len(events) > 0 && events[0].Title != "Тестовое событие" {
		t.Errorf("ожидалось название 'Тестовое событие', получено '%s'", events[0].Title)
	}
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	repo := storage.NewMemoryStorage()
	svc := NewCalendarService(repo)

	event := domain.Event{
		ID:     "test_delete",
		UserID: 2,
		Title:  "Событие для удаления",
		Date:   time.Now(),
	}

	_ = svc.CreateEvent(event)

	// Тест успешное удаление
	err := svc.DeleteEvent("test_delete")
	if err != nil {
		t.Fatalf("ожидалось успешное удаление, получена ошибка: %v", err)
	}

	events, _ := svc.GetEvents(2, "day")
	if len(events) != 0 {
		t.Errorf("ожидалось 0 событий после удаления, получено %d", len(events))
	}
}
