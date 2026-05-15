package main

import (
	"fmt"
	"log"
	"net/http"

	"gitstat/internal/server"
)

func main() {
	r := server.NewServer()

	fmt.Println("GitStat server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
