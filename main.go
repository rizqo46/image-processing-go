package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/handler"
	"github.com/rizqo46/image-processing-go/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.RequestBodyLimiter)

	handler.SetupImageRoute(r)

	r.Run(":10000")
}
