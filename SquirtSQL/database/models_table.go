package database

import (
	"sync"
)

type Record map[string]string

type Table struct {
	Name    string
	Fields  []string
	Records map[int]Record
	Mu      sync.RWMutex
	NextID  int
}
