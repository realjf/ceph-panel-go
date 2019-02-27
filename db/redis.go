package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"goframe/config"
	"goframe/exception"
	"goframe/utils"
	"goframe/middleware"
)

var RedisClient *redis.Client

type RedisDriver struct {
	Host string
	Port string
}

func NewRedis(config config.IConfig) *RedisDriver {
	configData := config.GetConfigData()
	return &RedisDriver{
		Host: configData.Redis.Host,
		Port: utils.ToString(configData.Redis.Port),
	}
}

func (r *RedisDriver) Init() {
	if r.Host == "" || r.Port == "" {
		exception.CheckError(exception.NewError("redis config is error"), 4001)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		exception.CheckError(exception.NewError(err.Error()), 4002)
	}

	middleware.Logger.Logger.Info("init redis...")
}
