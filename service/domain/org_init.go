package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

const larkDepartmentInitSql = consts.TemplateDirPrefix + "lark_department_init.template"
const larkUserInitSql = consts.TemplateDirPrefix + "lark_user_init.template"

func OrgOwnerInit(orgId int64, owner int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Owner = owner
	dbErr := store.Mysql.TransUpdate(tx, org)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	return nil
}

func GetSuiteTicket() (string, errs.SystemErrorInfo) {
	val, redisErr := store.Redis.Get(consts.CacheDingTalkSuiteTicket)
	if redisErr != nil {
		logger.Error(redisErr)
		return val, errs.RedisOperateError
	}
	return val, nil
}

func LarkDepartmentInit(orgId int64, sourceChannel, sourcePlatform string, orgName string, creator int64) (int64, errs.SystemErrorInfo) {
	departmentId := snowflake.Id()
	contextMap := map[string]interface{}{}
	contextMap["OrgId"] = orgId
	contextMap["OrgName"] = orgName
	contextMap["DepartmentId"] = departmentId
	//contextMap["OutDepartmentId"] = outDeparmentVo.Id
	contextMap["SourceChannel"] = sourceChannel
	contextMap["SourcePlatform"] = sourcePlatform
	contextMap["Creator"] = creator
	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		dbErr := util.ReadAndWrite(larkDepartmentInitSql, contextMap, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.MysqlOperateError
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	return departmentId, nil
}
