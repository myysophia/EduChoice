package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func SchoolList(page int) (*response.SchoolListResponse, error) {
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, err
	}
	var data = strings.NewReader(`{"keyword":"","page":` + strconv.Itoa(page) + `,"province_id":"","ranktype":"","request_type":1,"signsafe":"a6beb63405f371aece65cadfb263f006","size":20,"top_school_id":"[2461]","type":"","uri":"apidata/api/gkv3/school/lists"}`)
	req, err := http.NewRequest("POST", "https://api.zjzw.cn/web/api/?keyword=&page="+strconv.Itoa(page)+"&province_id=&ranktype=&request_type=1&size=20&top_school_id=\\[2461\\]&type=&uri=apidata/api/gkv3/school/lists&signsafe=a6beb63405f371aece65cadfb263f006", data)
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
	schoolListResp := &response.SchoolListResponse{}
	if err = json.Unmarshal(bodyText, schoolListResp); err != nil {
		return nil, err
	}
	return schoolListResp, nil
}

func SchoolInfo(schoolId int) (*response.SchoolInfoResponse, error) {
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://static-data.gaokao.cn/www/2.0/school/%d/info.json", schoolId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Microsoft Edge";v="121", "Chromium";v="121"`)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.gaokao.cn/")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	schoolInfoResp := &response.SchoolInfoResponse{}
	if err = json.Unmarshal(bodyText, schoolInfoResp); err != nil {
		return nil, err
	}
	return schoolInfoResp, nil
}

func JobDetail(schoolId int) (*response.JobDetailResponse, error) {
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://static-data.gaokao.cn/www/2.0/school/%d/pc_jobdetail.json", schoolId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Microsoft Edge";v="121", "Chromium";v="121"`)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.gaokao.cn/")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jobDetailResp := &response.JobDetailResponse{}
	if err = json.Unmarshal(bodyText, jobDetailResp); err != nil {
		common.LOG.Error("jobdetail.body:" + string(bodyText))
		jobDetailNullResp := &response.JobDetailNullResp{}
		if err = json.Unmarshal(bodyText, jobDetailNullResp); err != nil {
			common.LOG.Panic("job detail:" + err.Error())
			//fp := &response.FrequentResponse{}
			//if err = json.Unmarshal(bodyText, fp); err != nil {
			//	common.LOG.Panic("job detail:" + err.Error())
			//}
			//if fp.Message == "访问太过频繁，请稍后再试" {
			//	common.LOG.Panic("job detail:" + err.Error())
			//}
		}
		return jobDetailResp, nil
	}
	return jobDetailResp, nil
}

// typeId: 2文科，1理科，2074历史类，2073物理类
func ProvinceScore(schoolId int, provinceId int, typeId int, year int) (*response.ProvinceScoreResponse, error) {
	client, err := proxy.NewHttpClientWithProxy()
	if err != nil {
		return nil, err
	}
	var data = strings.NewReader(`{"e_sort":"zslx_rank,min","e_sorttype":"desc,desc","local_province_id":` + strconv.Itoa(provinceId) + `,"local_type_id":` + strconv.Itoa(typeId) + `,"page":1,"school_id":"` + strconv.Itoa(schoolId) + `","signsafe":"9b95450230606462623e5b42b78e464f","size":10,"uri":"apidata/api/gk/score/province","year":` + strconv.Itoa(year) + `}`)
	url := fmt.Sprintf("https://api.zjzw.cn/web/api/?e_sort=zslx_rank,min&e_sorttype=desc,desc&local_province_id=%d&local_type_id=%d&page=1&school_id=%d&size=10&uri=apidata/api/gk/score/province&year=%d&signsafe=9b95450230606462623e5b42b78e464f", provinceId, typeId, schoolId, year)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authority", "api.zjzw.cn")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://www.gaokao.cn")
	req.Header.Set("referer", "https://www.gaokao.cn/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	provinceScoreResp := &response.ProvinceScoreResponse{}
	if err = json.Unmarshal(bodyText, provinceScoreResp); err != nil {
		common.LOG.Error("Job Detail.Body: " + string(bodyText))
		fp := &response.FrequentResponse{}
		if err = json.Unmarshal(bodyText, fp); err != nil {
			common.LOG.Panic("不是访问过于频繁: " + err.Error())
		}
		if fp.Message == "访问太过频繁，请稍后再试" {
			common.LOG.Panic("是访问过于频繁: " + err.Error())
		}
	}
	return provinceScoreResp, nil
}

func GetPlan(year int) {

}
