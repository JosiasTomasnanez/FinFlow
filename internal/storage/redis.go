package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(redisURL string) (*RedisStore, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear url de redis: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("no se pudo conectar a Redis: %w", err)
	}

	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (r *RedisStore) SetWallet(wallet model.Wallet) error {
	key := fmt.Sprintf("wallet:%s", wallet.ID)

	data, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("error serializando wallet: %w", err)
	}

	err = r.client.Set(r.ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error guardando en redis: %w", err)
	}
	return nil
}

func (r *RedisStore) GetWallet(id string) (model.Wallet, bool) {
	key := fmt.Sprintf("wallet:%s", id)

	data, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return model.Wallet{}, false
	} else if err != nil {
		fmt.Printf("[REDIS ERROR] error obteniendo wallet %s: %v\n", id, err)
		return model.Wallet{}, false
	}

	var wallet model.Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return model.Wallet{}, false
	}

	return wallet, true
}
