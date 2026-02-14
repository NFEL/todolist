package domain

import (
	"graph-interview/internal/repository/enum"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name            string
	Description     string
	Status          enum.TaskStatus
	CreatedBy       *User `gorm:"foreignKey:CreatedByUserID"`
	CreatedByUserID uint
	UpdatedBy       *User `gorm:"foreignKey:UpdatedByUserID"`
	UpdatedByUserID uint
	Assignees       []*User `gorm:"many2many:user_tasks;"`
}
