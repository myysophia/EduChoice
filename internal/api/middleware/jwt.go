package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/pkg/jwt"
)

// JWT 认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 添加日志
		fmt.Printf("JWT中间件收到token: %s\n", token)

		claims, err := jwt.ParseToken(token)
		if err != nil {
			fmt.Printf("解析token失败: %v\n", err)
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 添加日志
		fmt.Printf("解析token成功, 用户ID: %d\n", claims.ID)

		// 将用户ID存入上下文
		c.Set("uid", claims.ID)
		c.Next()
	}
} 