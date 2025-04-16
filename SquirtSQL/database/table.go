package database

import (
	"errors"
)

var (
	ErrTableNotFound  = errors.New("таблица не найдена")
	ErrRecordNotFound = errors.New("запись не найдена")
	ErrMissFieldCount = errors.New("несоответствие количества полей")
)

func NewTable(name string, field []string) *Table {
	return &Table{
		Name:    name,
		Fields:  field,
		Records: make(map[int]Record),
		NextID:  1,
	}
}

func (t *Table) ValidateFields(field []string) bool {
	return len(field) == len(t.Fields)
}
