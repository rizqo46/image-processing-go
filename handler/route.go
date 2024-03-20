package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/usecase"
)

func SetupImageRoute(r *gin.Engine) {
	imageUsecase := usecase.NewImageUsecase()
	imageHandler := NewImageHandler(imageUsecase)

	r.
		POST("/", imageHandler.ProcessImage).
		POST("/png-to-jpeg", imageHandler.PngToJpeg).
		POST("/compress", imageHandler.CompressImages).
		POST("/resize", imageHandler.ResizeImages)
}
