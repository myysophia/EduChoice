package middleware

import (
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	uid, err := jwt.ParseGetUID(tokenString)
	if err != nil {
		c.Abort()
		response.FailMsg(c, "Authorization: token鉴权失败:"+err.Error())
		return
	}
	c.Set("uid", uid)
}
