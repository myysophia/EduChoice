package safe

import (
	"fmt"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

func MustGetSpecialScoresHis(schoolId, year, typeId, batchId, page int) *response.SpecialScoresHisResponse {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		resChan := make(chan *response.SpecialScoresHisResponse, 1)

		go func() {
			res, err := must.GetSpecialScoresHis(schoolId, year, typeId, batchId, page)
			if err != nil {
				common.LOG.Error(fmt.Sprintf(
					"GetSpecialScoresHis failed - SchoolId: %d, Year: %d, TypeId: %d, BatchId: %d, Page: %d, Error: %v",
					schoolId, year, typeId, batchId, page, err,
				))
				errChan <- err
				return
			}
			resChan <- res
		}()

		ticker := time.NewTicker(15 * time.Second)

		select {
		case res := <-resChan:
			return res
		case err := <-errChan:
			common.LOG.Error(fmt.Sprintf(
				"MustGetSpecialScoresHis: Time Out 15s TryCount:%d SchoolId:%d Year:%d TypeId:%d BatchId:%d Page:%d Error:%s",
				tryCount, schoolId, year, typeId, batchId, page, err.Error(),
			))
		case <-ticker.C:
			common.LOG.Error(fmt.Sprintf(
				"MustGetSpecialScoresHis: Time Out 15s TryCount:%d SchoolId:%d Year:%d TypeId:%d BatchId:%d Page:%d",
				tryCount, schoolId, year, typeId, batchId, page,
			))
		}
	}
}
