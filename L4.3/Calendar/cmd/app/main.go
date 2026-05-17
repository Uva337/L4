package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"calendar/internal/api"
	"calendar/internal/service"
	"calendar/internal/storage"
	"calendar/pkg/logger"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Инициализация слоев
	repo := storage.NewMemoryStorage()
	svc := service.NewCalendarService(repo)
	handler := api.NewHandler(svc)

	// Логгер с буфером на 100 сообщений
	asyncLogger := logger.NewAsyncLogger(100)

	// 2. Настройка роутера
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", api.LoggingMiddleware(handler.CreateEvent, asyncLogger))
	mux.HandleFunc("/update_event", api.LoggingMiddleware(handler.UpdateEvent, asyncLogger))
	mux.HandleFunc("/delete_event", api.LoggingMiddleware(handler.DeleteEvent, asyncLogger))

	mux.HandleFunc("/events_for_day", api.LoggingMiddleware(handler.GetEventsDay, asyncLogger))
	mux.HandleFunc("/events_for_week", api.LoggingMiddleware(handler.GetEventsWeek, asyncLogger))
	mux.HandleFunc("/events_for_month", api.LoggingMiddleware(handler.GetEventsMonth, asyncLogger))

	// 3. Запуск сервера
	fmt.Printf("🚀 Календарь запущен на порту :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Ошибка остановки сервера: %v", err)
	}
}
