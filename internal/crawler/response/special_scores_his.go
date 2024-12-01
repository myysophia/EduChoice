package response

// SpecialScoresHisResponse 专业历史分数响应结构
type SpecialScoresHisResponse struct {
	Code      string               `json:"code"`      // 响应码
	Message   string               `json:"message"`   // 响应消息
	Data      SpecialScoresHisData `json:"data"`      // 数据部分
	Location  string               `json:"location"`  // 位置信息
	Encrydata string               `json:"encrydata"` // 加密数据
}

type SpecialScoresHisData struct {
	Item      []SpecialScoresHisItem `json:"item"`      // 分数数据列表
	NumFound  int                    `json:"numFound"`  // 总记录数
	PageCount int                    `json:"pageCount"` // 总页数
	PageSize  int                    `json:"pageSize"`  // 每页大小
	Start     int                    `json:"start"`     // 起始位置
}

// SpecialScoresHisItem 专业历史分数数据项
type SpecialScoresHisItem struct {
	ID                string `json:"id"`                  // 记录ID，如 gkspecialscore2023295809
	SchoolID          int    `json:"school_id"`           // 学校ID
	SpecialID         int    `json:"special_id"`          // 专业ID
	SpeID             int    `json:"spe_id"`              // 特殊专业ID
	Year              int    `json:"year"`                // 年份
	SpName            string `json:"sp_name"`             // 专业名称
	SpInfo            string `json:"sp_info"`             // 专业信息
	Info              string `json:"info"`                // 额外信息
	LocalProvinceName string `json:"local_province_name"` // 省份
	LocalTypeName     string `json:"local_type_name"`     // 科类
	LocalBatchName    string `json:"local_batch_name"`    // 批次
	Level2Name        string `json:"level2_name"`         // 二级学科名称
	Level3Name        string `json:"level3_name"`         // 三级学科名称
	Average           int    `json:"average"`             // 平均分
	Max               int    `json:"max"`                 // 最高分
	Min               int    `json:"min"`                 // 最低分
	MinSection        int    `json:"min_section"`         // 最低位次
	Proscore          int    `json:"proscore"`            // 投档分
	DoubleHigh        int    `json:"doublehigh"`          // 是否双高
	IsTop             int    `json:"is_top"`              // 是否重点
	IsScoreRange      int    `json:"is_score_range"`      // 分数范围标识
	MinRange          string `json:"min_range"`           // 最低分范围
	MinRankRange      string `json:"min_rank_range"`      // 最低排名范围
	Remark            string `json:"remark"`              // 备注
	ZslxName          string `json:"zslx_name"`           // 招生类型名称
	DualClassName     string `json:"dual_class_name"`     // 双学位类名称
	FirstKm           int    `json:"first_km"`            // 首选科目
	SgFxk             string `json:"sg_fxk"`              // 选考科目要求
	SgSxk             string `json:"sg_sxk"`              // 首选科目要求
	SgInfo            string `json:"sg_info"`             // 选考科目信息
	SgName            string `json:"sg_name"`             // 选考科目组名称
	SgType            int    `json:"sg_type"`             // 选考科目类型
	SpFxk             string `json:"sp_fxk"`              // 专业选考科目要求
	SpSxk             string `json:"sp_sxk"`              // 专业首选科目要求
	SpType            int    `json:"sp_type"`             // 专业类型
	Single            string `json:"single"`              // 单科要求
	SpecialGroup      int    `json:"special_group"`       // 专业组
}
