package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository/enum"
	mockRepo "graph-interview/internal/repository/mock"
	"graph-interview/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupTaskRouter(taskSrv *services.TaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tasks := r.Group("/tasks")
	tasks.Use(func(c *gin.Context) {
		c.Set("userID", "1")
		c.Next()
	})
	tasks.POST("", CreateTask(taskSrv))
	tasks.GET("", ListTasks(taskSrv))
	tasks.GET("/:id", GetTask(taskSrv))
	tasks.PUT("/:id", UpdateTask(taskSrv))
	tasks.DELETE("/:id", DeleteTask(taskSrv))
	tasks.PATCH("/:id/archive", ArchiveTask(taskSrv))
	return r
}

func setupTaskRouterNoAuth(taskSrv *services.TaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tasks := r.Group("/tasks")
	// No userID set - simulates unauthenticated
	tasks.POST("", CreateTask(taskSrv))
	tasks.PUT("/:id", UpdateTask(taskSrv))
	tasks.PATCH("/:id/archive", ArchiveTask(taskSrv))
	return r
}

func TestCreateTaskHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
		Run(func(args mock.Arguments) {
			task := args.Get(1).(*domain.Task)
			task.ID = 1
		}).
		Return(uint(1), nil)

	body, _ := json.Marshal(dto.CreateTaskReq{
		Name:        "Test Task",
		Description: "Description",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	taskRepo.AssertExpectations(t)
}

func TestCreateTaskHandler_InvalidBody(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	body, _ := json.Marshal(map[string]string{"invalid": "body"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTaskHandler_Unauthorized(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouterNoAuth(taskSrv)

	body, _ := json.Marshal(dto.CreateTaskReq{
		Name:        "Test Task",
		Description: "Description",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateTaskHandler_RepoError(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
		Return(uint(0), errors.New("db error"))

	body, _ := json.Marshal(dto.CreateTaskReq{
		Name:        "Test Task",
		Description: "Description",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestGetTaskHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.Task{
			Name:        "Test Task",
			Description: "Description",
			Status:      enum.Created,
		}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	taskRepo.AssertExpectations(t)
}

func TestGetTaskHandler_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestGetTaskHandler_InvalidID(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListTasksHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	tasks := []domain.Task{
		{Name: "Task 1", Status: enum.Created},
		{Name: "Task 2", Status: enum.Started},
	}

	taskRepo.On("ListByFilter", mock.Anything, mock.Anything, 20, 0).
		Return(tasks, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks?limit=20&offset=0", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	taskRepo.AssertExpectations(t)
}

func TestListTasksHandler_EmptyResult(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("ListByFilter", mock.Anything, mock.Anything, 20, 0).
		Return([]domain.Task{}, int64(0), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestListTasksHandler_RepoError(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("ListByFilter", mock.Anything, mock.Anything, 20, 0).
		Return([]domain.Task(nil), int64(0), errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestUpdateTaskHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	existingTask := domain.Task{
		Name:   "Old Name",
		Status: enum.Created,
	}
	existingTask.ID = 1

	taskRepo.On("GetByID", mock.Anything, uint(1)).Return(existingTask, nil)
	taskRepo.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Task"), mock.Anything).Return(nil)

	newName := "New Name"
	body, _ := json.Marshal(dto.UpdateTaskReq{Name: &newName})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestUpdateTaskHandler_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	newName := "New Name"
	body, _ := json.Marshal(dto.UpdateTaskReq{Name: &newName})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestUpdateTaskHandler_Unauthorized(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouterNoAuth(taskSrv)

	newName := "New Name"
	body, _ := json.Marshal(dto.UpdateTaskReq{Name: &newName})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateTaskHandler_InvalidID(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	newName := "New Name"
	body, _ := json.Marshal(dto.UpdateTaskReq{Name: &newName})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTaskHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(1)).Return(domain.Task{}, nil)
	taskRepo.On("DeleteByID", mock.Anything, uint(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestDeleteTaskHandler_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestDeleteTaskHandler_InvalidID(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestArchiveTaskHandler(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	existingTask := domain.Task{
		Name:   "Task",
		Status: enum.Created,
	}
	existingTask.ID = 1

	taskRepo.On("GetByID", mock.Anything, uint(1)).Return(existingTask, nil)
	taskRepo.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Task"), []string{"status", "updated_by_user_id"}).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1/archive", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestArchiveTaskHandler_NotFound(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	taskRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.Task{}, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/999/archive", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	taskRepo.AssertExpectations(t)
}

func TestArchiveTaskHandler_Unauthorized(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouterNoAuth(taskSrv)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1/archive", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestArchiveTaskHandler_InvalidID(t *testing.T) {
	taskRepo := new(mockRepo.MockTaskRepo)
	taskSrv := services.NewTaskService(taskRepo)
	router := setupTaskRouter(taskSrv)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/abc/archive", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
