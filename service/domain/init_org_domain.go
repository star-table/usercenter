package domain

import (
	"strconv"
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//留着做对比方便的注释
func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	_, dbErr := GetOrgByOutOrgId(initOrgBo.OutOrgId, initOrgBo.SourceChannel)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
	} else {
		logger.ErrorF("组织已经存在，不需要初始化")
		return 0, errs.OrgNotNeedInitError
	}

	orgId, err := OrgInfoInit(initOrgBo, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	_, err = OrgOutInfoInit(initOrgBo, orgId, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	_, err = OrgConfigInfoInit(orgId, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	////权限、角色初始化
	//roleInitResp := rolefacade.RoleInit(rolevo.RoleInitReqVo{
	//	OrgId: orgId,
	//})
	//if roleInitResp.Failure() {
	//	logger.Error(roleInitResp.Message)
	//	return 0, roleInitResp.Error()
	//}
	//logger.Info("权限、角色初始化成功")
	//
	//logger.Info("优先级初始化成功")

	err = InitDepartment(
		orgId,
		initOrgBo.OutOrgId,
		initOrgBo.SourceChannel,
		0,
		0,
		tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return orgId, nil
}

//组织信息初始化
func OrgInfoInit(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgId := snowflake.Id()
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.SourceChannel = initOrgBo.SourceChannel
	org.Name = initOrgBo.OrgName
	org.LogoUrl = initOrgBo.OrgLogo
	org.Address = initOrgBo.OrgProvince + initOrgBo.OrgCity
	org.IsAuthenticated = isAuth
	dbErr := store.Mysql.TransInsert(tx, org)
	if dbErr != nil {
		logger.ErrorF("组织初始化，添加组织时异常: %s", strs.ObjectToString(dbErr))
		return 0, errs.MysqlOperateError
	}
	return orgId, nil
}

//组织外部信息初始化
func OrgOutInfoInit(initOrgBo bo.InitOrgBo, orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgOutInfoId := snowflake.Id()
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	orgOutInfo.Id = orgOutInfoId
	orgOutInfo.OrgId = orgId
	orgOutInfo.IsDelete = consts.AppIsNoDelete
	orgOutInfo.Status = consts.AppStatusEnable
	orgOutInfo.SourceChannel = initOrgBo.SourceChannel
	orgOutInfo.Name = initOrgBo.OrgName
	orgOutInfo.OutOrgId = initOrgBo.OutOrgId
	orgOutInfo.Industry = initOrgBo.Industry
	orgOutInfo.IsAuthenticated = isAuth
	orgOutInfo.AuthLevel = strconv.Itoa(initOrgBo.AuthLevel)

	dbErr := store.Mysql.TransInsert(tx, orgOutInfo)
	if dbErr != nil {
		logger.ErrorF("组织初始化，添加外部组织信息时异常: %s", strs.ObjectToString(dbErr))
		return 0, errs.MysqlOperateError
	}
	return orgOutInfoId, nil
}

//组织配置初始化
func OrgConfigInfoInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	sysConfig := &po.PpmOrcConfig{}

	payLevel := &po.PpmBasPayLevel{}
	dbErr := store.Mysql.SelectById(payLevel.TableName(), 1, payLevel)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	orgConfigId := snowflake.Id()

	sysConfig.Id = orgConfigId
	sysConfig.OrgId = orgId
	sysConfig.TimeZone = "Asia/Shanghai"
	sysConfig.TimeDifference = "+08:00"
	sysConfig.PayLevel = 1
	sysConfig.PayStartTime = time.Now()
	sysConfig.PayEndTime = time.Now().Add(time.Duration(payLevel.Duration) * time.Second)
	sysConfig.Language = "zh-CN"
	sysConfig.RemindSendTime = "09:00:00"
	sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
	sysConfig.PasswordLength = 6
	sysConfig.PasswordRule = 1
	sysConfig.MaxLoginFailCount = 0
	sysConfig.Status = consts.AppStatusEnable
	dbErr = store.Mysql.TransInsert(tx, sysConfig)
	if dbErr != nil {
		logger.ErrorF("组织初始化，添加组织配置信息时异常: %s" + strs.ObjectToString(dbErr))
		return 0, errs.MysqlOperateError
	}

	return orgConfigId, nil
}
