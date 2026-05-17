package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"dgrep/coordinator"
	"dgrep/worker"
)

func main() {
	// флаги командной строки
	mode := flag.String("mode", "client", "Режим работы: 'server' или 'client'")
	port := flag.String("port", "8081", "Порт для запуска сервера (только для mode=server)")
	serversList := flag.String("servers", "localhost:8081", "Список серверов через запятую (только для mode=client)")
	pattern := flag.String("pattern", "", "Строка для поиска (только для mode=client)")

	flag.Parse()

	// 2. Ветвление логики в зависимости от режима
	if *mode == "server" {
		// --- РЕЖИМ СЕРВЕРА ---
		worker.StartServer(*port)

	} else if *mode == "client" {

		if *pattern == "" {
			log.Fatal("Ошибка: необходимо указать строку для поиска через флаг -pattern")
		}

		args := flag.Args()
		if len(args) < 1 {
			log.Fatal("Ошибка: укажите путь к файлу для поиска. Пример: ./dgrep -mode=client -pattern=error text.log")
		}
		filename := args[0]

		lines, err := readLines(filename)
		if err != nil {
			log.Fatalf("Ошибка чтения файла: %v", err)
		}

		// Разбиваем строку с адресами серверов в массив
		servers := strings.Split(*serversList, ",")
		for i := range servers {
			servers[i] = strings.TrimSpace(servers[i])
		}

		fmt.Printf("🚀 Запуск распределенного поиска '%s' по файлу %s (%d строк)\n", *pattern, filename, len(lines))

		coordinator.RunCoordinator(servers, *pattern, lines)

	} else {
		log.Fatalf("Неизвестный режим: %s. Используйте 'server' или 'client'", *mode)
	}
}

// функция для чтения файла построчно
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
