package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
	out    *bufio.Writer
}

func Cache(c *gin.Context) {
	bw := &bodyLogWriter{
		Buffer: bytes.NewBuffer(nil),
		ResponseWriter: c.Writer,
	}
	c.Writer = bw
	c.Next()
	if c.Writer.Status() == http.StatusOK {
		fmt.Println(bw.Buffer.String(), c.Writer.Written(), "SD")
	}
}
