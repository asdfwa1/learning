package logic

import "errors"

func (tm *TaskManager) AddTask(name string) (Task, error) {
	if name == "" {
		return Task{}, errors.New("название задачи не может быть пустым")
	}

	task := Task{
		ID:   tm.NextID,
		Name: name,
	}
	tm.Tasks = append(tm.Tasks, task)
	tm.NextID++

	return task, nil
}

func (tm *TaskManager) ListTasks() []Task {
	return tm.Tasks
}

func (tm *TaskManager) UpdateTask(id int, newName string) error {
	if len(newName) < 3 {
		return errors.New("название задачи должно содержать минимум 3 символа")
	}

	for i := range tm.Tasks {
		if tm.Tasks[i].ID == id {
			if tm.Tasks[i].Name == newName {
				return errors.New("новое название не должно совпадать с предыдущим")
			}
			tm.Tasks[i].Name = newName
			return nil
		}
	}

	return errors.New("задача с таким ID не найдена")
}

func (tm *TaskManager) DeleteTask(id int) error {
	var foundIndex = -1
	for i, task := range tm.Tasks {
		if task.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return errors.New("задача с таким ID не найдена")
	}

	tm.Tasks = append(tm.Tasks[:foundIndex], tm.Tasks[foundIndex+1:]...)

	for i := foundIndex; i < len(tm.Tasks); i++ {
		tm.Tasks[i].ID = i + 1
	}

	tm.NextID = len(tm.Tasks) + 1

	return nil
}
