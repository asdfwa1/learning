package action

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"v0/logic"
)

const (
	Create = "create"
	Read   = "read"
	Update = "update"
	Delete = "delete"
)

type CLI struct {
	Reader      *bufio.Reader
	TaskManager *logic.TaskManager
}

func (c *CLI) Run() {
	fmt.Println("Welcome TO DO List")

	for {
		fmt.Println("\nВведите команду: (create, read, update, delete)")
		command, err := c.Reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения строки:", err)
			continue
		}
		command = strings.TrimSpace(command)

		switch command {
		case Create:
			c.createTask()
		case Read:
			c.listTasks()
		case Update:
			c.updateTask()
		case Delete:
			c.deleteTask()
		default:
			fmt.Println("Неверная команда. Попробуйте снова.")
		}
	}
}

func (c *CLI) createTask() {
	fmt.Println("Введите название новой задачи:")
	name, err := c.Reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения строки:", err)
		return
	}
	name = strings.TrimSpace(name)

	task, err := c.TaskManager.AddTask(name)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Printf("Задача добавлена: ID=%d, Название=%s\n", task.ID, task.Name)
}

func (c *CLI) listTasks() {
	tasks := c.TaskManager.ListTasks()
	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return
	}

	fmt.Println("Ваш список задач:")
	for _, task := range tasks {
		fmt.Printf("ID=%d, Название=%s\n", task.ID, task.Name)
	}
}

func (c *CLI) updateTask() {
	fmt.Println("Введите ID задачи, которая будет обновлена:")
	idStr, err := c.Reader.ReadString('\n')
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

	fmt.Println("Введите новое название задачи:")
	newName, err := c.Reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения строки:", err)
		return
	}
	newName = strings.TrimSpace(newName)

	err = c.TaskManager.UpdateTask(id, newName)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Printf("Задача обновлена: ID=%d, Новое название=%s\n", id, newName)
}

func (c *CLI) deleteTask() {
	fmt.Println("Введите ID задачи для удаления:")
	idStr, err := c.Reader.ReadString('\n')
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

	err = c.TaskManager.DeleteTask(id)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Printf("Задача с ID=%d удалена.\n", id)
}
