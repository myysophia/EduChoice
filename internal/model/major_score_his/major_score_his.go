package major_score_his

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

// MajorScoreHis 定义专业历史分数表的模型
type MajorScoreHis struct {
	ID                string `json:"id"`
	SchoolID          int    `json:"school_id"`
	SpecialID         int    `json:"special_id"`
	SpeID             int    `json:"spe_id"`
	Year              int    `json:"year"`
	SpName            string `json:"sp_name"`
	Spname            string `json:"spname"`
	Info              string `json:"info"`
	LocalProvinceName string `json:"local_province_name"`
	LocalTypeName     string `json:"local_type_name"`
	LocalBatchName    string `json:"local_batch_name"`
	Level2Name        string `json:"level2_name"`
	Level3Name        string `json:"level3_name"`
	Average           int    `json:"average"`
	Max               int    `json:"max"`
	Min               int    `json:"min"`
	MinSection        string `json:"min_section"`
	Proscore          int    `json:"proscore"`
	IsTop             int    `json:"is_top"`
	IsScoreRange      int    `json:"is_score_range"`
	MinRange          string `json:"min_range"`
	MinRankRange      string `json:"min_rank_range"`
	Remark            string `json:"remark"`
}

// CreateMajorScoresHis 创建历史分数记录
func CreateMajorScoresHis(records []*MajorScoreHis) error {
	// 按学校ID分组，减少索引竞争
	schoolGroups := make(map[int][]*MajorScoreHis)
	for _, record := range records {
		schoolGroups[record.SchoolID] = append(schoolGroups[record.SchoolID], record)
	}

	// 每个学校的数据在单独的事务中处理
	for _, schoolRecords := range schoolGroups {
		if err := createSchoolRecords(schoolRecords); err != nil {
			return err
		}
	}
	return nil
}

func createSchoolRecords(records []*MajorScoreHis) error {
	// 按主键ID去重
	uniqueRecords := make(map[string]*MajorScoreHis)
	for _, record := range records {
		uniqueRecords[record.ID] = record
	}

	// 转换回切片
	deduplicatedRecords := make([]*MajorScoreHis, 0, len(uniqueRecords))
	for _, record := range uniqueRecords {
		 deduplicatedRecords = append(deduplicatedRecords, record)
	}

	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.CreateInBatches(deduplicatedRecords, 100).Error; err != nil {
		tx.Rollback()
		if common.ErrMysqlDuplicate.Is(err) {
			return nil
		}
		return err
	}

	return tx.Commit().Error
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
