package service

import (
	"errors"
	"strings"

	"cron/internal/dao"
	"cron/internal/model"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(title, content string) (*model.Task, error)
	UpdateTask(id uint, title, content string) (*model.Task, error)
	UpdateTaskStatus(id uint, status string) (*model.Task, error)
	DeleteTask(id uint) error
	BulkDeleteTasks(ids []uint) error
	ListTasks(status string) ([]model.Task, error)
}

type taskService struct {
	taskDAO dao.TaskDAO
}

func NewTaskService(taskDAO dao.TaskDAO) TaskService {
	return &taskService{taskDAO: taskDAO}
}

func (s *taskService) CreateTask(title, content string) (*model.Task, error) {
	// 1. 校验：标题不能为空
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	// 2. 构建任务对象（默认 pending 状态）
	task := &model.Task{
		Title:   title,
		Content: strings.TrimSpace(content),
		Status:  model.TaskStatusPending,
	}

	// 3. 调用 DAO 存入数据库
	if err := s.taskDAO.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) UpdateTask(id uint, title, content string) (*model.Task, error) {
	// 校验标题
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	// 查询任务是否存在
	task, err := s.taskDAO.GetByID(id)
	if err != nil {
		return nil, mapDAOError(err)
	}
	// 修改内容
	task.Title = title
	task.Content = strings.TrimSpace(content)
	// 保存更新
	if err := s.taskDAO.Update(task); err != nil {
		return nil, mapDAOError(err)
	}

	return task, nil
}

func (s *taskService) UpdateTaskStatus(id uint, status string) (*model.Task, error) {
	status = strings.TrimSpace(status)
	if !model.IsValidTaskStatus(status) {
		return nil, ErrInvalidStatus
	}
	// 查询任务是否存在
	task, err := s.taskDAO.GetByID(id)
	if err != nil {
		return nil, mapDAOError(err)
	}

	task.Status = status
	// 保存
	if err := s.taskDAO.Update(task); err != nil {
		return nil, mapDAOError(err)
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uint) error {
	return mapDAOError(s.taskDAO.DeleteByID(id))
}

func (s *taskService) BulkDeleteTasks(ids []uint) error {
	// 业务校验：ID 列表不能为空
	if len(ids) == 0 {
		return ErrEmptyIDList
	}
	// 调用 DAO + 转换错误

	return mapDAOError(s.taskDAO.DeleteByIDs(ids))
}

func (s *taskService) ListTasks(status string) ([]model.Task, error) {
	status = normalizeFilter(status) // 空值 → all
	if !model.IsValidTaskFilter(status) {
		return nil, ErrInvalidStatus
	}

	return s.taskDAO.ListByStatus(status)
}

func normalizeFilter(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return model.TaskStatusAll
	}

	return status
}

func mapDAOError(err error) error {
	// 如果 DAO 返回「记录不存在」
	// 转换成业务错误：ErrTaskNotFound
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTaskNotFound
	}

	return err
}
