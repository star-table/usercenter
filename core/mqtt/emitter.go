package mqtt

import (
	emitter "github.com/emitter-io/go/v2"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/mqtt/emt"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/model/bo"
)

func onError(c *emitter.Client, e emitter.Error) {
	logger.Error(e.Message)
	if e.Status == 401 {
		logger.InfoF("mqtt root key失效 %s, 准备清理", e.Message)
		cacheErr := ClearRootKey()
		if cacheErr != nil {
			logger.Error(cacheErr)
		}
	}
}

func Publish(channel string, payload interface{}) errs.SystemErrorInfo {
	logger.InfoF("MQTT推送的channel %s", channel)

	key, err := GetRootKey()
	if err != nil {
		logger.Error(err)
		return err
	}
	mqttErr := emt.Publish(key, channel, payload, onError)
	if mqttErr != nil {
		logger.Error(mqttErr)
		return errs.MQTTPublishError
	}
	return nil
}

func GetRootKey() (string, errs.SystemErrorInfo) {
	key, redisErr := store.Redis.Get(consts.CacheMQTTRootKey)
	if redisErr != nil {
		logger.Error(redisErr)
		return "", errs.RedisOperateError
	}
	if key == "" {
		return GetRootNewKey()
	}
	return key, nil
}

func GetRootNewKey() (string, errs.SystemErrorInfo) {
	newKey, mqttErr := emt.GenerateKey(consts.MQTTChannelRoot, consts.MQTTDefaultRootPermissions, consts.MQTTDefaultTTL)
	if mqttErr != nil {
		logger.Error(mqttErr)
		return "", errs.MQTTKeyGenError
	}

	redisErr := store.Redis.Set(consts.CacheMQTTRootKey, newKey)
	if redisErr != nil {
		logger.Error(redisErr)
		return "", errs.RedisOperateError
	}
	return newKey, nil
}

func ClearRootKey() errs.SystemErrorInfo {
	_, redisErr := store.Redis.Del(consts.CacheMQTTRootKey)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func PushMQTTDataRefreshMsg(refreshInfo bo.MQTTDataRefreshNotice) errs.SystemErrorInfo {
	orgId := refreshInfo.OrgId
	projectId := refreshInfo.ProjectId

	channel := ""
	if projectId != 0 {
		channel = util.GetMQTTProjectChannel(orgId, projectId)
	} else {
		channel = util.GetMQTTOrgChannel(orgId)
	}

	pubErr := Publish(channel, json.ToJsonIgnoreError(bo.MQTTNoticeBo{
		Type: consts.MQTTNoticeTypeDataRefresh,
		Body: refreshInfo,
	}))
	if pubErr != nil {
		logger.Error(pubErr)
		return pubErr
	}
	return nil
}
