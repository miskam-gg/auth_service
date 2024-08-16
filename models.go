package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	CreatedAt time.Time `json:"created_at"`
	IPAddress string    `json:"ip_address"`
}

func SaveRefreshToken(db *sql.DB, userID string, tokenHash string, ipAddress string) error {
	query := `
        INSERT INTO refresh_tokens (id, user_id, token_hash, created_at, ip_address)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := db.Exec(query, uuid.New().String(), userID, tokenHash, time.Now(), ipAddress)
	return err
}

func GetRefreshToken(db *sql.DB, userID string) (*RefreshToken, error) {
	query := `
        SELECT id, user_id, token_hash, created_at, ip_address
        FROM refresh_tokens
        WHERE user_id = $1
    `
	var token RefreshToken
	row := db.QueryRow(query, userID)
	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.CreatedAt, &token.IPAddress)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
