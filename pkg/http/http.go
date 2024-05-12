package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, msg string, data interface{}) {
	NewResponse(c, http.StatusOK, msg, data)
}

func FailedResponse(c *gin.Context, msg string) {
	//NewResponse(c, http.StatusBadRequest, msg, nil)
	c.JSON(http.StatusOK, Response{
		Code: http.StatusBadRequest,
		Msg:  msg,
		Data: nil,
	})
}

func NewResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// NewServer returns a new server with sane defaults.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadHeaderTimeout: 32 * time.Second,
	}
}

type HttpServer struct {
}
