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

	walletService := service.NewWalletService(sqliteStore)
	authService := service.NewAuthService()
	server := api.NewServer(walletService, authService)

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
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			_ = os.Setenv(key, value)
		}
	}
}
