package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func InitSchoolRouter(r *gin.RouterGroup) {
	schoolGroup := r.Group("/school")
	{
		schoolGroup.POST("/cards", handler.GetSchoolCards)
	}
}
