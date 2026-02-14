package dto

import (
	"graph-interview/internal/repository/enum"
	"time"
)

// Auth DTOs

type LoginUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JWTResp struct {
	Access    string    `json:"access"`
	Refresh   string    `json:"refresh"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// User DTOs

type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResp struct {
	ID uint `json:"id"`
}

type UserProfileResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
}

// Task DTOs

type CreateTaskReq struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

type UpdateTaskReq struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Status      *enum.TaskStatus `json:"status,omitempty"`
}

type TaskResp struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedByID uint      `json:"created_by_id"`
	UpdatedByID uint      `json:"updated_by_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskListResp struct {
	Tasks  []TaskResp `json:"tasks"`
	Total  int64      `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// Filter DTOs

type UserListFilter struct {
}

type TaskListFilter struct {
	Status    *enum.TaskStatus `json:"status,omitempty" form:"status"`
	Assignee  uint             `json:"assignee,omitempty" form:"assignee"`
	CreatedAt time.Time        `json:"created_at,omitempty" form:"created_at"`
	UpdatedAt time.Time        `json:"updated_at,omitempty" form:"updated_at"`
}

// Pagination

type PaginationQuery struct {
	Limit  int `form:"limit,default=20" binding:"min=1,max=100"`
	Offset int `form:"offset,default=0" binding:"min=0"`
}
