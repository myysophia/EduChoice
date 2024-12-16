package zy

import (
	"github.com/big-dust/DreamBridge/internal/api/types"
	"github.com/big-dust/DreamBridge/internal/model/user"
)

func GetUserInfo(uid int) (*user.User, error) {
	return user.FindOne(uid)
}

func GetRecommend(u *user.User) (*types.ZYMockResp, error) {
	// 这里实现推荐逻辑
	// 暂时返回模拟数据
	return &types.ZYMockResp{
		ChongSchools: []types.School{
			{
				ID:   1,
				Name: "清华大学",
				HistoryInfos: map[string]types.HistoryInfo{
					"2023": {
						LowestScore:   680,
						LowestRank:    100,
						EnrollmentNum: 3000,
					},
				},
				Parts: map[string][]types.Major{
					"理工": {
						{
							ID:     1,
							Name:   "计算机科学与技术",
							Rate:   0.95,
							Weight: 1.0,
						},
					},
				},
			},
		},
		WenSchools:   []types.School{},
		BaoSchools:   []types.School{},
	}, nil
}
