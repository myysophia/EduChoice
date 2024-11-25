package router

import (
	"github.com/big-dust/DreamBridge/internal/api/handler/auth"
	"github.com/gin-gonic/gin"
)

func useAuth(r *gin.RouterGroup) {
	ag := r.Group("/auth")
	ag.POST("/login", auth.Login)
	ag.POST("/register", auth.Register)
	ag.GET("/email_code", auth.SendEmailVerificationCode)
}
