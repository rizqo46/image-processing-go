package dto

import "mime/multipart"

type ImageRequest struct {
	File *multipart.FileHeader `form:"file"`
}
