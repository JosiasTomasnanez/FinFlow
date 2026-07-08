package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type AuthService struct {
	users  map[string]string
	tokens map[string]string
	mu     sync.RWMutex
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: map[string]string{
			"admin": "password",
		},
		tokens: make(map[string]string),
	}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("username and password are required")
	}

	s.mu.RLock()
	expected, ok := s.users[username]
	s.mu.RUnlock()
	if !ok || expected != password {
		return "", fmt.Errorf("invalid credentials")
	}

	token := generateToken()
	s.mu.Lock()
	s.tokens[token] = username
	s.mu.Unlock()

	return token, nil
}

func (s *AuthService) ValidateToken(token string) (string, bool) {
	s.mu.RLock()
	username, ok := s.tokens[token]
	s.mu.RUnlock()
	return username, ok
}

func generateToken() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "token-fallback"
	}
	return hex.EncodeToString(bytes)
}
