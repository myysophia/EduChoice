package major

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

// SchoolScore 学校专业分数信息
type SchoolScore struct {
	SchoolName     string  `gorm:"column:学校名称"`
	MajorName      string  `gorm:"column:专业名称"`
	Province       string  `gorm:"column:招生省份"`
	Batch          string  `gorm:"column:招生批次"`
	Type           string  `gorm:"column:招生类别"`
	Year           int     `gorm:"column:年份"`
	MinScore       string  `gorm:"column:最低分"`
	MinRank        string  `gorm:"column:最低位次"`
	MaxScore       string  `gorm:"column:最高分"`
	AverageScore   float64 `gorm:"column:平均分"`
	ProScore       string  `gorm:"column:投档线"`
	PlanNum        string  `gorm:"column:计划录取人数"`
	Length         string  `gorm:"column:学制"`
	SchoolProvince string  `gorm:"column:学校所在省份"`
	Tuition        string  `gorm:"column:学费"`
}

// GetSchoolScores 获取学校分数信息
func GetSchoolScores(schoolName string, year int, studentType string) ([]SchoolScore, error) {
	var scores []SchoolScore
	query := `
        SELECT 
            sch.name AS '学校名称',
            scohis.sp_name AS '专业名称',
            scohis.local_province_name AS '招生省份',
            scohis.local_batch_name AS '招生批次',
            scohis.local_type_name AS '招生类别',
            scohis.year AS '年份',
            COALESCE(scohis.min, 0) AS '最低分',
            COALESCE(scohis.min_section, '-') AS '最低位次',
            COALESCE(scohis.max, 0) AS '最高分',
            COALESCE(scohis.average, 0) AS '平均分',
            COALESCE(scohis.proscore, 0) AS '投档线',
            COALESCE(plan.num, 0) AS '计划录取人数',
            COALESCE(plan.length, 0) AS '学制',
            plan.province_name AS '学校所在省份',
            COALESCE(plan.tuition, 0) AS '学费'
        FROM schools sch
        INNER JOIN major_score_his scohis ON sch.id = scohis.school_id
        INNER JOIN school_plan_his plan ON plan.school_id = sch.id 
            AND plan.sp_name = scohis.sp_name
            AND scohis.local_batch_name = plan.local_batch_name 
            AND scohis.local_type_name = plan.local_type_name
        WHERE scohis.year = ? AND sch.name = ?`

	// 打印完整SQL和参数
	fmt.Printf("GetSchoolScores 查询参数:\n")
	fmt.Printf("学校名称: %s\n", schoolName)
	fmt.Printf("年份: %d\n", year)
	fmt.Printf("考生类型: %s\n", studentType)

	// 先验证学校是否存在
	var schoolCount int64
	if err := common.DB.Table("schools").Where("name = ?", schoolName).Count(&schoolCount).Error; err != nil {
		return nil, fmt.Errorf("检查学校是否存在失败: %v", err)
	}
	fmt.Printf("找到 %d 所匹配的学校\n", schoolCount)

	// 执行主查询
	result := common.DB.Raw(query, year, schoolName).Debug().Scan(&scores)
	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %v", result.Error)
	}

	// 打印查询结果
	fmt.Printf("查询到 %d 条专业分数记录\n", result.RowsAffected)
	if len(scores) > 0 {
		fmt.Printf("第一条记录: %+v\n", scores[0])
	} else {
		// 如果没有结果，检查中间表数据
		var scoreCount int64
		common.DB.Table("major_score_his").
			Where("year = ? AND school_id IN (SELECT id FROM schools WHERE name = ?)",
				year, schoolName).
			Count(&scoreCount)
		fmt.Printf("major_score_his 表中找到 %d 条相关记录\n", scoreCount)

		var planCount int64
		common.DB.Table("school_plan_his").
			Where("school_id IN (SELECT id FROM schools WHERE name = ?)",
				schoolName).
			Count(&planCount)
		fmt.Printf("school_plan_his 表中找到 %d 条相关记录\n", planCount)
	}

	return scores, nil
}
