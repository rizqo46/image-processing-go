package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/handler"
	"github.com/rizqo46/image-processing-go/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.RequestBodyLimiter)

	handler.SetupImageRoute(r)

	port := flag.String("port", "8080", "define custom port")

	flag.Parse()
	_ = r.Run(":" + *port)
}
