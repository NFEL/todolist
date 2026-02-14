package handlers

import (
	"bytes"
	"encoding/json"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
	mockRepo "graph-interview/internal/repository/mock"
	"graph-interview/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupUserRouter(userSrv *services.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/register", Register(userSrv))

	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", "1")
		c.Next()
	})
	protected.GET("/profile", GetProfile(userSrv))

	return r
}

func setupAuthRouter(t *testing.T) (*gin.Engine, *mockRepo.MockUserRepo, *miniredis.Miniredis) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	userRepo := new(mockRepo.MockUserRepo)
	authSrv := services.NewAuthService(userRepo, rdb, []byte("test-secret"))

	r := gin.New()
	r.POST("/login", Login(authSrv))
	r.POST("/refresh", RefreshToken(authSrv))

	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", "1")
		c.Next()
	})
	protected.POST("/logout", Logout(authSrv))

	return r, userRepo, mr
}

func TestRegisterHandler_Success(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	userRepo.On("GetByField", mock.Anything, "username", "newuser").
		Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return(uint(1), nil)

	body, _ := json.Marshal(dto.CreateUserReq{
		Username: "newuser",
		Password: "password123",
		Email:    "new@example.com",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	userRepo.AssertExpectations(t)
}

func TestRegisterHandler_InvalidBody(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	body, _ := json.Marshal(map[string]string{"username": "ab"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterHandler_EmptyBody(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterHandler_DuplicateUsername(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	userRepo.On("GetByField", mock.Anything, "username", "existing").
		Return(domain.User{Username: "existing"}, nil)

	body, _ := json.Marshal(dto.CreateUserReq{
		Username: "existing",
		Password: "password123",
		Email:    "new@example.com",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	userRepo.AssertExpectations(t)
}

func TestGetProfileHandler_Success(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	userRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.User{
			Username: "testuser",
			Email:    "test@example.com",
		}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	userRepo.AssertExpectations(t)
}

func TestGetProfileHandler_Unauthorized(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/profile", GetProfile(userSrv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfileHandler_UserNotFound(t *testing.T) {
	userRepo := new(mockRepo.MockUserRepo)
	userSrv := services.NewUserService(userRepo)
	router := setupUserRouter(userSrv)

	userRepo.On("GetByID", mock.Anything, uint(1)).
		Return(domain.User{}, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	userRepo.AssertExpectations(t)
}

func TestLoginHandler_Success(t *testing.T) {
	router, userRepo, mr := setupAuthRouter(t)
	defer mr.Close()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := domain.User{Username: "testuser", Password: string(hashed)}
	user.ID = 1

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(user, nil)

	body, _ := json.Marshal(dto.LoginUserReq{
		Username: "testuser",
		Password: "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	userRepo.AssertExpectations(t)
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	router, _, mr := setupAuthRouter(t)
	defer mr.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	router, userRepo, mr := setupAuthRouter(t)
	defer mr.Close()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{Username: "testuser", Password: string(hashed)}
	user.ID = 1

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(user, nil)

	body, _ := json.Marshal(dto.LoginUserReq{
		Username: "testuser",
		Password: "wrong",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	userRepo.AssertExpectations(t)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	router, userRepo, mr := setupAuthRouter(t)
	defer mr.Close()

	userRepo.On("GetByField", mock.Anything, "username", "nonexistent").
		Return(domain.User{}, gorm.ErrRecordNotFound)

	body, _ := json.Marshal(dto.LoginUserReq{
		Username: "nonexistent",
		Password: "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	userRepo.AssertExpectations(t)
}

func TestLogoutHandler_WithBearerToken(t *testing.T) {
	router, userRepo, mr := setupAuthRouter(t)
	defer mr.Close()

	// First login to get a token
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := domain.User{Username: "testuser", Password: string(hashed)}
	user.ID = 1

	userRepo.On("GetByField", mock.Anything, "username", "testuser").
		Return(user, nil)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	authSrv := services.NewAuthService(userRepo, rdb, []byte("test-secret"))

	tokens, _ := authSrv.IssueTokens("1")
	_ = authSrv.Persist(t.Context(), tokens)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.Access)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogoutHandler_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr, _ := miniredis.Run()
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	userRepo := new(mockRepo.MockUserRepo)
	authSrv := services.NewAuthService(userRepo, rdb, []byte("test-secret"))

	r := gin.New()
	r.POST("/logout", Logout(authSrv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshTokenHandler_InvalidBody(t *testing.T) {
	router, _, mr := setupAuthRouter(t)
	defer mr.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserID_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "not-a-number")

	_, err := getUserID(c)
	assert.Error(t, err)
}

func TestGetUserID_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	_, err := getUserID(c)
	assert.Error(t, err)
}

func TestGetUserID_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "42")

	id, err := getUserID(c)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), id)
}
