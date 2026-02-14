package dto

import "time"

type LoginUserReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type JWTResp struct {
	Access    string    `json:"access,omitempty"`
	Refresh   string    `json:"refresh,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type CreateUserReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type CreateUserResp struct {
	ID uint `json:"id,omitempty"`
}

type TaskResp struct {
	ID          uint   `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
}

type TaskListResp struct {
	Tasks []TaskResp `json:"tasks,omitempty"`
}

type UserListFilter struct {
}

type TaskListFilter struct {
	Assignee  uint      `json:"assignee,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
