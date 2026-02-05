package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"todo-app/internal/todo"
)

func LoadJSON(path string) ([]todo.Task, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		SaveJSON(path, []todo.Task{})
		return []todo.Task{}, nil
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("ошибка: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	var tasks []todo.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования из JSON: %w", err)
	}
	return tasks, nil
}

func SaveJSON(path string, tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка преобразования в JSON: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файла: %w", err)
	}
	return nil
}
