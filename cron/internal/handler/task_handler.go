package handler

import (
	"errors"
	"net/http"
	"strconv"

	"cron/internal/service"
	"cron/pkg/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

type createTaskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateTaskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type bulkDeleteRequest struct {
	IDs []uint `json:"ids"`
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.taskService.CreateTask(req.Title, req.Content)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	// 从请求参数中获取筛选状态，默认值为 all 表示查询所有任务
	status := c.DefaultQuery("status", "all")

	tasks, err := h.taskService.ListTasks(status)
	if err != nil {
		// 统一错误处理
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(id, req.Title, req.Content)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, task)
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.taskService.UpdateTaskStatus(id, req.Status)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) BulkDeleteTasks(c *gin.Context) {
	var req bulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.taskService.BulkDeleteTasks(req.IDs); err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, nil)
}

func parseIDParam(c *gin.Context) (uint, bool) {
	idValue := c.Param("id")
	id64, err := strconv.ParseUint(idValue, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return 0, false
	}

	return uint(id64), true
}

func (h *TaskHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrTitleRequired),
		errors.Is(err, service.ErrEmptyIDList):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
