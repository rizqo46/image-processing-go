package usecase

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/rizqo46/image-processing-go/constants"
	"github.com/rizqo46/image-processing-go/dto"
)

func TestImageUsecase_ValidateAndProcessFilesRequest(t *testing.T) {
	createMultipartFileheaders := func(filePaths ...string) []*multipart.FileHeader {
		var buff bytes.Buffer
		buffWriter := io.Writer(&buff)
		formWriter := multipart.NewWriter(buffWriter)

		for _, filePath := range filePaths {
			file, err := os.Open(filePath)
			if err != nil {
				t.Errorf("failed to open test file%v", err)
				return nil
			}

			formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(formPart, file); err != nil {
				log.Fatal(err)
			}
			file.Close()
		}

		formWriter.Close()

		buffReader := bytes.NewReader(buff.Bytes())
		formReader := multipart.NewReader(buffReader, formWriter.Boundary())

		multipartForm, err := formReader.ReadForm(1 << 20)
		if err != nil {
			log.Fatal(err)
		}

		files := multipartForm.File["file"]
		return files
	}

	type args struct {
		files               []*multipart.FileHeader
		allowedContentTypes []string
	}
	tests := []struct {
		name    string
		args    args
		want    []dto.ImageData
		wantErr bool
	}{

		{
			name: "success validate and process files",
			args: args{
				files:               createMultipartFileheaders(".././imagetest/flower.png"),
				allowedContentTypes: []string{constants.ContentTypeImagePng},
			},
			want: []dto.ImageData{
				{ContentType: constants.ContentTypeImagePng},
			},
			wantErr: false,
		},

		{
			name: "failed content type not allowed",
			args: args{
				files:               createMultipartFileheaders(".././imagetest/cat.jpg"),
				allowedContentTypes: []string{constants.ContentTypeImagePng},
			},
			want:    []dto.ImageData{},
			wantErr: true,
		},

		{
			name: "failed error on open file",
			args: args{
				files:               []*multipart.FileHeader{{}},
				allowedContentTypes: []string{constants.ContentTypeImagePng},
			},
			want:    []dto.ImageData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := ImageUsecase{}
			got, err := uc.ValidateAndProcessFilesRequest(tt.args.files, tt.args.allowedContentTypes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageUsecase.ValidateAndProcessFilesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, res := range got {
				assert.Equal(t, res.ContentType, tt.want[i].ContentType)
			}
		})

	}
}

func generateImageDatas(t *testing.T, filePaths ...string) []dto.ImageData {
	images := make([]dto.ImageData, 0, len(filePaths))
	for _, filePath := range filePaths {
		contentType := ""
		if strings.HasSuffix(filePath, ".png") {
			contentType = constants.ContentTypeImagePng
		} else if strings.HasSuffix(filePath, ".jpg") {
			contentType = constants.ContentTypeImageJpeg
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("failed to open file %v", err)
			return nil
		}

		images = append(images, dto.ImageData{
			Filename:    filePath,
			ContentType: contentType,
			ImageBytes:  file,
		})

	}

	return images
}

func TestImageUsecase_ConvertPngToJpeg(t *testing.T) {
	type args struct {
		req []dto.ImageData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success convert image from png to jpeg",
			args: args{
				req: generateImageDatas(t, ".././imagetest/flower.png"),
			},
			wantErr: false,
		},
		{
			name: "failed on decode image",
			args: args{
				req: []dto.ImageData{{}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := ImageUsecase{}
			if err := uc.ConvertPngToJpeg(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ImageUsecase.ConvertPngToJpeg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageUsecase_CompressImages(t *testing.T) {
	type args struct {
		req []dto.ImageData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success compress image",
			args: args{
				req: generateImageDatas(t, ".././imagetest/flower.png"),
			},
			wantErr: false,
		},
		{
			name: "failed on decode image",
			args: args{
				req: []dto.ImageData{{}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := ImageUsecase{}
			if err := uc.CompressImages(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ImageUsecase.CompressImages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageUsecase_ResizeImages(t *testing.T) {
	type args struct {
		req dto.ImageDataResize
	}
	tests := []struct {
		name    string
		uc      ImageUsecase
		args    args
		wantErr bool
	}{
		{
			name: "success compress image",
			args: args{
				req: dto.ImageDataResize{
					ImageDatas: generateImageDatas(t, ".././imagetest/flower.png"),
					ResizeRequest: dto.ResizeRequest{
						Height: []int{30},
						Width:  []int{30},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed on decode image",
			args: args{
				req: dto.ImageDataResize{
					ImageDatas:    []dto.ImageData{{}},
					ResizeRequest: dto.ResizeRequest{Height: []int{1}, Width: []int{1}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := ImageUsecase{}
			if err := uc.ResizeImages(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ImageUsecase.ResizeImages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
