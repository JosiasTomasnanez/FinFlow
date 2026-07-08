package storage

import (
	"fmt"
	"sync"

	"github.com/josiastomasnanez/finflow/internal/model"
)

type Store interface {
	ListWallets() []model.Wallet
	GetWallet(id string) (model.Wallet, bool)
	SaveWallet(wallet model.Wallet) model.Wallet
	UpdateWallet(wallet model.Wallet) error
}

type MemoryStore struct {
	mu      sync.RWMutex
	wallets map[string]model.Wallet
	nextID  int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		wallets: make(map[string]model.Wallet),
		nextID:  1,
	}
}

func (m *MemoryStore) ListWallets() []model.Wallet {
	m.mu.RLock()
	defer m.mu.RUnlock()

	wallets := make([]model.Wallet, 0, len(m.wallets))
	for _, wallet := range m.wallets {
		wallets = append(wallets, wallet)
	}
	return wallets
}

func (m *MemoryStore) GetWallet(id string) (model.Wallet, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	wallet, ok := m.wallets[id]
	return wallet, ok
}

func (m *MemoryStore) SaveWallet(wallet model.Wallet) model.Wallet {
	m.mu.Lock()
	defer m.mu.Unlock()

	wallet.ID = fmt.Sprintf("w%d", m.nextID)
	m.nextID++
	m.wallets[wallet.ID] = wallet
	return wallet
}

func (m *MemoryStore) UpdateWallet(wallet model.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, found := m.wallets[wallet.ID]; !found {
		return fmt.Errorf("wallet %s not found", wallet.ID)
	}
	m.wallets[wallet.ID] = wallet
	return nil
}
