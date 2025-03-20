package app

import (
	"bufio"
	"os"
	"v0/TO_DO_LIST/action"
	"v0/TO_DO_LIST/logic"
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
