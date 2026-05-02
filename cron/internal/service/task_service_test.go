package service

import (
	"testing"
	"time"

	"cron/internal/model"
)

type fakeTaskDAO struct {
	tasks  map[uint]*model.Task
	nextID uint
}

func newFakeTaskDAO() *fakeTaskDAO {
	return &fakeTaskDAO{
		tasks:  make(map[uint]*model.Task),
		nextID: 1,
	}
}

func (f *fakeTaskDAO) Create(task *model.Task) error {
	task.ID = f.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	f.nextID++

	taskCopy := *task
	f.tasks[task.ID] = &taskCopy
	return nil
}

func (f *fakeTaskDAO) GetByID(id uint) (*model.Task, error) {
	task, ok := f.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	taskCopy := *task
	return &taskCopy, nil
}

func (f *fakeTaskDAO) Update(task *model.Task) error {
	if _, ok := f.tasks[task.ID]; !ok {
		return ErrTaskNotFound
	}

	task.UpdatedAt = time.Now()
	taskCopy := *task
	f.tasks[task.ID] = &taskCopy
	return nil
}

func (f *fakeTaskDAO) DeleteByID(id uint) error {
	if _, ok := f.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(f.tasks, id)
	return nil
}

func (f *fakeTaskDAO) DeleteByIDs(ids []uint) error {
	deleted := 0
	for _, id := range ids {
		if _, ok := f.tasks[id]; ok {
			delete(f.tasks, id)
			deleted++
		}
	}

	if deleted == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (f *fakeTaskDAO) ListByStatus(status string) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	for _, task := range f.tasks {
		if status == model.TaskStatusAll || task.Status == status {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

func TestCreateTaskDefaultsToPending(t *testing.T) {
	taskService := NewTaskService(newFakeTaskDAO())

	task, err := taskService.CreateTask("write docs", "finish the API guide")
	if err != nil {
		t.Fatalf("CreateTask returned error: %v", err)
	}

	if task.Status != model.TaskStatusPending {
		t.Fatalf("expected status %q, got %q", model.TaskStatusPending, task.Status)
	}
}

func TestListTasksRejectsInvalidStatus(t *testing.T) {
	taskService := NewTaskService(newFakeTaskDAO())

	if _, err := taskService.ListTasks("unknown"); err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestBulkDeleteTasksRejectsEmptyIDs(t *testing.T) {
	taskService := NewTaskService(newFakeTaskDAO())

	if err := taskService.BulkDeleteTasks(nil); err != ErrEmptyIDList {
		t.Fatalf("expected ErrEmptyIDList, got %v", err)
	}
}
