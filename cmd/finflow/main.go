package main

import (
    "log"
    "net/http"

    "github.com/josiastomasnanez/finflow/internal/api"
    "github.com/josiastomasnanez/finflow/internal/service"
    "github.com/josiastomasnanez/finflow/internal/storage"
)

func main() {
    store := storage.NewMemoryStore()
    service := service.NewWalletService(store)
    server := api.NewServer(service)

    log.Println("starting FinFlow API on http://0.0.0.0:8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server failed: %v", err)
    }
}
