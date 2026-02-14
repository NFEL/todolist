package handlers

import (
	"errors"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task for the authenticated user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateTaskReq  true  "Task data"
// @Success      201   {object}  dto.Response{data=dto.TaskResp}
// @Failure      400   {object}  dto.Response
// @Failure      401   {object}  dto.Response
// @Router       /v1/tasks [post]
func CreateTask(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserID(c)
		if err != nil {
			dto.ErrUnauthorized(c, api_error.ErrUnauthorized)
			return
		}

		req := dto.CreateTaskReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			dto.Err(c, err)
			return
		}

		resp, err := taskSrv.CreateTask(c, req, userID)
		if err != nil {
			dto.ErrInternal(c, err)
			return
		}
		dto.Created(c, "task created", resp)
	}
}

// GetTask godoc
// @Summary      Get a task by ID
// @Description  Retrieve a single task by its ID
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  dto.Response{data=dto.TaskResp}
// @Failure      404  {object}  dto.Response
// @Router       /v1/tasks/{id} [get]
func GetTask(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			dto.Err(c, err)
			return
		}

		resp, err := taskSrv.GetTask(c, uint(taskID))
		if err != nil {
			if errors.Is(err, api_error.ErrTaskNotFound) {
				dto.ErrNotFound(c, err)
				return
			}
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "task retrieved", resp)
	}
}

// ListTasks godoc
// @Summary      List tasks
// @Description  List tasks with optional filtering and pagination
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        limit     query     int     false  "Limit"     default(20)
// @Param        offset    query     int     false  "Offset"    default(0)
// @Param        status    query     int     false  "Status filter (0=Created,1=Started,2=Done,3=Failed,4=Delayed,5=Canceled)"
// @Param        assignee  query     int     false  "Assignee user ID"
// @Success      200       {object}  dto.Response{data=dto.TaskListResp}
// @Failure      400       {object}  dto.Response
// @Router       /v1/tasks [get]
func ListTasks(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := dto.PaginationQuery{Limit: 20, Offset: 0}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			dto.Err(c, err)
			return
		}

		filter := dto.TaskListFilter{}
		if err := c.ShouldBindQuery(&filter); err != nil {
			dto.Err(c, err)
			return
		}

		resp, err := taskSrv.ListTasks(c, filter, pagination.Limit, pagination.Offset)
		if err != nil {
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "tasks retrieved", resp)
	}
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update task fields (name, description, status)
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int              true  "Task ID"
// @Param        body  body      dto.UpdateTaskReq  true  "Fields to update"
// @Success      200   {object}  dto.Response{data=dto.TaskResp}
// @Failure      400   {object}  dto.Response
// @Failure      404   {object}  dto.Response
// @Router       /v1/tasks/{id} [put]
func UpdateTask(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserID(c)
		if err != nil {
			dto.ErrUnauthorized(c, api_error.ErrUnauthorized)
			return
		}

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			dto.Err(c, err)
			return
		}

		req := dto.UpdateTaskReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			dto.Err(c, err)
			return
		}

		resp, err := taskSrv.UpdateTask(c, uint(taskID), req, userID)
		if err != nil {
			if errors.Is(err, api_error.ErrTaskNotFound) {
				dto.ErrNotFound(c, err)
				return
			}
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "task updated", resp)
	}
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Soft-delete a task by ID
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  dto.Response
// @Failure      404  {object}  dto.Response
// @Router       /v1/tasks/{id} [delete]
func DeleteTask(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			dto.Err(c, err)
			return
		}

		if err := taskSrv.DeleteTask(c, uint(taskID)); err != nil {
			if errors.Is(err, api_error.ErrTaskNotFound) {
				dto.ErrNotFound(c, err)
				return
			}
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "task deleted", nil)
	}
}

// ArchiveTask godoc
// @Summary      Archive a task
// @Description  Set a task's status to Canceled (archived)
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  dto.Response{data=dto.TaskResp}
// @Failure      404  {object}  dto.Response
// @Router       /v1/tasks/{id}/archive [patch]
func ArchiveTask(taskSrv *services.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserID(c)
		if err != nil {
			dto.ErrUnauthorized(c, api_error.ErrUnauthorized)
			return
		}

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			dto.Err(c, err)
			return
		}

		resp, err := taskSrv.ArchiveTask(c, uint(taskID), userID)
		if err != nil {
			if errors.Is(err, api_error.ErrTaskNotFound) {
				dto.ErrNotFound(c, err)
				return
			}
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "task archived", resp)
	}
}
