package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/handler"
	"github.com/rizqo46/image-processing-go/middleware"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.RequestBodyLimiter)

	handler.SetupImageRoute(r)

	port := "8080"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	_ = r.Run(":" + port)
}
