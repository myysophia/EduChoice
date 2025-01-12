package handler

import (
	"github.com/big-dust/DreamBridge/internal/api/types"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/big-dust/DreamBridge/internal/api/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

// GetSchoolCards 获取学校卡片信息
func GetSchoolCards(c *gin.Context) {
	var req types.SchoolCardReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, "参数错误")
		return
	}

	// 构建查询SQL
	query := `SELECT DISTINCT
        sch.name as school_name,
        sch.brief_introduction,
        sch.school_code,
        sch.master_point,
        sch.phd_point,
        sch.title_985,
        sch.title_211,
        sch.region,
        sch.website,
        sch.recruitment_phone,
        sch.email,
        sch.double_first_class_disciplines,
        scores.tag,
        scores.year,
        scores.lowest,
        scores.lowest_rank,
        scores.batch_name,
        CASE WHEN scores.type_id = 1 THEN '理科'
             WHEN scores.type_id = 2 THEN '文科'
             ELSE '理科' END as subject_type
    FROM schools sch
    INNER JOIN scores ON sch.id = scores.school_id
    WHERE 1=1
        AND sch.id IN (?)
        AND scores.type_id = 1
        AND scores.tag = "普通类"
        AND scores.batch_name LIKE "%本科%"
    ORDER BY scores.lowest ASC`

	var results []struct {
		SchoolName            string `db:"school_name"`
		BriefIntroduction     string `db:"brief_introduction"`
		SchoolCode            string `db:"school_code"`
		MasterPoint           int    `db:"master_point"`
		PhdPoint              int    `db:"phd_point"`
		Title985              bool   `db:"title_985"`
		Title211              bool   `db:"title_211"`
		Region                string `db:"region"`
		Website               string `db:"website"`
		RecruitmentPhone      string `db:"recruitment_phone"`
		Email                 string `db:"email"`
		FirstClassDisciplines string `db:"double_first_class_disciplines"`
		Tag                   string `db:"tag"`
		Year                  int    `db:"year"`
		Lowest                int    `db:"lowest"`
		LowestRank            int    `db:"lowest_rank"`
		BatchName             string `db:"batch_name"`
		SubjectType           string `db:"subject_type"`
	}

	if err := common.DB.Raw(query, req.SchoolIds).Scan(&results).Error; err != nil {
		response.Error(c, "查询失败")
		return
	}

	// 转换为响应格式
	var resp []types.SchoolCardResp
	for _, r := range results {
		disciplines := strings.Split(r.FirstClassDisciplines, ",")
		resp = append(resp, types.SchoolCardResp{
			Name:                  r.SchoolName,
			BriefIntroduction:     r.BriefIntroduction,
			SchoolCode:            r.SchoolCode,
			MasterPoint:           r.MasterPoint,
			PhdPoint:              r.PhdPoint,
			Is985:                 r.Title985,
			Is211:                 r.Title211,
			Region:                r.Region,
			Website:               r.Website,
			RecruitmentPhone:      r.RecruitmentPhone,
			Email:                 r.Email,
			FirstClassDisciplines: disciplines,
			Tag:                   r.Tag,
			Year:                  r.Year,
			LowestScore:           r.Lowest,
			LowestRank:            r.LowestRank,
			Batch:                 r.BatchName,
			SubjectType:           r.SubjectType,
		})
	}

	response.Success(c, resp)
}
