package school

import (
	"fmt"
	"time"

	"github.com/big-dust/DreamBridge/internal/model/school_score"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

// 定义 Schools 表格的模型
type School struct {
	ID                          int
	Name                        string `gorm:"unique;not null"`
	BriefIntroduction           string
	SchoolCode                  string `gorm:"unique"`
	MasterPoint                 int
	PhdPoint                    int
	ResearchProject             int
	TitleDoubleFirstClass       bool
	Title_985                   bool
	Title_211                   bool
	TitleCollege                bool
	TitleUndergraduate          bool
	Region                      string
	Website                     string
	RecruitmentPhone            string
	Email                       string
	PromotionRate               string
	AbroadRate                  string
	EmploymentRate              string
	DoubleFirstClassDisciplines string
}

func FindIDsByLevelIn(IDs []int, isCollege bool) ([]int, error) {
	var ids []int
	if err := common.DB.Table("schools").Select("id").Find(&ids, "id in (?) and title_college = ?", IDs, isCollege).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func FindOne(id int) (*School, error) {
	school := &School{}
	if err := common.DB.First(school, id).Error; err != nil {
		return nil, err
	}
	return school, nil
}

// 查询记录
func GetSchoolIdList() ([]int, error) {
	var schoolIdList []int
	if err := common.DB.Model(&School{}).Select("id").
		//Where("id in (48,140,300,572)").
		Find(&schoolIdList).Error; err != nil {
		return nil, err
	}
	return schoolIdList, nil
}

func CreateSchoolScore(school *School, scores map[string]*school_score.Score) error {
	tx := common.DB.Begin()

	if err := tx.Create(school).Error; err != nil {
		common.LOG.Error("CreateSchoolScore: " + err.Error())
		tx.Rollback()
		if common.ErrMysqlDuplicate.Is(err) {
			return nil
		}
		return err
	}
	for _, score := range scores {
		if err := tx.Create(score).Error; err != nil {
			tx.Rollback()
			common.LOG.Error("CreateSchoolScore: " + err.Error())
			return err
		}
	}
	tx.Commit()
	return nil
}

func MustCreateSchoolScore(school *School, scores map[string]*school_score.Score) {
	tryCount := 0
this:
	for {
		errChan := make(chan error, 1)
		nilChan := make(chan error, 1)
		go func() {
			err := CreateSchoolScore(school, scores)
			if err != nil {
				errChan <- err
				return
			}
			nilChan <- nil
		}()
		ticker := time.NewTicker(10 * time.Second)
		select {
		case <-nilChan:
			break this
		case <-ticker.C:
			select {
			case err := <-errChan:
				common.LOG.Error("MustCreateSchoolScore: err:" + err.Error())
			default:
			}
			tryCount++
			common.LOG.Error(fmt.Sprintf("MustCreateSchoolScore:	Time Out 10s TryCount:%d schoolName:", tryCount, school.Name))
		}
	}
}
