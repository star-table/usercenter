package store

import (
	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/store/redis"
)

var (
	Redis *redis.Proxy
	Mysql *mysql.Client
)

func init() {
	redisConfig := conf.Cfg.Redis
	mysqlConfig := conf.Cfg.Mysql

	Redis = redis.NewProxy(redis.Config{
		Host:           redisConfig.Host,
		Port:           redisConfig.Port,
		Pwd:            redisConfig.Pwd,
		Database:       redisConfig.Database,
		MaxIdle:        redisConfig.MaxIdle,
		MaxActive:      redisConfig.MaxActive,
		MaxIdleTimeout: redisConfig.MaxIdleTimeout,
		IsSentinel:     redisConfig.IsSentinel,
		MasterName:     redisConfig.MasterName,
		SentinelPwd:    redisConfig.SentinelPwd,
	})

	Mysql = mysql.NewClient(mysql.Config{
		Host:         mysqlConfig.Host,
		Port:         mysqlConfig.Port,
		User:         mysqlConfig.User,
		Pwd:          mysqlConfig.Pwd,
		Database:     mysqlConfig.Database,
		MaxIdleConns: mysqlConfig.MaxIdleConns,
		MaxOpenConns: mysqlConfig.MaxOpenConns,
		MaxLifetime:  mysqlConfig.MaxLifetime,
	})
}
