package service

import (
	"context"
	"testing"

	"github.com/josiastomasnanez/finflow/internal/storage"
)

func TestCreateAndTransfer(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStore()
	service := NewWalletService(store, nil)

	walletA, err := service.CreateWallet(ctx, "alice", 1000)
	if err != nil {
		t.Fatalf("failed to create wallet A: %v", err)
	}
	walletB, err := service.CreateWallet(ctx, "bob", 500)
	if err != nil {
		t.Fatalf("failed to create wallet B: %v", err)
	}

	if walletA.Balance != 1000 {
		t.Fatalf("unexpected balance for wallet A: got %d", walletA.Balance)
	}
	if walletB.Balance != 500 {
		t.Fatalf("unexpected balance for wallet B: got %d", walletB.Balance)
	}

	result, err := service.Transfer(ctx, walletA.ID, walletB.ID, 300)
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	if result.Amount != 300 {
		t.Fatalf("expected transferred amount 300, got %d", result.Amount)
	}

	updatedA, _ := service.GetWallet(ctx, walletA.ID)
	updatedB, _ := service.GetWallet(ctx, walletB.ID)

	if updatedA.Balance != 700 {
		t.Fatalf("expected wallet A balance 700, got %d", updatedA.Balance)
	}
	if updatedB.Balance != 800 {
		t.Fatalf("expected wallet B balance 800, got %d", updatedB.Balance)
	}
}

func TestTransferInsufficientBalance(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStore()
	service := NewWalletService(store, nil)

	walletA, err := service.CreateWallet(ctx, "carla", 100)
	if err != nil {
		t.Fatalf("failed to create wallet A: %v", err)
	}
	walletB, err := service.CreateWallet(ctx, "diego", 100)
	if err != nil {
		t.Fatalf("failed to create wallet B: %v", err)
	}

	_, err = service.Transfer(ctx, walletA.ID, walletB.ID, 200)
	if err == nil {
		t.Fatal("expected transfer to fail due to insufficient funds")
	}
}

func TestCreateWalletValidation(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStore()
	service := NewWalletService(store, nil)

	_, err := service.CreateWallet(ctx, "", 100)
	if err == nil {
		t.Fatal("expected error when owner is empty")
	}

	_, err = service.CreateWallet(ctx, "elena", -10)
	if err == nil {
		t.Fatal("expected error when initial balance is negative")
	}
}
