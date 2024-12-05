package migration

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/big-dust/DreamBridge/internal/crawler/safe"
	"github.com/big-dust/DreamBridge/internal/model/major_score_his"
	"github.com/big-dust/DreamBridge/internal/model/school"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
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
		// 每个学校处理完后切换代理并等待，避免被反爬
		proxy.ChangeHttpProxyIP()
		time.Sleep(3000 * time.Millisecond)
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
					for _, item := range scores.Data.Item {
						score := &major_score_his.MajorScoreHis{
							ID:                item.ID,
							SchoolID:          item.SchoolID,
							SpecialID:         item.SpecialID,
							SpeID:             item.SpeID,
							Year:              item.Year,
							SpName:            item.SpName,
							Spname:            item.Spname,
							Info:              item.Info,
							LocalProvinceName: item.LocalProvinceName,
							LocalTypeName:     item.LocalTypeName,
							LocalBatchName:    item.LocalBatchName,
							Level2Name:        item.Level2Name,
							Level3Name:        item.Level3Name,
							Average:           parseIntWithDefault(item.Average),
							Max:               parseIntWithDefault(item.Max),
							Min:               parseIntWithDefault(item.Min),
							MinSection:        item.MinSection,
							Proscore:          item.Proscore,
							IsTop:             item.IsTop,
							IsScoreRange:      item.IsScoreRange,
							MinRange:          item.MinRange,
							MinRankRange:      item.MinRankRange,
							Remark:            item.Remark,
						}
						allScores = append(allScores, score)
					}

					// 如果当前页的数据量小于size，说明已经是最后一页
					if len(scores.Data.Item) < 10 {
						break
					}
					page++
					// 每页数据处理完后切换代理并等待
					// proxy.ChangeHttpProxyIP()
					// time.Sleep(500 * time.Millisecond)
				}
			}
		}
	}

	major_score_his.MustCreateMajorScoresHis(allScores)
}

// 添加辅助函数
func parseIntWithDefault(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// 添加新的辅助函数来处理ID
func processID(originalID string) string {
	// 移除"gkspecialscore"前缀
	return strings.TrimPrefix(originalID, "gkspecialscore")
}
