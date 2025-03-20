package logic

type Task struct {
	ID   int
	Name string
}

type TaskManager struct {
	Tasks  []Task
	NextID int
}
