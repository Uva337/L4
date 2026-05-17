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
	case <-time.After(100 * time.Millisecond):
		t.Error("Канал не закрылся при закрытии единственного входящего")
	}
}

func TestOr_Multiple(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	res := Or(ch1, ch2, ch3)
	close(ch2)

	select {
	case <-res:
	case <-time.After(100 * time.Millisecond):
		t.Error("Объединенный канал не отреагировал на закрытие ch2")
	}
}

func TestOr_GoroutineLeak(t *testing.T) {
	initialGoroutines := runtime.NumGoroutine()

	var channels []<-chan interface{}
	for i := 0; i < 20; i++ {
		channels = append(channels, make(chan interface{}))
	}
	triggerCh := make(chan interface{})
	channels = append(channels, triggerCh)

	res := Or(channels...)

	close(triggerCh)
	<-res
	time.Sleep(50 * time.Millisecond)

	finalGoroutines := runtime.NumGoroutine()
	if finalGoroutines > initialGoroutines+2 {
		t.Errorf("Обнаружена утечка горутин! Было: %d, Стало: %d", initialGoroutines, finalGoroutines)
	}
}
