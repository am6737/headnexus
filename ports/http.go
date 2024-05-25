package ports

import (
	"context"
	"errors"
	"fmt"
	"github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app"
	"github.com/am6737/headnexus/config"
	pkghttp "github.com/am6737/headnexus/pkg/http"
	"github.com/am6737/headnexus/pkg/http/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	oapimiddleware "github.com/oapi-codegen/gin-middleware"
	"log"
	"net/http"
	"os"
)

var _ v1.ServerInterface = &HttpHandler{}

//var _ interfaces.Runnable = &HttpHandler{}

func NewHttpHandler(c *config.Config, app *app.Application) *HttpHandler {
	return &HttpHandler{
		app: app,
		c:   c,
	}
}

type HttpHandler struct {
	c   *config.Config
	app *app.Application
}

func (h *HttpHandler) Start(ctx context.Context) error {
	// This is how you set up a basic gin router
	r := gin.Default()
	r.Use(middleware.RecoveryMiddleware(), middleware.CORSMiddleware())

	swagger, err := v1.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	validatorOptions := &oapimiddleware.Options{
		ErrorHandler: middleware.HandleOpenAPIError,
	}
	validatorOptions.Options.AuthenticationFunc = func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return middleware.HandleOpenApiAuthentication(ctx, h.app.JwtConfig, input)
		//return nil
	}

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

	v1.RegisterHandlers(r, h)

	//r.Use(middleware.OapiRequestValidator(swagger))
	server := pkghttp.NewServer(r)
	server.Addr = h.c.Http.Addr

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
