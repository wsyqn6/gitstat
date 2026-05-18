package server

import (
	"net/http"

	"gitstat/internal/handler"

	"github.com/go-chi/chi/v5"
)

func NewServer() *chi.Mux {
	r := chi.NewRouter()

	// CORS中间件
	r.Use(corsMiddleware)

	// 路由注册
	r.Post("/api/scan", handler.ScanHandler)
	r.Get("/api/repositories", handler.GetRepositoriesHandler)
	r.Get("/api/stats/overview", handler.GetOverviewStatsHandler)
	r.Get("/api/stats/daily", handler.GetDailyStatsHandler)
	r.Post("/api/export/json", handler.ExportDataHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
