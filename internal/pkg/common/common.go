package common

import (
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/gookit/config/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 数据库映射
var (
	DB     *gorm.DB
	CONFIG *config.Config
	LOG    *zap.Logger
	REDIS  *redis.Client
)

var (
	ShannXi   = 61 // 61为陕西
	Page      = 1
	Mu        = &sync.Mutex{}
	Count     = (Page - 1) * 5
	T_li      = 1
	T_wen     = 2
	T_Physics = 2073
	T_History = 2074
	// 21年后本科14, 专科10
	// 21年前本科二批8,本科一批7，专科批10
	BatchIds           = []int{7, 8, 10, 14}
	TypeIdsPH          = []string{"2074", "2073"}
	TypeIdsWL          = []int{1, 2}
	HollandMajorMap    = make(map[string]map[string]bool)
	HollandMajorSlice  = make(map[string][]string)
	NationalFocus      map[string]bool
	ConventionalSlice  = []string{"会计学", "工程学", "信息技术", "图书馆", "计算机科学", "行政管理", "统计学", "数学", "法律", "数据分析"}
	InvestigativeSlice = []string{"生物学", "物理学", "计算机科学", "社会科学研究", "心理学"}
	RealisticSlice     = []string{"工程学", "农业科学", "建筑学", "机械工程", "计算机科学", "医疗技术", "航空与航天工程", "物流与供应链管理"}
	EnterprisingSlice  = []string{"商业管理", "市场营销", "金融", "管理学", "企业管理", "商务管理", "国际关系", "创业管理"}
	ArtisticSlice      = []string{"美术", "音乐", "戏剧与表演", "艺术与设计", "摄影", "文学与创意写作", "影视制作"}
	SocialSlice        = []string{"社会工作", "心理学", "教育学", "医学与护理", "公共关系", "人力资源管理", "社会学"}
	Focus              = []string{"储能科学与工程", "密码科学与技术", "生物育种科学", "未来机器人", "医工学", "柔性电子学", "生物农药科学与工程", "数据科学与大数据技术", "集成电路设计与集成系统", "人工智能", "智能制造工程", "预防医学", "古文字学", "国际组织与全球治理", "儿科学", "数学", "物理学", "化学", "天文学", "地球科学", "生物学"}
)

const (
	Dream         = "edu_choice"
	Conventional  = "conventional"
	Investigative = "investigative"
	Realistic     = "realistic"
	Enterprising  = "enterprising"
	Artistic      = "artistic"
	Social        = "social"
)

func init() {
	HollandMajorMap[Conventional] = SliceToMap[string](ConventionalSlice)
	HollandMajorMap[Investigative] = SliceToMap[string](InvestigativeSlice)
	HollandMajorMap[Realistic] = SliceToMap[string](RealisticSlice)
	HollandMajorMap[Enterprising] = SliceToMap[string](EnterprisingSlice)
	HollandMajorMap[Artistic] = SliceToMap[string](ArtisticSlice)
	HollandMajorMap[Social] = SliceToMap[string](SocialSlice)

	HollandMajorSlice[Conventional] = ConventionalSlice
	HollandMajorSlice[Investigative] = InvestigativeSlice
	HollandMajorSlice[Realistic] = RealisticSlice
	HollandMajorSlice[Enterprising] = EnterprisingSlice
	HollandMajorSlice[Artistic] = ArtisticSlice
	HollandMajorSlice[Social] = SocialSlice

	NationalFocus = SliceToMap[string](Focus)
}

var (
	ErrMysqlDuplicate = &mysql.MySQLError{
		Number:   1062,
		SQLState: [5]byte{2, 3, 0, 0, 0},
		Message:  "",
	}
)

// typeId: 2-7文科，1-7理科，2074-14历史类，2073-14物理类
func Kelei(type_id string) string {
	if type_id[1] == '-' {
		type_id = type_id[:2]
	}
	switch type_id {
	case "2-":
		return "文科"
	case "1-":
		return "理科"
	case "2074-14":
		return "历史类"
	case "2073-14":
		return "物理类"
	}
	return ""
}

func IDConvKelei(typeID int) string {
	if typeID == 2073 {
		return "物理类"
	} else if typeID == 2074 {
		return "历史类"
	}
	return ""
}

func TypeID(ph bool, hi bool) int {
	if ph {
		return 2073
	} else if hi {
		return 2074
	}
	return 0
}

func Omit(typeID int) []string {
	if typeID == 2073 {
		return []string{"文科", "理科", "历史类"}
	} else if typeID == 2074 {
		return []string{"文科", "理科", "物理类"}
	}
	return nil
}

func SliceToMap[T int | string](s []T) map[T]bool {
	m := make(map[T]bool)
	for _, item := range s {
		m[item] = true
	}
	return m
}
