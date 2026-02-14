package services

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository/enum"
	mockRepo "graph-interview/internal/repository/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTask_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	taskRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
		Run(func(args mock.Arguments) {
			task := args.Get(1).(*domain.Task)
			task.ID = 1
		}).
		Return(uint(1), nil)

	req := dto.CreateTaskReq{
		Name:        "Test Task",
		Description: "A test task",
	}

	resp, err := svc.CreateTask(context.Background(), req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Test Task", resp.Name)
	assert.Equal(t, "Created", resp.Status)
	taskRepo.AssertExpectations(t)
}

func TestGetTask_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.Task{
			Name:        "Test Task",
			Description: "A test task",
			Status:      enum.Created,
		}, nil)

	resp, err := svc.GetTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test Task", resp.Name)
	taskRepo.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	resp, err := svc.GetTask(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, api_error.ErrTaskNotFound, err)
	taskRepo.AssertExpectations(t)
}

func TestListTasks_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	tasks := []domain.Task{
		{Name: "Task 1", Status: enum.Created},
		{Name: "Task 2", Status: enum.Started},
	}

	filter := dto.TaskListFilter{}
	taskRepo.On("ListByFilter", mock.Anything, filter, 20, 0).
		Return(tasks, int64(2), nil)

	resp, err := svc.ListTasks(context.Background(), filter, 20, 0)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Tasks))
	assert.Equal(t, int64(2), resp.Total)
	taskRepo.AssertExpectations(t)
}

func TestUpdateTask_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	existingTask := domain.Task{
		Name:        "Old Name",
		Description: "Old Desc",
		Status:      enum.Created,
	}
	existingTask.ID = 1

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(existingTask, nil)
	taskRepo.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Task"), mock.Anything).
		Return(nil)

	newName := "New Name"
	req := dto.UpdateTaskReq{Name: &newName}

	resp, err := svc.UpdateTask(context.Background(), 1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Name", resp.Name)
	taskRepo.AssertExpectations(t)
}

func TestUpdateTask_StatusChange(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	existingTask := domain.Task{
		Name:   "Task",
		Status: enum.Created,
	}
	existingTask.ID = 1

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(existingTask, nil)
	taskRepo.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Task"), mock.Anything).
		Return(nil)

	newStatus := enum.Started
	req := dto.UpdateTaskReq{Status: &newStatus}

	resp, err := svc.UpdateTask(context.Background(), 1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Started", resp.Status)
	taskRepo.AssertExpectations(t)
}

func TestDeleteTask_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.Task{}, nil)
	taskRepo.On("DeleteByID", mock.Anything, uint(1)).
		Return(nil)

	err := svc.DeleteTask(context.Background(), 1)

	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
}

func TestDeleteTask_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	err := svc.DeleteTask(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, api_error.ErrTaskNotFound, err)
	taskRepo.AssertExpectations(t)
}

func TestArchiveTask_Success(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	svc := NewTaskService(taskRepo)

	existingTask := domain.Task{
		Name:   "Task",
		Status: enum.Created,
	}
	existingTask.ID = 1

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(existingTask, nil)
	taskRepo.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Task"), []string{"status", "updated_by_user_id"}).
		Return(nil)

	resp, err := svc.ArchiveTask(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Canceled", resp.Status)
	taskRepo.AssertExpectations(t)
}
