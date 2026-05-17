package worker

import (
	"encoding/json"
	"log"
	"net/http"

	"dgrep/engine"
)

// описывает структуру данных, которую воркер ждет от координатора
type GrepRequest struct {
	Pattern string   `json:"pattern"`
	Lines   []string `json:"lines"`
}

func StartServer(port string) {
	http.HandleFunc("/grep", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var req GrepRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка чтения JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Получена задача: %d строк, ищем '%s'", len(req.Lines), req.Pattern)

		matches := engine.LocalGrep(req.Lines, req.Pattern, 4)
		if matches == nil {
			matches = []engine.Match{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	})

	log.Printf("🚀 Воркер запущен и слушает порт %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
