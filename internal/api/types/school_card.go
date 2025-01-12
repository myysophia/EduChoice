package types

// SchoolCardReq 学校卡片请求
type SchoolCardReq struct {
    SchoolIds []string `json:"school_ids" form:"school_ids"` // 学校ID列表
}

// SchoolCardResp 学校卡片响应
type SchoolCardResp struct {
    Name                    string   `json:"name"`                     // 学校名称
    BriefIntroduction      string   `json:"brief_introduction"`       // 学校简介
    SchoolCode             string   `json:"school_code"`              // 学校代码
    MasterPoint            int      `json:"master_point"`             // 硕士点
    PhdPoint              int      `json:"phd_point"`                // 博士点
    Is985                  bool     `json:"is_985"`                   // 是否985
    Is211                  bool     `json:"is_211"`                   // 是否211
    Region                 string   `json:"region"`                   // 城市名
    Website                string   `json:"website"`                  // 学校官网
    RecruitmentPhone       string   `json:"recruitment_phone"`        // 学校招生电话
    Email                  string   `json:"email"`                    // 学校招生邮箱
    FirstClassDisciplines  []string `json:"first_class_disciplines"`  // 一级学科
    Tag                    string   `json:"tag"`                      // 招生类型
    Year                   int      `json:"year"`                     // 年份
    LowestScore           int      `json:"lowest_score"`             // 最低分
    LowestRank            int      `json:"lowest_rank"`              // 最低位次
    Batch                  string   `json:"batch"`                    // 批次
    SubjectType           string   `json:"subject_type"`             // 科别
} 