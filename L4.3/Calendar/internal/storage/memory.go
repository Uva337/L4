package storage

import (
	"errors"
	"sync"
	"time"

	"calendar/internal/domain"
)

// хранилище событий в оперативной памяти
type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]domain.Event
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		events: make(map[string]domain.Event),
	}
}

func (s *MemoryStorage) Create(e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[e.ID]; exists {
		return errors.New("событие с таким ID уже существует")
	}
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStorage) Update(e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[e.ID]; !exists {
		return errors.New("событие не найдено")
	}
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[id]; !exists {
		return errors.New("событие не найдено")
	}
	delete(s.events, id)
	return nil
}

// возвращает все события пользователя за указанный период
func (s *MemoryStorage) GetEvents(userID int, start, end time.Time) []domain.Event {
	s.mu.RLock() // Используем Read-Lock, так как мы только читаем
	defer s.mu.RUnlock()

	var result []domain.Event
	for _, e := range s.events {
		if e.UserID == userID && !e.Date.Before(start) && !e.Date.After(end) {
			result = append(result, e)
		}
	}
	return result
}

// удаляет все события, дата которых раньше указанной (для архиватора)
func (s *MemoryStorage) DeleteBefore(t time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for id, e := range s.events {
		if e.Date.Before(t) {
			delete(s.events, id)
			count++
		}
	}
	return count
}
