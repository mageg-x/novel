package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mageg-x/novel/src/model"
)

var (
	// 密钥，用于签名JWT（在生产环境中应该从环境变量或配置文件中获取）
	jwtSecret = []byte("your-secret-key-here")
	// 过期时间
	expireDuration = 24 * time.Hour
)

// Claims 定义JWT的声明

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *model.User) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(expireDuration)
	
	// 创建声明
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Type:     user.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "novel-system",
			Subject:   "user-auth",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名令牌
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌并返回声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}