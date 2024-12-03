package response

type PlanHisResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      PlanHisData `json:"data"`
	Location  string      `json:"location"`
	Encrydata string      `json:"encrydata"`
}

type PlanHisData struct {
	Item     []PlanHisItem `json:"item"`
	NumFound int           `json:"numFound"`
}

type PlanHisItem struct {
	Length         string `json:"length"`
	Num            string `json:"num"`
	LocalBatchName string `json:"local_batch_name"`
	LocalTypeName  string `json:"local_type_name"`
	ProvinceName   string `json:"province_name"`
	SchoolID       string `json:"school_id"`
	SpName         string `json:"sp_name"`
	SpecialGroup   string `json:"special_group"`
	Spname         string `json:"spname"`
	Tuition        string `json:"tuition"`
	Year           string `json:"year"`
}
