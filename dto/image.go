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

func (r Resize) IsNoResize() bool {
	return r.Height == 0 && r.Width == 0
}

func (r Resize) Validate() error {
	if r.IsNoResize() {
		return nil
	}

	if r.Height < 0 || r.Width < 0 {
		return fmt.Errorf("resizeHeight or resizeHeight cannot be negative number")
	}

	if r.Height == 0 || r.Width == 0 {
		return fmt.Errorf("resizeHeight or resizeHeight should not zero when one of them is non zero")
	}

	return nil
}

func (r ImageRequest) Validate() error {
	return r.Resize.Validate()
}
