package server

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"gitstat/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(version string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)
	r.Use(middleware.Compress(5))

	r.Post("/api/scan/path", handler.SetScanPathHandler)
	r.Get("/api/scan/path", handler.GetScanPathHandler)

	r.Get("/api/stats/overview", handler.GetOverviewStatsHandler)
	r.Get("/api/stats/dashboard", handler.GetDashboardHandler)
	r.Get("/api/stats/daily", handler.GetStatsHandler("daily"))
	r.Get("/api/stats/daily-trend", handler.GetDailyTrendHandler)
	r.Get("/api/stats/weekly", handler.GetStatsHandler("weekly"))
	r.Get("/api/stats/monthly", handler.GetStatsHandler("monthly"))
	r.Get("/api/stats/yearly", handler.GetStatsHandler("yearly"))
	r.Get("/api/stats/authors", handler.GetAuthorRankHandler)
	r.Get("/api/stats/activity-heatmap", handler.GetActivityHeatmapHandler)
	r.Get("/api/stats/repo-comparison", handler.GetRepoComparisonHandler)
	r.Get("/api/stats/file-ranking", handler.GetFileRankingHandler)
	r.Get("/api/stats/compare", handler.GetComparisonHandler)
	r.Post("/api/export/json", handler.ExportDataHandler)
	r.Get("/api/repos/list", handler.GetReposListHandler)
	r.Get("/api/repos/info", handler.GetRepoInfoHandler)
	r.Get("/api/repos/stats", handler.GetRepoStatsHandler)
	r.Get("/api/repos/commits", handler.GetRepoCommitsHandler)
	r.Post("/api/repos/analyze", handler.GetRepoAnalyzeHandler)
	r.Get("/api/repos/tags", handler.GetRepoTagsHandler)
	r.Get("/api/repos/chart", handler.GetRepoChartHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": version})
	})

	return r
}

func NewServerWithStatic(staticFS fs.FS, version string) *chi.Mux {
	r := NewServer(version)

	if staticFS == nil {
		return r
	}

	r.Handle("/*", spaHandler(staticFS))

	return r
}

func spaHandler(staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// try to serve the requested static file
		f, err := staticFS.Open(path)
		if err == nil {
			defer f.Close()
			info, err := f.Stat()
			if err == nil && !info.IsDir() {
				http.ServeFileFS(w, r, staticFS, path)
				return
			}
		}

		// fallback to index.html for SPA routing
		_, err = staticFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.ServeFileFS(w, r, staticFS, "index.html")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%v)", r.Method, r.URL.Path, r.URL.RawQuery, time.Since(start))
	})
}
