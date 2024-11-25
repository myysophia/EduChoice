package types

type LoginReq struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"email"`
	Code     string `json:"code" form:"code" binding:"len=6"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type SendEmailVerificationCodeReq struct {
	Email string `json:"email" form:"email" binding:"email"`
}

type UserSetInfoReq struct {
	Province     string   `json:"province" form:"province"`
	ExamType     string   `json:"exam_type" form:"exam_type"`
	SchoolType   string   `json:"school_type" form:"school_type"`
	Subject      Subject  `json:"subject" form:"subject"`
	Score        int      `json:"score" form:"score"`
	ProvinceRank int      `json:"province_rank" form:"province_rank"`
	Holland      string   `json:"holland" form:"holland"`
	Interests    []string `json:"interests" form:"interests"`
}

type Subject struct {
	Physics   bool `json:"physics" form:"physics"`
	History   bool `json:"history" form:"history"`
	Chemistry bool `json:"chemistry" form:"chemistry"`
	Biology   bool `json:"biology" form:"biology"`
	Geography bool `json:"geography" form:"geography"`
	Politics  bool `json:"politics" form:"politics"`
}

type ZYMockResp struct {
	ChongSchools []*School `json:"chong_schools" form:"chong_schools"`
	WenSchools   []*School `json:"wen_schools" form:"wen_schools"`
	BaoSchools   []*School `json:"bao_schools" form:"bao_schools"`
}

type School struct {
	ID           int                  `json:"id" form:"id"`
	Name         string               `json:"name" form:"name"`
	Parts        map[int][]*Major     `json:"parts" form:"parts"`
	HistoryInfos map[int]*HistoryInfo `json:"history_infos" form:"history_infos"`
}

type Major struct {
	ID     int    `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	Rate   int    `json:"rate" form:"rate"`
	Weight int    `json:"weight" form:"weight"`
}

type HistoryInfo struct {
	//Year          int `json:"year" form:"year"`
	LowestScore   int `json:"lowest_score" form:"lowest_score"`
	LowestRank    int `json:"lowest_rank" form:"lowest_rank"`
	EnrollmentNum int `json:"enrollment_num" form:"enrollment_num"`
}
