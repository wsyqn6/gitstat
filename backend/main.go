package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gitstat/internal/scanner"
	"gitstat/internal/server"
	"gitstat/internal/store"
)

var Version = "dev"

//go:embed web/dist/*
var embeddedDist embed.FS

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "--debug":
		runDebug()
	case "--version", "-v":
		fmt.Println(Version)
	case "--help", "-h", "help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func scanAndRegister(scanPath string) {
	if scanPath == "" {
		return
	}
	repos, err := scanner.DiscoverRepos(scanPath)
	if err != nil {
		log.Printf("Warning: Failed to discover repos in %s: %v", scanPath, err)
		return
	}
	store.GlobalStore.SetScanPath(scanPath)
	store.GlobalStore.RegisterRepos(repos)
	log.Printf("Registered %d repos from %s", len(repos), scanPath)
}

func runDebug() {
	path := os.Getenv("GITSTAT_DEV_PATH")
	if path == "" {
		path = "D:/work/ems"
	}
	scanAndRegister(filepath.Clean(path))
	fmt.Println("GitStat server starting on :12580...")
	log.Fatal(http.ListenAndServe(":12580", server.NewServer(Version)))
}

func runServe(args []string) {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	port := serveCmd.Int("port", 12580, "Port number")
	serveCmd.Parse(args)

	var scanPath string
	if pos := serveCmd.Args(); len(pos) > 0 {
		scanPath = filepath.Clean(pos[0])
	} else if p, err := os.Getwd(); err == nil {
		scanPath = filepath.Clean(p)
	}
	scanAndRegister(scanPath)

	distFS, err := fs.Sub(embeddedDist, "web/dist")
	if err != nil {
		log.Printf("Warning: Failed to load embedded static files: %v", err)
		log.Println("Server will start with API only.")
	}

	r := server.NewServerWithStatic(distFS, Version)

	addr := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost:%d", *port)
	fmt.Println("GitStat Web Server")
	if scanPath != "" {
		fmt.Printf("Scan directory: %s\n", scanPath)
	}
	fmt.Printf("Listening on %s\n", url)
	fmt.Println("Open in browser to view analytics")

	go openBrowser(url)
	log.Fatal(http.ListenAndServe(addr, r))
}

func printHelp() {
	fmt.Println("GitStat - Git Repository Statistics Analyzer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitstat                     Show this help")
	fmt.Println("  gitstat serve [directory]   Start web server with embedded UI")
	fmt.Println("  gitstat --debug             Start API-only server (dev mode)")
	fmt.Println("  gitstat --version           Show version")
	fmt.Println("  gitstat --help              Show this help")
	fmt.Println()
	fmt.Println("Serve Options:")
	fmt.Println("  --port number       Port number (default: 12580)")
	fmt.Println()
	fmt.Println("Serve Arguments:")
	fmt.Println("  directory           Scan directory (default: current working directory)")
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Run()
}
