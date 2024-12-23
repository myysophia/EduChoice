package major

import (
	"context"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"sort"
	"strconv"
	"time"
)

// RecommendConfig 推荐配置
type RecommendConfig struct {
	// 分数线浮动范围
	ChongRange float64 // 冲刺院校分数线浮动
	WenRange   float64 // 稳妥院校分数线浮动
	BaoRange   float64 // 保底院校分数线浮动

	// 权重配置
	ScoreWeight    float64 // 分数权重
	LocationWeight float64 // 地理位置权重
	MajorWeight    float64 // 专业匹配权重
	CostWeight     float64 // 费用权重

	// 最大推荐数量
	MaxChongNum int
	MaxWenNum   int
	MaxBaoNum   int
}

// DefaultConfig 默认配置
var DefaultConfig = RecommendConfig{
	ChongRange:  1.05,
	WenRange:    0.98,
	BaoRange:    0.95,
	MaxChongNum: 3,
	MaxWenNum:   4,
	MaxBaoNum:   3,

	ScoreWeight:    0.4, // 分数占比40%
	LocationWeight: 0.2, // 地理位置占比20%
	MajorWeight:    0.3, // 专业匹配占比30%
	CostWeight:     0.1, // 费用占比10%
}

// RecommendSchools 根据用户信息推荐学校
func RecommendSchools(u *user.User, year string, config RecommendConfig) (*RecommendResp, error) {
	// 添加查询超时设置
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// 使用 WithContext 添加超时控制
	db := common.DB.WithContext(ctx)

	// 检查必要的用户信息
	if u.Score == nil || u.ExamType == nil {
		return nil, fmt.Errorf("用户信息不完整")
	}

	// 转换年份字符串为整数
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, fmt.Errorf("年份格式错误: %v", err)
	}

	// 1. 获取候选学校列表
	query := `
		SELECT DISTINCT 
			sch.name as '学校名称', 
			sch.brief_introduction as '学校简介',
			sch.school_code as '学校代码', 
			sch.master_point as '硕士点', 
			sch.phd_point as '博士点',
			sch.title_985, 
			sch.title_211, 
			sch.region as '所属省份', 
			sch.website as '学校官网',
			sch.recruitment_phone as '学校招生电话', 
			sch.email as '学校招生邮箱',
			sch.double_first_class_disciplines as '学校一级学科',
			scores.year as '年份',
			scores.lowest as '最低校分', 
			scores.lowest_rank as '最低校位次', 
			scores.batch_name as '批次'
		FROM schools sch
		INNER JOIN scores ON sch.id = scores.school_id
		WHERE scores.year >= ? AND scores.year <= ?
		AND (
			scores.lowest BETWEEN ? AND ?
			OR scores.lowest_rank BETWEEN ? AND ?
		)
		ORDER BY ABS(scores.lowest - ?) ASC`

	// 计算分数范围
	scoreMin := float64(*u.Score) * 0.9       // 下限：分数的90%
	scoreMax := float64(*u.Score) * 1.1       // 上限：分数的110%
	rankMin := float64(*u.ProvinceRank) * 0.5 // 下限：位次的50%
	rankMax := float64(*u.ProvinceRank) * 1.5 // 上限：位次的150%

	// 添加调试日志
	fmt.Printf("查询参数:\n")
	fmt.Printf("年份范围: %d - %d\n", yearInt-3, yearInt)
	fmt.Printf("分数范围: %.2f - %.2f\n", scoreMin, scoreMax)
	fmt.Printf("位次范围: %.2f - %.2f\n", rankMin, rankMax)
	fmt.Printf("用户分数: %d\n", *u.Score)

	// 执行查询
	var mainScores []SchoolMainScore
	result := db.Raw(query,
		yearInt-3, yearInt,
		scoreMin, scoreMax,
		rankMin, rankMax,
		*u.Score,
	).Debug().Scan(&mainScores) // 添加 Debug() 查看实际执行的 SQL

	if result.Error != nil {
		return nil, fmt.Errorf("查询候选学校失败: %v", result.Error)
	}

	// 打印查询结果
	fmt.Printf("查询到 %d 所候选学校\n", len(mainScores))
	if len(mainScores) > 0 {
		fmt.Printf("第一所学校: %+v\n", mainScores[0])
	} else {
		// 如果没有结果，执行一个简单的查询来验证数据库连接
		var count int64
		common.DB.Table("schools").Count(&count)
		fmt.Printf("schools 表中共有 %d 条记录\n", count)
	}

	// 2. 获取候选学校的专业详细信息
	var schoolScores []SchoolWithScore
	for _, mainScore := range mainScores {
		// 获取学校详细信息
		schoolScore, err := GetSchoolScores(mainScore.Name, yearInt-1, *u.ExamType)
		if err != nil {
			continue // 跳过获取失败的学校
		}

		// 计算综合评分
		for _, score := range schoolScore {
			totalScore := calculateSchoolScore(score, u, config)
			schoolScores = append(schoolScores, SchoolWithScore{
				School: score,
				Total:  totalScore,
			})
		}
	}

	// 3. 按综合评分排序
	sort.Slice(schoolScores, func(i, j int) bool {
		return schoolScores[i].Total > schoolScores[j].Total
	})

	// 4. 按分数线分组
	groups := groupSchoolsByScore(schoolScores, *u.Score, config)

	// 5. 构建响应
	resp := &RecommendResp{
		ChongSchools: make([]School, 0),
		WenSchools:   make([]School, 0),
		BaoSchools:   make([]School, 0),
	}

	// 6. 转换数据格式
	for _, s := range groups.Chong {
		if school := convertToSchool(s.School); school != nil {
			resp.ChongSchools = append(resp.ChongSchools, *school)
		}
	}
	for _, s := range groups.Wen {
		if school := convertToSchool(s.School); school != nil {
			resp.WenSchools = append(resp.WenSchools, *school)
		}
	}
	for _, s := range groups.Bao {
		if school := convertToSchool(s.School); school != nil {
			resp.BaoSchools = append(resp.BaoSchools, *school)
		}
	}

	return resp, nil
}

