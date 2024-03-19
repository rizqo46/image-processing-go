package handler

import (
	"bufio"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/dto"
	"github.com/rizqo46/image-processing-go/usecase"
)

type imageHandler struct {
	imageUc usecase.ImageUsecase
}

func NewImageHandler(imageUc usecase.ImageUsecase) imageHandler {
	return imageHandler{imageUc: imageUc}
}

func (h *imageHandler) ProcessImage(c *gin.Context) {
	var req dto.ImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)

	c.Status(http.StatusCreated)
	io.Copy(c.Writer, bufReader)
}
