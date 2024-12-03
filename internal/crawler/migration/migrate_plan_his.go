package migration

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/model/school"
	"github.com/big-dust/DreamBridge/internal/model/school_plan_his"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

func MigratePlanHis() {
	list, err := school.GetSchoolIdList()
	if err != nil {
		common.LOG.Panic("GetSchoolIdList: " + err.Error())
	}

	for _, schoolId := range list {
		wgPlan.Add(1)
		MigratePlanHisOneSafe(schoolId)
		time.Sleep(500 * time.Millisecond)
	}
	wgPlan.Wait()
}

func MigratePlanHisOneSafe(schoolId int) {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Info(fmt.Sprintf("Panic: MigratePlanHisOneSafe: %v SchoolId: %d", r, schoolId))
			common.LOG.Error(string(debug.Stack()))
			return
		}
		mu.Lock()
		SuccessCount++
		mu.Unlock()
		wgPlan.Done()
		mu.RLock()
		fmt.Printf("\rNumber: %d", SuccessCount)
		mu.RUnlock()
	}()

	var allPlans []*school_plan_his.SchoolPlanHis

	// 遍历年份
	for year := 2024; year <= 2024; year++ {
		// 遍历批次
		for _, batchId := range common.BatchIds {
			// 遍历科类
			for _, typeId := range common.TypeIdsWL {
				page := 1
				for {
					plans, err := must.GetPlanHis(schoolId, year, typeId, batchId, page)
					if err != nil {
						common.LOG.Error(fmt.Sprintf("GetPlanHis failed: %v", err))
						break
					}
					if len(plans.Data.Item) == 0 {
						break
					}

					for _, item := range plans.Data.Item {
						num, _ := strconv.Atoi(item.Num)
						specialGroup, _ := strconv.Atoi(item.SpecialGroup)
						schoolId, _ := strconv.Atoi(item.SchoolID)
						year, _ := strconv.Atoi(item.Year)

						plan := &school_plan_his.SchoolPlanHis{
							SchoolID:       schoolId,
							Year:           year,
							SpName:         item.SpName,
							Spname:         item.Spname,
							Num:            num,
							Length:         item.Length,
							Tuition:        item.Tuition,
							ProvinceName:   item.ProvinceName,
							SpecialGroup:   specialGroup,
							LocalBatchName: item.LocalBatchName,
							LocalTypeName:  item.LocalTypeName,
						}
						allPlans = append(allPlans, plan)
					}

					if len(plans.Data.Item) < 10 {
						break
					}

					page++
				}
			}
		}
	}

	school_plan_his.MustCreateSchoolPlanHis(allPlans)
}
