package service

import (
	"errors"

	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

var (
	errInsufficientBalance = errors.New("insufficient balance")
	errWalletNotFound      = errors.New("wallet not found")
	errInvalidAmount       = errors.New("amount must be greater than zero")
)

// WalletService contains the business rules for wallet operations.
type WalletService struct {
	store storage.Store
}

// NewWalletService constructs a wallet service.
func NewWalletService(store storage.Store) *WalletService {
	return &WalletService{store: store}
}

// CreateWallet registers a new wallet with an initial balance.
func (s *WalletService) CreateWallet(owner string, initialBalance int64) (model.Wallet, error) {
	if owner == "" {
		return model.Wallet{}, errors.New("owner is required")
	}
	if initialBalance < 0 {
		return model.Wallet{}, errors.New("initial balance cannot be negative")
	}

	wallet := model.Wallet{
		Owner:   owner,
		Balance: initialBalance,
	}
	return s.store.SaveWallet(wallet), nil
}

// GetWallet retrieves a wallet by ID.
func (s *WalletService) GetWallet(id string) (model.Wallet, bool) {
	return s.store.GetWallet(id)
}

// ListWallets returns all wallets in the system.
func (s *WalletService) ListWallets() []model.Wallet {
	return s.store.ListWallets()
}

// Transfer moves funds from one wallet to another.
func (s *WalletService) Transfer(fromID, toID string, amount int64) (model.PaymentResult, error) {
	if amount <= 0 {
		return model.PaymentResult{}, errInvalidAmount
	}
	if fromID == toID {
		return model.PaymentResult{}, errors.New("sender and receiver must differ")
	}

	fromWallet, ok := s.store.GetWallet(fromID)
	if !ok {
		return model.PaymentResult{}, errWalletNotFound
	}

	toWallet, ok := s.store.GetWallet(toID)
	if !ok {
		return model.PaymentResult{}, errWalletNotFound
	}

	if fromWallet.Balance < amount {
		return model.PaymentResult{}, errInsufficientBalance
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	if err := s.store.UpdateWallet(fromWallet); err != nil {
		return model.PaymentResult{}, err
	}
	if err := s.store.UpdateWallet(toWallet); err != nil {
		return model.PaymentResult{}, err
	}

	return model.PaymentResult{
		FromWalletID: fromID,
		ToWalletID:   toID,
		Amount:       amount,
		Status:       "completed",
	}, nil
}
