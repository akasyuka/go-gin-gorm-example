package main

import (
	"fmt"
	"log"
	"net"

	"github.com/akasyuka/service-a/config"
	"github.com/akasyuka/service-a/controller"
	"github.com/akasyuka/service-a/database"
	"github.com/akasyuka/service-a/metrics"
	"github.com/akasyuka/service-a/repository"
	"github.com/akasyuka/service-a/security"
	"github.com/akasyuka/service-a/service"

	userv1 "github.com/akasyuka/service-a/gen/user/v1"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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

	// ===== Initialize repositories / services =====
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// ===== gRPC controller =====
	userGrpcController := controller.NewUserGrpcController(userService)

	// =====================================================
	// ================= REST (Gin) ========================
	// =====================================================

	r := gin.Default()

	// ===== Prometheus metrics =====
	if cfg.Monitoring.Prometheus.Enabled {
		metrics.InitMetrics()
		r.Use(metrics.GinMetricsMiddleware())
		r.GET(
			cfg.Monitoring.Prometheus.MetricsPath,
			gin.WrapH(metrics.MetricsHandler()),
		)
		fmt.Printf(
			"Prometheus metrics enabled: path=%s\n",
			cfg.Monitoring.Prometheus.MetricsPath,
		)
	}

	// ===== Keycloak JWT middleware =====
	jwks, err := security.InitJWKS(cfg.Auth.Keycloak.JWKSURL)
	if err != nil {
		log.Fatalf("failed to initialize JWKS: %v", err)
	}

	private := r.Group("/api")
	private.Use(security.JWTMiddleware(jwks))

	// ===== Public routes =====
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// ===== Run HTTP server (async) =====
	go func() {
		addr := fmt.Sprintf(
			"%s:%d",
			cfg.Server.HTTP.Host,
			cfg.Server.HTTP.Port,
		)
		if err := r.Run(addr); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	// =====================================================
	// ================= gRPC ==============================
	// =====================================================

	grpcAddr := fmt.Sprintf(
		"%s:%d",
		cfg.Server.GRPC.Host,
		cfg.Server.GRPC.Port,
	)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer(
	// тут позже:
	// grpc.UnaryInterceptor(...)
	)

	userv1.RegisterUserServiceServer(
		grpcServer,
		userGrpcController,
	)

	fmt.Printf("gRPC server started at %s\n", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run gRPC server: %v", err)
	}
}
