package or

import (
	"runtime"
	"testing"
	"time"
)

func TestOr_Empty(t *testing.T) {
	if Or() != nil {
		t.Error("Or() без аргументов должен возвращать nil")
	}
}

func TestOr_Single(t *testing.T) {
	ch := make(chan interface{})
	res := Or(ch)

	close(ch)

	select {
	case <-res:
		// Успех, канал пропустил сигнал
	case <-time.After(100 * time.Millisecond):
		t.Error("Канал не закрылся при закрытии единственного входящего")
	}
}

func TestOr_Multiple(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	res := Or(ch1, ch2, ch3)

	// Закрываем второй канал (в середине цепочки)
	close(ch2)

	select {
	case <-res:
		// Успех
	case <-time.After(100 * time.Millisecond):
		t.Error("Объединенный канал не отреагировал на закрытие ch2")
	}
}

func TestOr_GoroutineLeak(t *testing.T) {
	initialGoroutines := runtime.NumGoroutine()

	// Создаем цепочку из 20 каналов
	var channels []<-chan interface{}
	for i := 0; i < 20; i++ {
		channels = append(channels, make(chan interface{}))
	}

	// Берем первый канал отдельно, чтобы его закрыть и триггернуть всю сеть
	triggerCh := make(chan interface{})
	channels = append(channels, triggerCh)

	res := Or(channels...)

	// Запускаем каскадное сворачивание
	close(triggerCh)
	<-res

	// Даем планировщику Go немного времени на сборку завершенных горутин
	time.Sleep(50 * time.Millisecond)

	finalGoroutines := runtime.NumGoroutine()
	// Допустима небольшая погрешность на системные процессы тестов,
	// но 20+ горутин остаться не должно.
	if finalGoroutines > initialGoroutines+2 {
		t.Errorf("Обнаружена утечка горутин! Было: %d, Стало: %d", initialGoroutines, finalGoroutines)
	}
}
