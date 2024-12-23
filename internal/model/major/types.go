package major

// MatchResult 匹配结果
type MatchResult struct {
	Score       float64  `json:"score"`       // 总分
	Factors     []Factor `json:"factors"`     // 影响因素
	Suggestions []string `json:"suggestions"` // 建议
}

// Factor 影响因素
type Factor struct {
	Name     string  `json:"name"`     // 因素名称
	Score    float64 `json:"score"`    // 得分
	Weight   float64 `json:"weight"`   // 权重
	Analysis string  `json:"analysis"` // 分析说明
}

// RecommendResp 推荐响应
type RecommendResp struct {
	ChongSchools []School `json:"chong_schools"` // 冲刺院校
	WenSchools   []School `json:"wen_schools"`   // 稳妥院校
	BaoSchools   []School `json:"bao_schools"`   // 保底院校
}

// School 学校信息
type School struct {
	Name         string                  `json:"name"`           // 学校名称
	HistoryInfos map[string]HistoryInfo `json:"history_infos"` // 历年录取信息
	Parts        map[string][]Major     `json:"parts"`         // 各批次专业
}

// SchoolWithScore 带评分的学校信息
type SchoolWithScore struct {
	School SchoolScore // 学校基础信息
	Total  float64    // 综合评分
}

// SchoolGroups 学校分组
type SchoolGroups struct {
	Chong []SchoolWithScore // 冲刺院校
	Wen   []SchoolWithScore // 稳妥院校
	Bao   []SchoolWithScore // 保底院校
}

// SchoolMainScore 学校基础信息和分数
type SchoolMainScore struct {
	Name                    string  `gorm:"column:学校名称"`
	BriefIntroduction      string  `gorm:"column:学校简介"`
	SchoolCode             string  `gorm:"column:学校代码"`
	MasterPoint            int     `gorm:"column:硕士点"`
	PhdPoint               int     `gorm:"column:博士点"`
	Title985               bool    `gorm:"column:title_985"`
	Title211               bool    `gorm:"column:title_211"`
	Region                 string  `gorm:"column:所属省份"`
	Website                string  `gorm:"column:学校官网"`
	RecruitmentPhone       string  `gorm:"column:学校招生电话"`
	Email                  string  `gorm:"column:学校招生邮箱"`
	DoubleFirstDisciplines string  `gorm:"column:学校一级学科"`
	Year                   int     `gorm:"column:年份"`
	LowestScore           int     `gorm:"column:最低校分"`
	LowestRank            int     `gorm:"column:最低校位次"`
	BatchName             string  `gorm:"column:批次"`
}
