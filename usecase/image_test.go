package usecase

import (
	"context"
	"io"
	"os"
	"testing"

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
