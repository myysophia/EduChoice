package major

import (
	"context"
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/internal/model/major"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// GetMajorRecommend 获取专业推荐
func GetMajorRecommend(c *gin.Context) {
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// 1. 获取用户ID和筛选条件
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "用户未登录")
		return
	}

	// 2. 获取筛选参数
	year := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	// 3. 获取用户信息
	u, err := user.FindOne(uid)
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	// 4. 使用推荐算法获取学校列表（添加上下文）
	respChan := make(chan *major.RecommendResp, 1)
	errChan := make(chan error, 1)
	
	go func() {
		resp, err := major.RecommendSchools(u, year, major.DefaultConfig)
		if err != nil {
			errChan <- err
			return
		}
		respChan <- resp
	}()

	// 等待结果或超时
	select {
	case <-ctx.Done():
		response.Error(c, "请求超时，请稍后重试")
		return
	case err := <-errChan:
		log.Printf("推荐学校失败: %v", err)
		response.Error(c, "推荐学校失败")
		return
	case resp := <-respChan:
		response.Success(c, resp)
	}
}
