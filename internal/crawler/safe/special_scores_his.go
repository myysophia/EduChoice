package safe

import (
	"fmt"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
)

func MustGetSpecialScoresHis(schoolId, year, typeId, batchId, page int) *response.SpecialScoresHisResponse {
	tryCount := 0
	var res *response.SpecialScoresHisResponse
this:
	for {
		tryCount++
		errChan := make(chan error, 1)
		resChan := make(chan *response.SpecialScoresHisResponse, 1)
		done := make(chan bool, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					common.LOG.Error(fmt.Sprintf("%v", r))
				}
				done <- true
			}()
			res, err := must.GetSpecialScoresHis(schoolId, year, typeId, batchId, page)
			if err != nil {
				errChan <- err
				return
			}
			resChan <- res
		}()

		ticker := time.NewTicker(15 * time.Second)
		//defer ticker.Stop()

		select {
		case res = <-resChan:
			break this
		case err := <-errChan:
			common.LOG.Error(fmt.Sprintf(
				"MustGetSpecialScoresHis: TryCount:%d SchoolId:%d Year:%d TypeId:%d BatchId:%d Page:%d Error:%s",
				tryCount, schoolId, year, typeId, batchId, page, err.Error(),
			))
			if err.Error() == "访问频率限制: 请求过多，请稍后再试" {
				time.Sleep(30 * time.Second)
			} else {
				time.Sleep(5 * time.Second)
			}
		case <-ticker.C:
			<-done
			common.LOG.Error("get score his time out 15s")
			proxy.ChangeHttpProxyIP()
		}
	}
	return res
}
