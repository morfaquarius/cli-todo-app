package todo

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
