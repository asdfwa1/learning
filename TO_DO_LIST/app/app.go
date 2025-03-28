package app

import (
	"bufio"
	"os"
	"v0/action"
	"v0/logic"
)

func NewTaskManager() *logic.TaskManager {
	return &logic.TaskManager{
		Tasks:  []logic.Task{},
		NextID: 1,
	}
}

func NewCLI() *action.CLI {
	return &action.CLI{
		Reader:      bufio.NewReader(os.Stdin),
		TaskManager: NewTaskManager(),
	}
}
