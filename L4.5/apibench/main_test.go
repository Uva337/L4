package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// замеряет скорость работы нашей плохой функции
func BenchmarkSlowHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/process-slow", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		SlowHandler(w, req)
	}
}

func BenchmarkFastHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/process-fast", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		FastHandler(w, req)
	}
}
