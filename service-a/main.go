package main

import (
	"fmt"
	"log"

	"github.com/akasyuka/service-a/config"
	"github.com/akasyuka/service-a/controller"
	"github.com/akasyuka/service-a/database"
	"github.com/akasyuka/service-a/metrics"
	"github.com/akasyuka/service-a/repository"
	"github.com/akasyuka/service-a/security"
	"github.com/akasyuka/service-a/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// ===== Load config =====
	cfg, err := config.Load("./application.yaml")
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
		metrics.InitMetrics()
		r.Use(metrics.GinMetricsMiddleware())
		r.GET(cfg.Monitoring.Prometheus.MetricsPath, gin.WrapH(metrics.MetricsHandler()))
		fmt.Printf("Prometheus metrics enabled: path=%s\n", cfg.Monitoring.Prometheus.MetricsPath)
	}

	// ===== Keycloak JWT middleware =====
	jwks, err := security.InitJWKS(cfg.Auth.Keycloak.JWKSURL)
	if err != nil {
		log.Fatalf("failed to initialize JWKS: %v", err)
	}

	// Приватные роуты через JWT middleware
	private := r.Group("/api")
	private.Use(security.JWTMiddleware(jwks))

	// ===== Register user routes на RouterGroup =====
	userController.RegisterRoutes(private)

	// ===== Optional public routes =====
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// ===== Run server =====
	addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
