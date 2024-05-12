package middleware

import (
	"fmt"
	pkghttp "github.com/am6737/headnexus/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recover error: %v\n", err)
				// 打印堆栈信息
				debug.PrintStack()
				pkghttp.NewResponse(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
