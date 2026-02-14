package services

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
	mockRepo "graph-interview/internal/repository/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	svc := NewUserService(userRepo)

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return(uint(1), nil)

	req := dto.CreateUserReq{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	resp, err := svc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	userRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	svc := NewUserService(userRepo)

	userRepo.On("GetByField", mock.Anything, "username", "existinguser").
		Return(domain.User{Username: "existinguser"}, nil)

	req := dto.CreateUserReq{
		Username: "existinguser",
		Password: "password123",
		Email:    "test@example.com",
	}

	resp, err := svc.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")
	userRepo.AssertExpectations(t)
}

func TestGetProfile_Success(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	svc := NewUserService(userRepo)

	userRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Avatar:   "avatar.png",
		}, nil)

	resp, err := svc.GetProfile(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	userRepo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	svc := NewUserService(userRepo)

	userRepo.On("GetByID", mock.Anything, uint(999)).
		Return(domain.User{}, gorm.ErrRecordNotFound)

	resp, err := svc.GetProfile(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}
