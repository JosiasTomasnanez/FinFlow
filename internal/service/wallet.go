package service

import (
	"context"
	"errors"

	"github.com/josiastomasnanez/finflow/internal/logging"
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

// Todos los métodos reciben ctx para poder loguear con el request_id que
// viene del handler HTTP (ver internal/api/middleware.go), así los logs
// de esta capa quedan correlacionados con los del handler que los llamó.

func (s *WalletService) CreateWallet(ctx context.Context, owner string, initialBalance int64) (model.Wallet, error) {
	log := logging.FromContext(ctx)

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
			log.Warn("redis_cache_write_failed", "wallet_id", savedWallet.ID, "error", err.Error())
		} else {
			log.Info("redis_cache_write", "wallet_id", savedWallet.ID)
		}
	}

	return savedWallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id string) (model.Wallet, bool) {
	log := logging.FromContext(ctx)

	if s.redisStore != nil {
		if wallet, found := s.redisStore.GetWallet(id); found {
			log.Info("redis_cache_hit", "wallet_id", id)
			return wallet, true
		}
	}

	log.Info("redis_cache_miss", "wallet_id", id)
	wallet, found := s.store.GetWallet(id)
	if !found {
		return model.Wallet{}, false
	}

	if s.redisStore != nil {
		_ = s.redisStore.SetWallet(wallet)
	}

	return wallet, true
}

func (s *WalletService) ListWallets(ctx context.Context) []model.Wallet {
	return s.store.ListWallets()
}

func (s *WalletService) Transfer(ctx context.Context, fromID, toID string, amount int64) (model.PaymentResult, error) {
	log := logging.FromContext(ctx)

	if amount <= 0 {
		return model.PaymentResult{}, errInvalidAmount
	}
	if fromID == toID {
		return model.PaymentResult{}, errors.New("sender and receiver must differ")
	}

	fromWallet, ok := s.GetWallet(ctx, fromID)
	if !ok {
		return model.PaymentResult{}, errWalletNotFound
	}

	toWallet, ok := s.GetWallet(ctx, toID)
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
		log.Info("redis_cache_update_post_transfer", "from_wallet_id", fromID, "to_wallet_id", toID)
	}

	return model.PaymentResult{
		FromWalletID: fromID,
		ToWalletID:   toID,
		Amount:       amount,
		Status:       "completed",
	}, nil
}
