package auth_test

import (
	"database/sql"
	"testing"

	"redix/pkg/auth"

	_ "github.com/mattn/go-sqlite3"
)

// mockDB is a mock implementation of sql.DB for testing
type mockDB struct {
	validTokens map[string]bool
}

// These methods are required by sql.DB but not used in our tests
func (m *mockDB) Close() error                                               { return nil }
func (m *mockDB) Ping() error                                                { return nil }
func (m *mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) { return nil, nil }
func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) { return nil, nil }
func (m *mockDB) Begin() (*sql.Tx, error)                                    { return nil, nil }
func (m *mockDB) Prepare(query string) (*sql.Stmt, error)                    { return nil, nil }
func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row        { return &sql.Row{} }

func TestNewValidator(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	validator := auth.NewValidator(db)
	if validator == nil {
		t.Error("NewValidator() returned nil")
	}
}

func TestIsValidToken(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		validToken bool
		want       bool
	}{
		{
			name:       "valid token",
			token:      "valid-token",
			validToken: true,
			want:       true,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			validToken: false,
			want:       false,
		},
		{
			name:       "empty token",
			token:      "",
			validToken: false,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a real sql.DB for testing
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Failed to create test database: %v", err)
			}
			defer db.Close()

			// Create the clients table
			_, err = db.Exec(`
				CREATE TABLE clients (
					token TEXT PRIMARY KEY,
					is_active INTEGER
				)
			`)
			if err != nil {
				t.Fatalf("Failed to create test table: %v", err)
			}

			// Insert test data
			if tt.validToken {
				_, err = db.Exec("INSERT INTO clients (token, is_active) VALUES (?, 1)", tt.token)
				if err != nil {
					t.Fatalf("Failed to insert test data: %v", err)
				}
			}

			validator := auth.NewValidator(db)
			got := validator.IsValidToken(tt.token)
			if got != tt.want {
				t.Errorf("IsValidToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMasterToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{
			name:  "master token",
			token: auth.MasterToken,
			want:  true,
		},
		{
			name:  "non-master token",
			token: "regular-token",
			want:  false,
		},
		{
			name:  "empty token",
			token: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := auth.IsMasterToken(tt.token)
			if got != tt.want {
				t.Errorf("IsMasterToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
