package must

import (
	"encoding/json"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"io"
	"net/http"
)

func GetSpecialScoresHis(schoolId, year, typeId, batchId, page int) (*response.SpecialScoresHisResponse, error) {
	url := fmt.Sprintf("https://api.zjzw.cn/web/api/?local_batch_id=%d&local_province_id=61&local_type_id=%d&page=%d&school_id=%d&size=10&sp_xuanke=&special_group=&uri=apidata/api/gk/score/special&year=%d",
		batchId, typeId, page, schoolId, year)

	headers := map[string]string{
		"accept":       "application/json, text/plain, */*",
		"content-type": "application/json",
		"origin":       "https://www.gaokao.cn",
		"referer":      "https://www.gaokao.cn/",
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	majScoresListResp := &response.SpecialScoresHisResponse{}
	if err = json.Unmarshal(bodyText, majScoresListResp); err != nil {
		return nil, err
	}
	return majScoresListResp, nil
}
