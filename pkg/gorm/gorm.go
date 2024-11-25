package gorm

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
)

func NewGorm() (*gorm.DB, error) {
	var (
		logout io.Writer
		err    error
	)

	logout = os.Stdout
	if common.CONFIG.String("gorm_log.outType") == "file" {
		logout, err = os.OpenFile(common.CONFIG.String("gorm_log.out"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("NewGorm: Cannot open file %s: %v", common.CONFIG.String("gorm_log.out"), err)
		}
	}

	newLogger := logger.New(
		log.New(logout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	config := common.CONFIG.StringMap("mysql")
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", config["user"], config["password"], config["addr"], config["port"], config["db"], config["config"])
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(common.CONFIG.Int("mysql.maxOpenConns"))
	sqlDB.SetMaxIdleConns(common.CONFIG.Int("mysql.maxIdleConns"))
	return db, nil
}
