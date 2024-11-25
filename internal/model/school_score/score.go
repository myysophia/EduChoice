package school_score

import (
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

// 定义 Scores 表格的模型
type Score struct {
	SchoolID   int
	Location   int    //地区
	Year       int    //年份
	TypeId     int    //科类
	BatchName  string //批次
	Tag        string //类型
	SgName     string //专业组
	Lowest     int
	LowestRank int
}

type HistoryScore struct {
	SchoolID   int
	Year       int
	Lowest     int
	LowestRank int
}

func SchoolIdsIn(startScore int, endScore int, typeID int) ([]int, error) {
	var ids []int
	if err := common.DB.Table("scores").
		Select("school_id").
		Group("school_id").
		Find(&ids, "type_id = ? and tag = \"普通类\" and year > 2020 and year <= 2023 and lowest != 0 and lowest >= ? and lowest < ?", typeID, startScore, endScore).
		Having("avg(lowest) >= ? and avg(lowest) < ?").Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func FindHistoryScore(sid int, typeId int) ([]*HistoryScore, error) {
	var hs []*HistoryScore
	if err := common.DB.Table("scores").
		Select("school_id,year,min(lowest) as lowest,min(lowest_rank) as lowest_rank").
		Group("school_id,year").
		Find(&hs, "school_id = ? and type_id = ? and tag = \"普通类\"", sid, typeId).Error; err != nil {
		return nil, err
	}
	return hs, nil
}
