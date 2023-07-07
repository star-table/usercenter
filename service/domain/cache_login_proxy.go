package domain

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/json"
)

const LoginCodeFreezeTimeBlankObj = "BLANK"

func CheckSMSLoginCodeFreezeTime(authType, addressType int, address string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       address,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	exist, redisErr := store.Redis.Exist(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	if exist {
		return errs.SMSSendTimeLimitError
	}
	return nil
}

func ClearSMSLoginCodeFreezeTime(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func SetSMSLoginCodeFreezeTime(authType, addressType int, phoneNumber string, minute int) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	redisErr := store.Redis.SetEx(key, LoginCodeFreezeTimeBlankObj, int64(60*minute))
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func ClearSMSLoginCode(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func ClearPwdLoginCode(loginName string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func SetSMSLoginCode(authType, addressType int, phoneNumber string, authCode string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	redisErr := store.Redis.SetEx(key, authCode, int64(60*5))
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetSMSLoginCode(authType, addressType int, phoneNumber string) (string, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return "", errs.TemplateRenderError
	}
	value, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return "", errs.RedisOperateError
	}
	if value == "" {
		logger.Error("登录验证码为空: " + phoneNumber)
		return "", errs.SMSLoginCodeInvalid
	}
	return value, nil
}

func GetPwdLoginCode(loginName string) (string, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(consts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return "", errs.TemplateRenderError
	}
	value, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return "", errs.RedisOperateError
	}
	if value == "" {
		logger.Error("登录验证码为空" + loginName)
		return "", errs.SMSLoginCodeInvalid
	}
	return value, nil
}

func SetPwdLoginCode(loginName, code string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	redisErr := store.Redis.SetEx(key, code, int64(60*5))
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	return nil
}

func SetSMSLoginCodeVerifyFailTimesIncrement(authType, addressType int, phoneNumber string) (int64, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return 0, errs.TemplateRenderError
	}
	c, redisErr := store.Redis.Incrby(key, 1)
	if redisErr != nil {
		logger.Error(redisErr)
		return 0, errs.RedisOperateError
	}
	return c, nil
}

func ClearSMSLoginCodeVerifyFailTimesIncrement(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheSmsLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func ClearSMSLoginCodeAllCache(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	err := ClearSMSLoginCode(authType, addressType, phoneNumber)
	if err != nil {
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	err = ClearSMSLoginCodeVerifyFailTimesIncrement(authType, addressType, phoneNumber)
	if err != nil {
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetPhoneNumberWhiteList() ([]string, errs.SystemErrorInfo) {
	key := consts.CachePhoneNumberWhiteList
	value, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if value == "" {
		//如果没有就设置
		defaultWhiteList, err := SetPhoneNumberWhiteList()
		if err != nil {
			return nil, err
		}

		return defaultWhiteList, nil
	}

	result := make([]string, 0)
	jsonErr := json.FromJson(value, &result)
	if jsonErr != nil {
		logger.Info(jsonErr)
		return nil, errs.JSONConvertError
	}

	return result, nil
}

func SetPhoneNumberWhiteList() ([]string, errs.SystemErrorInfo) {
	key := consts.CachePhoneNumberWhiteList
	whiteList := []string{
		"18891111111",
		"18891111112",
		"18891111113",
		"18891111114",
		"18891111115",
		"18616215755",
		"18117493299",
		"15618960078",
		"17621142248",
		"18516232262",
		"18917633966",
		"13037568259",
		"18917637631",
		"15618931187",
		"18221304331",
		"13621822254",
	}
	whiteListJson := json.ToJsonIgnoreError(whiteList)

	redisErr := store.Redis.SetEx(key, whiteListJson, 60*60*24*30)
	if redisErr != nil {
		logger.Info(redisErr)
		return whiteList, errs.RedisOperateError
	}

	return whiteList, nil
}

func AuthCodeVerify(authType, addressType int, contactAddress, authCode string) errs.SystemErrorInfo {
	localCode, err := GetSMSLoginCode(authType, addressType, contactAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	if !strings.EqualFold(localCode, authCode) {
		//计数，防止刷接口
		times, err := SetSMSLoginCodeVerifyFailTimesIncrement(authType, addressType, contactAddress)
		if err != nil {
			return err
		}
		if times > 4 {
			//如果失败五次，则清空其他缓存，让用户重新发送
			err := ClearSMSLoginCodeAllCache(authType, addressType, contactAddress)
			if err != nil {
				logger.Error(err)
			}
			return errs.BuildSystemErrorInfo(errs.SMSLoginCodeVerifyFailTimesOverLimit)
		}
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeNotMatch)
	}

	//登录成功
	//清理验证码相关缓存
	err = ClearSMSLoginCodeAllCache(authType, addressType, contactAddress)
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func ClearChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func GetChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	value, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	if value == "" {
		return errs.ChangeLoginNameInvalid
	}
	return nil
}

func SetChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	redisErr := store.Redis.SetEx(key, "true", int64(60*5))
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	return nil
}
