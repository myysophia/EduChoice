package zy

import (
	"context"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/api/types"
	"github.com/big-dust/DreamBridge/internal/model/major"
	"github.com/big-dust/DreamBridge/internal/model/major_score"
	"github.com/big-dust/DreamBridge/internal/model/school"
	"github.com/big-dust/DreamBridge/internal/model/school_num"
	"github.com/big-dust/DreamBridge/internal/model/school_score"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"github.com/big-dust/DreamBridge/internal/pkg/algo"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"sort"
	"strconv"
	"strings"
)

// [start:end) 前闭后开
func GetSchools(uid int, startSore int, endScore int) ([]*types.School, error) {
	// 获取考生信息
	u, err := user.FindOne(uid)
	if err != nil {
		return nil, err
	}
	// TypeID
	typeID := common.TypeID(u.Physics, u.History)
	if typeID == 0 {
		return nil, fmt.Errorf("不合法的选科: 物理: %v 历史: %v", u.Physics, u.History)
	}
	// 获取符合的所有学校id
	sids, err := school_score.SchoolIdsIn(u.Score+startSore, u.Score+endScore, typeID)
	if err != nil {
		return nil, err
	}
	// 专科、本科筛选
	if sids, err = school.FindIDsByLevelIn(sids, u.SchoolType == "本科"); err != nil {
		return nil, err
	}
	// id作为基础信息实例化schools
	schools := newSchools(sids)
	// 学校名称
	if err = setSchoolName(schools); err != nil {
		return nil, err
	}
	// 学校历年信息
	if err = setHistoryInfo(schools, typeID); err != nil {
		return nil, err
	}
	// 学校专业、按录取率分开
	if err = setMajorToSchool(schools, typeID, u.Score); err != nil {
		return nil, err
	}
	// 特征排序
	traitSort(schools, u)
	return schools, nil
}

// id作为基础信息实例化schools,没有添加其他信息
func newSchools(schoolIds []int) []*types.School {
	var schools []*types.School
	for _, id := range schoolIds {
		schools = append(schools,
			&types.School{
				ID:           id,
				Parts:        make(map[int][]*types.Major),
				HistoryInfos: make(map[int]*types.HistoryInfo),
			})
	}
	return schools
}

func setSchoolName(schools []*types.School) error {
	for _, s := range schools {
		sItem, err := school.FindOne(s.ID)
		if err != nil {
			return err
		}
		s.Name = sItem.Name
	}
	return nil
}

func setHistoryInfo(schools []*types.School, typeID int) error {
	for _, s := range schools {
		// 历年招生人数
		ens, err := school_num.FindHistoryEnrollmentNum(s.ID)
		if err != nil {
			return err
		}
		for _, en := range ens {
			info, ok := s.HistoryInfos[en.Year]
			if !ok {
				s.HistoryInfos[en.Year] = &types.HistoryInfo{}
				info = s.HistoryInfos[en.Year]
			}
			info.EnrollmentNum = en.Number
		}
		// 历年最低分数和最低排名
		hs, err := school_score.FindHistoryScore(s.ID, typeID)
		if err != nil {
			return err
		}
		for _, h := range hs {
			info, ok := s.HistoryInfos[h.Year]
			if !ok {
				s.HistoryInfos[h.Year] = &types.HistoryInfo{}
				info = s.HistoryInfos[h.Year]
			}
			info.LowestRank = h.LowestRank
			info.LowestScore = h.Lowest
		}
	}
	return nil
}

