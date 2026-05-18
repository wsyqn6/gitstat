package main

import (
	"fmt"
	"log"
	"net/http"

	"gitstat/internal/server"
)

func main() {
	r := server.NewServer()

	fmt.Println("GitStat server starting on :12580...")
	log.Fatal(http.ListenAndServe(":12580", r))
}
