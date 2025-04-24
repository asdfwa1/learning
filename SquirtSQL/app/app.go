package app

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"v4/database"
	"v4/database/actions"
	"v4/database/parser"
	"v4/storage"
)

type App struct {
	DB      *actions.Database
	Storage *storage.CSVStorage
}

func NewApp() *App {
	stor := storage.NewCSVStorage("data")
	db := actions.NewDatabase(stor)

	fmt.Println("SQUIRTSQL - простая база данных на основе CSV")
	fmt.Println("Введите /help для просмотра функционала")
	fmt.Println("Введите exit чтобы выйти")
	fmt.Println("---------------------------------------------")

	if err := db.LoadTables(); err != nil {
		return nil
	}
	return &App{
		DB:      db,
		Storage: stor,
	}
}

func (a *App) Run() {
	fmt.Println("---------------------------------------------")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("squirtsql>> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}

		a.handleQuery(input)
	}
}

func (a *App) handleQuery(input string) {
	query, err := parser.ParseQuery(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch query.Type {
	case parser.QueryCreateTable:
		a.HandleCreateTable(query)
	case parser.QuerySelect:
		a.handleSelect(query)
	case parser.QueryUpdate:
		a.handleUpdate(query)
	case parser.QueryInsert:
		a.handleInsert(query)
	case parser.QueryDelete:
		a.handleDelete(query)
	case parser.QueryHelp:
		a.handleHelp()
	default:
		fmt.Println("Error: неизвестный тип запроса")
	}
}

func (a *App) HandleCreateTable(query *parser.Query) error {
	if a.Storage.TableExist(query.Table) {
		return fmt.Errorf("Error: таблица %s уже существует", query.Table)
	}

	err := a.DB.CreateTable(query.Table, query.Fields)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	table, _ := a.DB.SelectAll(query.Table)
	err = a.Storage.SaveTable(&database.Table{
		NextID:  1,
		Name:    query.Table,
		Fields:  query.Fields,
		Records: table,
	})

	if err != nil {
		return fmt.Errorf("Error: таблица не сохранена: %v", err)
	}

	fmt.Println("Таблица успешно создана")
	return nil
}

func (a *App) handleSelect(query *parser.Query) {
	if !a.Storage.TableExist(query.Table) {
		fmt.Printf("Error: таблица %s не найдена\n", query.Table)
		return
	}

	if _, exist := a.DB.SelectAll(query.Table); exist != nil {
		table, err := a.Storage.LoadTable(query.Table)
		if err != nil {
			fmt.Printf("Error загрузки таблицы: %v\n", err)
			return
		}

		_ = a.DB.CreateTable(table.Name, table.Fields)
		for _, record := range table.Records {
			values := make([]string, len(table.Fields))
			for i, field := range table.Fields {
				values[i] = record[field]
			}
			_, _ = a.DB.Insert(table.Name, values)
		}
	}
	if query.ID == -1 {
		records, err := a.DB.SelectAll(query.Table)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(records) == 0 {
			fmt.Println("Записи не найдены")
			return
		}

		ids := make([]int, 0, len(records))
		for id := range records {
			ids = append(ids, id)
		}
		sort.Ints(ids)

		for _, id := range ids {
			record := records[id]
			fmt.Printf("%d: ", id)

			fields := make([]string, 0, len(record))
			for fieldName := range record {
				fields = append(fields, fieldName)
			}
			sort.Strings(fields)

			for i, fieldName := range fields {
				fmt.Printf("%s:%s", fieldName, record[fieldName])
				if i < len(fields)-1 {
					fmt.Printf(" ")
				}
			}
			fmt.Println()
		}
	} else {
		record, err := a.DB.Select(query.Table, query.ID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(record) == 0 {
			fmt.Println("Записи не найдены")
			return
		}
		fields := make([]string, 0, len(record))
		for fieldName := range record {
			fields = append(fields, fieldName)
		}
		sort.Strings(fields)
		fmt.Printf("%d: ", query.ID)
		for i, fieldName := range fields {
			fmt.Printf("%s: %s", fieldName, record[fieldName])
			if i < len(fields)-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func (a *App) handleUpdate(query *parser.Query) {
	if !a.Storage.TableExist(query.Table) {
		fmt.Printf("Error: таблица %s не найдена\n", query.Table)
		return
	}
	err := a.DB.Update(query.Table, query.ID, query.Fields)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	table := a.DB.Tables[query.Table]
	err = a.Storage.SaveTable(table)

	if err != nil {
		fmt.Printf("Error сохранения таблицы: %v\n", err)
	} else {
		fmt.Println("Таблица успешно сохранена")
	}
}

func (a *App) handleInsert(query *parser.Query) {
	if !a.Storage.TableExist(query.Table) {
		fmt.Printf("Error: таблица %s не найдена\n", query.Table)
		return
	}

	_, err := a.DB.Insert(query.Table, query.Fields)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	table := a.DB.Tables[query.Table]
	err = a.Storage.SaveTable(table)
	if err != nil {
		fmt.Printf("Error сохранения таблицы: %v\n", err)
	} else {
		fmt.Println("Данные успешно вставлены в таблицу")
	}
}

func (a *App) handleDelete(query *parser.Query) {
	if !a.Storage.TableExist(query.Table) {
		fmt.Printf("Error: таблица %s не найдена\n", query.Table)
		return
	}
	err := a.DB.Delete(query.Table, query.ID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	table, _ := a.DB.SelectAll(query.Table)
	err = a.Storage.SaveTable(&database.Table{
		Name:    query.Table,
		Fields:  a.DB.Tables[query.Table].Fields,
		Records: table,
	})
	if err != nil {
		fmt.Printf("Error сохранения таблицы: %v\n", err)
	} else {
		fmt.Println("Таблица успешно сохранена")
	}
}

func (a *App) handleHelp() {
	helpText := `
Доступные команды:

1. Создание таблицы:
   CREATE TABLE <имя_таблицы> <поле1>,<поле2>,...
   Пример: CREATE TABLE users name,email,age

2. Добавление данных:
   INSERT <имя_таблицы> <значение1>,<значение2>,...
   Пример: INSERT users kolya,test@mail.ru,22

3. Чтение данных:
   SELECT <имя_таблицы> <id|*>
   Примеры:
     SELECT users *       - все записи
     SELECT users 1       - запись с ID=1

4. Обновление данных:
   UPDATE <имя_таблицы> <id> <новое_значение1>,<новое_значение2>,...
   Пример: UPDATE users 1 NewName,new@email.com,23

5. Удаление данных:
   DELETE <имя_таблицы> <id>
   Пример: DELETE users 1

6. Справка:
   /help - вывести это сообщение

7. Выход:
   exit - завершить программу
`
	fmt.Println(helpText)
}
