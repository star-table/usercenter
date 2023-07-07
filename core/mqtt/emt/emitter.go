package emt

import (
	"errors"
	"strconv"
	"sync"
	"time"

	emitter "github.com/emitter-io/go/v2"
	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/lock"
)

var (
	mqttMutex     sync.Mutex
	disConnectErr                   = errors.New("mqtt disconnected")
	clients       []*emitter.Client = nil
)

func GetClient() (*emitter.Client, error) {
	client, selector, err := Connect(nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if client == nil {
		return nil, disConnectErr
	}

	//fmt.Printf("selector %d\n", selector)

	//断开连接，重试一次
	if !client.IsConnected() {
		logger.InfoF("连接断开，开始重连...")

		key := strconv.Itoa(selector)
		lock.Lock(key)
		defer lock.Unlock(key)
		client = clients[selector]
		if client == nil || !client.IsConnected() {
			clients[selector] = nil
			client, _, err = Connect(nil)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			if !client.IsConnected() {
				return nil, disConnectErr
			}
		}
	}
	return client, nil
}

func Connect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, int, error) {
	initPool()

	selector := int(time.Now().Unix() & int64(len(clients)-1))
	selectClient := clients[selector]

	if selectClient == nil {
		mqttMutex.Lock()
		defer mqttMutex.Unlock()

		selectClient = clients[selector]
		if selectClient == nil {
			mqttConfig := conf.Cfg.MQTT
			if mqttConfig == nil {
				panic("missing mqtt config.")
			}

			logger.InfoF("mqtt config %s", json.ToJsonIgnoreError(mqttConfig))

			if mqttConfig.Enable {
				var err error
				clients[selector], err = emitter.Connect(mqttConfig.Address, handler, options...)
				if err != nil {
					logger.Error(err)
					return nil, selector, err
				}
			} else {
				logger.Error("mqtt已被禁用")
			}
		}

	}
	return clients[selector], selector, nil
}

func initPool() {
	mqttConfig := conf.Cfg.MQTT
	if mqttConfig == nil {
		panic("missing mqtt config.")
	}

	if clients == nil {
		mqttMutex.Lock()
		defer mqttMutex.Unlock()
		if clients == nil {
			poolSize := mqttConfig.ConnectPoolSize
			if poolSize <= 0 {
				poolSize = 10
			}
			clients = make([]*emitter.Client, poolSize)
		}
	}
}

func GetNativeConnect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, error) {
	mqttConfig := conf.Cfg.MQTT
	if mqttConfig == nil {
		return nil, errors.New("missing mqtt config.")
	}

	logger.InfoF("mqtt config %s", json.ToJsonIgnoreError(mqttConfig))

	mc, err := emitter.Connect(mqttConfig.Address, handler, options...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return mc, nil
}

func GenerateKey(channel string, permissions string, ttl int) (string, error) {
	client, err := GetClient()
	if err != nil {
		logger.Error(err)
		return "", err
	}
	mqttConfig := conf.Cfg.MQTT
	if mqttConfig == nil {
		panic("missing mqtt config.")
	}
	secretKey := mqttConfig.SecretKey

	key, err := client.GenerateKey(secretKey, channel, permissions, ttl)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return key, nil
}

func Publish(key, channel string, payload interface{}, handler emitter.ErrorHandler) error {
	client, err := GetClient()
	if err != nil {
		logger.Error(err)
		return err
	}
	if handler != nil {
		client.OnError(handler)
	}
	return client.Publish(key, channel, payload, emitter.WithAtLeastOnce())
}
