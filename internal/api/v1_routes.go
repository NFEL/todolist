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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterV1Handlers(cfg *cfg.Config, r gin.IRouter) error {
	db, err := storage.NewDB(&cfg.DB)
	if err != nil {
		return err
	}
	cacheStore, err := cache.NewCache(cfg.Cache)
	if err != nil {
		return err
	}
	userRepo := storage_postgres.NewUserRepo(db)
	taskRepo := storage_postgres.NewTaskRepo(db)
	authSrv := services.NewAuthService(userRepo, cacheStore.Client, cfg.Server.JWT.Secret)
	userSrv := services.NewUserService(userRepo)
	taskSrv := services.NewTaskService(taskRepo)

	authMiddleware := middlewares.AuthMiddleware(authSrv, cacheStore.Client)

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	pubRoutes(userSrv, authSrv, r)
	authRoutes(userSrv, authSrv, taskSrv, r, authMiddleware)
	return nil
}

func pubRoutes(userSrv *services.UserService, authSrv *services.AuthService, r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(userSrv))
		auth.POST("/login", handlers.Login(authSrv))
		auth.POST("/refresh", handlers.RefreshToken(authSrv))
	}
}

func authRoutes(
	userSrv *services.UserService,
	authSrv *services.AuthService,
	taskSrv *services.TaskService,
	r gin.IRouter,
	authMiddleware gin.HandlerFunc,
) {
	protected := r.Group("")
	protected.Use(authMiddleware)
	{
		// Auth routes
		authGroup := protected.Group("/auth")
		authGroup.POST("/logout", handlers.Logout(authSrv))

		// User routes
		userGroup := protected.Group("/user")
		userGroup.GET("/profile", handlers.GetProfile(userSrv))

		// Task routes
		taskGroup := protected.Group("/tasks")
		taskGroup.POST("", handlers.CreateTask(taskSrv))
		taskGroup.GET("", handlers.ListTasks(taskSrv))
		taskGroup.GET("/:id", handlers.GetTask(taskSrv))
		taskGroup.PUT("/:id", handlers.UpdateTask(taskSrv))
		taskGroup.DELETE("/:id", handlers.DeleteTask(taskSrv))
		taskGroup.PATCH("/:id/archive", handlers.ArchiveTask(taskSrv))
	}
}
