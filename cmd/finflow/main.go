package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"time" // <-- Agregado para el manejo del Timeout

	"github.com/Unleash/unleash-client-go/v4"
	"github.com/josiastomasnanez/finflow/internal/api"
	"github.com/josiastomasnanez/finflow/internal/service"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

// Creamos una estructura simple para escuchar los eventos de Unleash sin llenarte de código
type unleashLogger struct{}

func (l *unleashLogger) OnReady() {
	log.Println("=== [UNLEASH] El cliente se sincronizó exitosamente y está LISTO ===")
}
func (l *unleashLogger) OnError(err error)                       { log.Printf("=== [UNLEASH ERROR] %v ===", err) }
func (l *unleashLogger) OnWarning(err error)                     { log.Printf("=== [UNLEASH WARNING] %v ===", err) }
func (l *unleashLogger) OnCount(name string, enabled bool)       {}
func (l *unleashLogger) OnSent(payload unleash.MetricsData)      {}
func (l *unleashLogger) OnRegistered(payload unleash.ClientData) {}

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

	unleashURL := os.Getenv("UNLEASH_URL")
	unleashToken := os.Getenv("UNLEASH_TOKEN")

	if unleashURL != "" && unleashToken != "" {
		log.Printf("Initializing Unleash with URL: %s", unleashURL)

		uLogger := &unleashLogger{}

		err := unleash.Initialize(
			unleash.WithAppName("finflow-backend"),
			unleash.WithUrl(unleashURL),
			unleash.WithCustomHeaders(http.Header{
				"Authorization": []string{unleashToken},
			}),
			unleash.WithListener(uLogger), // <-- Vinculamos nuestro escuchador de logs
		)
		if err != nil {
			log.Printf("Warning: failed to initialize Unleash: %v", err)
		} else {
			defer func() { _ = unleash.Close() }()

			// Esperamos 3 segundos a que el hilo asincrónico traiga las flags antes de arrancar Gin
			log.Println("Waiting for Unleash to fetch initial flags...")
			time.Sleep(3 * time.Second)
		}
	} else {
		log.Println("Warning: UNLEASH_URL or UNLEASH_TOKEN not found. Feature flags will default to false.")
	}

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
