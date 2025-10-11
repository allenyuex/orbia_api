package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"orbia_api/biz/infra/config"
)

// Claims JWT 声明结构
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64) (string, int64, error) {
	cfg := config.GlobalConfig.JWT
	expirationTime := time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)
	
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "orbia_api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(cfg.ExpireHours * 3600), nil // 返回秒数
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig.JWT
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证 token 是否有效
func ValidateToken(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	// 检查是否过期
	if claims.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("token expired")
	}

	return claims.UserID, nil
}