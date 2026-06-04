package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/josiastomasnanez/finflow/internal/api"
	"github.com/josiastomasnanez/finflow/internal/service"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

func main() {
	loadEnvFile()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/finflow.db"
	}

	sqliteStore, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize sqlite store: %v", err)
	}
	defer func() { _ = sqliteStore.Close() }()

	service := service.NewWalletService(sqliteStore)
	server := api.NewServer(service)

	log.Printf("starting FinFlow API on http://0.0.0.0:8080 using DB %s", dbPath)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}

	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "DB_PATH=") {
			value := strings.TrimPrefix(line, "DB_PATH=")
			if value != "" {
				_ = os.Setenv("DB_PATH", value)
			}
			return
		}
	}
}
