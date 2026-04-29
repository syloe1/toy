package repository

import (
	"encoding/json"
	"os"

	"cli_todo/internal/model"
)

type Repository struct {
	FilePath string
}

func (r *Repository) SaveAll(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.FilePath, data, 0644)
}

func (r *Repository) LoadAll() ([]model.Task, error) {
	data, err := os.ReadFile(r.FilePath)
	if err != nil {
		//file not exist
		if os.IsNotExist(err) {
			return []model.Task{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []model.Task{}, nil
	}
	var tasks []model.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
