package migration

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/safe"
	"github.com/big-dust/DreamBridge/internal/model/major_score_his"
	"github.com/big-dust/DreamBridge/internal/model/school"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

func MigrateSpecialScoresHis() {
	// 获取学校ID列表
	list, err := school.GetSchoolIdList()
	if err != nil {
		common.LOG.Panic("GetSchoolIdList: " + err.Error())
	}

	for _, schoolId := range list {
		wgSpecial.Add(1)
		MigrateSpecialScoresHisOneSafe(schoolId)
		time.Sleep(500 * time.Millisecond)
	}
	wgSpecial.Wait()
}

func MigrateSpecialScoresHisOneSafe(schoolId int) {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Info(fmt.Sprintf("Panic: MigrateSpecialScoresHisOneSafe: %v SchoolId: %d", r, schoolId))
			common.LOG.Error(string(debug.Stack()))
			return
		}
		mu.Lock()
		SuccessCount++
		mu.Unlock()
		wgSpecial.Done()
		mu.RLock()
		fmt.Printf("\rNumber: %d", SuccessCount)
		mu.RUnlock()
	}()

	var allScores []*major_score_his.MajorScoreHis

	// 遍历年份
	for year := 2020; year <= 2023; year++ {
		// 遍历批次
		for _, batchId := range common.BatchIds {
			// 遍历科类
			for _, typeId := range common.TypeIdsWL {
				page := 1
				for {
					scores := safe.MustGetSpecialScoresHis(schoolId, year, typeId, batchId, page)
					if len(scores.Data.Item) == 0 {
						break
					}

					for _, item := range scores.Data.Item {
						score := &major_score_his.MajorScoreHis{
							ID:                item.ID,
							SchoolID:          item.SchoolID,
							SpecialID:         item.SpecialID,
							SpeID:             item.SpeID,
							Year:              item.Year,
							SpName:            item.SpName,
							SpInfo:            item.SpInfo,
							Info:              item.Info,
							LocalProvinceName: item.LocalProvinceName,
							LocalTypeName:     item.LocalTypeName,
							LocalBatchName:    item.LocalBatchName,
							Level2Name:        item.Level2Name,
							Level3Name:        item.Level3Name,
							Average:           item.Average,
							Max:               item.Max,
							Min:               item.Min,
							MinSection:        item.MinSection,
							Proscore:          item.Proscore,
							DoubleHigh:        item.Doublehigh,
							IsTop:             item.IsTop,
							IsScoreRange:      item.IsScoreRange,
							MinRange:          item.MinRange,
							MinRankRange:      item.MinRankRange,
							Remark:            item.Remark,
							ZslxName:          item.ZslxName,
							DualClassName:     item.DualClassName,
							FirstKm:           item.FirstKm,
							SgFxk:             item.SgFxk,
							SgSxk:             item.SgSxk,
							SgInfo:            item.SgInfo,
							SgName:            item.SgName,
							SgType:            item.SgType,
							SpFxk:             item.SpFxk,
							SpSxk:             item.SpSxk,
							SpType:            item.SpType,
							Single:            item.Single,
							SpecialGroup:      item.SpecialGroup,
						}
						allScores = append(allScores, score)
					}

					// 如果当前页的数据量小于size，说明已经是最后一页
					if len(scores.Data.Item) < 10 {
						break
					}
					page++
				}
			}
		}
	}

	major_score_his.MustCreateMajorScoresHis(allScores)
}

// 添加辅助函数
func mustAtoi(n json.Number) int {
	if n == "" || n == "-" {
		return 0
	}
	i, err := n.Int64()
	if err != nil {
		return 0
	}
	return int(i)
}
