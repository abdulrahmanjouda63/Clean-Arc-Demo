package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

var DB *gorm.DB
var Redis *redis.Client

var Logger *zap.Logger

func InitLogger() error {
    var err error
    Logger, err = zap.NewProduction()
    return err
}
