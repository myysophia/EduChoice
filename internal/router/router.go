package router

import (
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/handler/auth"
	"github.com/big-dust/DreamBridge/internal/api/handler/user"
	"github.com/big-dust/DreamBridge/internal/api/handler/major"
	"github.com/big-dust/DreamBridge/internal/api/middleware"
)

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// 注册认证相关路由
		registerAuthRoutes(api)
		// 注册用户相关路由
		registerUserRoutes(api)
		// 注册专业相关路由
		registerMajorRoutes(api)
	}
}

func registerMajorRoutes(r *gin.RouterGroup) {
	zy := r.Group("/zy")
	{
		zy.GET("/recommend", middleware.JWT(), major.GetRecommend)
	}
}

// ... 其他路由注册函数 