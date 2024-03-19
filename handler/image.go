package handler

import (
	"bufio"
	"bytes"
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
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := req.Validate()
	if err != nil {
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
	sniff, err := bufReader.Peek(512)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to detect content type"})
		return
	}

	contentType := http.DetectContentType(sniff)
	if contentType != "image/png" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "file type not allowed, only support image/png"},
		)
		return
	}

	buf, err := io.ReadAll(bufReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imgResult, err := h.imageUc.ProcessImage(c.Request.Context(), buf, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
	io.Copy(c.Writer, bytes.NewReader(imgResult))
}
