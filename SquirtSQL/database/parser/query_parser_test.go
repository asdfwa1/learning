package parser

import (
	"strings"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *Query
		expectError bool
		errText     string
	}{
		// CREATE TABLE tests
		{
			name:  "valid CREATE TABLE",
			input: "CREATE TABLE users id, name, age",
			expected: &Query{
				Type:   QueryCreateTable,
				Table:  "users",
				Fields: []string{"id", "name", "age"},
			},
		},
		{
			name:        "CREATE TABLE missing name",
			input:       "CREATE TABLE",
			expectError: true,
			errText:     "формат: CREATE TABLE <table> <values>",
		},
		{
			name:        "CREATE TABLE missing fields",
			input:       "CREATE TABLE users",
			expectError: true,
			errText:     "не указаны поля таблицы",
		},

		// SELECT tests
		{
			name:  "valid SELECT with ID",
			input: "SELECT users 1",
			expected: &Query{
				Type:  QuerySelect,
				Table: "users",
				ID:    1,
			},
		},
		{
			name:  "valid SELECT all",
			input: "SELECT users *",
			expected: &Query{
				Type:  QuerySelect,
				Table: "users",
				ID:    -1,
			},
		},
		{
			name:        "SELECT missing table",
			input:       "SELECT",
			expectError: true,
			errText:     "формат: SELECT <table> <id> or <*>",
		},
		{
			name:        "SELECT invalid ID",
			input:       "SELECT users abc",
			expectError: true,
			errText:     "",
		},

		// INSERT tests
		{
			name:  "valid INSERT",
			input: "INSERT users John, 30, developer",
			expected: &Query{
				Type:   QueryInsert,
				Table:  "users",
				Fields: []string{"John", "30", "developer"},
			},
		},
		{
			name:        "INSERT missing values",
			input:       "INSERT users",
			expectError: true,
			errText:     "формат: INSERT <table> <values>",
		},
		{
			name:        "INSERT empty values",
			input:       "INSERT users   ",
			expectError: true,
			errText:     "формат: INSERT <table> <values>",
		},

		// UPDATE tests
		{
			name:  "valid UPDATE",
			input: "UPDATE users 1 Kolya, 23",
			expected: &Query{
				Type:   QueryUpdate,
				Table:  "users",
				ID:     1,
				Fields: []string{"Kolya", "23"},
			},
		},
		{
			name:        "UPDATE missing ID",
			input:       "UPDATE users name=John",
			expectError: true,
			errText:     "формат: UPDATE <table> <id> <values>",
		},
		{
			name:        "UPDATE invalid ID",
			input:       "UPDATE users abc name=John",
			expectError: true,
			errText:     "неподходящий ID в update",
		},

		// DELETE tests
		{
			name:  "valid DELETE",
			input: "DELETE users 1",
			expected: &Query{
				Type:  QueryDelete,
				Table: "users",
				ID:    1,
			},
		},
		{
			name:        "DELETE missing ID",
			input:       "DELETE users",
			expectError: true,
			errText:     "формат: DELETE <table> <id>",
		},
		{
			name:        "DELETE invalid ID",
			input:       "DELETE users abc",
			expectError: true,
			errText:     "неподходящий ID в delete",
		},

		// HELP test
		{
			name:  "valid HELP",
			input: "/HELP",
			expected: &Query{
				Type: QueryHelp,
			},
		},

		// Invalid queries
		{
			name:        "empty query",
			input:       "",
			expectError: true,
			errText:     "неизвестный тип запроса",
		},
		{
			name:        "unknown command",
			input:       "RANDOM COMMAND",
			expectError: true,
			errText:     "неизвестный тип запроса",
		},

		// Normalization tests
		{
			name:  "query with extra spaces",
			input: "  CREATE   TABLE   users   id,  name  ",
			expected: &Query{
				Type:   QueryCreateTable,
				Table:  "users",
				Fields: []string{"id", "name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseQuery(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errText, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if actual.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", actual.Type, tt.expected.Type)
			}

			if actual.Table != tt.expected.Table {
				t.Errorf("Table = %v, want %v", actual.Table, tt.expected.Table)
			}

			if actual.ID != tt.expected.ID {
				t.Errorf("ID = %v, want %v", actual.ID, tt.expected.ID)
			}

			if len(actual.Fields) != len(tt.expected.Fields) {
				t.Errorf("Fields length = %v, want %v", len(actual.Fields), len(tt.expected.Fields))
			} else {
				for i := range actual.Fields {
					if actual.Fields[i] != tt.expected.Fields[i] {
						t.Errorf("Fields[%d] = %v, want %v", i, actual.Fields[i], tt.expected.Fields[i])
					}
				}
			}
		})
	}
}
