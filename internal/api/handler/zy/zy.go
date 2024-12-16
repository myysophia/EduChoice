package zy

import (
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/internal/api/service/zy"
)

// GetRecommend 获取专业推荐
func GetRecommend(c *gin.Context) {
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "用户未登录")
		return
	}

	// 获取用户信息
	userInfo, err := zy.GetUserInfo(uid)
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	// 检查必要信息
	if !userInfo.IsComplete() {
		response.Error(c, "请先完善个人信息")
		return
	}

	// 获取推荐结果
	result, err := zy.GetRecommend(userInfo)
	if err != nil {
		response.Error(c, "获取推荐失败: " + err.Error())
		return
	}

	response.Success(c, result)
}
