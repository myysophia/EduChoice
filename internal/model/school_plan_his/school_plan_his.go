package school_plan_his

import (
	"fmt"
	"time"

	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

type SchoolPlanHis struct {
	ID             int64  `json:"id"`
	SchoolID       int    `json:"school_id"`
	Year           int    `json:"year"`
	SpName         string `json:"sp_name"`
	Spname         string `json:"spname"`
	Num            int    `json:"num"`
	Length         string `json:"length"`
	Tuition        string `json:"tuition"`
	ProvinceName   string `json:"province_name"`
	SpecialGroup   int    `json:"special_group"`
	LocalBatchName string `json:"local_batch_name"`
	LocalTypeName  string `json:"local_type_name"`
}

func CreateSchoolPlanHis(records []*SchoolPlanHis) error {
	// 1. 首先按学校ID分组
	schoolGroups := make(map[int][]*SchoolPlanHis)
	for _, record := range records {
		schoolGroups[record.SchoolID] = append(schoolGroups[record.SchoolID], record)
	}

	// 2. 每个学校的所有数据在一个事务中处理
	for _, schoolRecords := range schoolGroups {
		if err := createSchoolRecords(schoolRecords); err != nil {
			return err
		}
	}
	return nil
}

func createSchoolRecords(records []*SchoolPlanHis) error {
	tx := common.DB.Begin() // 开启事务
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// CreateInBatches 只是为了优化内存使用，将大量记录分批插入
	// 但这些批次都在同一个事务中
	if err := tx.CreateInBatches(records, 100).Error; err != nil {
		if common.ErrMysqlDuplicate.Is(err) {
			return nil
		}
		return err
	}

	// 所有批次插入成功后，一次性提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true
	return nil
}

func MustCreateSchoolPlanHis(records []*SchoolPlanHis) {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		nilChan := make(chan struct{}, 1)

		go func() {
			if err := CreateSchoolPlanHis(records); err != nil {
				errChan <- err
				return
			}
			nilChan <- struct{}{}
		}()

		ticker := time.NewTicker(15 * time.Second)
		select {
		case <-nilChan:
			ticker.Stop()
			return
		case err := <-errChan:
			ticker.Stop()
			common.LOG.Error(fmt.Sprintf(
				"MustCreateSchoolPlanHis: TryCount:%d SchoolID:%d Error:%s",
				tryCount, records[0].SchoolID, err.Error(),
			))
		case <-ticker.C:
			common.LOG.Error(fmt.Sprintf(
				"MustCreateSchoolPlanHis: Timeout 15s TryCount:%d SchoolID:%d",
				tryCount, records[0].SchoolID,
			))
		}
		time.Sleep(5 * time.Second) // 失败后等待5秒再重试
	}
}
