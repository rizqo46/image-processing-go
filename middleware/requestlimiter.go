package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestBodyLimiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(2<<20))
}
