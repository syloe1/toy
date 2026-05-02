package model

import "time"

const (
	TaskStatusPending = "pending"
	TaskStatusDone    = "done"
	TaskStatusAll     = "all"
)

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    string    `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 校验任务状态是否有效（创建/更新时用）
func IsValidTaskStatus(status string) bool {
	return status == TaskStatusPending || status == TaskStatusDone
}

// 校验列表筛选状态是否有效（查询时用）
func IsValidTaskFilter(status string) bool {
	return status == TaskStatusAll || IsValidTaskStatus(status)
}
