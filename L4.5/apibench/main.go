package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "net/http/pprof"
)

// Структура для строгого и быстрого парсинга JSON
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// медленный код
func SlowHandler(w http.ResponseWriter, r *http.Request) {
	data := generateJSONPayload(1000)

	var users []map[string]interface{}
	if err := json.Unmarshal(data, &users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := ""
	for _, u := range users {
		name := u["name"].(string)
		result += name + ", "
	}

	w.Write([]byte(result))
}

// оптимизированная версия
func FastHandler(w http.ResponseWriter, r *http.Request) {
	data := generateJSONPayload(1000)

	// 1: Парсим в строгую структуру (никаких интерфейсов)
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2: Используем strings.Builder для конкатенации
	var builder strings.Builder
	builder.Grow(len(users) * 15)

	for _, u := range users {
		builder.WriteString(u.Name)
		builder.WriteString(", ")
	}

	w.Write([]byte(builder.String()))
}

func generateJSONPayload(count int) []byte {
	jsonStr := "["
	for i := 0; i < count; i++ {
		jsonStr += fmt.Sprintf(`{"id": %d, "name": "User%d", "isActive": true}`, i, i)
		if i < count-1 {
			jsonStr += ","
		}
	}
	jsonStr += "]"
	return []byte(jsonStr)
}

func main() {
	http.HandleFunc("/process-slow", SlowHandler)
	http.HandleFunc("/process-fast", FastHandler)

	fmt.Println("🚀 Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
