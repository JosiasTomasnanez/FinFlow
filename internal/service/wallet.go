package service

import (
	"errors"
	"fmt" // <-- Agregado para los prints en consola

	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/storage"
)

var (
	errInsufficientBalance = errors.New("insufficient balance")
	errWalletNotFound      = errors.New("wallet not found")
	errInvalidAmount       = errors.New("amount must be greater than zero")
)

// WalletService contiene las reglas de negocio y ahora coordina DB y Caché
type WalletService struct {
	store      storage.Store       // Nuestra DB principal (Postgres)
	redisStore *storage.RedisStore // Nuestra Caché (Redis)
}

// NewWalletService ahora recibe también el storage de Redis
func NewWalletService(store storage.Store, redisStore *storage.RedisStore) *WalletService {
	return &WalletService{
		store:      store,
		redisStore: redisStore,
	}
}

// CreateWallet registra en Postgres y hace Write-Through en Redis
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

	// 1. Guardamos en la base de datos Postgres obligatoriamente
	savedWallet := s.store.SaveWallet(wallet)
	if savedWallet.ID == "" {
		return model.Wallet{}, errors.New("failed to save wallet in postgres")
	}

	// 2. Guardamos en Redis de forma transparente (Write-Through)
	if s.redisStore != nil {
		if err := s.redisStore.SetWallet(savedWallet); err != nil {
			fmt.Printf("[REDIS ERROR] No se pudo cachear la nueva wallet: %v\n", err)
		} else {
			fmt.Printf("[⚡ REDIS WRITE] Wallet %s cacheada al crearse\n", savedWallet.ID)
		}
	}

	return savedWallet, nil
}

// GetWallet implementa Cache-Aside con los logs distintivos que pediste
func (s *WalletService) GetWallet(id string) (model.Wallet, bool) {
	// 1. Intentar buscar en Redis primero
	if s.redisStore != nil {
		if wallet, found := s.redisStore.GetWallet(id); found {
			fmt.Printf("🟢 [CACHE HIT] La wallet %s se obtuvo desde REDIS\n", id)
			return wallet, true
		}
	}

	// 2. Si no está en Redis, es un Cache Miss -> Vamos a Postgres
	fmt.Printf("🔴 [CACHE MISS] La wallet %s NO estaba en Redis. Buscando en POSTGRES...\n", id)
	wallet, found := s.store.GetWallet(id)
	if !found {
		return model.Wallet{}, false
	}

	// 3. Lo encontramos en Postgres, así que lo guardamos en Redis para la próxima consulta
	if s.redisStore != nil {
		_ = s.redisStore.SetWallet(wallet)
	}

	return wallet, true
}

// ListWallets (Podríamos cachear la lista completa, pero por simplicidad de este paso, lee directo de DB)
func (s *WalletService) ListWallets() []model.Wallet {
	return s.store.ListWallets()
}

// Transfer mueve fondos, actualiza DB e invalida o actualiza la caché
func (s *WalletService) Transfer(fromID, toID string, amount int64) (model.PaymentResult, error) {
	if amount <= 0 {
		return model.PaymentResult{}, errInvalidAmount
	}
	if fromID == toID {
		return model.PaymentResult{}, errors.New("sender and receiver must differ")
	}

	// Usamos el método interno del servicio para que aproveche la lógica de logs
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

	// Si las actualizaciones en Postgres salieron bien, actualizamos Redis
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
