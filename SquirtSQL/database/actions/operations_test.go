package actions

import (
	"os"
	"testing"
	"v4/storage"
)

func setupTestDB(t *testing.T) (*Database, string) {
	tempDir, err := os.MkdirTemp("", "db_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storage := storage.NewCSVStorage(tempDir)
	db := NewDatabase(storage)

	return db, tempDir
}

func cleanupTestDB(tempDir string) {
	_ = os.RemoveAll(tempDir)
}

func TestCreateTable(t *testing.T) {
	db, tempDir := setupTestDB(t)
	defer cleanupTestDB(tempDir)

	tests := []struct {
		name       string
		tableName  string
		fields     []string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "Valid table creation",
			tableName: "users",
			fields:    []string{"name", "email"},
			wantErr:   false,
		},
		{
			name:       "Duplicate table",
			tableName:  "users",
			fields:     []string{"name", "email"},
			wantErr:    true,
			errMessage: "таблица users уже существует",
		},
		{
			name:       "Reserved field name",
			tableName:  "products",
			fields:     []string{"id", "price"},
			wantErr:    true,
			errMessage: "поле 'id' зарезервированно системой",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.CreateTable(tt.tableName, tt.fields)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if err.Error() != tt.errMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if _, exists := db.Tables[tt.tableName]; !exists {
					t.Errorf("Table '%s' was not created", tt.tableName)
				}
			}
		})
	}
}

func TestInsertAndSelect(t *testing.T) {
	db, tempDir := setupTestDB(t)
	defer cleanupTestDB(tempDir)

	tableName := "users"
	fields := []string{"name", "email"}
	err := db.CreateTable(tableName, fields)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	tests := []struct {
		name     string
		values   []string
		wantErr  bool
		wantID   int
		wantData map[string]string
	}{
		{
			name:     "Valid insert",
			values:   []string{"John Doe", "john@example.com"},
			wantErr:  false,
			wantID:   1,
			wantData: map[string]string{"name": "John Doe", "email": "john@example.com"},
		},
		{
			name:    "Invalid field count",
			values:  []string{"Jane Doe"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := db.Insert(tableName, tt.values)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if id != tt.wantID {
				t.Errorf("Expected ID %d, got %d", tt.wantID, id)
			}

			record, err := db.Select(tableName, id)
			if err != nil {
				t.Errorf("Failed to select record: %v", err)
			}

			for field, wantValue := range tt.wantData {
				if record[field] != wantValue {
					t.Errorf("Field %s: expected '%s', got '%s'", field, wantValue, record[field])
				}
			}
		})
	}
}

func TestSelectAll(t *testing.T) {
	db, tempDir := setupTestDB(t)
	defer cleanupTestDB(tempDir)

	tableName := "products"
	fields := []string{"name", "price"}
	err := db.CreateTable(tableName, fields)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	testData := []struct {
		values []string
	}{
		{[]string{"Laptop", "999.99"}},
		{[]string{"Phone", "599.99"}},
		{[]string{"Tablet", "399.99"}},
	}

	for _, data := range testData {
		_, err := db.Insert(tableName, data.values)
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	records, err := db.SelectAll(tableName)
	if err != nil {
		t.Fatalf("Failed to select all records: %v", err)
	}

	if len(records) != len(testData) {
		t.Errorf("Expected %d records, got %d", len(testData), len(records))
	}

	lastID := 0
	for id, record := range records {
		if id <= lastID {
			t.Error("Records are not sorted by ID")
		}
		lastID = id

		if len(record) != len(fields) {
			t.Errorf("Record %d has wrong field count", id)
		}
	}
}

func TestUpdate(t *testing.T) {
	db, tempDir := setupTestDB(t)
	defer cleanupTestDB(tempDir)

	tableName := "employees"
	fields := []string{"name", "position"}
	err := db.CreateTable(tableName, fields)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	id, err := db.Insert(tableName, []string{"Alice", "Developer"})
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	tests := []struct {
		name     string
		id       int
		values   []string
		wantErr  bool
		wantData map[string]string
	}{
		{
			name:     "Valid update",
			id:       id,
			values:   []string{"Alice Smith", "Senior Developer"},
			wantErr:  false,
			wantData: map[string]string{"name": "Alice Smith", "position": "Senior Developer"},
		},
		{
			name:    "Invalid field count",
			id:      id,
			values:  []string{"Only name"},
			wantErr: true,
		},
		{
			name:    "Non-existent record",
			id:      999,
			values:  []string{"Bob", "Manager"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Update(tableName, tt.id, tt.values)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			record, err := db.Select(tableName, tt.id)
			if err != nil {
				t.Errorf("Failed to select record: %v", err)
			}

			for field, wantValue := range tt.wantData {
				if record[field] != wantValue {
					t.Errorf("Field %s: expected '%s', got '%s'", field, wantValue, record[field])
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, tempDir := setupTestDB(t)
	defer cleanupTestDB(tempDir)

	tableName := "customers"
	fields := []string{"name", "email"}
	err := db.CreateTable(tableName, fields)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	id, err := db.Insert(tableName, []string{"Bob", "bob@example.com"})
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Valid delete",
			id:      id,
			wantErr: false,
		},
		{
			name:    "Non-existent record",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Delete(tableName, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			_, err = db.Select(tableName, tt.id)
			if err == nil {
				t.Error("Record still exists after deletion")
			}
		})
	}
}
