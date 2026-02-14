package storage_postgres

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository/storage"
	"math"
	"reflect"

	"gorm.io/gorm"
)

type taskImp struct {
	db *gorm.DB
}

func NewTaskRepo(db *storage.DB) *taskImp {
	return &taskImp{
		db: db.DB,
	}
}

func (i *taskImp) Create(ctx context.Context, task *domain.Task) (uint, error) {
	err := gorm.G[domain.Task](i.db).Create(ctx, task)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (i *taskImp) GetByID(ctx context.Context, ID uint) (domain.Task, error) {
	return gorm.G[domain.Task](i.db).Where("id = ?", ID).Take(ctx)
}

func (i *taskImp) List(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	r := make([]domain.Task, int(math.Abs(float64(limit-offset))))
	err := gorm.G[domain.Task](i.db).Select("*").Limit(limit).Offset(offset).Scan(ctx, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (i *taskImp) ListByFilter(ctx context.Context, filter dto.TaskListFilter, limit, offset int) ([]domain.Task, int64, error) {
	q := i.db.Model(&domain.Task{})

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.Assignee != 0 {
		q = q.Where("id IN (SELECT task_id FROM user_tasks WHERE user_id = ?)", filter.Assignee)
	}
	if !reflect.ValueOf(filter.CreatedAt).IsZero() {
		q = q.Where("created_at >= ?", filter.CreatedAt)
	}
	if !reflect.ValueOf(filter.UpdatedAt).IsZero() {
		q = q.Where("updated_at >= ?", filter.UpdatedAt)
	}

	var total int64
	q.Count(&total)

	var tasks []domain.Task
	if err := q.Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (i *taskImp) UpdateByID(ctx context.Context, task *domain.Task, fields []string) error {
	_, err := gorm.G[domain.Task](i.db).Where("id = ?", task.ID).Select(fields[0], fields[1:]).Updates(ctx, *task)
	return err
}

func (i *taskImp) DeleteByID(ctx context.Context, ID uint) error {
	_, err := gorm.G[domain.Task](i.db).Where("id = ?", ID).Delete(ctx)
	return err
}
