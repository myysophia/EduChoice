package router

import (
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/handler/auth"
	"github.com/big-dust/DreamBridge/internal/api/handler/user"
	"github.com/big-dust/DreamBridge/internal/api/handler/major"
	"github.com/big-dust/DreamBridge/internal/api/middleware"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Cors())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", auth.Register)
			auth.POST("/login", auth.Login)
			// ...其他认证路由
		}

		// 用户相关路由
		user := api.Group("/user")
		user.Use(middleware.JWT())
		{
			user.GET("/info", user.GetUserInfo)
			user.PUT("/info", user.UpdateUserInfo)
			// ...其他用户路由
		}

		// 专业推荐相关路由
		zy := api.Group("/zy")
		zy.Use(middleware.JWT())
		{
			zy.GET("/recommend", major.GetMajorRecommend) // 使用新的推荐处理器
		}
	}

	return r
} 