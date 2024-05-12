package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTConfig struct {
	SecretKey      string
	ExpiryDuration time.Duration
}

func (c *JWTConfig) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(c.ExpiryDuration).Unix()

	signedToken, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseTokenWithKey 解析带有自定义秘钥的 JWT
func (c *JWTConfig) ParseTokenWithKey(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateTokenWithExpiryAndKey(claims jwt.Claims, expiry time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(expiry).Unix()

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseTokenWithKey 解析带有自定义秘钥的 JWT
func ParseTokenWithKey(tokenString string, key []byte) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
