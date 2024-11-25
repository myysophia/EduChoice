// must 必须获取到指定的信息，但是获取不到并没有采取panic而是无限循环
package must

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/crawler/scraper"
	"github.com/big-dust/DreamBridge/internal/model/school_num"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"time"
)

func GetPlanInfo(schoolId int, provinceId int, year int, typeId string, batchId int, page int, size int) *response.PlanInfoResp {
	for {
		planCh := make(chan *response.PlanInfoResp)
		for i := 0; i < 2; i++ {
			go func() {
				plan, err := scraper.PlanInfo(schoolId, provinceId, year, typeId, batchId, page, size)
				if err != nil {
					common.LOG.Error("scraper.PlanInfo: " + err.Error())
					return
				}
				planCh <- plan
			}()
		}

		ticker := time.NewTicker(5 * time.Second)
		select {
		case plan := <-planCh:
			return plan
		case <-ticker.C:
			// 超时
			common.LOG.Info("scraper.PlanInfo: 请求超时: 5s")
			proxy.ChangeHttpProxyIP()
		}
	}
}

func SchoolNumCreate(sn []*school_num.SchoolNum) {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Error(fmt.Sprint(r))
		}
	}()
	for {
		err := school_num.SchoolNumCreate(sn)
		if err != nil {
			common.LOG.Error("model.SchoolNumCreate(sn)" + err.Error())
			common.LOG.Info("sn: SchoolID: " + fmt.Sprint(sn[0].SchoolID))
			if common.ErrMysqlDuplicate.Is(err) {
				return
			}
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}
