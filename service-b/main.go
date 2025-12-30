package main

import (
	"fmt"
	"sync"
)

//	func main() {
//		r := gin.Default()
//
//		// Простая проверка здоровья
//		r.GET("/health", func(c *gin.Context) {
//			c.JSON(200, gin.H{
//				"status": "UP",
//			})
//		})
//
//		// Запуск сервера на 8081
//		r.Run(":8081")
//	}
func main() {
	group := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			fmt.Println("Goroutine", i)
		}(i)
	}
	group.Wait()
}
