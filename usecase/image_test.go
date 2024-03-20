package usecase

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/rizqo46/image-processing-go/constants"
	"github.com/rizqo46/image-processing-go/dto"
)

func TestImageUsecase_ProcessImage(t *testing.T) {
	// Load test file
	file, err := os.Open("./image/flower.png")
	if err != nil {
		t.Errorf("Test ImageUsecase.ProcessImage() failed to open file, error = %v", err)
		return
	}
	defer file.Close()

	imageByte, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("Test ImageUsecase.ProcessImage() failed to open file, error = %v", err)
		return
	}

	type args struct {
		ctx       context.Context
		imageByte []byte
		req       dto.ImageRequest
	}
	tests := []struct {
		name    string
		uc      ImageUsecase
		args    args
		wantErr bool
	}{
		{
			name: "success with resize",
			uc:   ImageUsecase{},
			args: args{
				ctx:       context.Background(),
				imageByte: imageByte,
				req: dto.ImageRequest{
					Resize: dto.Resize{
						Height: 50,
						Width:  50,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error empty byte image",
			uc:   ImageUsecase{},
			args: args{
				ctx:       context.Background(),
				imageByte: []byte{},
				req:       dto.ImageRequest{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.uc.ProcessImage(tt.args.ctx, tt.args.imageByte, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageUsecase.ProcessImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

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
				files:               createMultipartFileheaders("./image/flower.png"),
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
				files:               createMultipartFileheaders("./image/cat.jpg"),
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
