package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/handler"
	"github.com/rizqo46/image-processing-go/middleware"
	"github.com/rizqo46/image-processing-go/usecase"
)

func main() {
	imageUsecase := usecase.NewImageUsecase()
	imageHandler := handler.NewImageHandler(imageUsecase)

	r := gin.Default()
	r.Use(middleware.RequestBodyLimiter)
	r.POST("/", imageHandler.ProcessImage)
	r.POST("/png-to-jpeg", imageHandler.PngToJpeg)
	r.POST("/compress", imageHandler.CompressImages)
	r.POST("/resize", imageHandler.ResizeImages)

	r.Run(":10000")
}