func setMajorToSchool(schools []*types.School, typeID int, studentScore int) error {
	for _, s := range schools {
		// 获取所有的专业
		majors, err := major.FindBySchoolID(s.ID)
		if err != nil {
			return err
		}
		mIDs, err := major.FindIDListBySchoolID(s.ID)
		if err != nil {
			return err
		}
		//获取不考虑的科类
		omitKelei := common.Omit(typeID)
		//过滤掉不能上的 typeID 以外的
		omitIDs, err := major_score.FindByKeleiIn(mIDs, omitKelei...)
		if err != nil {
			return err
		}
		// 获取必须要的 IDs
		kl := common.IDConvKelei(typeID)
		neededIDs, err := major_score.FindByKeleiIn(mIDs, kl)
		// omit - needIDs
		algo.RemoveFromSlice(omitIDs, neededIDs)
		omit := common.SliceToMap[int](omitIDs)
		for i, m := range majors {
			specialID, err := strconv.Atoi(m.SpecialId)
			if err != nil {
				return err
			}
			if omit[specialID] {
				majors = append(majors[:i], majors[i+1:]...)
			}
		}
		// 每一个专业计算录取率、存入
		err = setMajorToParts(majors, s.Parts, kl, studentScore)
		if err != nil {
			return err
		}
	}
	return nil
}

func setMajorToParts(majors []*major.Major, parts map[int][]*types.Major, kelei string, studentScore int) error {
	for _, m := range majors {
		// 获取平均分
		avg, err := major_score.FindScoreAvg(m.ID, kelei)
		if err != nil {
			return err
		}
		if !avg.Valid {
			parts[0] = append(parts[0], &types.Major{
				ID:   m.ID,
				Name: m.Name,
				Rate: 0,
			})
			continue
		}
		// 计算录取率
		r := rate(studentScore, int(avg.Float64))
		//存入map
		parts[r] = append(parts[r], &types.Major{
			ID:   m.ID,
			Name: m.Name,
			Rate: r,
		})
	}
	return nil
}

func rate(studentScore int, targetScore int) int {
	c := targetScore - studentScore
	switch c {
	case 9, 10:
		return 20
	case 7, 8:
		return 25
	case 5, 6:
		return 30
	case 1, 2, 3, 4:
		return 45
	case 0:
		return 50
	case -1, -2, -3, -4, -5:
		return 55
	case -6:
		return 80
	case -7, -8:
		return 85
	case -9, -10:
		return 90
	}
	if c > 0 {
		return 1
	}
	return 95
}

func traitSort(schools []*types.School, userInfo *user.User) {
	interests := strings.Split(userInfo.Interests, " ")
	// 喜欢的专业
	like := common.SliceToMap[string](interests)
	// 性格推荐的专业
	ho := common.HollandMajorMap[userInfo.Holland]
	for _, sc := range schools {
		parts := sc.Parts
		for _, majors := range parts {
			// 为每一个专业设置权重
			for _, m := range majors {
				if like[m.Name] {
					m.Rate += 10
				} else {
					// 模糊匹配
					ok := SlicesContainsFunc[string](interests, m.Name, func(target string, s string) bool {
						return strings.Contains(s, m.Name) || strings.Contains(m.Name, s)
					})
					if ok {
						m.Rate += 10
					}
				}
				if ho[m.Name] {
					m.Rate += 8
				} else {
					hos := common.HollandMajorSlice[userInfo.Holland]
					ok := SlicesContainsFunc[string](hos, m.Name, func(target string, s string) bool {
						return strings.Contains(s, m.Name) || strings.Contains(m.Name, s)
					})
					if ok {
						m.Rate += 8
					}
				}
				if common.NationalFocus[m.Name] {
					m.Rate += 3
				} else {
					ok := SlicesContainsFunc[string](common.Focus, m.Name, func(target string, s string) bool {
						return strings.Contains(s, m.Name) || strings.Contains(m.Name, s)
					})
					if ok {
						m.Rate += 3
					}
				}
			}
			// 排序
			sort.Slice(majors, func(i int, j int) bool {
				return majors[i].Rate > majors[j].Rate
			})
		}
	}
}

func SlicesContainsFunc[T string](slice []T, target T, f func(target T, s T) bool) bool {
	for _, s := range slice {
		if f(target, s) {
			return true
		}
	}
	return false
}

func SetResp(owner string, value any) error {
	//缓存
	return common.REDIS.Set(context.Background(), owner, value, 0).Err()
}

func DelResp(owner string) error {
	return common.REDIS.Del(context.Background(), owner).Err()
}

func GetResp(owner string) ([]byte, error) {
	get := common.REDIS.Get(context.Background(), owner)
	if get.Err() != nil {
		return nil, get.Err()
	}
	bytes, err := get.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
