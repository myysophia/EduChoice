package major

import "github.com/big-dust/DreamBridge/internal/pkg/common"

// 定义 Majors 表格的模型
type Major struct {
	ID                 int
	Name               string `gorm:"not null"`
	NationalFeature    bool
	Level              string
	DisciplineCategory string
	MajorCategory      string
	LimitYear          string
	SchoolID           int
	SpecialId          string
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
