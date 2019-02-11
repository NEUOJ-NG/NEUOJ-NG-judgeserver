package redis

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	goRedis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	Client *goRedis.Client
)

// initialize the redis client
func InitRedisClient() {
	log.Info("connecting to redis")
	Client = goRedis.NewClient(&goRedis.Options{
		Addr:     config.GetConfig().Redis.Addr,
		Password: config.GetConfig().Redis.Password,
		DB:       config.GetConfig().Redis.DB,
	})
}
