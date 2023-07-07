package domain

import (
	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

func GetUserIdBatchByEmpId(sourceChannel string, orgId int64, empIds []string) ([]int64, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(empIds))
	for i, empId := range empIds {
		key, keyErr := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyOutUserIdConstName:     empId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return nil, errs.TemplateRenderError
		}
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, redisErr := store.Redis.MGet(keys...)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
		resultList = list
	}
	userIds := make([]int64, 0)
	validEmpIds := make([]string, 0)
	for _, empInfoJson := range resultList {
		empIdInfo := &bo.UserEmpIdInfo{}
		jsonErr := json.FromJson(empInfoJson, empIdInfo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		userIds = append(userIds, empIdInfo.UserId)
		validEmpIds = append(validEmpIds, empIdInfo.EmpId)
	}
	//找不存在的
	if len(empIds) != len(validEmpIds) {
		for _, empId := range empIds {
			exist, _ := slice.Contain(validEmpIds, empId)
			if !exist {
				userId, err := GetUserIdByEmpId(sourceChannel, orgId, empId)
				if err != nil {
					logger.Error(err)
					continue
				}
				userIds = append(userIds, userId)
			}
		}
	}
	return userIds, nil
}

func GetUserIdByEmpId(sourceChannel string, orgId int64, empId string) (int64, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
		consts.CacheKeyOutUserIdConstName:     empId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return 0, errs.TemplateRenderError
	}

	empInfoJson, err := store.Redis.Get(key)
	if err != nil {
		logger.Error(err)
		return 0, errs.RedisOperateError
	}
	if empInfoJson != "" {
		empIdInfo := &bo.UserEmpIdInfo{}
		jsonErr := json.FromJson(empInfoJson, empIdInfo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return 0, errs.JSONConvertError
		}
		return empIdInfo.UserId, nil
	} else {
		userOutInfo := &po.PpmOrgUserOutInfo{}
		dbErr := store.Mysql.SelectOneByCond(userOutInfo.TableName(), db.Cond{
			consts.TcOutUserId:     empId,
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcIsDelete:      consts.AppIsNoDelete,
			consts.TcStatus:        consts.AppStatusEnable,
		}, userOutInfo)
		if dbErr != nil {
			logger.Error(dbErr)
			return 0, errs.UserNotExist
		}
		empIdInfo := bo.UserEmpIdInfo{
			EmpId:  empId,
			UserId: userOutInfo.UserId,
		}
		redisErr := store.Redis.Set(key, json.ToJsonIgnoreError(empIdInfo))
		if redisErr != nil {
			logger.Error(redisErr)
			return 0, errs.RedisOperateError
		}
		return userOutInfo.UserId, nil
	}
}

func GetDingTalkBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(consts.AppSourceChannelDingTalk, orgId, empId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return GetDingTalkBaseUserInfo(orgId, userId)
}

func GetBaseUserInfoByEmpId(sourceChannel string, orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(sourceChannel, orgId, empId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return GetBaseUserInfo(sourceChannel, orgId, userId)
}

func GetFeiShuBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(consts.AppSourceChannelFeiShu, orgId, empId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return GetFeiShuBaseUserInfo(orgId, userId)
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(consts.AppSourceChannelDingTalk, orgId, userId)
}

func GetFeiShuBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(consts.AppSourceChannelFeiShu, orgId, userId)
}
