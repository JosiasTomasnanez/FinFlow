package storage

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/josiastomasnanez/finflow/internal/model"
)

type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore connects to Postgres using a DSN and ensures schema exists.
func NewPostgresStore(dsn string) (*PostgresStore, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres db: %w", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS wallets (
        id TEXT PRIMARY KEY,
        owner TEXT NOT NULL,
        balance BIGINT NOT NULL
    );`
	if _, err := db.Exec(createTable); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create wallets table: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) ListWallets() []model.Wallet {
	rows, err := s.db.Query("SELECT id, owner, balance FROM wallets")
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	wallets := make([]model.Wallet, 0)
	for rows.Next() {
		var w model.Wallet
		if err := rows.Scan(&w.ID, &w.Owner, &w.Balance); err != nil {
			continue
		}
		wallets = append(wallets, w)
	}
	return wallets
}

func (s *PostgresStore) GetWallet(id string) (model.Wallet, bool) {
	var w model.Wallet
	row := s.db.QueryRow("SELECT id, owner, balance FROM wallets WHERE id = $1", id)
	if err := row.Scan(&w.ID, &w.Owner, &w.Balance); err != nil {
		if err == sql.ErrNoRows {
			return model.Wallet{}, false
		}
		return model.Wallet{}, false
	}
	return w, true
}

func (s *PostgresStore) SaveWallet(wallet model.Wallet) model.Wallet {
	if wallet.ID == "" {
		wallet.ID = generateWalletID()
	}
	_, err := s.db.Exec("INSERT INTO wallets (id, owner, balance) VALUES ($1, $2, $3)", wallet.ID, wallet.Owner, wallet.Balance)
	if err != nil {
		return model.Wallet{}
	}
	return wallet
}

func (s *PostgresStore) UpdateWallet(wallet model.Wallet) error {
	result, err := s.db.Exec("UPDATE wallets SET owner = $1, balance = $2 WHERE id = $3", wallet.Owner, wallet.Balance, wallet.ID)
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
