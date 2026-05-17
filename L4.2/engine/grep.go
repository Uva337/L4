package engine

import (
	"strings"
	"sync"
)

// описывает одно найденное совпадение
type Match struct {
	LineNum int    `json:"line_num"` // Номер строки
	Text    string `json:"text"`     // Сама строка
}

// параллельно ищет подстроку в массиве строк
func LocalGrep(lines []string, pattern string, workers int) []Match {
	if len(lines) == 0 {
		return nil
	}

	type job struct {
		index int
		line  string
	}

	jobs := make(chan job, len(lines))
	results := make(chan Match, len(lines))

	// 1. Заполние очереди задач
	for i, line := range lines {
		jobs <- job{index: i + 1, line: line}
	}
	close(jobs)

	// 2. Запуск Worker Pool
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if strings.Contains(j.line, pattern) {
					results <- Match{
						LineNum: j.index,
						Text:    j.line,
					}
				}
			}
		}()
	}

	// 3. Ожидание
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Сбор резов
	var matches []Match
	for match := range results {
		matches = append(matches, match)
	}

	return matches
}
