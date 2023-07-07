package redis

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

func newRedisPool(conf Config) *redis.Pool {
	maxIdle := 10
	if conf.MaxIdle > 0 {
		maxIdle = conf.MaxIdle
	}

	maxActive := 10
	if conf.MaxActive > 0 {
		maxActive = conf.MaxActive
	}

	maxIdleTimeout := 60
	if conf.MaxIdleTimeout > 0 {
		maxIdleTimeout = conf.MaxIdleTimeout
	}

	timeout := time.Duration(5)

	if conf.IsSentinel {
		sntnl := &Sentinel{
			Addrs:      []string{conf.Host + ":" + strconv.Itoa(conf.Port)},
			MasterName: conf.MasterName,
			Dial: func(addr string) (redis.Conn, error) {
				timeout := 500 * time.Millisecond
				c, err := redis.Dial("tcp", addr,
					redis.DialPassword(conf.SentinelPwd),
					redis.DialDatabase(conf.Database),
					redis.DialConnectTimeout(timeout*time.Second),
					redis.DialReadTimeout(timeout*time.Second),
					redis.DialWriteTimeout(timeout*time.Second))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		return &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				masterAddr, err := sntnl.MasterAddr()
				if err != nil {
					return nil, err
				}
				c, err := redis.Dial("tcp", masterAddr,
					redis.DialPassword(conf.Pwd),
					redis.DialDatabase(conf.Database),
					redis.DialConnectTimeout(timeout*time.Second),
					redis.DialReadTimeout(timeout*time.Second),
					redis.DialWriteTimeout(timeout*time.Second))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
	}

	// 建立连接池
	redisClient := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", conf.Host+":"+strconv.Itoa(conf.Port),
				redis.DialPassword(conf.Pwd),
				redis.DialDatabase(conf.Database),
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	return redisClient
}
