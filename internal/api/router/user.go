package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler/middleware"
	"github.com/big-dust/DreamBridge/internal/api/handler/user"
	"github.com/gin-gonic/gin"
)

func useUser(r *gin.RouterGroup) {
	u := r.Group("/user")
	u.Use(middleware.Authorization)
	u.POST("/info", user.SetInfo)
}
