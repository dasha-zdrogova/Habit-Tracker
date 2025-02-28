package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type tokenData struct {
	UserID    int
	ExpiredAt time.Time
}

type TokenManager struct {
	tokens map[string]tokenData
	mu     sync.RWMutex
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: map[string]tokenData{},
		mu:     sync.RWMutex{},
	}
}

func (tm *TokenManager) AddToken(token string, userID int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tokens[token] = tokenData{
		UserID:    userID,
		ExpiredAt: time.Now().Add(72 * time.Hour),
	}
}

func (tm *TokenManager) ValidateToken(token string) (int, error) {
	tm.mu.RLock()
	data, exist := tm.tokens[token]
	tm.mu.RUnlock()

	if !exist || time.Now().After(data.ExpiredAt) {
		return 0, ErrInvalidToken
	}

	data.ExpiredAt = time.Now().Add(72 * time.Hour)
	tm.mu.Lock()
	tm.tokens[token] = data
	tm.mu.Unlock()

	return data.UserID, nil
}

func (tm *TokenManager) StartCleanUp () {
	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		tm.cleanUp()
	}
}

func (tm *TokenManager) cleanUp () {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	now := time.Now()
	for token, data := range tm.tokens {
		if now.After(data.ExpiredAt) {
			delete(tm.tokens, token)
		}
	}
}

func (*TokenManager) GenerateToken () string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}