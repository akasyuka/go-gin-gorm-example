package main

import (
	"fmt"
	"github.com/akasyuka/service-a/config"
	"github.com/akasyuka/service-a/controller"
	"github.com/akasyuka/service-a/database"
	"github.com/akasyuka/service-a/repository"
	"github.com/akasyuka/service-a/service"
	"log"

	"github.com/akasyuka/service-a/metrics"
	"github.com/gin-gonic/gin"
)

func main() {
	// ===== Load config =====
	cfg, err := config.Load("config/application.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ===== Initialize database =====
	db, err := database.NewPostgres(cfg.Database.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// ===== Initialize repositories / services / controllers =====
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// ===== Setup Gin router =====
	r := gin.Default()

	// ===== Prometheus metrics =====
	if cfg.Monitoring.Prometheus.Enabled {
		// Инициализация метрик
		metrics.InitMetrics()

		// Middleware для REST маршрутов
		r.Use(metrics.GinMetricsMiddleware())

		// Endpoint /metrics
		r.GET(cfg.Monitoring.Prometheus.MetricsPath, gin.WrapH(metrics.MetricsHandler()))

		fmt.Printf("Prometheus metrics enabled: path=%s\n", cfg.Monitoring.Prometheus.MetricsPath)
	}

	// ===== Register user routes =====
	userController.RegisterRoutes(r)

	// ===== Run server =====
	addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
