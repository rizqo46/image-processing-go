package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func Test_imageHandler_ProcessImage(t *testing.T) {
	generateBodyForTest := func(fileImagePath string, resizeHeight, resizeWidth int) *http.Request {
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		if fileImagePath != "" {
			file, err := os.Open(fileImagePath)
			if err != nil {
				t.Fatal(err)
			}

			ws, err := mw.CreateFormFile("file", "flower.png")
			if err != nil {
				t.Fatal(err)
			}

			if _, err := io.Copy(ws, file); err != nil {
				t.Fatal(err)
			}
		}

		resizeH, _ := mw.CreateFormField("resizeHeight")
		resizeH.Write([]byte(fmt.Sprintf("%d", resizeHeight)))

		resizeW, _ := mw.CreateFormField("resizeWidth")
		resizeW.Write([]byte(fmt.Sprintf("%d", resizeWidth)))

		mw.Close()

		req, err := http.NewRequest(http.MethodPost, "/", body)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Content-Type", mw.FormDataContentType())

		return req
	}

	type field struct {
		fileImagePath string
		resizeHeight  int
		resizeWidth   int
	}
	var tests = []struct {
		name           string
		field          field
		wantStatusCode int
	}{
		{
			name:           "success process image",
			field:          field{"./image/flower.png", 50, 50},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "error resize request not valid",
			field:          field{"./image/flower.png", -50, 50},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "error image request not provided",
			field:          field{"", 50, 50},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "error image is not png",
			field:          field{"./image/cat.jpg", -50, 50},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := generateBodyForTest(
				tt.field.fileImagePath,
				tt.field.resizeHeight,
				tt.field.resizeWidth,
			)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h := &imageHandler{}
			h.ProcessImage(c)

			result := w.Result()
			assert.Equal(t, result.StatusCode, tt.wantStatusCode)
		})
	}
}

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
			ff.Write([]byte(formData.value))
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
	req, _ := http.NewRequest(http.MethodPost, "/png-to-jpeg", body)
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
					value:      "./image/flower.png",
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
					value:      "./image/cat.jpg",
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
