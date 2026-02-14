package repository

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) (uint, error)
	GetByID(ctx context.Context, ID uint) (domain.User, error)
	GetByField(ctx context.Context, field string, value any) (domain.User, error)
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
	ListByFilter(ctx context.Context, filter dto.UserListFilter, limit, offset int) ([]domain.User, error)
	UpdateByID(ctx context.Context, user *domain.User, fields []string) error
}

type TaskRepo interface {
	Create(ctx context.Context, user *domain.Task) (uint, error)
	List(ctx context.Context, limit, offset int) ([]domain.Task, error)
	ListByFilter(ctx context.Context, filter dto.TaskListFilter, limit, offset int) ([]domain.Task, error)
	UpdateByID(ctx context.Context, user *domain.Task, fields []string) error
}
