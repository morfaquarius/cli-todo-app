package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"todo-app/internal/todo"
)

// LoadCSV загружает список задач из файла формата CSV
func LoadCSV(path string) ([]todo.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия CSV файла: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Ошибка закрытия файла: %v\n", err)
		}
	}()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV файла: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("файл CSV пуст")
	}

	if len(records) == 1 {
		return []todo.Task{}, nil
	}

	var tasks []todo.Task

	for _, record := range records[1:] {
		if len(record) != 3 {
			return nil, fmt.Errorf("неверный формат строки CSV: %v", record)
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования строки: %w", err)
		}
		done, err := strconv.ParseBool(record[2])
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования строки: %w", err)
		}
		task := todo.Task{
			ID:          id,
			Description: record[1],
			Done:        done,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// SaveCSV сохраняет список задач в файл формата CSV
func SaveCSV(path string, tasks []todo.Task) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Ошибка закрытия файла: %v\n", err)
		}
	}()
	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"ID", "Description", "Done"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка CSV: %w", err)
	}
	for _, task := range tasks {
		if err := writer.Write([]string{
			strconv.Itoa(task.ID),
			task.Description,
			strconv.FormatBool(task.Done),
		}); err != nil {
			return fmt.Errorf("ошибка записи строки CSV: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("ошибка записи в CSV: %w", err)
	}
	return nil
}
