package dto

import "mime/multipart"

type ImageRequest struct {
	File *multipart.FileHeader `form:"file"`
	Resize
}

type Resize struct {
	Height int `form:"resizeHeight"`
	Width  int `form:"resizeWidth"`
}
