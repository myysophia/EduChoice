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
	ID                string `json:"id"`
	Info              string `json:"info"`
	IsScoreRange      int    `json:"is_score_range"`
	IsTop             int    `json:"is_top"`
	LocalBatchName    string `json:"local_batch_name"`
	LocalProvinceName string `json:"local_province_name"`
	LocalTypeName     string `json:"local_type_name"`
	Max               string `json:"max"`
	Min               string `json:"min"`
	MinRange          string `json:"min_range"`
	MinRankRange      string `json:"min_rank_range"`
	MinSection        string `json:"min_section"`
	Proscore          int    `json:"proscore"`
	Remark            string `json:"remark"`
	SchoolID          int    `json:"school_id"`
	SpName            string `json:"sp_name"`
	SpType            int    `json:"sp_type"`
	SpeID             int    `json:"spe_id"`
	SpecialID         int    `json:"special_id"`
	Spname            string `json:"spname"`
	Year              int    `json:"year"`
}
