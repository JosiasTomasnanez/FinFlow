package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"time" 

	"github.com/Unleash/unleash-client-go/v4"
	"github.com/josiastomasnanez/finflow/internal/api"
	"github.com/josiastomasnanez/finflow/internal/service"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

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

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DB_PATH")
	}
	if dbURL == "" {
		log.Fatalf("DATABASE_URL or DB_PATH environment variable must be set to a Postgres DSN")
	}

	pgStore, err := storage.NewPostgresStore(dbURL)
	if err != nil {
		log.Fatalf("failed to initialize postgres store: %v", err)
	}
	defer func() { _ = pgStore.Close() }()

	var redisStore *storage.RedisStore
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		log.Printf("Connecting to Redis on: %s", redisURL)
		redisStore, err = storage.NewRedisStore(redisURL)
		if err != nil {
			log.Printf("Warning: failed to initialize redis store: %v. Proceeding without cache.", err)
		} else {
			log.Println("🟢 Redis conectado exitosamente.")
		}
	} else {
		log.Println("Warning: REDIS_URL not found. App will run without caching mechanism.")
	}

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
			unleash.WithListener(uLogger),
		)
		if err != nil {
			log.Printf("Warning: failed to initialize Unleash: %v", err)
		} else {
			defer func() { _ = unleash.Close() }()

			log.Println("Waiting for Unleash to fetch initial flags...")
			time.Sleep(3 * time.Second)
		}
	} else {
		log.Println("Warning: UNLEASH_URL or UNLEASH_TOKEN not found. Feature flags will default to false.")
	}

	walletService := service.NewWalletService(pgStore, redisStore)
	authService := service.NewAuthService()
	server := api.NewServer(walletService, authService)

	log.Printf("starting FinFlow API on http://0.0.0.0:8080 using DB %s", dbURL)
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
