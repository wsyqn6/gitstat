package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitstat/internal/scanner"
	"gitstat/internal/server"
	"gitstat/internal/store"
)

const Version = "v0.2.0"

func main() {
	// 默认扫描路径（开发调试用）
	defaultPath := "D:/work"
	repos, err := scanner.DiscoverRepos(defaultPath)
	if err != nil {
		log.Printf("Warning: Failed to discover repos in %s: %v", defaultPath, err)
	} else {
		store.GlobalStore.RegisterRepos(repos)
		log.Printf("Registered %d repos from %s", len(repos), defaultPath)
	}

	r := server.NewServer()

	// 版本接口
	r.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": Version})
	})

	fmt.Println("GitStat server starting on :12580...")
	log.Fatal(http.ListenAndServe(":12580", r))
}
