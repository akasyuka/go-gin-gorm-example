package main

import (
	"log"

	"github.com/akasyuka/go-gin-gorm-example/config"
	"github.com/akasyuka/go-gin-gorm-example/controller"
	"github.com/akasyuka/go-gin-gorm-example/database"
	"github.com/akasyuka/go-gin-gorm-example/repository"
	"github.com/akasyuka/go-gin-gorm-example/service"
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
	userController.RegisterRoutes(r)

	// ===== Run server =====
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
