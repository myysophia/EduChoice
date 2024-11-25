package api

import (
	"github.com/big-dust/DreamBridge/internal/api/router"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/config"
	"github.com/big-dust/DreamBridge/pkg/gorm"
	"github.com/big-dust/DreamBridge/pkg/redis"
	"github.com/big-dust/DreamBridge/pkg/zap"
	"github.com/gin-gonic/gin"
)

func Run() {
	common.CONFIG = config.New("./config/config.toml")
	common.LOG = zap.AddZap()
	rd, err := redis.NewRedisClient()
	if err != nil {
		common.LOG.Panic("Connect Redis Error: " + err.Error())
	}
	common.REDIS = rd
	db, err := gorm.NewGorm()
	if err != nil {
		common.LOG.Panic("Open DB Error: " + err.Error())
	}
	common.DB = db
	e := gin.Default()
	router.Load(e)
	addr := common.CONFIG.String("server.host") + ":" + common.CONFIG.String("server.port")
	if err = e.Run(addr); err != nil {
		common.LOG.Panic("Run Error: " + err.Error())
	}
}
