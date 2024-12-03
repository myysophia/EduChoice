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

func MustCreateSchoolPlanHis(records []*SchoolPlanHis) {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		nilChan := make(chan error, 1)
		go func() {
			err := CreateSchoolPlanHis(records)
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
				common.LOG.Error(fmt.Sprintf("MustCreateSchoolPlanHis: Time Out 15s TryCount:%d SchoolID:%v Error: %s",
					tryCount, records[0].SchoolID, err.Error()))
			default:
				common.LOG.Error(fmt.Sprintf("MustCreateSchoolPlanHis: Time Out 15s TryCount:%d SchoolID:%v",
					tryCount, records[0].SchoolID))
			}
		}
	}
}

func CreateSchoolPlanHis(records []*SchoolPlanHis) error {
	tx := common.DB.Begin()
	for _, record := range records {
		if err := tx.Create(record).Error; err != nil {
			tx.Rollback()
			if common.ErrMysqlDuplicate.Is(err) {
				return nil
			}
			return err
		}
	}
	return tx.Commit().Error
}
