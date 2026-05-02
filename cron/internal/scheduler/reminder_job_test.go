package scheduler

import (
	"testing"

	"cron/internal/model"
)

type fakeTaskService struct {
	tasks []model.Task
	err   error
}

func (f *fakeTaskService) ListTasks(status string) ([]model.Task, error) {
	return f.tasks, f.err
}

type fakeReminderMailer struct {
	sendCalled bool
	taskCount  int
}

func (f *fakeReminderMailer) SendPendingTasksReminder(tasks []model.Task) error {
	f.sendCalled = true
	f.taskCount = len(tasks)
	return nil
}

func TestReminderJobSkipsEmptyTaskList(t *testing.T) {
	taskService := &fakeTaskService{tasks: []model.Task{}}
	mailer := &fakeReminderMailer{}
	job := NewReminderJob(taskService, mailer)

	job.Run()

	if mailer.sendCalled {
		t.Fatalf("expected reminder email not to be sent")
	}
}

func TestReminderJobSendsPendingTasks(t *testing.T) {
	taskService := &fakeTaskService{
		tasks: []model.Task{
			{ID: 1, Title: "task 1", Status: model.TaskStatusPending},
			{ID: 2, Title: "task 2", Status: model.TaskStatusPending},
		},
	}
	mailer := &fakeReminderMailer{}
	job := NewReminderJob(taskService, mailer)

	job.Run()

	if !mailer.sendCalled {
		t.Fatalf("expected reminder email to be sent")
	}

	if mailer.taskCount != 2 {
		t.Fatalf("expected 2 tasks in email, got %d", mailer.taskCount)
	}
}
