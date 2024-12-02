package response

// SpecialScoresHisResponse 结构体
type SpecialScoresHisResponse struct {
	Code      string               `json:"code"`
	Message   string               `json:"message"`
	Data      SpecialScoresHisData `json:"data"`
	Location  string               `json:"location"`
	Encrydata string               `json:"encrydata"`
}

type SpecialScoresHisData struct {
	Item     []SpecialScoresHisItem `json:"item"`
	NumFound int                    `json:"numFound"`
}

type SpecialScoresHisItem struct {
	Average           string `json:"average"`
	Doublehigh        int    `json:"doublehigh"`
	DualClassName     string `json:"dual_class_name"`
	FirstKm           int    `json:"first_km"`
	ID                string `json:"id"`
	Info              string `json:"info"`
	IsScoreRange      int    `json:"is_score_range"`
	IsTop             int    `json:"is_top"`
	Level2Name        string `json:"level2_name"`
	Level3Name        string `json:"level3_name"`
	LocalBatchName    string `json:"local_batch_name"`
	LocalProvinceName string `json:"local_province_name"`
	LocalTypeName     string `json:"local_type_name"`
	Max               string `json:"max"`
	Min               string `json:"min"`
	MinRange          string `json:"min_range"`
	MinRankRange      string `json:"min_rank_range"`
	MinSection        int    `json:"min_section"`
	Name              string `json:"name"`
	Proscore          int    `json:"proscore"`
	Remark            string `json:"remark"`
	SchoolID          int    `json:"school_id"`
	SgFxk             string `json:"sg_fxk"`
	SgInfo            string `json:"sg_info"`
	SgName            string `json:"sg_name"`
	SgSxk             string `json:"sg_sxk"`
	SgType            int    `json:"sg_type"`
	Single            string `json:"single"`
	SpFxk             string `json:"sp_fxk"`
	SpInfo            string `json:"sp_info"`
	SpName            string `json:"sp_name"`
	SpSxk             string `json:"sp_sxk"`
	SpType            int    `json:"sp_type"`
	SpeID             int    `json:"spe_id"`
	SpecialGroup      int    `json:"special_group"`
	SpecialID         int    `json:"special_id"`
	Spname            string `json:"spname"`
	Year              int    `json:"year"`
	ZslxName          string `json:"zslx_name"`
}
