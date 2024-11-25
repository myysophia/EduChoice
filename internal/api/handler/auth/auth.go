package auth

import (
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/internal/api/service/auth"
	"github.com/big-dust/DreamBridge/internal/api/types"
	"github.com/big-dust/DreamBridge/pkg/email"
	"github.com/gin-gonic/gin"
)

// @Summary		注册
// @Description	邮箱验证码注册
// @Accept			json;multipart/form-data
// @Produce		json
// @Param			email		formData		string					true	"email"
// @Param			username	formData		string					true	"username"
// @Param			password	formData		string					true	"password"
// @Success		200			{object}	response.OkMsgResp		"注册成功"
// @Failure		400			{object}	response.FailMsgResp	"注册失败"
// @Router			/api/v1/auth/resgister [post]
func Register(c *gin.Context) {
	req := &types.RegisterReq{}
	if err := c.ShouldBind(req); err != nil {
		response.FailMsg(c, "注册失败: "+err.Error())
		return
	}
	if !auth.OkEmailCode(req.Email, req.Code) {
		response.FailMsg(c, "邮箱验证失败")
		return
	}
	if err := auth.Register(req.Username, req.Email, req.Password); err != nil {
		response.FailMsg(c, "注册失败: "+err.Error())
		return
	}
	response.OkMsg(c, "注册成功")
}

// @Summary		登录
// @Description	邮箱，密码登录
// @Accept			json;multipart/form-data
// @Produce		json
// @Param			email		formData		string									true	"email"
// @Param			password	formData		string									true	"password"
// @Success		200			{object}	response.OkMsgDataResp[types.LoginResp]	"登录成功，返回token"
// @Failure		400			{object}	response.FailMsgResp					"登录失败"
// @Router			/api/v1/auth/login [post]
func Login(c *gin.Context) {
	req := &types.LoginReq{}
	if err := c.ShouldBind(req); err != nil {
		response.FailMsg(c, "登录失败: "+err.Error())
		return
	}
	token, err := auth.LoginGetToken(req.Email, req.Password)
	if err != nil {
		response.FailMsg(c, "登录失败: "+err.Error())
		return
	}
	response.OkMsgData(c, "登录成功", &types.LoginResp{Token: token})
}

// @Summary	发送邮箱验证码
// @Description
// @Accept		json;multipart/form-data
// @Produce	json
// @Param		email	formData		string					true	"email"
// @Success	200		{object}	response.OkMsgResp		"发送成功"
// @Failure	400		{object}	response.FailMsgResp	"发送失败"
// @Router		/api/v1/auth/email_code [get]
func SendEmailVerificationCode(c *gin.Context) {
	req := &types.SendEmailVerificationCodeReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.FailMsg(c, "发送失败: "+err.Error())
		return
	}
	err := email.Send(req.Email)
	if err != nil {
		response.FailMsg(c, "发送失败: "+err.Error())
		return
	}
	response.OkMsg(c, "发送成功")
}
