package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"time"
)

func GenerateToken(userId int64, ttl time.Duration, scope string) *domain.Token {
	var token domain.Token
	token.Plaintext = rand.Text()
	token.UserId = userId
	token.Expiry = time.Now().Add(ttl)
	token.Scope = scope
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return &token
}
