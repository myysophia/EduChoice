package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/model/auth"
	"github.com/big-dust/DreamBridge/internal/api/response"
	authService "github.com/big-dust/DreamBridge/internal/api/service/auth"
	"fmt"
)

// Register 处理注册请求
func Register(c *gin.Context) {
	var req auth.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误")
		return
	}

	if !authService.OkEmailCode(req.Email, req.Code) {
		response.Error(c, "验证码错误")
		return
	}

	if err := authService.Register(req.Username, req.Email, req.Password); err != nil {
		response.Error(c, "注册失败: " + err.Error())
		return
	}

	response.Success(c, nil)
}

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("登录请求参数解析失败: %v, 原始数据: %v\n", err, c.Request.Body)
		response.Error(c, "账号和密码不能为空")
		return
	}

	fmt.Printf("收到登录请求: account=%s, password长度=%d\n", req.Account, len(req.Password))

	token, err := authService.LoginGetToken(req.Account, req.Password)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		response.Error(c, err.Error())
		return
	}

	fmt.Printf("登录成功，生成token: %s\n", token)

	response.Success(c, gin.H{
		"token": token,
	})
}

// SendEmailVerificationCode 发送邮箱验证码
func SendEmailVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "邮箱格式错误")
		return
	}

	if err := authService.SendEmailCode(req.Email); err != nil {
		response.Error(c, "发送验证码失败")
		return
	}

	response.Success(c, nil)
}
