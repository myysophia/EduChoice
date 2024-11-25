package safe

import (
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/crawler/scraper"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"time"
)

func MustGetSpecialInfoSafe(schoolId int) *response.SpecialResponse[response.SpecialList] {
	for {
		specialInfosChan := make(chan *response.SpecialResponse[response.SpecialList], 1)
		for i := 0; i < 2; i++ {
			go func() {
				specialInfos, err := scraper.SpecialInfo(schoolId)
				if err != nil {
					common.LOG.Error("SpecialInfo: " + err.Error())
					return
				}
				specialInfosChan <- specialInfos
			}()
		}
		var specialInfos *response.SpecialResponse[response.SpecialList]
		ticker := time.NewTicker(5 * time.Second)
		select {
		case specialInfos = <-specialInfosChan:
			return specialInfos
		case <-ticker.C:
			common.LOG.Error("MustGetSpecialInfoSafe: Time Out: 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
}

func MustGetHistoryRecruit(schoolId int) *response.HistoryRecruitResponse {
	for {
		recruitChan := make(chan *response.HistoryRecruitResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				recruit, err := scraper.HistoryRecruit(schoolId, common.HuBei)
				if err != nil {
					common.LOG.Error("HistoryRecruit:" + err.Error())
					return
				}
				recruitChan <- recruit
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case recruit := <-recruitChan:
			return recruit
		case <-ticker.C:
			common.LOG.Error("MustGetHistoryRecruit: Time Out: 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
}

func MustGetHistoryAdmission(schoolId int) *response.HistoryAdmissionResponse {
	for {
		admissionChan := make(chan *response.HistoryAdmissionResponse, 1)
		for i := 0; i < 2; i++ {
			go func() {
				admission, err := scraper.HistoryAdmission(schoolId, common.HuBei)
				if err != nil {
					common.LOG.Error("HistoryAdmission: Err:" + err.Error())
					return
				}
				admissionChan <- admission
			}()
		}
		ticker := time.NewTicker(5 * time.Second)
		select {
		case admission := <-admissionChan:
			return admission
		case <-ticker.C:
			common.LOG.Error("MustGetHistoryAdmission: Time Out: 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
}
