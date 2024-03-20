package handler

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rizqo46/image-processing-go/constants"
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
	ErrFilesCannotBeEmpty        = fmt.Errorf("files cannot be empty")
)

func (h *imageHandler) PngToJpeg(c *gin.Context) {
	var req dto.FilesRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	err := req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	images, err := h.imageUc.ValidateAndProcessFilesRequest(req.Files, constants.ContentTypeImagePng)
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	err = h.imageUc.ConvertPngToJpeg(images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, parseResponseError(err))
		return
	}

	sendImagesRespAsZip(c, images)
}

func (h *imageHandler) CompressImages(c *gin.Context) {
	var req dto.FilesRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	err := req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	images, err := h.imageUc.ValidateAndProcessFilesRequest(
		req.Files, constants.ContentTypeImagePng, constants.ContentTypeImageJpeg,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	err = h.imageUc.CompressImages(images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, parseResponseError(err))
		return
	}

	sendImagesRespAsZip(c, images)
}

func sendImagesRespAsZip(c *gin.Context, images []dto.ImageData) {
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	now := time.Now()
	for _, image := range images {
		w, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:     image.Filename,
			Method:   zip.Deflate,
			Modified: now,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, parseResponseError(err))
			return
		}

		if _, err := io.Copy(w, bytes.NewReader(image.ImageBytes)); err != nil {
			c.JSON(http.StatusInternalServerError, parseResponseError(err))
			return
		}
	}

	c.Status(http.StatusCreated)
}

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
