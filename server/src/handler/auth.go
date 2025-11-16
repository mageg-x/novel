package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mageg-x/novel/src/util"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Authorization头格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		// 解析JWT令牌
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userType", claims.Type)

		c.Next()
	}
}

// AuthorMiddleware 作者身份中间件
func AuthorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户类型
		userType, exists := c.Get("userType")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供用户类型"})
			c.Abort()
			return
		}

		// 检查用户是否为作者
		if userType != "author" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要作者身份"})
			c.Abort()
			return
		}

		c.Next()
	}
}