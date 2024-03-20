package dto

import (
	"fmt"
	"mime/multipart"
)

type ImageData struct {
	Filename    string
	ContentType string
	ImageBytes  []byte
}

type FilesRequest struct {
	Files []*multipart.FileHeader `form:"files[]"`
}

func (r FilesRequest) Validate() error {
	if len(r.Files) == 0 {
		return fmt.Errorf("files[] cannot be enmpy")
	}

	return nil
}

type FilesResizeRequest struct {
	ResizeRequest
	FilesRequest
}

func (r FilesResizeRequest) Validate() error {
	err := r.FilesRequest.Validate()
	if err != nil {
		return err
	}

	if len(r.Files) != len(r.Height) || len(r.Files) != len(r.Width) {
		return fmt.Errorf("len of files and resize param must be the same")
	}

	return r.ResizeRequest.Validate()
}

type ImageDataResize struct {
	ResizeRequest
	ImageDatas []ImageData
}

type ResizeRequest struct {
	Height []int `form:"height[]"`
	Width  []int `form:"width[]"`
}

func (r ResizeRequest) Validate() error {
	for _, v := range append(r.Height, r.Width...) {
		if v > 0 {
			continue
		}

		return fmt.Errorf("height and width must be large than zero")
	}

	return nil
}

type ImageRequest struct {
	File *multipart.FileHeader `form:"file"`
	Resize
}

type Resize struct {
	Height int `form:"resizeHeight"`
	Width  int `form:"resizeWidth"`
}

var (
	ErrFormFiledNamedFileShouldBeProvided                  = fmt.Errorf("form filed named file should be provided")
	ErrResizeRequestCannotBeNegativeNumber                 = fmt.Errorf("resizeHeight or resizeHeight cannot be negative number")
	ErrOneOfResizeParamShouldNotZeroWhenOneOfThemIsNonZero = fmt.Errorf("resizeHeight or resizeHeight should not zero when one of them is non zero")
)

func (r Resize) IsNoResize() bool {
	return r.Height == 0 && r.Width == 0
}

func (r Resize) Validate() error {
	if r.IsNoResize() {
		return nil
	}

	if r.Height < 0 || r.Width < 0 {
		return ErrResizeRequestCannotBeNegativeNumber
	}

	if r.Height == 0 || r.Width == 0 {
		return ErrOneOfResizeParamShouldNotZeroWhenOneOfThemIsNonZero
	}

	return nil
}

func (r ImageRequest) Validate() error {
	if r.File == nil {
		return ErrFormFiledNamedFileShouldBeProvided
	}

	return r.Resize.Validate()
}
