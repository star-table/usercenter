package domain

import (
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

const OrcConfigSql = consts.TemplateDirPrefix + "ppm_orc_config.template"

func OrgSysConfigInit(tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {

	payLevel := &po.PpmBasPayLevel{}
	dbErr := store.Mysql.SelectById(payLevel.TableName(), 1, payLevel)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	id := snowflake.Id()

	contextMap := map[string]interface{}{}
	contextMap["Id"] = id
	contextMap["OrgId"] = orgId
	contextMap["TimeZone"] = "Asia/Shanghai"
	contextMap["TimeDifference"] = "+08:00"
	contextMap["PayLevel"] = 1
	contextMap["PayStartTime"] = time.Now().Format(consts.AppTimeFormat)
	contextMap["PayEndTime"] = time.Now().Add(time.Duration(payLevel.Duration) * time.Second).Format(consts.AppTimeFormat)
	contextMap["Language"] = "zh-CN"
	contextMap["RemindSendTime"] = "09:00:00"
	contextMap["DatetimeFormat"] = "yyyy-MM-dd HH:mm:ss"
	contextMap["PasswordLength"] = 6
	contextMap["PasswordRule"] = 1
	contextMap["MaxLoginFailCount"] = 0
	contextMap["Status"] = 1
	dbErr = util.ReadAndWrite(OrcConfigSql, contextMap, tx)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	//sysConfig.Id = respVo.Id
	//sysConfig.OrgId = orgId
	//sysConfig.TimeZone = "Asia/Shanghai"
	//sysConfig.TimeDifference = "+08:00"
	//sysConfig.PayLevel = 1
	//sysConfig.PayStartTime = time.Now()
	//sysConfig.PayEndTime = time.Now().Add(time.Duration(payLevel.Duration) * time.Second)
	//sysConfig.Language = "zh-CN"
	//sysConfig.RemindSendTime = "09:00:00"
	//sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
	//sysConfig.PasswordLength = 6
	//sysConfig.PasswordRule = 1
	//sysConfig.MaxLoginFailCount = 0
	//sysConfig.Status = 1
	//sysConfig.IsDelete = consts.AppIsNoDelete
	//
	//err = store.Mysql.TransInsert(tx, sysConfig)
	//if err != nil {
	//	logger.Error(err)
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	return nil
}
