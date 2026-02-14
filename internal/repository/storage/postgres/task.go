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

func (i *taskImp) List(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	r := make([]domain.Task, int(math.Abs(float64(limit-offset))))
	err := gorm.G[domain.Task](i.db).Select("*").Limit(limit).Offset(offset).Scan(ctx, &r)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (i *taskImp) ListByFilter(ctx context.Context, filter dto.TaskListFilter, limit, offset int) ([]domain.Task, error) {
	q := gorm.G[domain.Task](i.db).Select("*").Limit(limit).Offset(offset)
	if !reflect.ValueOf(filter.CreatedAt).IsZero() {
		q = q.Where("created_at >= ?", filter.CreatedAt)
	}
	return q.Find(ctx)
}

func (i *taskImp) UpdateByID(ctx context.Context, user *domain.Task, fields []string) error {
	_, err := gorm.G[domain.Task](i.db).Where("id = ?", user.ID).Select(fields[0], fields[1:]).Updates(ctx, *user)
	return err
}
