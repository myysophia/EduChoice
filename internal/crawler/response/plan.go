package response

type PlanInfoResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Item []struct {
			//FirstKm        string `json:"first_km"`
			//Length         string `json:"length"`
			//Level2Name     string `json:"level2_name"`
			//LocalBatchName string `json:"local_batch_name"`
			//LocalTypeName  string `json:"local_type_name"`
			//Name           string `json:"name"`
			Num string `json:"num"`
			//ProvinceName   string `json:"province_name"`
			//SchoolID       string `json:"school_id"`
			//SgFxk          string `json:"sg_fxk"`
			//SgInfo         string `json:"sg_info"`
			//SgName         string `json:"sg_name"`
			//SgSxk          string `json:"sg_sxk"`
			//SgType         string `json:"sg_type"`
			//SpFxk          string `json:"sp_fxk"`
			//SpInfo         string `json:"sp_info"`
			//SpSxk          string `json:"sp_sxk"`
			//SpType         string `json:"sp_type"`
			//SpXuanke       string `json:"sp_xuanke"`
			//Spcode         string `json:"spcode"`
			//SpecialGroup   string `json:"special_group"`
			//Spname         string `json:"spname"`
			//Tuition        string `json:"tuition"`
			//Year           string `json:"year"`
		} `json:"item"`
		NumFound int `json:"numFound"`
	} `json:"data"`
	//Location  string `json:"location"`
	//Encrydata string `json:"encrydata"`
}
