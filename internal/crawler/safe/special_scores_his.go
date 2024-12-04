package safe

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
)

func MustGetSpecialScoresHis(schoolId, year, typeId, batchId, page int) *response.SpecialScoresHisResponse {
	tryCount := 0
	for {
		tryCount++
		errChan := make(chan error, 1)
		resChan := make(chan *response.SpecialScoresHisResponse, 1)

		go func() {
			client, err := proxy.NewHttpClientWithProxy()
			if err != nil {
				errChan <- fmt.Errorf("创建代理客户端失败: %v", err)
				return
			}

			transport, ok := client.Transport.(*http.Transport)
			proxyInfo := "direct"
			if ok && transport.Proxy != nil {
				if proxyURL, err := transport.Proxy(&http.Request{URL: &url.URL{Scheme: "https"}}); err == nil && proxyURL != nil {
					proxyInfo = proxyURL.String()
				}
			}

			res, err := must.GetSpecialScoresHis(schoolId, year, typeId, batchId, page)
			if err != nil {
				common.LOG.Error(fmt.Sprintf(
					"GetSpecialScoresHis failed - SchoolId: %d, Year: %d, TypeId: %d, BatchId: %d, Page: %d, Proxy: %s, Error: %v",
					schoolId, year, typeId, batchId, page, proxyInfo, err,
				))
				errChan <- err
				return
			}
			resChan <- res
		}()

		ticker := time.NewTicker(5 * time.Second)
		select {
		case res := <-resChan:
			return res
		case err := <-errChan:
			common.LOG.Error(fmt.Sprintf(
				"MustGetSpecialScoresHis: TryCount:%d SchoolId:%d Year:%d TypeId:%d BatchId:%d Page:%d Error:%s",
				tryCount, schoolId, year, typeId, batchId, page, err.Error(),
			))
		case <-ticker.C:
			proxy.ChangeHttpProxyIP()
		}
	}
}
