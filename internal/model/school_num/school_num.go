package school_num

import "github.com/big-dust/DreamBridge/internal/pkg/common"

type SchoolNum struct {
	ID       int
	SchoolID int
	Year     int
	TypeID   string
	Number   int
}

type EnrollmentNum struct {
	SchoolID int
	Year     int
	Number   int
}

func SchoolNumCreate(sn []*SchoolNum) error {
	return common.DB.Create(sn).Error
}

func NumList() ([]int, error) {
	var schoolIdList []int
	if err := common.DB.Model(&SchoolNum{}).Select("school_id").Group("school_id").Find(&schoolIdList).Error; err != nil {
		return nil, err
	}
	return schoolIdList, nil
}

func FindHistoryEnrollmentNum(sid int) ([]*EnrollmentNum, error) {
	var ens []*EnrollmentNum
	if err := common.DB.Table("school_nums").
		Group("school_id,year").
		Select("school_id,year,sum(number) as number").
		Find(&ens, "school_id = ?", sid).Error; err != nil {
		return nil, err
	}
	return ens, nil
}
