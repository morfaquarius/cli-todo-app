package main

import (
	"fmt"
	"todo-app/internal/todo"
	"todo-app/internal/todo/storage"
)

func main() {
	// Тест 1: Сохраняем задачи
	tasks := []todo.Task{
		{ID: 1, Description: "Купить молоко", Done: false},
		{ID: 2, Description: "Изучить Go", Done: true},
	}

	err := storage.SaveJSON("test.json", tasks)
	if err != nil {
		fmt.Println("Ошибка сохранения:", err)
		return
	}
	fmt.Println("Файл сохранён!")

	// Тест 2: Загружаем обратно
	loadedTasks, err := storage.LoadJSON("test.json")
	if err != nil {
		fmt.Println("Ошибка загрузки:", err)
		return
	}

	fmt.Println("Загруженные задачи:")
	for _, task := range loadedTasks {
		fmt.Printf("ID: %d, Описание: %s, Выполнена: %v\n",
			task.ID, task.Description, task.Done)
	}
}
