package storage

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"v4/database"
)

type CSVStorage struct {
	BasePath string
	Mu       sync.Mutex
}

func NewCSVStorage(basePath string) *CSVStorage {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		_ = os.MkdirAll(basePath, 0755)
	}
	return &CSVStorage{
		BasePath: basePath,
	}
}

func (s *CSVStorage) SaveTable(table *database.Table) (err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	filePath := s.BasePath + "/" + table.Name + ".csv"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := append([]string{"id"}, table.Fields...)
	if err := writer.Write(header); err != nil {
		return err
	}

	for id, record := range table.Records {
		row := make([]string, len(header))
		row[0] = strconv.Itoa(id)
		for i, field := range table.Fields {
			row[i+1] = record[field]
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func (s *CSVStorage) LoadTable(name string) (*database.Table, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	filePath := s.BasePath + "/" + name + ".csv"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, errors.New("файл с таблицей пуст")
	}

	userFields := records[0][1:]
	table := database.NewTable(name, userFields)
	maxID := 0

	for _, row := range records[1:] {
		if len(row) != len(records[0]) {
			continue
		}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		record := make(database.Record)
		for i, field := range userFields {
			record[field] = row[i+1]
		}

		table.Records[id] = record
		if id > maxID {
			maxID = id
		}
	}
	table.NextID = maxID + 1
	return table, nil
}

func (s *CSVStorage) TableExist(name string) bool {
	filePath := s.BasePath + "/" + name + ".csv"
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (s *CSVStorage) ListTables() ([]string, error) {
	files, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, err
	}

	var tables []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			tables = append(tables, strings.TrimSuffix(file.Name(), ".csv"))
		}
	}
	return tables, nil
}
