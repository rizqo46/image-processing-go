package handler

import (
	"bufio"
	"bytes"
	"fmt"
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

func parseResponseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}

var (
	ErrFailedToDetectContentType = fmt.Errorf("failed to detect content type")
	ErrFileTypeNotAllowed        = fmt.Errorf("file type not allowed, only support image/png")
)

func (h *imageHandler) ProcessImage(c *gin.Context) {
	var req dto.ImageRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	err := req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)
	sniff, err := bufReader.Peek(512)
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(ErrFailedToDetectContentType))
		return
	}

	contentType := http.DetectContentType(sniff)
	if contentType != "image/png" {
		c.JSON(http.StatusBadRequest, parseResponseError(ErrFileTypeNotAllowed))
		return
	}

	buf, err := io.ReadAll(bufReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	imgResult, err := h.imageUc.ProcessImage(c.Request.Context(), buf, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, parseResponseError(err))
		return
	}

	c.Status(http.StatusCreated)
	io.Copy(c.Writer, bytes.NewReader(imgResult))
}
