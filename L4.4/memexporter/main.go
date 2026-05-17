package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"

	_ "net/http/pprof"
)

// собирает метрики памяти и отдает их в формате Prometheus
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Формируем ответ в формате Prometheus (Plain Text)
	metrics := `
# HELP go_memstats_alloc_bytes Bytes of allocated heap objects.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes %d

# HELP go_memstats_sys_bytes Total bytes of memory obtained from the OS.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes %d

# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total %d

# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total %d

# HELP go_memstats_num_gc Number of completed GC cycles.
# TYPE go_memstats_num_gc counter
go_memstats_num_gc %d

# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds %f
`
	// Prometheus ожидает время в секундах, а LastGC возвращает наносекунды
	lastGCSeconds := float64(m.LastGC) / 1e9

	// Устанавливаем правильный заголовок, который ждет сервер Prometheus
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, metrics, m.Alloc, m.Sys, m.Mallocs, m.Frees, m.NumGC, lastGCSeconds)
}

// позволяет на лету менять агрессивность сборщика мусора
func setGCHandler(w http.ResponseWriter, r *http.Request) {
	valStr := r.URL.Query().Get("val")
	if valStr == "" {
		http.Error(w, "Укажите параметр val, например ?val=50", http.StatusBadRequest)
		return
	}

	val, err := strconv.Atoi(valStr)
	if err != nil || val < 0 {
		http.Error(w, "Неверный формат числа. Ожидается положительное целое число.", http.StatusBadRequest)
		return
	}

	// возвращает предыдущее значение
	oldVal := debug.SetGCPercent(val)
	fmt.Fprintf(w, "GC Percent успешно изменен с %d на %d\n", oldVal, val)
}

func main() {
	// Устанавливаем дефолтное значение сборщика мусора (100 - стандарт в Go)
	debug.SetGCPercent(100)

	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/set_gc", setGCHandler)

	port := ":8080"
	fmt.Println("🚀 Сервер экспортера запущен на порту", port)
	fmt.Println("📊 Метрики Prometheus: http://localhost:8080/metrics")
	fmt.Println("⚙️  Управление GC:      http://localhost:8080/set_gc?val=50")
	fmt.Println("🔍 Профилирование:    http://localhost:8080/debug/pprof/")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
