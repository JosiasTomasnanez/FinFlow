package service

import (
	"errors"
	"fmt" 

	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

var (
	errInsufficientBalance = errors.New("insufficient balance")
	errWalletNotFound      = errors.New("wallet not found")
	errInvalidAmount       = errors.New("amount must be greater than zero")
)

type WalletService struct {
	store      storage.Store       
	redisStore *storage.RedisStore 
}

func NewWalletService(store storage.Store, redisStore *storage.RedisStore) *WalletService {
	return &WalletService{
		store:      store,
		redisStore: redisStore,
	}
}

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

	savedWallet := s.store.SaveWallet(wallet)
	if savedWallet.ID == "" {
		return model.Wallet{}, errors.New("failed to save wallet in postgres")
	}

	if s.redisStore != nil {
		if err := s.redisStore.SetWallet(savedWallet); err != nil {
			fmt.Printf("[REDIS ERROR] No se pudo cachear la nueva wallet: %v\n", err)
		} else {
			fmt.Printf("[⚡ REDIS WRITE] Wallet %s cacheada al crearse\n", savedWallet.ID)
		}
	}

	return savedWallet, nil
}

func (s *WalletService) GetWallet(id string) (model.Wallet, bool) {
	if s.redisStore != nil {
		if wallet, found := s.redisStore.GetWallet(id); found {
			fmt.Printf("🟢 [CACHE HIT] La wallet %s se obtuvo desde REDIS\n", id)
			return wallet, true
		}
	}

	fmt.Printf("🔴 [CACHE MISS] La wallet %s NO estaba en Redis. Buscando en POSTGRES...\n", id)
	wallet, found := s.store.GetWallet(id)
	if !found {
		return model.Wallet{}, false
	}

	if s.redisStore != nil {
		_ = s.redisStore.SetWallet(wallet)
	}

	return wallet, true
}

func (s *WalletService) ListWallets() []model.Wallet {
	return s.store.ListWallets()
}

func (s *WalletService) Transfer(fromID, toID string, amount int64) (model.PaymentResult, error) {
	if amount <= 0 {
		return model.PaymentResult{}, errInvalidAmount
	}
	if fromID == toID {
		return model.PaymentResult{}, errors.New("sender and receiver must differ")
	}

	fromWallet, ok := s.GetWallet(fromID)
	if !ok {
		return model.PaymentResult{}, errWalletNotFound
	}

	toWallet, ok := s.GetWallet(toID)
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

	if s.redisStore != nil {
		_ = s.redisStore.SetWallet(fromWallet)
		_ = s.redisStore.SetWallet(toWallet)
		fmt.Printf("[⚡ REDIS UPDATE] Saldos actualizados en caché post-transferencia\n")
	}

	return model.PaymentResult{
		FromWalletID: fromID,
		ToWalletID:   toID,
		Amount:       amount,
		Status:       "completed",
	}, nil
}
