package storage

import (
	"os"
	"path/filepath"
	"testing"
	"v4/database"
)

func TestNewCSVStorage(t *testing.T) {
	tempDir := t.TempDir()

	testPath := filepath.Join(tempDir, "newdir")
	storage := NewCSVStorage(testPath)

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("NewCSVStorage() didn't create directory")
	}

	if storage.BasePath != testPath {
		t.Errorf("NewCSVStorage() BasePath = %v, want %v", storage.BasePath, testPath)
	}
}

func TestCSVStorage_SaveAndLoadTable(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewCSVStorage(tempDir)

	fields := []string{"name", "age"}
	table := database.NewTable("users", fields)

	table.Records[1] = database.Record{"name": "Alice", "age": "30"}
	table.Records[2] = database.Record{"name": "Bob", "age": "25"}
	table.NextID = 3

	err := storage.SaveTable(table)
	if err != nil {
		t.Fatalf("SaveTable() error = %v", err)
	}

	filePath := filepath.Join(tempDir, "users.csv")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("SaveTable() didn't create file")
	}

	loadedTable, err := storage.LoadTable("users")
	if err != nil {
		t.Fatalf("LoadTable() error = %v", err)
	}

	if loadedTable.Name != "users" {
		t.Errorf("LoadTable() table name = %v, want %v", loadedTable.Name, "users")
	}

	if len(loadedTable.Fields) != len(fields) {
		t.Errorf("LoadTable() fields count = %v, want %v", len(loadedTable.Fields), len(fields))
	}

	if len(loadedTable.Records) != 2 {
		t.Errorf("LoadTable() records count = %v, want %v", len(loadedTable.Records), 2)
	}

	if loadedTable.NextID != 3 {
		t.Errorf("LoadTable() NextID = %v, want %v", loadedTable.NextID, 3)
	}
}

func TestCSVStorage_TableExist(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewCSVStorage(tempDir)

	if storage.TableExist("nonexistent") {
		t.Errorf("TableExist() returned true for nonexistent table")
	}

	table := database.NewTable("test", []string{"field"})
	err := storage.SaveTable(table)
	if err != nil {
		t.Fatal(err)
	}

	if !storage.TableExist("test") {
		t.Errorf("TableExist() returned false for existing table")
	}
}

func TestCSVStorage_ListTables(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewCSVStorage(tempDir)

	tables := []string{"users", "products", "orders"}
	for _, name := range tables {
		table := database.NewTable(name, []string{"field"})
		if err := storage.SaveTable(table); err != nil {
			t.Fatal(err)
		}
	}

	list, err := storage.ListTables()
	if err != nil {
		t.Fatalf("ListTables() error = %v", err)
	}

	if len(list) != len(tables) {
		t.Errorf("ListTables() returned %d tables, want %d", len(list), len(tables))
	}

	for _, name := range tables {
		found := false
		for _, tableName := range list {
			if tableName == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListTables() missing table %s", name)
		}
	}
}
