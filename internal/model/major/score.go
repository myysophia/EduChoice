package major

import (
	"github.com/big-dust/DreamBridge/internal/model/user"
	"math"
	"strconv"
	"strings"
)

// SchoolScore 学校评分

// calculateSchoolScore 计算学校评分
func calculateSchoolScore(school SchoolScore, user *user.User, config RecommendConfig) float64 {
	if user.Score == nil {
		return 0
	}

	// 1. 计算分数匹配度 (0-1)
	minScore, err := strconv.Atoi(school.MinScore)
	if err != nil {
		// 处理转换错误
		return 0
	}
	scoreMatch := calculateScoreMatch(minScore, *user.Score)

	// 2. 计算地理位置匹配度 (0-1)
	var locationMatch float64 = 0.5 // 默认中等匹配度
	if user.Province != nil {
		locationMatch = calculateLocationMatch(school.SchoolProvince, *user.Province)
	}

	// 3. 计算专业匹配度 (0-1)
	var majorMatch float64 = 0.5 // 默认中等匹配度
	if user.Interests != nil && user.Holland != nil {
		interests := strings.Split(*user.Interests, ",") // 将字符串转换为字符串数组
		majorMatch = calculateMajorMatch(school.MajorName, interests, *user.Holland)
	}

	// 4. 计算费用匹配度 (0-1)
	tuition, err := strconv.Atoi(school.Tuition)
	if err != nil {
		return 0
	}
	costMatch := calculateCostMatch(float64(tuition))

	// 5. 计算加权总分
	totalScore := scoreMatch*config.ScoreWeight +
		locationMatch*config.LocationWeight +
		majorMatch*config.MajorWeight +
		costMatch*config.CostWeight

	return totalScore
}

// calculateScoreMatch 计算分数匹配度
func calculateScoreMatch(schoolScore, userScore int) float64 {
	diff := math.Abs(float64(schoolScore - userScore))
	maxDiff := float64(schoolScore) * 0.2 // 允许20%的分差

	if diff > maxDiff {
		return 0
	}
	return 1 - (diff / maxDiff)
}

// calculateLocationMatch 计算地理位置匹配度
func calculateLocationMatch(schoolProvince, userProvince string) float64 {
	// 同省分数最高
	if schoolProvince == userProvince {
		return 1.0
	}

	// 相邻省份次之
	if isNeighborProvince(schoolProvince, userProvince) {
		return 0.8
	}

	// 同区域再次之
	if isSameRegion(schoolProvince, userProvince) {
		return 0.6
	}

	// 其他省份基础分
	return 0.4
}

// calculateMajorMatch 计算专业匹配度
func calculateMajorMatch(majorName string, interests []string, holland string) float64 {
	// 1. 兴趣匹配度 (30%)
	interestScore := calculateInterestMatch(majorName, interests)

	// 2. Holland码匹配度 (70%)
	hollandScore := calculateHollandMatch(majorName, holland)

	// 加权计算总分
	return interestScore*0.3 + hollandScore*0.7
}

// calculateInterestMatch 计算兴趣匹配度
func calculateInterestMatch(majorName string, interests []string) float64 {
	if len(interests) == 0 {
		return 0.5 // 未填写兴趣返回中等匹配度
	}

	// 将专业名称分词
	majorKeywords := splitMajorName(majorName)

	// 计算关键词匹配度
	var matchCount int
	for _, interest := range interests {
		for _, keyword := range majorKeywords {
			if strings.Contains(interest, keyword) || strings.Contains(keyword, interest) {
				matchCount++
				break
			}
		}
	}

	// 计算匹配度
	return float64(matchCount) / float64(len(interests))
}

// splitMajorName 专业名称分词
func splitMajorName(majorName string) []string {
	// 这里可以使用分词库，这里简单处理
	keywords := []string{
		"计算机", "软件", "人工智能", "数据", "电子", "通信",
		"机械", "土木", "建筑", "设计", "金融", "会计",
		"市场", "心理", "教育", "医学", "护理", "法学",
		"新闻", "英语",
	}

	var result []string
	for _, keyword := range keywords {
		if strings.Contains(majorName, keyword) {
			result = append(result, keyword)
		}
	}
	return result
}

// calculateCostMatch 计算费用匹配度
func calculateCostMatch(tuition float64) float64 {
	// 基于学费区间计算匹配度
	switch {
	case tuition <= 5000:
		return 1.0
	case tuition <= 8000:
		return 0.8
	case tuition <= 12000:
		return 0.6
	case tuition <= 15000:
		return 0.4
	default:
		return 0.2
	}
}
