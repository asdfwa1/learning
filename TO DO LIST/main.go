package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Tasks struct {
	ID   int
	Name string
}

var (
	tasks  []Tasks
	nextID = 1
)

const (
	Create = "create"
	Read   = "read"
	Update = "update"
	Delete = "delete"
)

func main() {
	fmt.Println("Welcome TO DO List")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nВведиье команду: (create, read, update, delete)")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения строки:", err)
			continue
		}
		command = strings.TrimSpace(command)

		switch command {
		case Create:
			createTask(reader)
		case Read:
			readTasks(reader)
		case Update:
			updateTask(reader)
		case Delete:
			deleteTask(reader)
		default:
			fmt.Println("Неверная команда. Попробуйте снова.")
		}
	}
}

func createTask(reader *bufio.Reader) {
	fmt.Println("Введите название новой задачи: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения строки:", err)
		return
	}
	name = strings.TrimSpace(name)

	if name == "" {
		fmt.Println("Название задачи не может быть пустым.")
		return
	}

	task := Tasks{
		ID:   nextID,
		Name: name,
	}
	tasks = append(tasks, task)
	nextID++
	fmt.Printf("Задача добавлена: ID=%d, Название=%s\n", task.ID, task.Name)
}

func readTasks(reader *bufio.Reader) {
	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return
	}

	fmt.Println("Ваш список задач:")
	for _, task := range tasks {
		fmt.Printf("ID=%d, Название=%s\n", task.ID, task.Name)
	}
}

func updateTask(reader *bufio.Reader) {

}

func deleteTask(reader *bufio.Reader) {

}
