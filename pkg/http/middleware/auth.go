package middleware

//
//import (
//	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strings"
//)
//
//func AuthMiddleware(jwt2 *pkgjwt.JWTConfig) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 从请求头中获取 token
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
//			return
//		}
//
//		// 解析 token
//		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
//		token, err := jwt2.ParseTokenWithKey(tokenString)
//		if err != nil || token.Valid == nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//			return
//		}
//
//		// 将用户信息存储到上下文中
//		claims := token.(jwt.MapClaims)
//		c.Set("user_id", claims["user_id"])
//
//		// 继续处理请求
//		c.Next()
//	}
//}
