package or

import (
	"fmt"
	"time"
)

// sig создает канал, который автоматически закроется через заданное время.
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func ExampleOr() {
	start := time.Now()

	// Передаем несколько каналов с разным временем жизни.
	// Самый быстрый сработает через 100 миллисекунд.
	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(100*time.Millisecond), // <-- Победитель
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// Проверяем, что выполнение заняло доли секунды, а не часы
	elapsed := time.Since(start)
	if elapsed < 500*time.Millisecond {
		fmt.Println("Успешно завершено по самому быстрому каналу!")
	}

	// Output:
	// Успешно завершено по самому быстрому каналу!
}
