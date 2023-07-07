package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
)

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetCtxParameters(ctx context.Context, key string) (string, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}
	v := gc.GetString(key)
	return v, nil
}

func GetCurrentUser(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, true)
}

func GetCurrentUserWithoutOrgVerify(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, false)
}

func GetCurrentUserWithCond(ctx context.Context, orgVerify bool) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	token, err := GetCtxParameters(ctx, consts.AppHeaderTokenName)

	if err != nil || token == "" {
		return nil, errs.BuildSystemErrorInfo(errs.TokenNotExist)
	} else {

		cacheUserInfoJson, redisErr := store.Redis.Get(sconsts.CacheUserToken + token)
		if redisErr != nil {
			logger.Error(err)
			return nil, errs.RedisOperateError
		}
		if cacheUserInfoJson == "" {
			logger.Error("token失效")
			return nil, errs.TokenExpires
		}
		cacheUserInfo := &bo.CacheUserInfoBo{}
		jsonErr := json.FromJson(cacheUserInfoJson, cacheUserInfo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		_, redisErr = store.Redis.Expire(sconsts.CacheUserToken+token, consts.CacheUserTokenExpire)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.TokenExpires
		}
		//判断用户组织状态
		if cacheUserInfo.OrgId != 0 && orgVerify {
			baseUserInfo, err := GetBaseUserInfo("", cacheUserInfo.OrgId, cacheUserInfo.UserId)
			if err != nil {
				return nil, err
			}
			err = BaseUserInfoOrgStatusCheck(*baseUserInfo)
			if err != nil {
				return nil, err
			}
		}
		return cacheUserInfo, nil
	}
}

//用户信息所在组织状态监测
func BaseUserInfoOrgStatusCheck(baseUserInfo bo.BaseUserInfoBo) errs.SystemErrorInfo {
	baseOrgOutInfo, errSys := domain.GetBaseOrgOutInfo(baseUserInfo.OrgId)
	if errSys != nil {
		logger.Error(errSys)
		return errSys
	}
	if baseUserInfo.OrgUserStatus != consts.AppStatusEnable {
		if baseOrgOutInfo.OutOrgId != "" {
			return errs.OrgUserInvalidErr
		}
		return errs.OrgUserUnabledErr
	}
	if baseUserInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
		return errs.OrgUserCheckStatusUnabledErr
	}
	if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted {
		return errs.OrgUserDeletedErr
	}
	return nil
}

func UpdateCacheUserInfoOrgId(token string, orgId int64) errs.SystemErrorInfo {
	cacheUserInfoJson, redisErr := store.Redis.Get(sconsts.CacheUserToken + token)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	if cacheUserInfoJson == "" {
		return errs.BuildSystemErrorInfo(errs.TokenExpires)
	}
	cacheUserInfo := &bo.CacheUserInfoBo{}
	_ = json.FromJson(cacheUserInfoJson, cacheUserInfo)

	//更新缓存用户的orgId
	cacheUserInfo.OrgId = orgId
	cacheUserInfoJson = json.ToJsonIgnoreError(cacheUserInfo)
	redisErr = store.Redis.SetEx(sconsts.CacheUserToken+token, cacheUserInfoJson, consts.CacheUserTokenExpire)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}