// groupSchoolsByScore 按分数线对学校分组
func groupSchoolsByScore(scores []SchoolWithScore, userScore int, config RecommendConfig) SchoolGroups {
	var groups SchoolGroups

	// 分组
	for _, s := range scores {
		minScore, err := strconv.ParseFloat(s.School.MinScore, 64)
		if err != nil {
			continue // 如果转换失败，跳过当前学校
		}
		scoreRatio := float64(userScore) / minScore

		switch {
		case scoreRatio >= config.ChongRange && len(groups.Chong) < config.MaxChongNum:
			groups.Chong = append(groups.Chong, s)
		case scoreRatio >= config.WenRange && len(groups.Wen) < config.MaxWenNum:
			groups.Wen = append(groups.Wen, s)
		case scoreRatio >= config.BaoRange && len(groups.Bao) < config.MaxBaoNum:
			groups.Bao = append(groups.Bao, s)
		}
	}

	return groups
}

// convertToSchool 将分数信息转换为学校信息
func convertToSchool(score SchoolScore) *School {
	// 转换位次字符串为整数
	minRank := 0
	MinScore := 0
	PlanNum := 0
	if score.MinRank != "-" {
		if rank, err := strconv.Atoi(score.MinRank); err == nil {
			minRank = rank
		}
	}

	if score.MinScore != "-" {
		if rank, err := strconv.Atoi(score.MinScore); err == nil {
			MinScore = rank
		}
	}

	if score.PlanNum != "-" {
		if rank, err := strconv.Atoi(score.PlanNum); err == nil {
			PlanNum = rank
		}
	}
	school := &School{
		Name: score.SchoolName,
		HistoryInfos: map[string]HistoryInfo{
			strconv.Itoa(score.Year): {
				LowestScore:   MinScore,
				LowestRank:    minRank,
				EnrollmentNum: PlanNum,
			},
		},
		Parts: map[string][]Major{
			score.Batch: {
				{
					Name:   score.MajorName,
					Rate:   0.8,
					Weight: 0.8,
				},
			},
		},
	}
	return school
}

// 如果在其他地方也使用了 MinRank，需要添加辅助函数
func getMinRank(score SchoolScore) int {
	if score.MinRank == "-" {
		return 0
	}
	rank, err := strconv.Atoi(score.MinRank)
	if err != nil {
		return 0
	}
	return rank
}
