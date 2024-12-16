package major

import (
	"github.com/gin-gonic/gin"
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/internal/model/major"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"fmt"
)

// GetRecommend 获取专业推荐
func GetRecommend(c *gin.Context) {
	// 获取用户ID
	uid := c.GetInt("uid")
	if uid == 0 {
		response.Error(c, "用户未登录")
		return
	}

	// 获取用户信息
	u, err := user.FindOne(uid)
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	fmt.Printf("用户信息: %+v\n", u)

	// 检查用户信息是否完整
	if !u.IsComplete() {
		response.Error(c, "请先完善个人信息")
		return
	}

	// 模拟数据
	resp := &major.RecommendResp{
		ChongSchools: []major.School{
			{
				ID:   1,
				Name: "清华大学",
				HistoryInfos: map[string]major.HistoryInfo{
					"2023": {
						LowestScore:   680,
						LowestRank:    100,
						EnrollmentNum: 100,
					},
					"2022": {
						LowestScore:   675,
						LowestRank:    120,
						EnrollmentNum: 95,
					},
				},
				Parts: map[string][]major.Major{
					"提前批": {
						{
							ID:     1,
							Name:   "计算机科学与技术",
							Rate:   0.95,
							Weight: 0.9,
						},
						{
							ID:     2,
							Name:   "人工智能",
							Rate:   0.92,
							Weight: 0.85,
						},
					},
				},
			},
		},
		WenSchools: []major.School{
			{
				ID:   2,
				Name: "北京大学",
				HistoryInfos: map[string]major.HistoryInfo{
					"2023": {
						LowestScore:   670,
						LowestRank:    200,
						EnrollmentNum: 120,
					},
					"2022": {
						LowestScore:   665,
						LowestRank:    220,
						EnrollmentNum: 115,
					},
				},
				Parts: map[string][]major.Major{
					"普通批": {
						{
							ID:     3,
							Name:   "软件工程",
							Rate:   0.88,
							Weight: 0.85,
						},
						{
							ID:     4,
							Name:   "数据科学",
							Rate:   0.86,
							Weight: 0.82,
						},
					},
				},
			},
		},
		BaoSchools: []major.School{
			{
				ID:   3,
				Name: "浙江大学",
				HistoryInfos: map[string]major.HistoryInfo{
					"2023": {
						LowestScore:   660,
						LowestRank:    300,
						EnrollmentNum: 150,
					},
					"2022": {
						LowestScore:   655,
						LowestRank:    320,
						EnrollmentNum: 145,
					},
				},
				Parts: map[string][]major.Major{
					"普通批": {
						{
							ID:     5,
							Name:   "电子信息工程",
							Rate:   0.82,
							Weight: 0.8,
						},
						{
							ID:     6,
							Name:   "通信工程",
							Rate:   0.80,
							Weight: 0.78,
						},
					},
				},
			},
		},
	}

	fmt.Printf("返回推荐数据: %+v\n", resp)
	response.Success(c, resp)
} 