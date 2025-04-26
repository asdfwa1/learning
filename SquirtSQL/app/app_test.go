package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
		name             string
		input            string
		wantErr          bool
		errMessage       string
		expectParseError bool
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
			name:             "Invalid syntax - missing fields",
			input:            "CREATE TABLE products",
			wantErr:          true,
			errMessage:       "не указаны поля таблицы",
			expectParseError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)
			if tt.expectParseError {
				if parseErr == nil {
					t.Error("Expected parse error, got nil")
				} else if !strings.Contains(parseErr.Error(), tt.errMessage) {
					t.Errorf("Expected parse error to contain '%s', got '%s'",
						tt.errMessage, parseErr.Error())
				}
				return
			}

			if parseErr != nil {
				t.Fatalf("Unexpected parse error: %v", parseErr)
			}

			app.HandleCreateTable(query)

			if !tt.wantErr {
				if !app.Storage.TableExist(query.Table) {
					t.Errorf("Table %s should exist", query.Table)
				}
			}
		})
	}
}

func TestHandleSelect(t *testing.T) {
	app, tempDir := setupTestApp(t)
	defer cleanupTestApp(tempDir)

	createQuery := &parser.Query{
		Type:   parser.QueryCreateTable,
		Table:  "users",
		Fields: []string{"name", "email"},
	}
	app.HandleCreateTable(createQuery)

	insertQuery := &parser.Query{
		Type:   parser.QueryInsert,
		Table:  "users",
		Fields: []string{"Kolya", "kolya@mail.ru"},
	}
	_, _ = app.DB.Insert(insertQuery.Table, insertQuery.Fields)

	tests := []struct {
		name        string
		input       string
		wantError   bool
		errContains string
	}{
		{
			name:      "Select all records - valid",
			input:     "SELECT users *",
			wantError: false,
		},
		{
			name:      "Select by ID - valid",
			input:     "SELECT users 1",
			wantError: false,
		},
		{
			name:        "Non-existent table",
			input:       "SELECT unknown *",
			wantError:   true,
			errContains: "таблица unknown не найдена",
		},
		{
			name:        "Non-existent record",
			input:       "SELECT users 999",
			wantError:   true,
			errContains: "Записи не найдены",
		},
		{
			name:        "Empty table",
			input:       "SELECT empty *",
			wantError:   true,
			errContains: "Записи не найдены",
		},
		{
			name:        "Invalid query syntax",
			input:       "SELECT users",
			wantError:   true,
			errContains: "SELECT <table> <id> or <*>",
		},
	}

	app.HandleCreateTable(&parser.Query{
		Type:   parser.QueryCreateTable,
		Table:  "empty",
		Fields: []string{"field"},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)

			if parseErr != nil {
				if !tt.wantError {
					t.Fatalf("Unexpected parse error: %v", parseErr)
				}
				if !strings.Contains(parseErr.Error(), tt.errContains) {
					t.Errorf("Expected parse error to contain '%s', got '%s'",
						tt.errContains, parseErr.Error())
				}
				return
			}

			app.handleSelect(query)

			if !tt.wantError {
				if !app.Storage.TableExist(query.Table) {
					t.Errorf("Table %s should exist", query.Table)
				}

				if query.ID != -1 {
					_, err := app.DB.Select(query.Table, query.ID)
					if err != nil {
						t.Errorf("Record %d should exist in table %s", query.ID, query.Table)
					}
				}
			}
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	app, tempDir := setupTestApp(t)
	defer cleanupTestApp(tempDir)

	createQuery := &parser.Query{
		Type:   parser.QueryCreateTable,
		Table:  "users",
		Fields: []string{"name", "email"},
	}
	app.HandleCreateTable(createQuery)

	insertQuery := &parser.Query{
		Type:   parser.QueryInsert,
		Table:  "users",
		Fields: []string{"kolya", "kolya@mail.ru"},
	}
	insertedID, _ := app.DB.Insert(insertQuery.Table, insertQuery.Fields)

	tests := []struct {
		name        string
		input       string
		wantError   bool
		errContains string
	}{
		{
			name:      "Valid update",
			input:     fmt.Sprintf("UPDATE users %d Alice,alice@mail.ru", insertedID),
			wantError: false,
		},
		{
			name:        "Non-existent table",
			input:       "UPDATE unknown 1 New,Value",
			wantError:   true,
			errContains: "таблица unknown не найдена",
		},
		{
			name:        "Non-existent record",
			input:       "UPDATE users 999 New,Value",
			wantError:   true,
			errContains: "запись не найдена",
		},
		{
			name:        "Invalid field count",
			input:       fmt.Sprintf("UPDATE users %d OnlyName", insertedID),
			wantError:   true,
			errContains: "несоответствие количества полей",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)
			if parseErr != nil && !tt.wantError {
				t.Fatalf("Parse error: %v", parseErr)
			}
			app.handleUpdate(query)

			if !tt.wantError {

				record, err := app.DB.Select(query.Table, query.ID)
				if err != nil {
					t.Errorf("Failed to select updated record: %v", err)
				}
				for i, field := range app.DB.Tables[query.Table].Fields {
					if record[field] != query.Fields[i] {
						t.Errorf("Field %s not updated, got %s, want %s",
							field, record[field], query.Fields[i])
					}
				}
			}
		})
	}
}

