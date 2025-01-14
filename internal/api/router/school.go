package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler"
	"github.com/big-dust/DreamBridge/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitSchoolRouter(r *gin.RouterGroup) {
	schoolGroup := r.Group("/school")
	schoolGroup.Use(middleware.JWT())
	{
		schoolGroup.GET("/cards", handler.GetSchoolCards)
	}
}
