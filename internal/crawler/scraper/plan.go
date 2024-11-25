package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func PlanInfo(schoolId int, provinceId int, year int, typeId string, batchId int, page int, size int) (*response.PlanInfoResp, error) {
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, err
	}
	// body
	dataMap := map[string]any{
		"local_batch_id":    batchId,
		"local_province_id": provinceId,
		"local_type_id":     typeId,
		"page":              page,
		"school_id":         strconv.Itoa(schoolId),
		"signsafe":          "9f01e17aad5d03a4837de21dea88bb31",
		"size":              size,
		"special_group":     "",
		"uri":               "apidata/api/gkv3/plan/school",
		"year":              year,
	}
	b, _ := json.Marshal(dataMap)
	var data = strings.NewReader(string(b))

	// query
	values := url.Values{}
	values.Set("local_batch_id", strconv.Itoa(batchId))
	values.Set("local_province_id", strconv.Itoa(provinceId))
	values.Set("local_type_id", typeId)
	values.Set("page", strconv.Itoa(page))
	values.Set("school_id", strconv.Itoa(schoolId))
	values.Set("signsafe", "9f01e17aad5d03a4837de21dea88bb31")
	values.Set("size", strconv.Itoa(size))
	values.Set("special_group", "")
	values.Set("uri", "apidata/api/gkv3/plan/school")
	values.Set("year", strconv.Itoa(year))

	req, err := http.NewRequest("POST", "https://api.zjzw.cn/web/api/?"+values.Encode(), data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authority", "api.zjzw.cn")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://www.gaokao.cn")
	req.Header.Set("referer", "https://www.gaokao.cn/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Microsoft Edge";v="121", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	plResp := &response.PlanInfoResp{}
	if err = json.Unmarshal(bodyText, plResp); err != nil {
		common.LOG.Info("URL Values: " + fmt.Sprint(dataMap))
		common.LOG.Info("PlanInfo: " + string(bodyText))
		common.LOG.Error("PlanInfo: " + err.Error())
		return nil, err
	}
	return plResp, nil
}
