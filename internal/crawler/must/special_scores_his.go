package must

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	//"github.com/big-dust/DreamBridge/internal/crawler/common"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
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
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 使用 interface{} 接收所有字段
	var rawResponse map[string]interface{}
	if err = json.Unmarshal(bodyText, &rawResponse); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	// 手动构造最终响应对象
	finalResponse := &response.SpecialScoresHisResponse{
		Code:      fmt.Sprint(rawResponse["code"]),
		Message:   fmt.Sprint(rawResponse["message"]),
		Location:  fmt.Sprint(rawResponse["location"]),
		Encrydata: fmt.Sprint(rawResponse["encrydata"]),
	}

	// 解析 data 字段
	dataMap, ok := rawResponse["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data字段类型不匹配")
	}

	// 解析 item 数组
	items, ok := dataMap["item"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("item字段类型不匹配")
	}

	for _, item := range items {
		itemMap := item.(map[string]interface{})
		specialItem := response.SpecialScoresHisItem{
			Average:           fmt.Sprint(itemMap["average"]),
			Doublehigh:        safeToInt(itemMap["doublehigh"]),
			DualClassName:     fmt.Sprint(itemMap["dual_class_name"]),
			FirstKm:           safeToInt(itemMap["first_km"]),
			ID:                fmt.Sprint(itemMap["id"]),
			Info:              fmt.Sprint(itemMap["info"]),
			IsScoreRange:      safeToInt(itemMap["is_score_range"]),
			IsTop:             safeToInt(itemMap["is_top"]),
			Level2Name:        fmt.Sprint(itemMap["level2_name"]),
			Level3Name:        fmt.Sprint(itemMap["level3_name"]),
			LocalBatchName:    fmt.Sprint(itemMap["local_batch_name"]),
			LocalProvinceName: fmt.Sprint(itemMap["local_province_name"]),
			LocalTypeName:     fmt.Sprint(itemMap["local_type_name"]),
			Max:               fmt.Sprint(itemMap["max"]),
			Min:               fmt.Sprint(itemMap["min"]),
			MinRange:          fmt.Sprint(itemMap["min_range"]),
			MinRankRange:      fmt.Sprint(itemMap["min_rank_range"]),
			MinSection:        safeToInt(itemMap["min_section"]),
			Name:              fmt.Sprint(itemMap["name"]),
			Proscore:          safeToInt(itemMap["proscore"]),
			Remark:            fmt.Sprint(itemMap["remark"]),
			SchoolID:          safeToInt(itemMap["school_id"]),
			SgFxk:             fmt.Sprint(itemMap["sg_fxk"]),
			SgInfo:            fmt.Sprint(itemMap["sg_info"]),
			SgName:            fmt.Sprint(itemMap["sg_name"]),
			SgSxk:             fmt.Sprint(itemMap["sg_sxk"]),
			SgType:            safeToInt(itemMap["sg_type"]),
			Single:            fmt.Sprint(itemMap["single"]),
			SpFxk:             fmt.Sprint(itemMap["sp_fxk"]),
			SpInfo:            fmt.Sprint(itemMap["sp_info"]),
			SpName:            fmt.Sprint(itemMap["sp_name"]),
			SpSxk:             fmt.Sprint(itemMap["sp_sxk"]),
			SpType:            safeToInt(itemMap["sp_type"]),
			SpeID:             safeToInt(itemMap["spe_id"]),
			SpecialGroup:      safeToInt(itemMap["special_group"]),
			SpecialID:         safeToInt(itemMap["special_id"]),
			Spname:            fmt.Sprint(itemMap["spname"]),
			Year:              safeToInt(itemMap["year"]),
			ZslxName:          fmt.Sprint(itemMap["zslx_name"]),
		}
		finalResponse.Data.Item = append(finalResponse.Data.Item, specialItem)
	}

	finalResponse.Data.NumFound = int(dataMap["numFound"].(float64))

	return finalResponse, nil
}

// safeToInt 尝试将 interface{} 转换为 int，处理可能的类型错误
func safeToInt(value interface{}) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}
