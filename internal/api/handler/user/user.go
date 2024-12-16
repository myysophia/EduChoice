package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/response"
	userService "github.com/big-dust/DreamBridge/internal/api/service/user"
	userModel "github.com/big-dust/DreamBridge/internal/model/user"
)

// SetInfo 设置用户信息
func SetInfo(c *gin.Context) {
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "用户未登录")
		return
	}

	var req userService.UserSetInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("解析请求参数失败: %v\n原始数据: %v\n", err, c.Request.Body)
		response.Error(c, "参数错误")
		return
	}

	fmt.Printf("收到设置用户信息请求: uid=%d, data=%+v\n", uid, req)

	if err := userService.SetInfo(uid, &req); err != nil {
		fmt.Printf("设置用户信息失败: %v\n", err)
		response.Error(c, "设置失败: "+err.Error())
		return
	}

	fmt.Printf("设置用户信息成功: uid=%d\n", uid)
	response.Success(c, nil)
}

// GetInfo 获取用户信息
func GetInfo(c *gin.Context) {
	// 从上下文获取用户ID
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "未登录")
		return
	}

	// 添加日志
	fmt.Printf("获取用户信息请求, uid: %d\n", uid)

	// 查询用户信息
	u, err := userModel.FindOne(uid)
	if err != nil {
		fmt.Printf("查询用户信息失败: %v\n", err)
		response.Error(c, "获取失败: "+err.Error())
		return
	}

	// 添加日志
	fmt.Printf("查询到用户信息: %+v\n", u)

	response.Success(c, u)
}

// GetRecommend 获取用户推荐信息
func GetRecommend(c *gin.Context) {
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "用户未登录")
		return
	}

	info, err := userService.GetRecommend(uid)
	if err != nil {
		response.Error(c, "获取推荐失败: "+err.Error())
		return
	}

	response.Success(c, info)
}
