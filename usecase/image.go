package usecase

import (
	"bufio"
	"context"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/rizqo46/image-processing-go/dto"
	"gocv.io/x/gocv"
)

type ImageUsecase struct{}

func NewImageUsecase() ImageUsecase {
	return ImageUsecase{}
}

func convretFilenameFromPngToJpeg(name string) string {
	return strings.TrimSuffix(name, "png") + "jpeg"
}

func (uc ImageUsecase) ValidateAndParseBeforeConvertPngToJpeg(files []*multipart.FileHeader) ([]dto.ConvertImagePngToJpegParam, error) {
	results := make([]dto.ConvertImagePngToJpegParam, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		bufReader := bufio.NewReader(file)
		sniff, err := bufReader.Peek(512)
		if err != nil {
			return nil, err
		}

		contentType := http.DetectContentType(sniff)
		if contentType != "image/png" {
			return nil, err
		}

		bytes, err := io.ReadAll(bufReader)
		if err != nil {
			return nil, err
		}

		results = append(results, dto.ConvertImagePngToJpegParam{
			Filename:   convretFilenameFromPngToJpeg(fileHeader.Filename),
			ImageBytes: bytes,
		})
	}

	return results, nil
}

func (uc ImageUsecase) ConvertPngToJpeg(req []dto.ConvertImagePngToJpegParam) error {
	for i := range req {
		img, err := gocv.IMDecode(req[i].ImageBytes, gocv.IMReadAnyColor)
		if err != nil {
			return err
		}

		params := []int{gocv.IMWriteJpegQuality, 95}
		nativeBuffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, img, params)
		if err != nil {
			return err
		}

		req[i].ImageBytes = nativeBuffer.GetBytes()
	}

	return nil
}

func (uc ImageUsecase) ProcessImage(ctx context.Context, imageByte []byte, req dto.ImageRequest) ([]byte, error) {
	img, err := gocv.IMDecode(imageByte, gocv.IMReadAnyColor)
	if err != nil {
		return nil, err
	}

	interpolationMethod := gocv.InterpolationCubic

	if !req.Resize.IsNoResize() {
		newImage := gocv.NewMat()
		gocv.Resize(img, &newImage, image.Pt(req.Width, req.Height), 0, 0, interpolationMethod)
		img = newImage
	}

	params := []int{gocv.IMWriteJpegQuality, 96}
	nativeBuffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, img, params)
	if err != nil {
		return nil, err
	}

	newImgByte := nativeBuffer.GetBytes()
	return newImgByte, nil
}
