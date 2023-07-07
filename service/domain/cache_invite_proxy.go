package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/model/bo"
)

func SetUserInviteCodeInfo(inviteCode string, inviteInfo bo.InviteInfoBo) errs.SystemErrorInfo {
	inviteInfoJson := json.ToJsonIgnoreError(inviteInfo)

	redisErr := store.Redis.SetEx(consts.CacheUserInviteCode+inviteCode, inviteInfoJson, int64(consts.CacheUserInviteCodeExpire))
	if redisErr != nil {
		logger.Info(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func GetUserInviteCodeInfo(inviteCode string) (*bo.InviteInfoBo, errs.SystemErrorInfo) {
	inviteInfoJson, redisErr := store.Redis.Get(consts.CacheUserInviteCode + inviteCode)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	logger.InfoF("邀请code %s 对应的邀请信息 %s", inviteCode, inviteInfoJson)
	if inviteInfoJson == "" {
		return nil, errs.InviteCodeInvalid
	}
	inviteInfo := &bo.InviteInfoBo{}
	jsonErr := json.FromJson(inviteInfoJson, inviteInfo)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return nil, errs.JSONConvertError
	}
	return inviteInfo, nil
}
