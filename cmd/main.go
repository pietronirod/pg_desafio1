package main

import (
	"fmt"
	"rate-limiter/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	config.InitLogger()
	defer config.Logger.Sync()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Rate Limiter funcionando!",
		})
	})

	port := config.Cfg.ServerPort
	fmt.Println("Servidor rodando na porta:", port)
	if err := r.Run(":" + port); err != nil {
		config.Logger.Error("Erro ao iniciar servidor", err)
	}
}
