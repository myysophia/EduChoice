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
		for i := 1; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()
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
		}

		ticker := time.NewTicker(5 * time.Second)

		select {
		case res = <-resChan:
			break this
		case err := <-errChan:
			common.LOG.Error(fmt.Sprintf(
				"MustGetSpecialScoresHis: TryCount:%d SchoolId:%d Year:%d TypeId:%d BatchId:%d Page:%d Error:%s",
				tryCount, schoolId, year, typeId, batchId, page, err.Error(),
			))

			// 如果是访问频率限制错误，增加更长的等待时间
			if err.Error() == "访问频率限制: 请求过多，请稍后再试" {
				time.Sleep(30 * time.Second) // 访问频率限制时等待30秒
			} else {
				time.Sleep(5 * time.Second) // 其他错误等待5秒
			}

		case <-ticker.C:
			common.LOG.Error("get score his time out 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
	return res
}
