package major

import "github.com/big-dust/DreamBridge/internal/pkg/common"

// HistoryInfo 历年录取信息
type HistoryInfo struct {
	LowestScore   int `json:"lowestScore"`   // 最低分
	LowestRank    int `json:"lowestRank"`    // 最低位次
	EnrollmentNum int `json:"enrollmentNum"` // 录取人数
}

// Major 专业信息
type Major struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Rate   float64 `json:"rate"`   // 匹配度
	Weight float64 `json:"weight"` // 权重
}

func FindBySchoolID(schoolID int) ([]*Major, error) {
	var majors []*Major
	if err := common.DB.Distinct("*").Find(&majors, "school_id = ?", schoolID).Error; err != nil {
		return nil, err
	}
	return majors, nil
}

func FindIDListBySchoolID(schoolID int) ([]int, error) {
	var majorIDs []int
	if err := common.DB.Table("majors").Distinct("id").Find(&majorIDs, "school_id = ?", schoolID).Error; err != nil {
		return nil, err
	}
	return majorIDs, nil
}
