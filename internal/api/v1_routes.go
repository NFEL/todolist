package api

import (
	"graph-interview/internal/api/handlers"
	"graph-interview/internal/api/middlewares"
	"graph-interview/internal/cfg"
	"graph-interview/internal/repository/cache"
	"graph-interview/internal/repository/storage"
	storage_postgres "graph-interview/internal/repository/storage/postgres"
	"graph-interview/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterV1Handlers(cfg *cfg.Config, r gin.IRouter) error {
	db, err := storage.NewDB(&cfg.DB)
	if err != nil {
		return err
	}
	cache, err := cache.NewCache(cfg.Cache)
	if err != nil {
		return err
	}
	userRepo := storage_postgres.NewUserRepo(db)
	taskRepo := storage_postgres.NewTaskRepo(db)
	authSrv := services.NewAuthService(userRepo, cache.Client, cfg.Server.JWT.Secret)
	userSrv := services.NewUserService(userRepo)
	taskSrv := services.NewTaskService(taskRepo)

	authMiddleware := middlewares.AuthMiddleware(authSrv, cache.Client)

	pubRoutes(userSrv, r)
	r.Use(authMiddleware)
	authRoutes(userSrv, authSrv, taskSrv, r)
	return nil
}

func pubRoutes(userSrv *services.UserService, r gin.IRouter) {
	{
		auth := r.Group("/auth/")
		auth.GET("/login", handlers.CreateUser(userSrv))
	}
}

func authRoutes(
	userSrv *services.UserService,
	authSrv *services.AuthService,
	taskSrv *services.TaskService,
	r gin.IRouter) {
	{
		auth := r.Group("/auth/")
		auth.GET("/logout", handlers.Logout(authSrv))
	}
	{
		task := r.Group("/task/")
		task.GET("/", handlers.TaskList(taskSrv))

	}
}
