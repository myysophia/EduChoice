package migration

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/crawler/must"
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/crawler/safe"
	"github.com/big-dust/DreamBridge/internal/model/major"
	"github.com/big-dust/DreamBridge/internal/model/major_score"
	"github.com/big-dust/DreamBridge/internal/model/school"
	"github.com/big-dust/DreamBridge/internal/model/school_num"
	"github.com/big-dust/DreamBridge/internal/model/school_score"
	"github.com/big-dust/DreamBridge/internal/pkg/algo"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/proxy"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

var (
	wgSchool     = &sync.WaitGroup{}
	wgSpecial    = &sync.WaitGroup{}
	wgPlan       = &sync.WaitGroup{}
	mu           = &sync.RWMutex{}
	SuccessCount = 0
)

func Migrate() {
	//MigrateSchoolScores()
	MigrateSpecialScores()
	//MigratePlanNum()
}

func MigrateSchoolScores() {
	defer LOGPageCount()
	for common.Page <= 147 {
		//学校基础信息
		list := safe.GetSchoolListSafe(common.Page)
		common.Page++ //下一页
		proxy.ChangeHttpProxyIP()
		time.Sleep(2 * time.Second)
		for i, item := range list.Data.Item {
			wgSchool.Add(1)
			go MigrateSchoolScoresOneSafe(i, item)
		}
	}
	wgSchool.Wait()
}

func MigrateSchoolScoresOneSafe(i int, item response.Item) {
	defer func() {
		wgSchool.Done()
		if r := recover(); r != nil {
			common.LOG.Error(fmt.Sprintf("MigrateSchoolScores: Panic: %v", r))
			return
		}
		common.Mu.Lock()
		common.Count++
		common.Mu.Unlock()
		fmt.Printf("\rNumber: %d", common.Count)
	}()
	common.LOG.Info(fmt.Sprintf("当前的爬取第[ %d ]所大学：%s", (common.Page-2)*5+i+1, item.Name))
	// 学校具体信息
	info := safe.GetSchoolInfoSafe(item.SchoolID)
	// 三个指标
	master, _ := strconv.Atoi(info.Data.NumMaster)
	doctor, _ := strconv.Atoi(info.Data.NumDoctor)
	gbh, _ := strconv.Atoi(info.Data.GbhNum)
	// 毕业去向
	promotion, abroad, job := safe.ToWhereSafe(item.SchoolID)
	// 双一流学科text
	text := TextDualClass(info)

	s := &school.School{
		ID:                          item.SchoolID,
		Name:                        item.Name,
		BriefIntroduction:           info.Data.Content,
		SchoolCode:                  item.CodeEnroll,
		MasterPoint:                 master,
		PhdPoint:                    doctor,
		ResearchProject:             gbh,
		TitleDoubleFirstClass:       item.DualClassName == "双一流",
		Title_985:                   item.F985 == 1,
		Title_211:                   item.F211 == 1,
		TitleCollege:                info.Data.LevelName == "普通本科",
		TitleUndergraduate:          info.Data.LevelName == "专科（高职）",
		Region:                      info.Data.CityName,
		Website:                     info.Data.Site,
		RecruitmentPhone:            info.Data.Phone,
		Email:                       info.Data.Email,
		PromotionRate:               promotion,
		AbroadRate:                  abroad,
		EmploymentRate:              job,
		DoubleFirstClassDisciplines: text,
	}
	var scores []*school_score.Score
	// 学校各省历年分数（这里只针对陕西省）
	for year := 2021; year <= 2023; year++ {
		// 物理类
		scores = append(scores, safe.GetScoresSafe(s.ID, common.ShannXi, common.T_Physics, year)...)
		// 历史类
		scores = append(scores, safe.GetScoresSafe(s.ID, common.ShannXi, common.T_History, year)...)
	}
	for year := 2018; year <= 2020; year++ {
		// 理科
		scores = append(scores, safe.GetScoresSafe(s.ID, common.ShannXi, common.T_li, year)...)
		// 文科
		scores = append(scores, safe.GetScoresSafe(s.ID, common.ShannXi, common.T_wen, year)...)
	}
	// 去重
	scoresDistinct := make(map[string]*school_score.Score, 30)
	for _, score := range scores {
		key := fmt.Sprintf("%d-%d-%d-%d-%s-%s-%s", score.SchoolID, score.Location, score.TypeId, score.Year, score.Tag, score.SgName, score.BatchName)
		scoresDistinct[key] = score
	}
	// Must！
	school.MustCreateSchoolScore(s, scoresDistinct)
}

