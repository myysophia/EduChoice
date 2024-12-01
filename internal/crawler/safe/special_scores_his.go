package safe

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"time"
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
				errChan <- err
				return
			}
			resChan <- res
		}()

		ticker := time.NewTicker(15 * time.Second)
		select {
		case res := <-resChan:
			return res
		case <-ticker.C:
			select {
			case err := <-errChan:
				common.LOG.Error(fmt.Sprintf("MustGetSpecialScoresHis: Time Out 15s TryCount:%d SchoolId:%d Error:%s", 
					tryCount, schoolId, err.Error()))
			default:
				common.LOG.Error(fmt.Sprintf("MustGetSpecialScoresHis: Time Out 15s TryCount:%d SchoolId:%d", 
					tryCount, schoolId))
			}
		}
	}
}