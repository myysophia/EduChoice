package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler/middleware"
	"github.com/big-dust/DreamBridge/internal/api/handler/zy"
	"github.com/gin-gonic/gin"
)

func useZY(r *gin.RouterGroup) {
	z := r.Group("/zy")
	z.Use(middleware.Authorization)
	z.GET("/mock", zy.Mock)
}
