// 可以安全的用于并发，safe中的函数不会导致父协程panic
package safe

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/crawler/scraper"
	"github.com/big-dust/DreamBridge/internal/model/school_score"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"strconv"
	"time"
)

func ToWhereSafe(schoolId int) (promotion string, abroad string, job string) {
	// 学生毕业去向
	var detail *response.JobDetailResponse
this:
	for {
		detailChan := make(chan *response.JobDetailResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()
				detail, err := scraper.JobDetail(schoolId)
				if err != nil || detail == nil {
					return
				}
				detailChan <- detail
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case detail = <-detailChan:
			break this
		case <-ticker.C:
			common.LOG.Error("get job detail: time out 5 s")
			proxy.ChangeHttpProxyIP()
		}
	}
	promotion = detail.Data.Jobrate.Postgraduate.One
	abroad = detail.Data.Jobrate.Abroad.One
	job = detail.Data.Jobrate.Job.One
	return
}

func GetScoresSafe(schooldId int, provinceId int, type_id int, year int) []*school_score.Score {
	var score *response.ProvinceScoreResponse
this:
	for {
		scoreChan := make(chan *response.ProvinceScoreResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()
				score, err := scraper.ProvinceScore(schooldId, provinceId, type_id, year)
				if err != nil || score == nil {
					return
				}
				scoreChan <- score
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case score = <-scoreChan:
			break this
		case <-ticker.C:
			common.LOG.Error("get score: time out 5 s")
			proxy.ChangeHttpProxyIP()
		}
	}
	var scores []*school_score.Score
	for _, item := range score.Data.Item {
		rank := fmt.Sprint(item.MinSection)
		r, _ := strconv.Atoi(rank)
		score := &school_score.Score{
			SchoolID:   schooldId,
			Location:   common.ShannXi,
			Year:       year,
			TypeId:     type_id,
			BatchName:  item.LocalBatchName,
			Tag:        item.ZslxName, // 专项...
			SgName:     item.SgName,
			Lowest:     item.Min,
			LowestRank: r,
		}
		scores = append(scores, score)
	}
	return scores
}

func GetSchoolListSafe(page int) *response.SchoolListResponse {
	var list *response.SchoolListResponse
this:
	for {
		listChan := make(chan *response.SchoolListResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()
				list, err := scraper.SchoolList(page)
				if err != nil || list == nil {
					return
				}
				listChan <- list
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case list = <-listChan:
			break this
		case <-ticker.C:
			common.LOG.Error("get list: time out 5 s")
			proxy.ChangeHttpProxyIP()
		}
	}
	return list
}

func GetSchoolInfoSafe(schoolId int) *response.SchoolInfoResponse {
	var info *response.SchoolInfoResponse
this:
	for {
		infoChan := make(chan *response.SchoolInfoResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.LOG.Error(fmt.Sprintf("%v", r))
					}
				}()
				info, err := scraper.SchoolInfo(schoolId)
				if err != nil || info == nil {
					return
				}
				infoChan <- info
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case info = <-infoChan:
			break this
		case <-ticker.C:
			common.LOG.Error("get school info: time out 5 s")
			proxy.ChangeHttpProxyIP()
		}
	}
	return info
}
