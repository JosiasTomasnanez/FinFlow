package storage

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/josiastomasnanez/finflow/internal/model"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens or creates a SQLite database at dbPath and prepares the schema.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH is required")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS wallets (
        id TEXT PRIMARY KEY,
        owner TEXT NOT NULL,
        balance INTEGER NOT NULL
    );`
	if _, err := db.Exec(createTable); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create wallets table: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) ListWallets() []model.Wallet {
	rows, err := s.db.Query("SELECT id, owner, balance FROM wallets")
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	wallets := make([]model.Wallet, 0)
	for rows.Next() {
		var wallet model.Wallet
		if err := rows.Scan(&wallet.ID, &wallet.Owner, &wallet.Balance); err != nil {
			continue
		}
		wallets = append(wallets, wallet)
	}
	return wallets
}

func (s *SQLiteStore) GetWallet(id string) (model.Wallet, bool) {
	var wallet model.Wallet
	row := s.db.QueryRow("SELECT id, owner, balance FROM wallets WHERE id = ?", id)
	if err := row.Scan(&wallet.ID, &wallet.Owner, &wallet.Balance); err != nil {
		if err == sql.ErrNoRows {
			return model.Wallet{}, false
		}
		return model.Wallet{}, false
	}
	return wallet, true
}

func (s *SQLiteStore) SaveWallet(wallet model.Wallet) model.Wallet {
	if wallet.ID == "" {
		wallet.ID = generateWalletID()
	}
	_, err := s.db.Exec("INSERT INTO wallets (id, owner, balance) VALUES (?, ?, ?)", wallet.ID, wallet.Owner, wallet.Balance)
	if err != nil {
		return model.Wallet{}
	}
	return wallet
}

func (s *SQLiteStore) UpdateWallet(wallet model.Wallet) error {
	result, err := s.db.Exec("UPDATE wallets SET owner = ?, balance = ? WHERE id = ?", wallet.Owner, wallet.Balance, wallet.ID)
	if err != nil {
		return fmt.Errorf("update wallet: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check update result: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("wallet %s not found", wallet.ID)
	}
	return nil
}

func generateWalletID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("w%v", os.Getpid())
	}
	return "w" + hex.EncodeToString(bytes)
}
