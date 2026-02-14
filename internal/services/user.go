package services

import (
	"context"
	"errors"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserReq) (*dto.CreateUserResp, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(hashed)
	{
		_, err := s.UserRepo.GetByField(ctx, "username", req.Username)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api_error.UsernameExists(req.Username)
		}
	}
	ID, err := s.UserRepo.Create(ctx, &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Avatar:   "",
	})
	if err != nil {
		return nil, err
	}
	return &dto.CreateUserResp{
		ID: ID,
	}, nil
}

func (s *UserService) UserLogin(ctx context.Context, req dto.LoginUserReq) (*dto.JWTResp, error) {
	return nil, nil
}
