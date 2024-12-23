package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler/auth"
	"github.com/big-dust/DreamBridge/internal/api/handler/major"
	"github.com/big-dust/DreamBridge/internal/api/handler/user"
	"github.com/big-dust/DreamBridge/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func Load(e *gin.Engine) {
	// 使用CORS中间件
	e.Use(middleware.Cors())

	// API路由组
	api := e.Group("/api/v1")
	{
		// 认证相关
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/email_code", auth.SendEmailVerificationCode)

		// 需要认证的路由
		authed := api.Group("")
		authed.Use(middleware.JWT())
		{
			// 用户相关
			authed.POST("/user/info", user.SetInfo)
			authed.GET("/user/info", user.GetInfo)

			// 专业推荐相关
			// 专业推荐相关路由
			zy := api.Group("/zy")
			zy.Use(middleware.JWT())
			{
				zy.GET("/recommend", major.GetMajorRecommend) // 使用新的推荐处理器
			}
		}
	}
}
