package user

import (
	"github.com/big-dust/DreamBridge/internal/model/user"
)

// UserSetInfoReq 设置用户信息请求
type UserSetInfoReq struct {
	Province     *string `json:"province" binding:"omitempty"`
	ExamType     *string `json:"exam_type" binding:"omitempty"`
	SchoolType   *string `json:"school_type" binding:"omitempty"`
	Score        *int    `json:"score" binding:"omitempty"`
	ProvinceRank *int    `json:"province_rank" binding:"omitempty"`
	Physics      *bool   `json:"physics" binding:"omitempty"`
	History      *bool   `json:"history" binding:"omitempty"`
	Chemistry    *bool   `json:"chemistry" binding:"omitempty"`
	Biology      *bool   `json:"biology" binding:"omitempty"`
	Geography    *bool   `json:"geography" binding:"omitempty"`
	Politics     *bool   `json:"politics" binding:"omitempty"`
	Holland      *string `json:"holland" binding:"omitempty"`
	Interests    *string `json:"interests" binding:"omitempty"`
}

// SetInfo 设置用户信息
func SetInfo(uid int, req *UserSetInfoReq) error {
	u := &user.User{
		Province:     req.Province,
		ExamType:     req.ExamType,
		SchoolType:   req.SchoolType,
		Score:        req.Score,
		ProvinceRank: req.ProvinceRank,
		Physics:      req.Physics,
		History:      req.History,
		Chemistry:    req.Chemistry,
		Biology:      req.Biology,
		Geography:    req.Geography,
		Politics:     req.Politics,
		Holland:      req.Holland,
		Interests:    req.Interests,
	}
	return user.UpdateOne(uid, u)
}

// GetInfo 获取用户信息
func GetInfo(uid int) (*user.User, error) {
	return user.FindOne(uid)
}

// GetRecommend 获取用户推荐信息
type RecommendInfo struct {
	HasBasicInfo bool `json:"has_basic_info"` // 是否已填写基本信息
	HasSubjects  bool `json:"has_subjects"`   // 是否已选择科目
	HasHolland   bool `json:"has_holland"`    // 是否已完成霍兰德测试
	HasInterests bool `json:"has_interests"`  // 是否已填写兴趣爱好
}

func GetRecommend(uid int) (*RecommendInfo, error) {
	u, err := user.FindOne(uid)
	if err != nil {
		return nil, err
	}

	info := &RecommendInfo{
		HasBasicInfo: u.Province != nil && u.ExamType != nil && u.SchoolType != nil && 
					 u.Score != nil && u.ProvinceRank != nil,
		HasSubjects:  u.Physics != nil && u.History != nil && 
					 u.Chemistry != nil && u.Biology != nil && 
					 u.Geography != nil && u.Politics != nil,
		HasHolland:   u.Holland != nil,
		HasInterests: u.Interests != nil,
	}

	return info, nil
}
