package parser

import (
	"errors"
	"strconv"
	"strings"
)

type QueryType int

const (
	QueryCreateTable QueryType = iota
	QuerySelect
	QueryInsert
	QueryUpdate
	QueryDelete
	QueryHelp
)

const (
	CREATE = "CREATE TABLE"
	SELECT = "SELECT"
	INSERT = "INSERT"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	HELP   = "/HELP"
)

type Query struct {
	Type   QueryType
	Table  string
	Fields []string
	ID     int
}

func ParseQuery(input string) (*Query, error) {
	normalized := strings.Join(strings.Fields(input), " ")
	parts := strings.SplitN(normalized, " ", 3)
	if len(parts) < 1 {
		return nil, errors.New("неверный формат запроса")
	}

	query := &Query{}
	if len(parts) > 1 {
		twoWordCommand := strings.ToUpper(parts[0] + " " + parts[1])
		switch twoWordCommand {
		case CREATE:
			query.Type = QueryCreateTable
			if len(parts) < 3 {
				return nil, errors.New("формат: CREATE TABLE <table> <values>")
			}
			nameAndFields := strings.SplitN(parts[2], " ", 2)
			if len(nameAndFields) < 1 {
				return nil, errors.New("не указано имя таблицы")
			}
			query.Table = nameAndFields[0]

			if len(nameAndFields) > 1 {
				fields := strings.Split(nameAndFields[1], ",")
				for i := range fields {
					fields[i] = strings.TrimSpace(fields[i])
				}
				query.Fields = fields
			}
			return query, nil
		}
	}
	command := strings.ToUpper(parts[0])
	switch command {
	case SELECT:
		if len(parts) < 3 {
			return nil, errors.New("формат: SELECT <table> <id> or <*>")
		}
		query.Type = QuerySelect
		query.Table = parts[1]
		if len(parts) > 2 {
			if parts[2] == "*" {
				query.ID = -1
			} else {
				id, err := strconv.Atoi(parts[2])
				if err != nil {
					return nil, errors.New("неподходящий ID в select")
				}
				query.ID = id
			}
		}
	case INSERT:
		if len(parts) < 3 {
			return nil, errors.New("формат: INSERT <table> <values>")
		}
		query.Type = QueryInsert
		query.Table = parts[1]

		valuesPart := strings.Join(parts[2:], " ")
		query.Fields = strings.Split(valuesPart, ",")

		for i := range query.Fields {
			query.Fields[i] = strings.TrimSpace(query.Fields[i])
		}

		if len(query.Fields) == 0 {
			return nil, errors.New("для вставки необходимо указать значения")
		}
	case UPDATE:
		newParts := strings.SplitN(normalized, " ", 4)
		query.Type = QueryUpdate

		if len(newParts) < 4 {
			return nil, errors.New("формат: UPDATE <table> <id> <values>")
		}
		query.Table = newParts[1]
		id, err := strconv.Atoi(newParts[2])
		if err != nil {
			return nil, errors.New("неподходящий ID в update")
		}
		query.ID = id

		valuesPart := strings.Join(newParts[3:], " ")
		query.Fields = strings.Split(valuesPart, ",")
		for i := range query.Fields {
			query.Fields[i] = strings.TrimSpace(query.Fields[i])
		}
	case DELETE:
		if len(parts) < 3 {
			return nil, errors.New("формат: DELETE <table> <id>")
		}
		query.Type = QueryDelete
		query.Table = parts[1]

		id, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.New("неподходящий ID в delete")
		}
		query.ID = id
	case HELP:
		query.Type = QueryHelp
		return query, nil
	default:
		return nil, errors.New("неизвестный тип запроса")
	}
	return query, nil
}
