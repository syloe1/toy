package main

import (
	"fmt"
	"os"
	"strconv"

	"cli_todo/internal/repository"
	"cli_todo/internal/service"
)

func main() {
	repo := &repository.Repository{FilePath: "tasks.json"}
	svc := &service.TaskService{Repo: repo}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide a task description")
			return
		}
		if err := svc.AddTask(os.Args[2]); err != nil {
			fmt.Printf("failed to add task: %v\n", err)
			return
		}
		fmt.Println("task added successfully")
	case "list":
		showAll := len(os.Args) >= 3 && (os.Args[2] == "-a" || os.Args[2] == "--all")
		if err := svc.List(showAll); err != nil {
			fmt.Printf("failed to list tasks: %v\n", err)
		}
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("error: task ID must be an integer")
			return
		}
		if err := svc.CompleteTask(id); err != nil {
			fmt.Printf("failed to complete task: %v\n", err)
			return
		}
		fmt.Println("task completed successfully")
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("error: task ID must be an integer")
			return
		}
		if err := svc.DeleteTask(id); err != nil {
			fmt.Printf("failed to delete task: %v\n", err)
			return
		}
		fmt.Println("task deleted successfully")
	case "help", "h", "--help":
		printUsage()
	default:
		fmt.Printf("unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  tasks add <content>       add a task")
	fmt.Println("  tasks list                list incomplete tasks")
	fmt.Println("  tasks list -a             list all tasks")
	fmt.Println("  tasks complete <id>       complete a task")
	fmt.Println("  tasks delete <id>         delete a task")
}
