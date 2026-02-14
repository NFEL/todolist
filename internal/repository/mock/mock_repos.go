package mock

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock of UserRepo interface
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) (uint, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, ID uint) (domain.User, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByField(ctx context.Context, field string, value any) (domain.User, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepo) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) ListByFilter(ctx context.Context, filter dto.UserListFilter, limit, offset int) ([]domain.User, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) UpdateByID(ctx context.Context, user *domain.User, fields []string) error {
	args := m.Called(ctx, user, fields)
	return args.Error(0)
}

// MockTaskRepo is a mock of TaskRepo interface
type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) Create(ctx context.Context, task *domain.Task) (uint, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockTaskRepo) GetByID(ctx context.Context, ID uint) (domain.Task, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepo) List(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepo) ListByFilter(ctx context.Context, filter dto.TaskListFilter, limit, offset int) ([]domain.Task, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]domain.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepo) UpdateByID(ctx context.Context, task *domain.Task, fields []string) error {
	args := m.Called(ctx, task, fields)
	return args.Error(0)
}

func (m *MockTaskRepo) DeleteByID(ctx context.Context, ID uint) error {
	args := m.Called(ctx, ID)
	return args.Error(0)
}
