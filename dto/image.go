package dto

import (
	"fmt"
	"mime/multipart"
)

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
