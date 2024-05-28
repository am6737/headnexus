package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/am6737/headnexus/pkg/code"
	pkghttp "github.com/am6737/headnexus/pkg/http"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	omiddleware "github.com/oapi-codegen/gin-middleware"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")
)

func HandleOpenAPIError(c *gin.Context, message string, statusCode int) {
	if strings.Contains(message, "security requirements failed: authorization failed") {
		statusCode = http.StatusUnauthorized
		message = code.Unauthorized.Message()
	}
	if strings.Contains(message, "request body has an error: doesn't match schema") {
		index := strings.Index(message, "Error at")
		if index != -1 {
			message = strings.TrimSpace(message[index+len("Error at "):])
		} else {
			//pkghttp.NewResponse(c, statusCode, message, nil)
			//return
			//message = code.InvalidParameter.Message()
		}
	}
	pkghttp.NewResponse(c, statusCode, message, nil)
}

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	queryToken := req.URL.Query().Get("token")
	if queryToken != "" {
		return queryToken, nil
	}
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func NewAuthenticator(jwt2 *pkgjwt.JWTConfig) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, jwt2, input)
	}
}

func Authenticate(ctx context.Context, jwt2 *pkgjwt.JWTConfig, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" && input.SecuritySchemeName != "bearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return err
	}

	token, err := jwt2.ParseTokenWithKey(jws)
	if err != nil {
		return err
	}
	if token.Valid() != nil {
		return token.Valid()
	}

	// 将用户信息存储到上下文中
	claims := token.(jwt.MapClaims)

	//claims, err := authService.Access(ctx, &authv1.AccessRequest{Token: jws})
	//if err != nil {
	//	return code.Unauthorized
	//}

	gctx := omiddleware.GetGinContext(ctx)
	gctx.Set("user_id", claims["user_id"])
	fmt.Println(gctx.Value("user_id"), gctx.Value("user_id"))
	return nil
}

func HandleOpenApiAuthentication(ctx context.Context, jwt2 *pkgjwt.JWTConfig, input *openapi3filter.AuthenticationInput) error {
	if err := Authenticate(ctx, jwt2, input); err != nil {
		//gx := omiddleware.GetGinContext(ctx)
		//gx.JSON(http.StatusUnauthorized, gin.H{
		//	"code": 401,
		//	"msg":  err.Error(),
		//})
		//gx.Abort()
		return input.NewError(err)
	}

	return nil
}
