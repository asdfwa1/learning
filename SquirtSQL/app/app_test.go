package app

import (
	"os"
	"strings"
	"testing"
	"v4/database/actions"
	"v4/database/parser"
	"v4/storage"
)

func setupTestApp(t *testing.T) (*App, string) {
	tempDir, err := os.MkdirTemp("", "app_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storage := storage.NewCSVStorage(tempDir)
	db := actions.NewDatabase(storage)

	return &App{
		DB:      db,
		Storage: storage,
	}, tempDir
}

func cleanupTestApp(tempDir string) {
	_ = os.RemoveAll(tempDir)
}

func TestHandleCreateTable(t *testing.T) {
	app, tempDir := setupTestApp(t)
	defer cleanupTestApp(tempDir)

	tests := []struct {
		name       string
		input      string
		wantErr    bool
		errMessage string
	}{
		{
			name:    "Valid table creation",
			input:   "CREATE TABLE users name,email",
			wantErr: false,
		},
		{
			name:       "Duplicate table",
			input:      "CREATE TABLE users name,email",
			wantErr:    true,
			errMessage: "таблица users уже существует",
		},
		{
			name:    "Invalid syntax - missing fields",
			input:   "CREATE TABLE products",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)
			if parseErr != nil && !tt.wantErr {
				t.Fatalf("Parse error: %v", parseErr)
			}
			err := app.HandleCreateTable(query)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMessage) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