func TextDualClass(info *response.SchoolInfoResponse) string {
	var text string
	for _, dualclass := range info.Data.Dualclass {
		text += dualclass.Class + " "
	}
	return text
}

func LOGPageCount() {
	if r := recover(); r != nil {
		common.LOG.Error(fmt.Sprintf("%v", r))
		common.LOG.Info(string(debug.Stack()))
	}
	common.LOG.Info(fmt.Sprintf("page: %d ,count: %d", common.Page, common.Count))
}

func MigrateSpecialScores() {
	// shoolId列表
	list, err := school.GetSchoolIdList()
	if err != nil {
		common.LOG.Panic("GetSchoolIdList: " + err.Error())
	}
	for _, schoolId := range list {
		wgSpecial.Add(1)
		//go MigrateSpecialScoresOneSafe(schoolId)
		MigrateSpecialScoresOneSafe(schoolId)
		time.Sleep(500 * time.Millisecond)
	}
	wgSpecial.Wait()
}

func MigrateSpecialScoresOneSafe(schoolId int) {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Info(fmt.Sprintf("Panic: MigrateSpecialScoresOneSafe: %v SchoolId: %d", r, schoolId))
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

	specialInfos := safe.MustGetSpecialInfoSafe(schoolId)

	// 拿到school的所有专业
	var specials []*major.Major
	for _, info := range specialInfos.Data {
		id, _ := strconv.Atoi(info.ID)
		special := &major.Major{
			ID:                 id,
			Name:               info.SpecialName,
			NationalFeature:    info.NationFeature == "1", //guo
			Level:              info.TypeName,
			DisciplineCategory: info.Level2Name,
			MajorCategory:      info.Level3Name,
			LimitYear:          info.LimitYear,
			SchoolID:           schoolId,
			SpecialId:          info.SpecialID,
		}
		specials = append(specials, special)
	}

	// 所有专业的招生信息和录取信息
	scores := make(map[string]*major_score.MajorScore, 60)

	recruit := safe.MustGetHistoryRecruit(schoolId)

	for major_id, major := range recruit.Data {
		for kelei_id, kelei := range major {
			for _, yearInfo := range kelei {
				major_id, _ := strconv.Atoi(major_id)
				year, _ := strconv.Atoi(yearInfo.Year)
				num, _ := strconv.Atoi(yearInfo.Num)
				keleiStr := common.Kelei(kelei_id)
				score := &major_score.MajorScore{
					SpecialID:         major_id,
					Location:          "陕西",
					Year:              year,
					Kelei:             keleiStr,
					Batch:             yearInfo.Batch,
					RecruitmentNumber: num,
					LowestScore:       0,
					LowestRank:        0,
				}
				key := fmt.Sprintf("%d-%s-%d-%s-%s", major_id, "陕西", year, keleiStr, yearInfo.Batch)
				if _, exist := scores[key]; !exist {
					scores[key] = score
					continue
				}
				scores[key].RecruitmentNumber += num
				common.LOG.Info(fmt.Sprintf("key exist: %s schoolId: %d", key, schoolId))
			}
		}
	}

	admission := safe.MustGetHistoryAdmission(schoolId)

	for major_id, major := range admission.Data {
		for kelei_id, kelei := range major {
			for _, yearInfo := range kelei {
				major_id, _ := strconv.Atoi(major_id)
				year, _ := strconv.Atoi(yearInfo.Year)
				keleiStr := common.Kelei(kelei_id)
				key := fmt.Sprintf("%d-%s-%d-%s-%s", major_id, "陕西", year, keleiStr, yearInfo.Batch)
				lowestScore, _ := strconv.Atoi(yearInfo.Min)
				rank, _ := strconv.Atoi(fmt.Sprintf("%v", yearInfo.MinSection))
				score, exist := scores[key]
				if !exist {
					score := &major_score.MajorScore{
						SpecialID:         major_id,
						Location:          "陕西",
						Year:              year,
						Kelei:             keleiStr,
						Batch:             yearInfo.Batch,
						RecruitmentNumber: 0,
						LowestScore:       lowestScore,
						LowestRank:        rank,
					}
					scores[key] = score
					continue
				}
				score.LowestScore = lowestScore
				score.LowestRank = rank
			}
		}
	}
	major_score.MustCreateMajorScores(specials, scores)
}

