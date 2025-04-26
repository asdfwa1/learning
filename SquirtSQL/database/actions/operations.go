package actions

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"sort"
	"sync"
	"v4/database"
	"v4/storage"
)

type Database struct {
	Tables  map[string]*database.Table
	Mu      sync.RWMutex
	Storage *storage.CSVStorage
}

func NewDatabase(storage *storage.CSVStorage) *Database {
	db := &Database{
		Tables:  make(map[string]*database.Table),
		Storage: storage,
	}
	return db
}

func (db *Database) CreateTable(name string, userFields []string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	if _, exist := db.Tables[name]; exist == true {
		return fmt.Errorf("таблица %s уже существует", name)
	}

	for _, field := range userFields {
		if field == "id" {
			return errors.New("поле 'id' зарезервированно системой")
		}
	}

	db.Tables[name] = &database.Table{
		Name:    name,
		Fields:  userFields,
		Records: make(map[int]database.Record),
		NextID:  1,
	}
	return nil
}

func (db *Database) Insert(tableName string, values []string) (int, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	table, ok := db.Tables[tableName]
	if !ok {
		return 0, database.ErrTableNotFound
	}
	if !table.ValidateFields(values) {
		return 0, database.ErrMissFieldCount
	}

	table.Mu.Lock()
	defer table.Mu.Unlock()

	id := table.NextID
	record := make(database.Record)

	for i, field := range table.Fields {
		record[field] = values[i]
	}
	table.Records[id] = record
	table.NextID++
	return id, nil
}

func (db *Database) Select(tableName string, id int) (database.Record, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	table, ok := db.Tables[tableName]
	if !ok {
		return nil, database.ErrTableNotFound
	}

	table.Mu.RLock()
	defer table.Mu.RUnlock()

	if id == -1 {
		return nil, nil
	}

	record, exist := table.Records[id]
	if !exist {
		return nil, database.ErrRecordNotFound
	}

	return record, nil
}

func (db *Database) SelectAll(tableName string) (map[int]database.Record, error) {
	db.Mu.RLock()
	defer db.Mu.RUnlock()

	table, exist := db.Tables[tableName]
	if !exist {
		return nil, database.ErrTableNotFound
	}

	table.Mu.RLock()
	defer table.Mu.RUnlock()

	sortedCopy := make(map[int]database.Record)
	ids := make([]int, 0, len(table.Records))
	for id := range table.Records {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		sortedCopy[id] = table.Records[id]
	}

	return sortedCopy, nil
}

func (db *Database) Update(tableName string, id int, values []string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	table, exist := db.Tables[tableName]
	if !exist {
		return database.ErrTableNotFound
	}

	if !table.ValidateFields(values) {
		return database.ErrMissFieldCount
	}

	table.Mu.Lock()
	defer table.Mu.Unlock()

	record, exists := table.Records[id]
	if !exists {
		return database.ErrRecordNotFound
	}

	for i := 0; i < len(table.Fields); i++ {
		fieldName := table.Fields[i]
		record[fieldName] = values[i]
	}

	return nil
}

func (db *Database) Delete(tableName string, id int) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	table, exist := db.Tables[tableName]
	if !exist {
		return database.ErrTableNotFound
	}

	table.Mu.Lock()
	defer table.Mu.Unlock()

	if _, exists := table.Records[id]; !exists {
		return database.ErrRecordNotFound
	}

	delete(table.Records, id)
	return nil
}

func (db *Database) LoadTables() error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	tableNames, err := db.Storage.ListTables()
	if err != nil {
		return err
	}

	for _, name := range tableNames {
		table, err := db.Storage.LoadTable(name)
		if err != nil {
			fmt.Printf("Ошибка загрузки таблицы %s : %v\n", name, err)
			continue
		}
		db.Tables[name] = table
		TableColor := color.New(color.FgBlue).SprintFunc()
		valid := fmt.Sprintf("Таблица %s загружена", name)
		fmt.Println(TableColor(valid))
	}

	return nil
}
