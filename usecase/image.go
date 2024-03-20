package usecase

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"

	"github.com/rizqo46/image-processing-go/constants"
	"github.com/rizqo46/image-processing-go/dto"
	"gocv.io/x/gocv"
)

type ImageUsecase struct{}

func NewImageUsecase() ImageUsecase {
	return ImageUsecase{}
}

var (
	ErrOpenFile          = fmt.Errorf("failed to open a file")
	ErrReadFile          = fmt.Errorf("failed to read a file")
	ErrDetectContentType = fmt.Errorf("failed to detect content type")
)

func (uc ImageUsecase) ValidateAndProcessFilesRequest(files []*multipart.FileHeader, allowedContentTypes ...string) ([]dto.ImageData, error) {
	images := make([]dto.ImageData, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, ErrOpenFile
		}
		defer file.Close()

		bufReader := bufio.NewReader(file)
		sniff, err := bufReader.Peek(512)
		if err != nil {
			return nil, ErrDetectContentType
		}

		contentType := http.DetectContentType(sniff)
		if !slices.Contains(allowedContentTypes, contentType) {
			return nil, fmt.Errorf("filetype not allowed, only allow %+v", allowedContentTypes)
		}

		bytes, err := io.ReadAll(bufReader)
		if err != nil {
			return nil, ErrReadFile
		}

		images = append(images, dto.ImageData{
			Filename:    fileHeader.Filename,
			ContentType: contentType,
			ImageBytes:  bytes,
		})
	}

	return images, nil
}

func convretFilenameFromPngToJpeg(name string) string {
	return strings.TrimSuffix(name, "png") + "jpeg"
}

func (uc ImageUsecase) ConvertPngToJpeg(req []dto.ImageData) error {
	for i := range req {
		img, err := gocv.IMDecode(req[i].ImageBytes, gocv.IMReadAnyColor)
		if err != nil {
			return err
		}

		params := []int{gocv.IMWriteJpegQuality, 100}
		nativeBuffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, img, params)
		if err != nil {
			return err
		}

		req[i].Filename = convretFilenameFromPngToJpeg(req[i].Filename)
		req[i].ImageBytes = nativeBuffer.GetBytes()
	}

	return nil
}

func (uc ImageUsecase) CompressImages(req []dto.ImageData) error {
	imWriteContentTypeMapping := map[string]gocv.FileExt{
		constants.ContentTypeImagePng:  gocv.PNGFileExt,
		constants.ContentTypeImageJpeg: gocv.JPEGFileExt,
	}

	encodeParamFileExtMapping := map[gocv.FileExt][]int{
		gocv.PNGFileExt:  {gocv.IMWritePngCompression, 3},
		gocv.JPEGFileExt: {gocv.IMWriteJpegQuality, 95},
	}

	for i := range req {
		img, err := gocv.IMDecode(req[i].ImageBytes, gocv.IMReadUnchanged)
		if err != nil {
			return err
		}

		fileExt := imWriteContentTypeMapping[req[i].ContentType]
		params := encodeParamFileExtMapping[fileExt]
		nativeBuffer, err := gocv.IMEncodeWithParams(fileExt, img, params)
		if err != nil {
			return err
		}

		req[i].ImageBytes = nativeBuffer.GetBytes()
	}

	return nil
}

func (uc ImageUsecase) ResizeImages(req dto.ImageDataResize) error {
	imWriteContentTypeMapping := map[string]gocv.FileExt{
		constants.ContentTypeImagePng:  gocv.PNGFileExt,
		constants.ContentTypeImageJpeg: gocv.JPEGFileExt,
	}

	for i := range req.ImageDatas {
		img, err := gocv.IMDecode(req.ImageDatas[i].ImageBytes, gocv.IMReadUnchanged)
		if err != nil {
			return err
		}

		interpolationMethod := gocv.InterpolationCubic
		newImage := gocv.NewMat()
		gocv.Resize(img, &newImage, image.Pt(req.Width[i], req.Height[i]), 0, 0, interpolationMethod)

		fileExt := imWriteContentTypeMapping[req.ImageDatas[i].ContentType]
		nativeBuffer, err := gocv.IMEncode(fileExt, newImage)
		if err != nil {
			return err
		}

		req.ImageDatas[i].ImageBytes = nativeBuffer.GetBytes()
	}

	return nil
}

func (uc ImageUsecase) ProcessImages(req dto.ImageDataResize) error {
	for i := range req.ImageDatas {
		img, err := gocv.IMDecode(req.ImageDatas[i].ImageBytes, gocv.IMReadUnchanged)
		if err != nil {
			return err
		}

		interpolationMethod := gocv.InterpolationCubic
		newImage := gocv.NewMat()
		gocv.Resize(img, &newImage, image.Pt(req.Width[i], req.Height[i]), 0, 0, interpolationMethod)

		params := []int{gocv.IMWriteJpegQuality, 100}
		nativeBuffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, img, params)
		if err != nil {
			return err
		}

		req.ImageDatas[i].ImageBytes = nativeBuffer.GetBytes()
	}

	return nil
}