func MigratePlanNum() {
	list, err := school.GetSchoolIdList()
	if err != nil {
		common.LOG.Panic("GetSchoolIdList: " + err.Error())
	}
	elist, err := school_num.NumList()
	if err != nil {
		common.LOG.Panic("NumList: " + err.Error())
	}
	SuccessCount += len(elist)
	list = algo.RemoveFromSlice(list, elist)
	limitCh := make(chan struct{}, 10)
	for _, schoolId := range list {
		limitCh <- struct{}{}
		wgPlan.Add(1)
		go MigrateOneSchoolPlan(schoolId, limitCh)
	}
	wgPlan.Wait()
}

func MigrateOneSchoolPlan(schoolId int, limitCh <-chan struct{}) {
	defer func() {
		<-limitCh
		if r := recover(); r != nil {
			common.LOG.Error("MigrateOneSchoolPlan: Panic" + fmt.Sprint(r))
			common.LOG.Info(string(debug.Stack()))
			return
		}
		mu.Lock()
		SuccessCount++
		mu.Unlock()
		mu.RLock()
		fmt.Printf("\rNumber: %d", SuccessCount)
		mu.RUnlock()
		wgPlan.Done()
	}()
	var schoolNums []*school_num.SchoolNum
	for year := 2021; year <= 2023; year++ {
		for _, typeId := range common.TypeIdsPH {
			sum := 0
			for _, batchId := range common.BatchIds {
				for page := 1; page < 10; page++ {
					plan := must.GetPlanInfo(schoolId, 61, year, typeId, batchId, page, 10)
					if len(plan.Data.Item) == 0 {
						break
					}
					for _, item := range plan.Data.Item {
						num, err := strconv.Atoi(item.Num)
						if err != nil {
							common.LOG.Error("strconv.Atoi(item.Num): " + err.Error())
						}
						sum += num
					}
				}
			}
			// 科目招生人数入slice
			sn := &school_num.SchoolNum{
				SchoolID: schoolId,
				Year:     year,
				TypeID:   typeId,
				Number:   sum,
			}
			schoolNums = append(schoolNums, sn)
		}
	}
	for year := 2018; year <= 2020; year++ {
		for _, typeId := range common.TypeIdsWL {
			sum := 0
			for _, batchId := range common.BatchIds {
				for page := 1; page < 10; page++ {
					plan := must.GetPlanInfo(schoolId, 61, year, typeId, batchId, page, 10)
					if len(plan.Data.Item) == 0 {
						break
					}
					for _, item := range plan.Data.Item {
						num, err := strconv.Atoi(item.Num)
						if err != nil {
							common.LOG.Error("strconv.Atoi(item.Num): " + err.Error())
						}
						sum += num
					}
				}
			}
			sn := &school_num.SchoolNum{
				SchoolID: schoolId,
				Year:     year,
				TypeID:   typeId,
				Number:   sum,
			}
			schoolNums = append(schoolNums, sn)
		}
	}
	must.SchoolNumCreate(schoolNums)
}
