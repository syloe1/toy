package dao

import (
	"cron/internal/model"

	"gorm.io/gorm"
)

type TaskDAO interface {
	Create(task *model.Task) error
	GetByID(id uint) (*model.Task, error)
	Update(task *model.Task) error
	DeleteByID(id uint) error
	DeleteByIDs(ids []uint) error
	ListByStatus(status string) ([]model.Task, error)
}

type taskDAO struct {
	db *gorm.DB
}

func NewTaskDAO(db *gorm.DB) TaskDAO {
	return &taskDAO{db: db}
}

// Create(task *model.Task) error
func (d *taskDAO) Create(task *model.Task) error {
	return d.db.Create(task).Error
}

// GetByID(id uint) (*model.Task, error)
func (d *taskDAO) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	//func (db *DB) First(dest interface{}, conds ...interface{}) (tx *DB)
	if err := d.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// Update(task *model.Task) error
func (d *taskDAO) Update(task *model.Task) error {
	return d.db.Save(task).Error
}

// DeleteByID(id uint) error
func (d *taskDAO) DeleteByID(id uint) error {
	result := d.db.Delete(&model.Task{}, id)
	// 判断是否发生数据库错误
	if result.Error != nil {
		return result.Error
	}
	//判断是否真的删除了数据
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteByIDs(ids []uint) error
func (d *taskDAO) DeleteByIDs(ids []uint) error {
	result := d.db.Where("id IN ?", ids).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (d *taskDAO) ListByStatus(status string) ([]model.Task, error) {
	var tasks []model.Task
	query := d.db.Order("created_at DESC")

	if status != model.TaskStatusAll {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
