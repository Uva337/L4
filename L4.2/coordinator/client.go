package coordinator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"dgrep/engine"
	"dgrep/worker"
)

// читает строки, бьет их на батчи и раскидывает по серверам
func RunCoordinator(servers []string, pattern string, allLines []string) {
	if len(servers) == 0 {
		log.Fatal("Нет доступных серверов для работы")
	}

	numServers := len(servers)
	chunkSize := (len(allLines) + numServers - 1) / numServers

	type resultMsg struct {
		matches []engine.Match
		err     error
	}
	resultsCh := make(chan resultMsg, numServers)

	var wg sync.WaitGroup

	// Раскидываем задачи по серверам
	for i, serverURL := range servers {
		wg.Add(1)

		// Вычисляем, какой кусок текста (срез) достанется этому серверу
		start := i * chunkSize
		end := start + chunkSize
		if end > len(allLines) {
			end = len(allLines)
		}

		var chunk []string
		if start < len(allLines) {
			chunk = allLines[start:end]
		}

		// Запускаем горутину-отправщика
		go func(url string, linesBatch []string) {
			defer wg.Done()

			// Если кусок пустой, возвращаем успех
			if len(linesBatch) == 0 {
				resultsCh <- resultMsg{matches: []engine.Match{}, err: nil}
				return
			}

			// Упаковываем данные
			reqData := worker.GrepRequest{
				Pattern: pattern,
				Lines:   linesBatch,
			}
			jsonData, _ := json.Marshal(reqData)

			// Отправляем HTTP POST запрос воркеру
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post("http://"+url+"/grep", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				resultsCh <- resultMsg{err: fmt.Errorf("ошибка запроса к %s: %v", url, err)}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				resultsCh <- resultMsg{err: fmt.Errorf("сервер %s вернул статус %d", url, resp.StatusCode)}
				return
			}

			// Распаковываем ответ
			var matches []engine.Match
			if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
				resultsCh <- resultMsg{err: fmt.Errorf("ошибка парсинга ответа от %s: %v", url, err)}
				return
			}

			resultsCh <- resultMsg{matches: matches, err: nil}
		}(serverURL, chunk)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var totalMatches []engine.Match
	successCount := 0

	for msg := range resultsCh {
		if msg.err != nil {
			log.Printf("⚠️ Воркер отвалился: %v\n", msg.err)
		} else {
			successCount++
			totalMatches = append(totalMatches, msg.matches...)
		}
	}

	// Проверяем кворум: N/2 + 1
	quorum := (numServers / 2) + 1
	if successCount >= quorum {
		fmt.Printf("✅ Кворум достигнут! Успешных ответов: %d из %d\n", successCount, numServers)
		fmt.Printf("🎯 Найдено совпадений: %d\n", len(totalMatches))
		fmt.Println("--------------------------------------------------")

		for _, m := range totalMatches {
			// номер_строки: текст
			fmt.Printf("%d: %s\n", m.LineNum, m.Text)
		}
	} else {
		log.Fatalf("❌ Кворум НЕ достигнут. Успешно ответило: %d из %d. Требуется: %d\n", successCount, numServers, quorum)
	}
}
