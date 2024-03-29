package handler

import (
	"archive/zip"
	"bytes"
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

	c.Status(http.StatusCreated)
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

	c.Status(http.StatusCreated)
	sendImagesRespAsZip(c, images)
}

func (h *imageHandler) ResizeImages(c *gin.Context) {
	var req dto.FilesResizeRequest
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

	imageDataResize := dto.ImageDataResize{
		ResizeRequest: req.ResizeRequest,
		ImageDatas:    images,
	}

	err = h.imageUc.ResizeImages(imageDataResize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, parseResponseError(err))
		return
	}

	c.Status(http.StatusCreated)
	sendImagesRespAsZip(c, imageDataResize.ImageDatas)
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
}

func (h *imageHandler) ProcessImage(c *gin.Context) {
	var req dto.FilesResizeRequest
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
		req.Files, constants.ContentTypeImagePng,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, parseResponseError(err))
		return
	}

	imageDataResize := dto.ImageDataResize{
		ResizeRequest: req.ResizeRequest,
		ImageDatas:    images,
	}
	err = h.imageUc.ProcessImages(imageDataResize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, parseResponseError(err))
		return
	}

	c.Status(http.StatusCreated)
	sendImagesRespAsZip(c, imageDataResize.ImageDatas)
}
