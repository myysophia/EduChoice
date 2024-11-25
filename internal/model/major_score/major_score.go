package major_score

import (
	"database/sql"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/model/major"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"time"
)

// 定义 MajorScores 表格的模型
type MajorScore struct {
	ID                int
	SpecialID         int
	Location          string
	Year              int
	Kelei             string
	Batch             string
	RecruitmentNumber int
	LowestScore       int
	LowestRank        int
}

func FindScoreAvg(specialID int, kl string) (*sql.NullFloat64, error) {
	avgScore := &sql.NullFloat64{}
	if err := common.DB.Table("major_scores").
		Select("avg(lowest_score) as avg_score").
		Find(avgScore, "special_id = ? and kelei = ? and year > 2020 and year <= 2023 and lowest_score != 0", specialID, kl).Error; err != nil {
		return nil, err
	}
	return avgScore, nil
}

func FindByKeleiIn(majorIds []int, kl ...string) ([]int, error) {
	var mids []int
	if err := common.DB.Table("major_scores").Distinct("special_id").Find(&mids, "special_id in (?) and kelei in (?)", majorIds, kl).Error; err != nil {
		return nil, err
	}
	return mids, nil
}

func CreateMajorScores(major []*major.Major, majorScores map[string]*MajorScore) error {
	tx := common.DB.Begin()
	if err := tx.Create(major).Error; err != nil {
		tx.Rollback()
		if common.ErrMysqlDuplicate.Is(err) {
			return nil
		}
		common.LOG.Error("MustCreateMajorScores: " + err.Error())
		return err
	}
	for _, score := range majorScores {
		if err := tx.Create(score).Error; err != nil {
			tx.Rollback()
			if common.ErrMysqlDuplicate.Is(err) {
				return nil
			}
			common.LOG.Error("MustCreateMajorScores: " + err.Error())
			return err
		}
	}
	tx.Commit()
	return nil
}

func MustCreateMajorScores(major []*major.Major, majorScores map[string]*MajorScore) {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		nilChan := make(chan error, 1)
		go func() {
			err := CreateMajorScores(major, majorScores)
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
				common.LOG.Error(fmt.Sprintf("MustCreateSchoolScore:	Time Out 15s TryCount:%d Major.SchoolId: %v Error: %s", tryCount, major[0].SchoolID, err.Error()))
			default:
				common.LOG.Error(fmt.Sprintf("MustCreateSchoolScore:	Time Out 15s TryCount:%d Major.SchoolId: %v", tryCount, major[0].SchoolID))
			}
		}
	}
}
