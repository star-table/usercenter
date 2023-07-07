package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetOrgIdByOutOrgId(sourceChannel string, outOrgId string) (int64, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(consts.CacheOutOrgIdRelationId, map[string]interface{}{
		consts.CacheKeyOutOrgIdConstName:      outOrgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return 0, errs.TemplateRenderError
	}
	orgIdInfoJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return 0, errs.RedisOperateError
	}
	if orgIdInfoJson != "" {
		orgIdInfo := &bo.OrgIdInfo{}
		jsonErr := json.FromJson(orgIdInfoJson, orgIdInfo)
		if jsonErr != nil {
			return 0, errs.JSONConvertError
		}
		return orgIdInfo.OrgId, nil
	} else {
		orgBo, dbErr := GetOrgByOutOrgId(sourceChannel, outOrgId)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				return 0, errs.OrgOutInfoNotExist
			}
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
		orgIdInfo := bo.OrgIdInfo{
			OutOrgId: outOrgId,
			OrgId:    orgBo.Id,
		}
		orgIdInfoJson = json.ToJsonIgnoreError(orgIdInfo)
		redisErr := store.Redis.SetEx(key, orgIdInfoJson, consts.GetCacheBaseExpire())
		if redisErr != nil {
			logger.Error(redisErr)
			return 0, errs.RedisOperateError
		}
		return orgIdInfo.OrgId, nil
	}
}

func GetBaseOrgOutInfo(orgId int64) (*bo.BaseOrgOutInfoBo, errs.SystemErrorInfo) {
	key, sysErr := util.ParseCacheKey(consts.CacheBaseOrgOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if sysErr != nil {
		logger.Error(sysErr)
		return nil, errs.TemplateRenderError
	}
	outOrgInfoJson, err := store.Redis.Get(key)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if outOrgInfoJson != "" {
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{}
		err := json.FromJson(outOrgInfoJson, orgOutInfoBo)
		if err != nil {
			logger.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return orgOutInfoBo, nil
	} else {
		var orgOutInfos []*po.PpmOrgOrganizationOutInfo
		err = store.Mysql.SelectAllByCond("ppm_org_organization_out_info", db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: db.In([]int{0, consts.AppIsNoDelete}),
		}, &orgOutInfos)
		if err != nil || len(orgOutInfos) == 0 {
			logger.InfoF("get org out info orgId: %d err: %v", orgId, err)
			return nil, errs.OrgOutInfoNotExist
		}
		var orgOutInfo *po.PpmOrgOrganizationOutInfo
		for _, outInfo := range orgOutInfos {
			if outInfo.OutOrgId != "" { // 优先拿有out_org_id的
				orgOutInfo = outInfo
				break
			}
		}
		if orgOutInfo == nil { // 没有就拿第一个
			orgOutInfo = orgOutInfos[0]
		}
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{
			OrgId:         orgId,
			OutOrgId:      orgOutInfo.OutOrgId,
			SourceChannel: orgOutInfo.SourceChannel,
		}
		outOrgInfoJson = json.ToJsonIgnoreError(orgOutInfoBo)
		//err = store.Redis.SetEx(key, outOrgInfoJson, consts.GetCacheBaseExpire())
		//if err != nil {
		//	logger.Error(err)
		//	return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		//}
		return orgOutInfoBo, nil
	}
}

func GetBaseOrgInfoByOutOrgId(sourceChannel string, outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	orgId, err := GetOrgIdByOutOrgId(sourceChannel, outOrgId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return GetBaseOrgInfo(sourceChannel, orgId)
}

func ClearCacheBaseOrgInfo(sourceChannel string, orgId int64) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(consts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
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

func GetBaseOrgInfo(sourceChannel string, orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	baseOrgInfo := &bo.BaseOrgInfoBo{}

	key, keyErr := util.ParseCacheKey(consts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return nil, errs.TemplateRenderError
	}

	if value, err := store.Redis.Get(key); err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	} else if value != "" {
		if err = json.FromJson(value, baseOrgInfo); err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	}

	// 没拿到缓存或者缓存有问题，从DB拿数据
	if baseOrgInfo.OrgId == 0 || baseOrgInfo.SourceChannel == "" {
		orgOutInfo, sysErr := GetBaseOrgOutInfo(orgId)
		if sysErr != nil {
			logger.ErrorF("[GetBaseOrgInfo] GetBaseOrgOutInfo orgId:%v, err:%v", orgId, sysErr.Error())
			return nil, sysErr
		}

		var orgInfo po.PpmOrgOrganization
		if err := store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
			consts.TcId:       orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &orgInfo); err != nil {
			logger.ErrorF("[GetBaseOrgInfo] get ppm_org_organization orgId:%v, err:%v", orgId, strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
		}

		baseOrgInfo.OrgId = orgId
		baseOrgInfo.OrgName = orgInfo.Name
		baseOrgInfo.OrgOwnerId = orgInfo.Owner
		baseOrgInfo.Creator = orgInfo.Creator
		baseOrgInfo.OutOrgId = orgOutInfo.OutOrgId
		baseOrgInfo.SourceChannel = orgOutInfo.SourceChannel // 这个才是目前正确的source_channel

		logger.InfoF("[GetBaseOrgInfo] getOrgInfo: %v", json.ToJsonIgnoreError(baseOrgInfo))
	}
	return baseOrgInfo, nil
}

func GetOutDeptAndInnerDept(orgId int64, tx sqlbuilder.Tx) (map[string]int64, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(consts.CacheDeptRelation, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return nil, errs.TemplateRenderError
	}

	deptRelationListJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	var deptRelationList = make(map[string]int64)
	if deptRelationListJson != "" {
		jsonErr := json.FromJson(deptRelationListJson, &deptRelationList)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		return deptRelationList, nil
	} else {
		deptOurInfoList, err := queryDepartmentOutInfWithTx(tx, orgId)

		if err != nil {
			logger.Error(err)
			return deptRelationList, errs.MysqlOperateError
		}
		logger.Info("部门关联关系: " + strs.ObjectToString(deptOurInfoList))

		for _, v := range deptOurInfoList {
			deptRelationList[v.OutOrgDepartmentId] = v.DepartmentId
		}
		deptRelationListJson := json.ToJsonIgnoreError(deptRelationList)

		redisErr = store.Redis.SetEx(key, deptRelationListJson, consts.GetCacheBaseExpire())
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}

		return deptRelationList, nil
	}
}

func queryDepartmentOutInfWithTx(tx sqlbuilder.Tx, orgId int64) ([]po.PpmOrgDepartmentOutInfo, error) {
	var deptOurInfoList []po.PpmOrgDepartmentOutInfo
	if tx != nil {
		//TODO 未定义TransSelectAllByCond，先使用SelectAllByCond(不使用事务不会有问题)
		dbErr := store.Mysql.TransSelectAllByCond(tx, consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &deptOurInfoList)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, dbErr
		}
	} else {
		dbErr := store.Mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &deptOurInfoList)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, dbErr
		}
	}
	return deptOurInfoList, nil
}
