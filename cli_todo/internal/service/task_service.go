package service

import (
	"errors"
	"fmt"
	"time"

	"cli_todo/internal/model"
	"cli_todo/internal/repository"
)

type TaskService struct {
	Repo *repository.Repository
}

func (s *TaskService) AddTask(content string) error {
	tasks, err := s.Repo.LoadAll()
	if err != nil {
		return err
	}

	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}

	newTask := model.Task{
		ID:        newID,
		Content:   content,
		Done:      false,
		CreatedAt: time.Now(),
	}

	tasks = append(tasks, newTask)
	return s.Repo.SaveAll(tasks)
}

func (s *TaskService) List(showAll bool) error {
	tasks, err := s.Repo.LoadAll()
	if err != nil {
		return err
	}

	fmt.Println("ID\tTask\tCreated At\tDone")

	count := 0
	for _, task := range tasks {
		if !showAll && task.Done {
			continue
		}
		fmt.Printf("%d\t%s\t%s\t%v\n", task.ID, task.Content, task.CreatedAt.Format(time.RFC3339), task.Done)
		count++
	}

	fmt.Printf("Total: %d task(s)\n", count)
	return nil
}

func (s *TaskService) CompleteTask(id int) error {
	tasks, err := s.Repo.LoadAll()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			return s.Repo.SaveAll(tasks)
		}
	}

	return errors.New("task not found")
}

func (s *TaskService) DeleteTask(id int) error {
	tasks, err := s.Repo.LoadAll()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.Repo.SaveAll(tasks)
		}
	}

	return errors.New("task not found")
}
