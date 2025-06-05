package auth

import (
	"database/sql"
)

const (
	// MasterToken is the special token that has administrative privileges
	MasterToken = "MASTER_TOKEN"
)

// Validator handles token validation
type Validator struct {
	db *sql.DB
}

// NewValidator creates a new token validator
func NewValidator(db *sql.DB) *Validator {
	return &Validator{db: db}
}

// IsValidToken checks if a token is valid
func (v *Validator) IsValidToken(token string) bool {
	var count int
	_ = v.db.QueryRow("SELECT COUNT(*) FROM clients WHERE token = ? AND is_active = 1", token).Scan(&count)
	return count > 0
}

// IsMasterToken checks if a token is the master token
func IsMasterToken(token string) bool {
	return token == MasterToken
}
