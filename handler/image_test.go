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

func Test_imageHandler_ProcessImage(t *testing.T) {

	t.Run("success process image", func(t *testing.T) {
		// generate body
		bodyWithFile := new(bytes.Buffer)
		mwWithFile := multipart.NewWriter(bodyWithFile)

		file, err := os.Open("./image/flower.png")
		if err != nil {
			t.Fatal(err)
		}

		ws, err := mwWithFile.CreateFormFile("file", "flower.png")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(ws, file); err != nil {
			t.Fatal(err)
		}

		resizeHeight, _ := mwWithFile.CreateFormField("resizeHeight")
		resizeHeight.Write([]byte("50"))

		resizeWidth, _ := mwWithFile.CreateFormField("resizeWidth")
		resizeWidth.Write([]byte("50"))

		mwWithFile.Close()
		// end generate body

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqWithFile, err := http.NewRequest(http.MethodPost, "/", bodyWithFile)
		if err != nil {
			t.Fatal(err)
		}

		reqWithFile.Header.Add("Content-Type", mwWithFile.FormDataContentType())

		c.Request = reqWithFile

		h := &imageHandler{}
		h.ProcessImage(c)

		result := w.Result()
		assert.Equal(t, result.StatusCode, http.StatusCreated)
	})

	t.Run("error resize request not valid", func(t *testing.T) {
		// generate body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		resizeHeight, _ := mw.CreateFormField("resizeHeight")
		resizeHeight.Write([]byte("-50"))

		resizeWidth, _ := mw.CreateFormField("resizeWidth")
		resizeWidth.Write([]byte("50"))

		file, err := os.Open("./image/flower.png")
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

		mw.Close()
		// end generate body

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqWithFile, err := http.NewRequest(http.MethodPost, "/", body)
		if err != nil {
			t.Fatal(err)
		}

		reqWithFile.Header.Add("Content-Type", mw.FormDataContentType())

		c.Request = reqWithFile

		h := &imageHandler{}
		h.ProcessImage(c)

		result := w.Result()
		assert.Equal(t, result.StatusCode, http.StatusBadRequest)
	})

	t.Run("error image request not provided", func(t *testing.T) {
		// generate body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		mw.Close()
		// end generate body

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqWithFile, err := http.NewRequest(http.MethodPost, "/", body)
		if err != nil {
			t.Fatal(err)
		}

		reqWithFile.Header.Add("Content-Type", mw.FormDataContentType())

		c.Request = reqWithFile

		h := &imageHandler{}
		h.ProcessImage(c)

		result := w.Result()
		assert.Equal(t, result.StatusCode, http.StatusBadRequest)
	})

	t.Run("error image is not png", func(t *testing.T) {
		// generate body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		resizeHeight, _ := mw.CreateFormField("resizeHeight")
		resizeHeight.Write([]byte("50"))

		resizeWidth, _ := mw.CreateFormField("resizeWidth")
		resizeWidth.Write([]byte("50"))

		file, err := os.Open("./image/cat.jpg")
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

		mw.Close()
		// end generate body

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqWithFile, err := http.NewRequest(http.MethodPost, "/", body)
		if err != nil {
			t.Fatal(err)
		}

		reqWithFile.Header.Add("Content-Type", mw.FormDataContentType())

		c.Request = reqWithFile

		h := &imageHandler{}
		h.ProcessImage(c)

		result := w.Result()
		assert.Equal(t, result.StatusCode, http.StatusBadRequest)
	})

}
