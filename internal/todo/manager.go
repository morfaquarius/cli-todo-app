package todo

import "fmt"

// Add возвращает обновлённый срез с новой задачей
func Add(tasks []Task, desc string) []Task {
	var maxID int
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	newTask := Task{
		ID:          maxID + 1,
		Description: desc,
		Done:        false,
	}
	tasks = append(tasks, newTask)
	return tasks
}

// List возвращает отфильтрованный срез по заданному параметру
func List(tasks []Task, filter string) []Task {
	switch filter {
	case "all":
		return tasks
	case "done":
		doneTasks := []Task{}
		for _, task := range tasks {
			if task.Done {
				doneTasks = append(doneTasks, task)
			}
		}
		return doneTasks
	case "pending":
		pendingTasks := []Task{}
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}
		return pendingTasks
	default:
		return []Task{}
	}
}

// Complete отмечает задачу выполненной по её ID
func Complete(tasks []Task, id int) ([]Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			return tasks, nil
		}
	}
	return tasks, fmt.Errorf("задача с ID %d не найдена", id)
}

// Delete удаляет задачу по её ID
func Delete(tasks []Task, id int) ([]Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks, nil
		}
	}
	return tasks, fmt.Errorf("задача с ID %d не найдена", id)
}
