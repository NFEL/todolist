package services

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/repository"
)

type TaskService struct {
	TaskRepo repository.TaskRepo
}

func NewTaskService(taskrepo repository.TaskRepo) *TaskService {
	return &TaskService{
		TaskRepo: taskrepo,
	}
}

func (s *TaskService) TaskList(ctx context.Context, userID uint) ([]dto.TaskListResp, error) {
	// s.TaskRepo.
	return nil, nil
}
