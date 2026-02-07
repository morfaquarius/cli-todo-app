package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"todo-app/internal/storage"
	"todo-app/internal/todo"
)

const taskFile = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: не указана команда")
		os.Exit(1)
	}

	command := os.Args[1]

	flags := os.Args[2:]

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	desc := addCmd.String("desc", "", "Описание задачи.")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	filter := listCmd.String("filter", "all", "Фильтрация задач done(только выполненные), pending(только невыполненные), all(все).")

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	idComplete := completeCmd.Int("id", 0, "Отметить задачу выполненной по её ID.")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	idDelete := deleteCmd.Int("id", 0, "Удалить задачу по её ID.")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	format := exportCmd.String("format", "csv", "Формат файла для экспорта текущего списка задач (пример - json/csv)")
	outFile := exportCmd.String("out", "backup.csv", "Имя файла для экспорта текущего списка (пример - backup.csv/backup.json)")

	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
	loadFile := loadCmd.String("file", "", "Имя файла для импорта задач (пример - tasks.csv/tasks.json)")

	switch command {
	case "add":
		if err := addCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			addCmd.Usage()
			os.Exit(1)
		}
		if *desc == "" {
			fmt.Println("Ошибка: описание не может быть пустым")
			addCmd.Usage()
			os.Exit(1)
		}
		tasks := loadTasksOrExit(taskFile)

		tasks = todo.Add(tasks, *desc)

		saveTasksOrExit(taskFile, tasks)

		fmt.Println("Задача успешно добавлена!")
	case "list":
		if err := listCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			listCmd.Usage()
			os.Exit(1)
		}
		validFilters := map[string]bool{
			"all":     true,
			"done":    true,
			"pending": true,
		}

		if !validFilters[*filter] {
			fmt.Println("Неизвестный фильтр", *filter)
			listCmd.Usage()
			os.Exit(1)
		}

		tasks := loadTasksOrExit(taskFile)

		if len(tasks) == 0 {
			fmt.Println("Нет задач")
			return
		}

		filteredTasks := todo.List(tasks, *filter)

		if len(filteredTasks) == 0 {
			fmt.Println("Нет задач с фильтром", *filter)
			return
		}

		fmt.Printf("%-4s %-10s %s\n", "ID", "Статус", "Описание")
		fmt.Printf("%-4s %-10s %s\n", "--", "------", "-----------")

		for _, task := range filteredTasks {
			status := "[ ]"
			if task.Done {
				status = "[✓]"
			}
			fmt.Printf("%-4d %-10s %s\n", task.ID, status, task.Description)
		}
	case "complete":
		if err := completeCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			completeCmd.Usage()
			os.Exit(1)
		}

		if *idComplete <= 0 {
			fmt.Println("Ошибка: ID должен быть положительным числом (1, 2, 3, ...)")
			completeCmd.Usage()
			os.Exit(1)
		}

		tasks := loadTasksOrExit(taskFile)

		tasks, err := todo.Complete(tasks, *idComplete)
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
		saveTasksOrExit(taskFile, tasks)

		fmt.Printf("Задача с ID - %d помечена выполненной!\n", *idComplete)
	case "delete":
		if err := deleteCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			deleteCmd.Usage()
			os.Exit(1)
		}

		if *idDelete <= 0 {
			fmt.Println("Ошибка: ID должен быть положительным числом")
			deleteCmd.Usage()
			os.Exit(1)
		}
		tasks := loadTasksOrExit(taskFile)

		tasks, err := todo.Delete(tasks, *idDelete)
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
		saveTasksOrExit(taskFile, tasks)

		fmt.Printf("Задача с ID - %d успешно удалена!\n", *idDelete)
	case "export":
		if err := exportCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			exportCmd.Usage()
			os.Exit(1)
		}

		validFormats := map[string]bool{
			"csv":  true,
			"json": true,
		}

		if !validFormats[*format] {
			fmt.Println("Неизвестный формат", *format)
			exportCmd.Usage()
			os.Exit(1)
		}

		fileExt := filepath.Ext(*outFile)

		expectedExt := ""

		if *format == "json" {
			expectedExt = ".json"
		} else {
			expectedExt = ".csv"
		}

		if fileExt != expectedExt {
			fmt.Printf("Ошибка: расширение файла '%s' не соответствует формату '%s'\n", *outFile, *format)
			os.Exit(1)
		}

		tasks := loadTasksOrExit(taskFile)

		if *format == "json" {
			err := storage.SaveJSON(*outFile, tasks)
			if err != nil {
				fmt.Println("Ошибка:", err)
				os.Exit(1)
			}
		}
		if *format == "csv" {
			err := storage.SaveCSV(*outFile, tasks)
			if err != nil {
				fmt.Println("Ошибка:", err)
				os.Exit(1)
			}
		}
		fmt.Println("Задачи успешно экспортированы в", *outFile)
	case "load":
		if err := loadCmd.Parse(flags); err != nil {
			fmt.Println("Ошибка парсинга аргументов:", err)
			loadCmd.Usage()
			os.Exit(1)
		}

		if *loadFile == "" {
			fmt.Println("Ошибка: файл для импорта задач не указан")
			loadCmd.Usage()
			os.Exit(1)
		}

		fileExt := filepath.Ext(*loadFile)

		switch fileExt {
		case ".json":
			tasks, err := storage.LoadJSON(*loadFile)
			if err != nil {
				fmt.Println("Ошибка:", err)
				os.Exit(1)
			}
			saveTasksOrExit(taskFile, tasks)

		case ".csv":
			tasks, err := storage.LoadCSV(*loadFile)
			if err != nil {
				fmt.Println("Ошибка:", err)
				os.Exit(1)
			}
			saveTasksOrExit(taskFile, tasks)
		default:
			fmt.Println("Ошибка: неподдерживаемое расширение файла", *loadFile)
			os.Exit(1)
		}
	default:
		fmt.Println("Ошибка: неизвестная команда -", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Доступные команды:
  add
  list
  complete
  delete
  export
  load`)
}

func loadTasksOrExit(path string) []todo.Task {
	tasks, err := storage.LoadJSON(path)
	if err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
	return tasks
}

func saveTasksOrExit(path string, tasks []todo.Task) {
	if err := storage.SaveJSON(path, tasks); err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
}
