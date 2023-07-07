package service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"upper.io/db.v3"
)

// ApiKeyAuth  验证OpenApi
func ApiKeyAuth(apiKey string) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	// 获取机构
	org, dbErr := domain.GetOrgByApiKey(apiKey)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.ApiKeyAuthErr
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	uBo := &bo.CacheUserInfoBo{
		OrgId:  org.Id,
		UserId: org.Owner,
	}
	return uBo, nil
}
