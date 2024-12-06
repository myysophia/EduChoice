package must

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
)

func GetPlanHis(schoolId, year, typeId, batchId, page int) (*response.PlanHisResponse, error) {
	var plan *response.PlanHisResponse
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, fmt.Errorf("proxy.NewHttpClientWithProxy: %w", err)
	}
this:
	for {
		planChan := make(chan *response.PlanHisResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()

				url := fmt.Sprintf("https://api.zjzw.cn/web/api/?local_batch_id=%d&local_province_id=61&local_type_id=%d&page=%d&school_id=%d&sg_xuanke=&size=10&special_group=&uri=apidata/api/gkv3/plan/school&year=%d",
					batchId, typeId, page, schoolId, year)

				headers := map[string]string{
					"accept":       "application/json, text/plain, */*",
					"content-type": "application/json",
					"origin":       "https://www.gaokao.cn",
					"referer":      "https://www.gaokao.cn/",
					"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
				}

				//client := &http.Client{}
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return
				}

				for key, value := range headers {
					req.Header.Add(key, value)
				}

				resp, err := client.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}

				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					common.LOG.Error(fmt.Sprintf("解析JSON失败: %v", err))
					return
				}

				//common.LOG.Info(fmt.Sprintf("Response Type: %T", result))
				//common.LOG.Info(fmt.Sprintf("Response Content: %+v", result))

				resultMap, ok := result.(map[string]interface{})
				if !ok {
					common.LOG.Error("响应格式不正确")
					return
				}

				planResp := &response.PlanHisResponse{
					Code:      fmt.Sprint(resultMap["code"]),
					Message:   fmt.Sprint(resultMap["message"]),
					Location:  fmt.Sprint(resultMap["location"]),
					Encrydata: fmt.Sprint(resultMap["encrydata"]),
				}

				if data, ok := resultMap["data"].(map[string]interface{}); ok {
					planResp.Data.NumFound = int(data["numFound"].(float64))
					if items, ok := data["item"].([]interface{}); ok {
						for _, item := range items {
							if itemMap, ok := item.(map[string]interface{}); ok {
								planItem := response.PlanHisItem{
									Length:         fmt.Sprint(itemMap["length"]),
									Num:            fmt.Sprint(itemMap["num"]),
									LocalBatchName: fmt.Sprint(itemMap["local_batch_name"]),
									LocalTypeName:  fmt.Sprint(itemMap["local_type_name"]),
									ProvinceName:   fmt.Sprint(itemMap["province_name"]),
									SchoolID:       fmt.Sprint(itemMap["school_id"]),
									SpName:         fmt.Sprint(itemMap["sp_name"]),
									SpecialGroup:   fmt.Sprint(itemMap["special_group"]),
									Spname:         fmt.Sprint(itemMap["spname"]),
									Tuition:        fmt.Sprint(itemMap["tuition"]),
									Year:           fmt.Sprint(itemMap["year"]),
								}
								planResp.Data.Item = append(planResp.Data.Item, planItem)
							}
						}
					}
				}

				planChan <- planResp
			}()
		}

		ticker := time.NewTicker(5 * time.Second)
		select {
		case plan = <-planChan:
			break this
		case <-ticker.C:
			common.LOG.Error("GetPlanHis: time out 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
	return plan, nil
}
