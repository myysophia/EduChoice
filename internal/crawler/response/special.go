package response

type SpecialResponse[T ZhuanSpecialData | BenSpecialData | SpecialList] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Md5     string `json:"md5"`
}

type SpecialList []DetailInfo

type ZhuanSpecialData struct {
	//NationFeature []NationFeature `json:"nation_feature"`
	SpecialDetail ZhuanSpecialDetail `json:"special_detail"`
}
type BenSpecialData struct {
	//NationFeature []NationFeature `json:"nation_feature"`
	SpecialDetail BenSpecialDetail `json:"special_detail"`
}

type ZhuanSpecialDetail struct {
	//Num1 interface{} `json:"1"`
	Num2 SpecialList `json:"2"`
}

type BenSpecialDetail struct {
	Num1 SpecialList `json:"1"`
	//Num2 interface{} `json:"2"`
}

type DetailInfo struct {
	ID string `json:"id"`
	//SchoolID         string `json:"school_id"`
	SpecialID     string `json:"special_id"`
	NationFeature string `json:"nation_feature"`
	//ProvinceFeature  string `json:"province_feature"`
	//IsImportant      string `json:"is_important"`
	LimitYear string `json:"limit_year"`
	Year      string `json:"year"`
	//Level3Weight     string `json:"level3_weight"`
	//NationFirstClass string `json:"nation_first_class"`
	//XuekeRankScore   string `json:"xueke_rank_score"`
	//IsVideo          int    `json:"is_video"`
	SpecialName string `json:"special_name"`
	//SpecialType      string `json:"special_type"`
	TypeName   string `json:"type_name"`
	Level3Name string `json:"level3_name"`
	//Level3Code       string `json:"level3_code"`
	Level2Name string `json:"level2_name"`
	//Level2ID         string `json:"level2_id"`
	//Level2Code string `json:"level2_code"`
	//Code string `json:"code"`
}

type NationFeature struct {
	ID               string `json:"id"`
	SchoolID         string `json:"school_id"`
	SpecialID        string `json:"special_id"`
	NationFeature    string `json:"nation_feature"`
	ProvinceFeature  string `json:"province_feature"`
	IsImportant      string `json:"is_important"`
	LimitYear        string `json:"limit_year"`
	Year             string `json:"year"`
	Level3Weight     string `json:"level3_weight"`
	NationFirstClass string `json:"nation_first_class"`
	XuekeRankScore   string `json:"xueke_rank_score"`
	IsVideo          int    `json:"is_video"`
	SpecialName      string `json:"special_name"`
	SpecialType      string `json:"special_type"`
	TypeName         string `json:"type_name"`
	Level3Name       string `json:"level3_name"`
	Level3ID         string `json:"level3_id"`
	Level3Code       string `json:"level3_code"`
	Level2Name       string `json:"level2_name"`
	Level2ID         string `json:"level2_id"`
	Level2Code       string `json:"level2_code"`
	Code             string `json:"code"`
}

type HistoryRecruitResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    HistoryRecruitData `json:"data"`
	Md5     string             `json:"md5"`
	Time    string             `json:"time"`
}

// 2-7文科，1-7理科，2074-14历史类，2073-14物理类
type HistoryRecruitData map[string]map[string][]RecruitInfo

type RecruitInfo struct {
	Year         string `json:"year"`
	Type         string `json:"type"`
	Batch        string `json:"batch"`
	Name         string `json:"name"`
	Remark       string `json:"remark"`
	ElectiveInfo string `json:"elective_info"`
	Num          string `json:"num"`
}

type HistoryAdmissionResponse struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Data    HistoryAdmissionData `json:"data"`
	Md5     string               `json:"md5"`
	Time    string               `json:"time"`
}

// 2-7文科，1-7理科，2074-14历史类，2073-14物理类
type HistoryAdmissionData map[string]map[string][]AdmissionInfo

type AdmissionInfo struct {
	Year         string      `json:"year"`
	Type         string      `json:"type"`
	Batch        string      `json:"batch"`
	Name         string      `json:"name"`
	Remark       string      `json:"remark"`
	ElectiveInfo string      `json:"elective_info"`
	Min          string      `json:"min"`
	MinSection   interface{} `json:"min_section"`
}
