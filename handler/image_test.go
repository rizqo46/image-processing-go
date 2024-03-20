package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

type formData struct {
	isTypeFile bool
	label      string
	value      string
}

func httpRequestWithFormData(t *testing.T, httpMethod string, path string, formsData ...formData) *http.Request {
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	for _, formData := range formsData {
		if !formData.isTypeFile {
			ff, _ := mw.CreateFormField(formData.label)
			_, _ = ff.Write([]byte(formData.value))
			continue
		}

		file, err := os.Open(formData.value)
		if err != nil {
			t.Fatal(err)
		}

		ws, err := mw.CreateFormFile(formData.label, file.Name())
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(ws, file); err != nil {
			t.Fatal(err)
		}
	}

	mw.Close()
	req, _ := http.NewRequest(httpMethod, path, body)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	return req
}

func Test_imageHandler_PngToJpeg(t *testing.T) {
	router := gin.Default()
	SetupImageRoute(router)

	var tests = []struct {
		name           string
		field          []formData
		wantStatusCode int
	}{
		{
			name: "success process image",
			field: []formData{
				{
					isTypeFile: true,
					label:      "files[]",
					value:      ".././imagetest/flower.png",
				},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "error image request not provided",
			field:          []formData{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "error image is not png",
			field: []formData{
				{
					isTypeFile: true,
					label:      "files[]",
					value:      ".././imagetest/cat.jpg",
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httpRequestWithFormData(t, http.MethodPost, "/png-to-jpeg", tt.field...)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}

func Test_imageHandler_CompressImages(t *testing.T) {
	router := gin.Default()
	SetupImageRoute(router)

	var tests = []struct {
		name           string
		field          []formData
		wantStatusCode int
	}{
		{
			name: "success process image",
			field: []formData{
				{
					isTypeFile: true,
					label:      "files[]",
					value:      ".././imagetest/flower.png",
				},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "error image request not provided",
			field:          []formData{},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httpRequestWithFormData(t, http.MethodPost, "/compress", tt.field...)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}

func Test_imageHandler_Resize(t *testing.T) {
	router := gin.Default()
	SetupImageRoute(router)

	var tests = []struct {
		name           string
		field          []formData
		wantStatusCode int
	}{
		{
			name: "success process image",
			field: []formData{
				{
					isTypeFile: true,
					label:      "files[]",
					value:      ".././imagetest/flower.png",
				},
				{
					isTypeFile: false,
					label:      "width[]",
					value:      "70",
				},
				{
					isTypeFile: false,
					label:      "height[]",
					value:      "70",
				},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "error image request not provided",
			field:          []formData{},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httpRequestWithFormData(t, http.MethodPost, "/resize", tt.field...)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}
