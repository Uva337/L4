package or

import (
	"fmt"
	"time"
)

// создает канал, который автоматически закроется через время
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

	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(100*time.Millisecond),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	elapsed := time.Since(start)
	if elapsed < 500*time.Millisecond {
		fmt.Println("Успешно завершено по самому быстрому каналу!")
	}

}
