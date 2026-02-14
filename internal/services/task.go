package services

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository"
	"graph-interview/internal/repository/enum"
)

type TaskService struct {
	TaskRepo repository.TaskRepo
}

func NewTaskService(taskrepo repository.TaskRepo) *TaskService {
	return &TaskService{
		TaskRepo: taskrepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskReq, userID uint) (*dto.TaskResp, error) {
	task := &domain.Task{
		Name:            req.Name,
		Description:     req.Description,
		Status:          enum.Created,
		CreatedByUserID: userID,
		UpdatedByUserID: userID,
	}

	id, err := s.TaskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return &dto.TaskResp{
		ID:          id,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status.String(),
		CreatedByID: task.CreatedByUserID,
		UpdatedByID: task.UpdatedByUserID,
	}, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID uint) (*dto.TaskResp, error) {
	task, err := s.TaskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, api_error.ErrTaskNotFound
	}
	return taskToResp(&task), nil
}

func (s *TaskService) ListTasks(ctx context.Context, filter dto.TaskListFilter, limit, offset int) (*dto.TaskListResp, error) {
	tasks, total, err := s.TaskRepo.ListByFilter(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	taskResps := make([]dto.TaskResp, len(tasks))
	for i, t := range tasks {
		taskResps[i] = *taskToResp(&t)
	}

	return &dto.TaskListResp{
		Tasks:  taskResps,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID uint, req dto.UpdateTaskReq, userID uint) (*dto.TaskResp, error) {
	task, err := s.TaskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, api_error.ErrTaskNotFound
	}

	var fields []string
	task.UpdatedByUserID = userID
	fields = append(fields, "updated_by_user_id")

	if req.Name != nil {
		task.Name = *req.Name
		fields = append(fields, "name")
	}
	if req.Description != nil {
		task.Description = *req.Description
		fields = append(fields, "description")
	}
	if req.Status != nil {
		task.Status = *req.Status
		fields = append(fields, "status")
	}

	if err := s.TaskRepo.UpdateByID(ctx, &task, fields); err != nil {
		return nil, err
	}

	return taskToResp(&task), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uint) error {
	_, err := s.TaskRepo.GetByID(ctx, taskID)
	if err != nil {
		return api_error.ErrTaskNotFound
	}
	return s.TaskRepo.DeleteByID(ctx, taskID)
}

func (s *TaskService) ArchiveTask(ctx context.Context, taskID uint, userID uint) (*dto.TaskResp, error) {
	task, err := s.TaskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, api_error.ErrTaskNotFound
	}

	task.Status = enum.Canceled
	task.UpdatedByUserID = userID
	fields := []string{"status", "updated_by_user_id"}

	if err := s.TaskRepo.UpdateByID(ctx, &task, fields); err != nil {
		return nil, err
	}

	return taskToResp(&task), nil
}

func taskToResp(task *domain.Task) *dto.TaskResp {
	return &dto.TaskResp{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status.String(),
		CreatedByID: task.CreatedByUserID,
		UpdatedByID: task.UpdatedByUserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
