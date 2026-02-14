package services

import (
	"context"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/domain"
	mockRepo "graph-interview/internal/repository/mock"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T) (*AuthService, *mockRepo.MockUserRepo, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	userRepo := new(mockRepo.MockUserRepo)
	authSrv := NewAuthService(userRepo, rdb, "test-secret")
	if rdb == nil {
		t.FailNow()
	}
	return authSrv, userRepo, mr
}

func TestLoginUser_Success(t *testing.T) {
	authSrv, userRepo, mr := setupAuthTest(t)
	defer mr.Close()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := domain.User{
		Username: "testuser",
		Password: string(hashed),
	}
	user.ID = 1

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(user, nil)

	req := dto.LoginUserReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := authSrv.LoginUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Access)
	assert.NotEmpty(t, resp.Refresh)
	userRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidUsername(t *testing.T) {
	authSrv, userRepo, mr := setupAuthTest(t)
	defer mr.Close()

	userRepo.On("GetByField", mock.Anything, "username", "nonexistent").
		Return(domain.User{}, gorm.ErrRecordNotFound)

	req := dto.LoginUserReq{
		Username: "nonexistent",
		Password: "password123",
	}

	resp, err := authSrv.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, api_error.ErrInvalidCredentials, err)
	userRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	authSrv, userRepo, mr := setupAuthTest(t)
	defer mr.Close()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := domain.User{
		Username: "testuser",
		Password: string(hashed),
	}
	user.ID = 1

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(user, nil)

	req := dto.LoginUserReq{
		Username: "testuser",
		Password: "wrongpassword",
	}

	resp, err := authSrv.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, api_error.ErrInvalidCredentials, err)
	userRepo.AssertExpectations(t)
}

func TestIssueTokens(t *testing.T) {
	authSrv, _, mr := setupAuthTest(t)
	defer mr.Close()

	tokens, err := authSrv.IssueTokens("1")

	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.Access)
	assert.NotEmpty(t, tokens.Refresh)
	assert.Equal(t, "1", tokens.UserID)
}

func TestParseToken(t *testing.T) {
	authSrv, _, mr := setupAuthTest(t)
	defer mr.Close()

	tokens, err := authSrv.IssueTokens("42")
	assert.NoError(t, err)

	claims, err := authSrv.ParseToken(tokens.Access)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "42", claims.Subject)
}

func TestParseToken_Invalid(t *testing.T) {
	authSrv, _, mr := setupAuthTest(t)
	defer mr.Close()

	claims, err := authSrv.ParseToken("invalid-token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestPersistAndRevoke(t *testing.T) {
	authSrv, _, mr := setupAuthTest(t)
	defer mr.Close()

	tokens, err := authSrv.IssueTokens("1")
	assert.NoError(t, err)

	// Persist
	err = authSrv.Persist(context.Background(), tokens)
	assert.NoError(t, err)

	// Verify token exists in Redis
	assert.True(t, mr.Exists("access:"+tokens.JTIAcc))
	assert.True(t, mr.Exists("refresh:"+tokens.JTIRef))

	// Revoke
	err = authSrv.RevokeToken(context.Background(), tokens)
	assert.NoError(t, err)

	// Verify token removed from Redis
	assert.False(t, mr.Exists("access:"+tokens.JTIAcc))
	assert.False(t, mr.Exists("refresh:"+tokens.JTIRef))
}

func TestRefreshToken_Success(t *testing.T) {
	authSrv, _, mr := setupAuthTest(t)
	defer mr.Close()

	tokens, err := authSrv.IssueTokens("1")
	assert.NoError(t, err)

	err = authSrv.Persist(context.Background(), tokens)
	assert.NoError(t, err)

	resp, err := authSrv.RefreshToken(context.Background(), tokens.Refresh)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Access)
	assert.NotEmpty(t, resp.Refresh)
	// Old refresh token should be gone
	assert.False(t, mr.Exists("refresh:"+tokens.JTIRef))
}
