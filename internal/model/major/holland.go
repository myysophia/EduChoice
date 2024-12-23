package major

import "sort"

// HollandCode 霍兰德类型代码
type HollandCode string

const (
	Realistic     HollandCode = "R" // 现实型
	Investigative HollandCode = "I" // 研究型
	Artistic      HollandCode = "A" // 艺术型
	Social        HollandCode = "S" // 社会型
	Enterprising  HollandCode = "E" // 企业型
	Conventional  HollandCode = "C" // 常规型
)

// HollandMajorMap 专业-霍兰德类型映射
//var HollandMajorMap = map[string][]HollandCode{
//	"计算机科学与技术": {"I", "R", "C"},
//	"软件工程":      {"I", "R", "E"},
//	"人工智能":      {"I", "R"},
//	"数据科学":      {"I", "C"},
//	"电子信息工程":    {"R", "I"},
//	"通信工程":      {"R", "I"},
//	"机械工程":      {"R", "I"},
//	"土木工程":      {"R", "C"},
//	"建筑学":       {"A", "R"},
//	"工业设计":      {"A", "R"},
//	"金融学":       {"E", "C"},
//	"会计学":       {"C", "E"},
//	"市场营销":      {"E", "S"},
//	"心理学":       {"S", "I"},
//	"教育学":       {"S", "A"},
//	"医学":        {"I", "S"},
//	"护理学":       {"S", "R"},
//	"法学":        {"E", "S"},
//	"新闻学":       {"A", "S"},
//	"英语":        {"S", "A"},
//	// 可以继续添加更多专业
//}

// calculateHollandMatch 计算霍兰德类型匹配度
func calculateHollandMatch(majorName string, userHolland string) float64 {
	// 获取专业对应的霍兰德类型
	majorCodes, exists := HollandMajorMap[majorName]
	if !exists {
		return 0.5 // 未知专业返回中等匹配度
	}

	// 用户的霍兰德类型
	userCodes := parseHollandCode(userHolland)
	if len(userCodes) == 0 {
		return 0.5 // 用户未填写霍兰德类型
	}

	// 计算匹配度
	var totalScore float64
	var maxScore float64

	for i, uCode := range userCodes {
		weight := 1.0 - float64(i)*0.2 // 权重递减：1.0, 0.8, 0.6
		if weight < 0.2 {
			break // 只考虑前三个类型
		}

		for j, mCode := range majorCodes {
			mWeight := 1.0 - float64(j)*0.2
			if mWeight < 0.2 {
				break
			}

			if uCode == mCode {
				totalScore += weight * mWeight
			}
		}
		maxScore += weight
	}

	if maxScore == 0 {
		return 0.5
	}
	return totalScore / maxScore
}

// parseHollandCode 解析霍兰德代码字符串
func parseHollandCode(code string) []HollandCode {
	var result []HollandCode
	for _, c := range code {
		switch HollandCode(string(c)) {
		case Realistic, Investigative, Artistic, Social, Enterprising, Conventional:
			result = append(result, HollandCode(string(c)))
		}
	}
	return result
}

// GetMajorsByHolland 根据霍兰德类型获取推荐专业
func GetMajorsByHolland(hollandCode string) []string {
	var recommendMajors []string
	userCodes := parseHollandCode(hollandCode)
	if len(userCodes) == 0 {
		return recommendMajors
	}

	// 计算每个专业的匹配度
	type majorScore struct {
		name  string
		score float64
	}
	var scores []majorScore

	for major := range HollandMajorMap {
		score := calculateHollandMatch(major, hollandCode)
		if score >= 0.6 { // 只推荐匹配度大于60%的专业
			scores = append(scores, majorScore{
				name:  major,
				score: score,
			})
		}
	}

	// 按匹配度排序
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// 取前10个专业
	for i := 0; i < len(scores) && i < 10; i++ {
		recommendMajors = append(recommendMajors, scores[i].name)
	}

	return recommendMajors
}

// GetHollandDescription 获取霍兰德类型描述
func GetHollandDescription(code HollandCode) string {
	descriptions := map[HollandCode]string{
		Realistic:     "现实型(R)：偏好具体的、有序的、系统性的操作活动，喜欢与机器、工具、物品等打交道。",
		Investigative: "研究型(I)：偏好研究性、分析性的工作，喜欢观察、学习、调查、分析、评估或解决问题。",
		Artistic:      "艺术型(A)：偏好艺术性、创造性的工作，喜欢表达自我、追求个性化和创新。",
		Social:        "社会型(S)：偏好与人交往的工作，喜欢帮助、教育、咨询、服务他人。",
		Enterprising:  "企业型(E)：偏好影响他人、说服他人的工作，喜欢领导、管理、经营等活动。",
		Conventional:  "常规型(C)：偏好有条理、规范性的工作，喜欢按规则办事，做事细心、有计划。",
	}
	return descriptions[code]
}
