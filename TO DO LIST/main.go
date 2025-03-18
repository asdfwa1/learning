package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	fmt.Println("Введите ID задачи, которая будет обновлена:")
	idStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ID:", err)
		return
	}
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка преобразования ID:", err)
		return
	}

	var foundTask *Tasks
	var lastName string
	for i := range tasks {
		if tasks[i].ID == id {
			foundTask = &tasks[i]
			lastName = tasks[i].Name
			break
		}
	}

	if foundTask == nil {
		fmt.Println("Задача с таким ID не найдена.")
		return
	}

	fmt.Println("Введите новое название задачи:")
	newName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения строки:", err)
		return
	}
	newName = strings.TrimSpace(newName)

	if len(newName) < 2 || newName == lastName {
		fmt.Println("Название задачи должно содержать минимум 3 символа\nи новое название не должно совпадать с предыдущим ")
		return
	}
	foundTask.Name = newName
	fmt.Printf("Задача обновлена: ID=%d, Новое название=%s\n", foundTask.ID, foundTask.Name)
}

func deleteTask(reader *bufio.Reader) {
	fmt.Println("Введите ID задачи для удаления:")
	idStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ID:", err)
		return
	}
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка преобразования ID:", err)
		return
	}

	var foundIndex = -1
	for i, task := range tasks {
		if task.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		fmt.Println("Задача с таким ID не найдена.")
		return
	}

	tasks = append(tasks[:foundIndex], tasks[foundIndex+1:]...)

	for i := foundIndex; i < len(tasks); i++ {
		tasks[i].ID = i + 1
	}

	nextID = len(tasks) + 1
	fmt.Printf("Задача с ID=%d удалена.\n", id)
}
