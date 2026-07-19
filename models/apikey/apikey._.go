package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	SecretPrefix = "hup_"
	secretBytes  = 32
)

type APIKey struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	SecretHash string    `json:"-"`

	LastUsedAt pgtype.Timestamp `json:"last_used_at"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
	DeletedAt  pgtype.Timestamp `json:"deleted_at"`
}

func New(userID uuid.UUID, name string) (key *APIKey, secret string, err error) {
	raw := make([]byte, secretBytes)
	if _, err = rand.Read(raw); err != nil {
		return nil, "", err
	}

	secret = SecretPrefix + base64.RawURLEncoding.EncodeToString(raw)

	key = new(APIKey)
	key.UserID = userID
	key.Name = name
	key.SecretHash = HashSecret(secret)

	return key, secret, nil
}

func HashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func IsSecret(s string) bool {
	return strings.HasPrefix(s, SecretPrefix)
}
