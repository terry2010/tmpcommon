package Common

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)
import "go.uber.org/zap"

var RedisClient *redis.Client

var Config = viper.New()
var Logger, _ = zap.NewProduction()
var RunPath, _ = GetCurrentPath()

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

var Master ServerInfo

var SelfRole string
