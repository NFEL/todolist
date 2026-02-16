package api

import (
	"context"
	"fmt"
	"graph-interview/internal/api/middlewares"
	"graph-interview/internal/cfg"
	"graph-interview/pkg/logger"
	"os/signal"
	"syscall"

	_ "graph-interview/docs"
	_ "net/http/pprof"

	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func genSwagHandler(baseURL string) gin.HandlerFunc {
	var conf func(*ginSwagger.Config)
	switch cfg.Cfg.Environment {
	case cfg.Local:
		conf = ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.Cfg.Server.Port))
	default:
		conf = ginSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", baseURL))
	}
	return ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		conf,
		ginSwagger.PersistAuthorization(true),
	)
}

// @title           Task Manager API
// @version         1.0
// @description     A RESTful Task Manager API with JWT authentication
// @host            localhost:3154
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func ServeREST(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize tracing
	tp, err := initTracer()
	if err != nil {
		logger.Logger.Warn("failed to initialize tracer, continuing without tracing", "err", err)
	} else {
		defer func() {
			_ = tp.Shutdown(context.Background())
		}()
	}

	listeningUrl := fmt.Sprintf("%s:%d", cfg.Cfg.Server.Host, cfg.Cfg.Server.Port)
	g, err := graceful.Default(graceful.WithAddr(listeningUrl))
	if err != nil {
		return err
	}
	defer g.Close()

	// Middleware stack
	g.Use(middlewares.PrometheusMiddleware())
	g.Use(middlewares.CorsMiddleware(cfg.Cfg.Server.Cors))
	g.Use(middlewares.JSONMiddleware())
	g.Use(middlewares.I18nMiddleware())
	g.Use(middlewares.AccessLogMiddleware(logger.Logger))

	// OpenTelemetry tracing middleware
	// if tp != nil {
	// 	g.Use(otelgin.Middleware("task-manager"))
	// }

	_ = tp

	// Health check
	g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"msg": "ok"})
	})

	// Swagger
	g.GET("/swagger/*any", genSwagHandler(""))

	// Register API routes
	if err := RegisterV1Handlers(cfg.Cfg, g.Group("/v1")); err != nil {
		return err
	}

	logger.Logger.Info("started gin", "url", listeningUrl)
	if err := g.RunWithContext(ctx); err != nil && err != context.Canceled {
		return err
	}
	return nil
}
