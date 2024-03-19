package usecase

import (
	"context"
	"image"

	"github.com/rizqo46/image-processing-go/dto"
	"gocv.io/x/gocv"
)

type ImageUsecase struct{}

func NewImageUsecase() ImageUsecase {
	return ImageUsecase{}
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
