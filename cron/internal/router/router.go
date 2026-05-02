package router

import (
	"cron/internal/handler"
	"cron/pkg/response"

	"github.com/gin-gonic/gin"
)

func New(taskHandler *handler.TaskHandler) *gin.Engine {
	engine := gin.Default()

	engine.GET("/healthz", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	api := engine.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.ListTasks)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.DELETE("", taskHandler.BulkDeleteTasks)
		}
	}

	return engine
}
