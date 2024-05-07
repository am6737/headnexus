package ports

import (
	"context"
	"errors"
	"fmt"
	"github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app"
	"github.com/am6737/nexus/api/interfaces"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var _ v1.ServerInterface = &HttpHandler{}

var _ interfaces.Runnable = &HttpHandler{}

func NewHttpHandler(app *app.Application) *HttpHandler {
	return &HttpHandler{
		app: app,
	}
}

type HttpHandler struct {
	app *app.Application
}

func (h *HttpHandler) Start(ctx context.Context) error {
	swagger, err := v1.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic gin router
	r := gin.Default()
	r.Use(RecoveryMiddleware(), CORSMiddleware())

	validatorOptions := &middleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			FailedResponse(c, message)
			//FailedResponse(c, "参数校验失败")
		},
	}
	validatorOptions.Options.AuthenticationFunc = func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// TODO 验证token
		return nil
	}

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

	//r.Use(middleware.OapiRequestValidator(swagger))
	server := newServer(r)
	server.Addr = ":7777"

	v1.RegisterHandlers(r, h)

	serverShutdown := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(serverShutdown)
	}()
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("启动http服务失败：%v", err)
		}
		return nil
	}

	//<-serverShutdown

	fmt.Println("Server exited properly")
	return nil
}

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

// newServer returns a new server with sane defaults.
func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadHeaderTimeout: 32 * time.Second,
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recover error: %v\n", err)
				// 打印堆栈信息
				debug.PrintStack()
				NewResponse(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		}

		ctx.Request.Header.Del("Origin")
		ctx.Next()

	}
}

type HttpServer struct {
}
