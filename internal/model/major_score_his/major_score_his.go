package major_score_his

import (
	"database/sql"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"time"
)

// MajorScoreHis 定义专业历史分数表的模型
type MajorScoreHis struct {
	ID                string `gorm:"primary_key"` // 如 gkspecialscore2023295809
	SchoolID          int    `gorm:"index"`       // 学校ID
	SpecialID         int    `gorm:"index"`       // 专业ID
	SpeID             int    // 特殊专业ID
	Year              int    `gorm:"index"` // 年份
	SpName            string // 专业名称
	SpInfo            string // 专业信息
	Info              string // 额外信息（如：钱学森班本研一体）
	LocalProvinceName string // 省份
	LocalTypeName     string // 科类
	LocalBatchName    string // 批次
	Level2Name        string // 二级学科名称
	Level3Name        string // 三级学科名称
	Average           int    // 平均分
	Max               int    // 最高分
	Min               int    // 最低分
	MinSection        int    // 最低位次
	Proscore          int    // 投档分
	DoubleHigh        int    // 是否双高
	IsTop             int    // 是否重点
	IsScoreRange      int    // 分数范围标识
	MinRange          string // 最低分范围
	MinRankRange      string // 最低排名范围
	Remark            string // 备注
	ZslxName          string // 招生类型名称
	DualClassName     string // 双学位类名称
	FirstKm           int    // 首选科目
	SgFxk             string // 选考科目要求
	SgSxk             string // 首选科目要求
	SgInfo            string // 选考科目信息
	SgName            string // 选考科目组名称
	SgType            int    // 选考科目类型
	SpFxk             string // 专业选考科目要求
	SpSxk             string // 专业首选科目要求
	SpType            int    // 专业类型
	Single            string // 单科要求
	SpecialGroup      int    // 专业组
}

// CreateMajorScoresHis 创建历史分数记录
func CreateMajorScoresHis(records []*MajorScoreHis) error {
	tx := common.DB.Begin()
	for _, record := range records {
		if err := tx.Create(record).Error; err != nil {
			tx.Rollback()
			if common.ErrMysqlDuplicate.Is(err) {
				return nil
			}
			common.LOG.Error("CreateMajorScoresHis: " + err.Error())
			return err
		}
	}
	tx.Commit()
	return nil
}

// MustCreateMajorScoresHis 确保创建成功的历史分数记录
func MustCreateMajorScoresHis(records []*MajorScoreHis) {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		nilChan := make(chan error, 1)
		go func() {
			err := CreateMajorScoresHis(records)
			if err != nil {
				errChan <- err
				return
			}
			nilChan <- nil
		}()
		ticker := time.NewTicker(15 * time.Second)
		select {
		case <-nilChan:
			return
		case <-ticker.C:
			select {
			case err := <-errChan:
				common.LOG.Error(fmt.Sprintf("MustCreateMajorScoresHis: Time Out 15s TryCount:%d SchoolID:%v Error: %s",
					tryCount, records[0].SchoolID, err.Error()))
			default:
				common.LOG.Error(fmt.Sprintf("MustCreateMajorScoresHis: Time Out 15s TryCount:%d SchoolID:%v",
					tryCount, records[0].SchoolID))
			}
		}
	}
}

// FindScoreAvg 计算某专业在特定科类的平均分
func FindScoreAvg(specialID int, localTypeName string) (*sql.NullFloat64, error) {
	avgScore := &sql.NullFloat64{}
	if err := common.DB.Table("major_scores_his").
		Select("avg(average) as avg_score").
		Find(avgScore, "special_id = ? and local_type_name = ? and year > 2020 and year <= 2023 and average != 0",
			specialID, localTypeName).Error; err != nil {
		return nil, err
	}
	return avgScore, nil
}

// FindByLocalTypeNameIn 根据科类查询专业ID列表
func FindByLocalTypeNameIn(majorIds []int, localTypeNames ...string) ([]int, error) {
	var mids []int
	if err := common.DB.Table("major_scores_his").
		Distinct("special_id").
		Find(&mids, "special_id in (?) and local_type_name in (?)", majorIds, localTypeNames).Error; err != nil {
		return nil, err
	}
	return mids, nil
}

// FindBySpName 根据专业名称查询历史分数记录
func FindBySpName(spName string) ([]*MajorScoreHis, error) {
	var records []*MajorScoreHis
	if err := common.DB.Find(&records, "sp_name = ?", spName).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// FindBySchoolID 根据学校ID查询历史分数记录
func FindBySchoolID(schoolID int) ([]*MajorScoreHis, error) {
	var records []*MajorScoreHis
	if err := common.DB.Find(&records, "school_id = ?", schoolID).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// FindBySpecialID 根据专业ID查询历史分数记录
func FindBySpecialID(specialID int) ([]*MajorScoreHis, error) {
	var records []*MajorScoreHis
	if err := common.DB.Find(&records, "special_id = ?", specialID).Error; err != nil {
		return nil, err
	}
	return records, nil
}
