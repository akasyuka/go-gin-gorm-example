package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Простая проверка здоровья
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Запуск сервера на 8081
	r.Run(":8081")
}
