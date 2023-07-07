package open_service

import (
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/service"
)

// AuthToken 对token进行认证
func AuthToken(accessToken string) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	if accessToken == "" {
		return nil, errs.TokenAuthError
	}
	loginUserJson, redisErr := store.Redis.Get(sconsts.CacheUserToken + accessToken)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if loginUserJson == "" {
		return nil, errs.TokenAuthError
	}

	logger.InfoF("token %s 认证信息: %s", accessToken, loginUserJson)
	loginUser := &bo.CacheUserInfoBo{}
	_ = json.FromJson(loginUserJson, loginUser)
	//校验用户
	baseUserInfo, err := domain.GetBaseUserInfo("", loginUser.OrgId, loginUser.UserId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = service.BaseUserInfoOrgStatusCheck(*baseUserInfo)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return loginUser, nil
}
