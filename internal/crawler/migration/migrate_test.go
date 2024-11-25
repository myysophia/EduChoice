package migration

import (
	"github.com/big-dust/DreamBridge/internal/crawler/response"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/config"
	"github.com/big-dust/DreamBridge/pkg/gorm"
	"github.com/big-dust/DreamBridge/pkg/zap"
	"testing"
)

func TestMigrateSchoolScores(t *testing.T) {
	// 初始化配置
	common.CONFIG = config.New("./config/config.toml")
	// 日志配置
	common.LOG = zap.AddZap()
	// 连接数据库
	DB, err := gorm.NewGorm()
	if err != nil {
		panic("gorm:" + err.Error())
	}
	common.DB = DB
	item := &response.Item{
		SchoolID:      2561,
		Name:          "西安工商学院",
		CodeEnroll:    "1368200",
		CityName:      "西安市",
		DualClassName: "",
		F211:          0,
		F985:          0,
		Level:         "普通本科",
	}
	MigrateSchoolScoresOneSafe(1, *item)
}

func TestMigrateSpecialScoresOneSafe(t *testing.T) {
	// 初始化配置
	common.CONFIG = config.New("./config/config.toml")
	// 日志配置
	common.LOG = zap.AddZap()
	// 连接数据库
	DB, err := gorm.NewGorm()
	if err != nil {
		panic("gorm:" + err.Error())
	}
	common.DB = DB
	MigrateSpecialScoresOneSafe(602)
}