func TestHandleInsert(t *testing.T) {
	app, tempDir := setupTestApp(t)
	defer cleanupTestApp(tempDir)

	createQuery := &parser.Query{
		Type:   parser.QueryCreateTable,
		Table:  "users",
		Fields: []string{"name", "email"},
	}
	app.HandleCreateTable(createQuery)

	tests := []struct {
		name        string
		input       string
		wantError   bool
		errContains string
	}{
		{
			name:      "Valid insert",
			input:     "INSERT users kolya,kolya@mail.ru",
			wantError: false,
		},
		{
			name:        "Non-existent table",
			input:       "INSERT unknown Value1,Value2",
			wantError:   true,
			errContains: "таблица unknown не найдена",
		},
		{
			name:        "Invalid field count",
			input:       "INSERT users OnlyName",
			wantError:   true,
			errContains: "несоответствие количества полей",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)
			if parseErr != nil && !tt.wantError {
				t.Fatalf("Parse error: %v", parseErr)
			}

			initialCount := 0
			if app.Storage.TableExist(query.Table) {
				records, _ := app.DB.SelectAll(query.Table)
				initialCount = len(records)
			}

			app.handleInsert(query)

			if !tt.wantError {
				records, err := app.DB.SelectAll(query.Table)
				if err != nil {
					t.Errorf("Failed to get records: %v", err)
				}
				if len(records) != initialCount+1 {
					t.Errorf("Expected %d records, got %d", initialCount+1, len(records))
				}
			} else {
				if app.Storage.TableExist(query.Table) {
					records, _ := app.DB.SelectAll(query.Table)
					if len(records) != initialCount {
						t.Errorf("Record count should not change on error, got %d, want %d",
							len(records), initialCount)
					}
				}
			}
		})
	}
}

func TestHandleDelete(t *testing.T) {
	app, tempDir := setupTestApp(t)
	defer cleanupTestApp(tempDir)

	createQuery := &parser.Query{
		Type:   parser.QueryCreateTable,
		Table:  "users",
		Fields: []string{"name", "email"},
	}
	app.HandleCreateTable(createQuery)

	insertQuery := &parser.Query{
		Type:   parser.QueryInsert,
		Table:  "users",
		Fields: []string{"kolya", "kolya@mail.ru"},
	}
	insertedID, _ := app.DB.Insert(insertQuery.Table, insertQuery.Fields)

	tests := []struct {
		name        string
		input       string
		wantError   bool
		errContains string
	}{
		{
			name:      "Valid delete",
			input:     fmt.Sprintf("DELETE users %d", insertedID),
			wantError: false,
		},
		{
			name:        "Non-existent table",
			input:       "DELETE unknown 1",
			wantError:   true,
			errContains: "таблица unknown не найдена",
		},
		{
			name:        "Non-existent record",
			input:       "DELETE users 999",
			wantError:   true,
			errContains: "запись не найдена",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, parseErr := parser.ParseQuery(tt.input)
			if parseErr != nil && !tt.wantError {
				t.Fatalf("Parse error: %v", parseErr)
			}
			initialRecords, _ := app.DB.SelectAll("users")
			app.handleDelete(query)

			if !tt.wantError {
				_, err := app.DB.Select(query.Table, query.ID)
				if err == nil {
					t.Errorf("Record %d should be deleted", query.ID)
				}

				currentRecords, _ := app.DB.SelectAll("users")
				if len(initialRecords)-1 != len(currentRecords) {
					t.Errorf("Expected %d records after delete, got %d",
						len(initialRecords)-1, len(currentRecords))
				}
			} else if tt.errContains != "Error сохранения таблицы" {
				currentRecords, _ := app.DB.SelectAll("users")
				if len(initialRecords) != len(currentRecords) {
					t.Error("Record count should not change on error")
				}
			}
		})
	}
}

func TestHandleHelp(t *testing.T) {
	app, _ := setupTestApp(t)

	notPanics := func() {
		app.handleHelp()
	}

	assert.NotPanics(t, notPanics)
}
