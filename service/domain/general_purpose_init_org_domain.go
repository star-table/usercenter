package domain

import (
	"strconv"
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//通用版
func GeneralInitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {

	isDingTalk := initOrgBo.SourceChannel == consts.AppSourceChannelDingTalk

	//判断这个组织是否已经存在 内含存在的时候更新
	_, isUpdateFail, err := GetGeneralOrgInfoByOutOrgId(initOrgBo.OutOrgId, initOrgBo.PermanentCode, initOrgBo.SourceChannel, isDingTalk, tx)

	//不需要初始化并且已经更新过了
	if err == nil {
		logger.ErrorF("组织已经存在，不需要初始化，初始化信息为：%s", json.ToJsonIgnoreError(initOrgBo))
		return 0, errs.BuildSystemErrorInfo(errs.OrgNotNeedInitError)
	}
	//走这里说明组织存在 且更新时候异常 直接返回出去
	if isUpdateFail != nil && *isUpdateFail {
		logger.ErrorF("组织初始化，更新组织时异常：%s", strs.ObjectToString(err))
		return 0, err
	}

	//组织信息初始化
	orgId, err := GeneralOrgInfoInit(initOrgBo, isDingTalk, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	//处理外部信息
	_, err = GeneralOrgOutInfoInit(initOrgBo, orgId, isDingTalk, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	_, err = GeneralOrgConfigInfoInit(orgId, tx)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

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

	//钉钉单独初始化的内容
	//dErr := DingDingInit(initOrgBo.OutOrgId, orgId, roleInitResp.RoleInitResp, tx)

	//if dErr != nil {
	//	return 0, errs.BuildSystemErrorInfo(errs.OrgInitError, dErr)
	//}

	return orgId, nil
}

func GetGeneralOrgInfoByOutOrgId(outOrgId string, permanentCode string, sourceChannel string, isDingTalk bool, tx sqlbuilder.Tx) (baseOrgInfo *bo.BaseOrgInfoBo, updateFlag *bool, returnErr errs.SystemErrorInfo) {
	conds := db.Cond{
		consts.TcOutOrgId:      outOrgId,
		consts.TcSourceChannel: sourceChannel,
	}

	//飞书的拼接未删除的条件
	if !isDingTalk {
		conds[consts.TcIsDelete] = consts.AppIsNoDelete
	}
	var outOrgInfo po.PpmOrgOrganizationOutInfo
	dbErr := store.Mysql.SelectOneByCond(consts.TableOrganizationOutInfo, conds, &outOrgInfo)
	//err 不为空说明组织不存在 查不到
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, nil, errs.OrgOutInfoNotExist
		}
		logger.Error(dbErr)
		return nil, nil, errs.MysqlOperateError
	}

	//获取原本的组织信息 不存在就出去初始化
	orgInfo, dbErr := GetOrgById(outOrgInfo.OrgId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, nil, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return nil, nil, errs.MysqlOperateError
	}

	//组装数据返回
	return &bo.BaseOrgInfoBo{
		OrgId:         orgInfo.Id,
		OrgName:       orgInfo.Name,
		OutOrgId:      outOrgId,
		SourceChannel: sourceChannel,
	}, nil, nil
}

//组织信息初始化
func GeneralOrgInfoInit(initOrgBo bo.InitOrgBo, isDingTalk bool, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {

	isAuth := 0

	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgId := snowflake.Id()
	org := &po.PpmOrgOrganization{}
	//统一处理的id,状态,名字,sourceChannel
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.SourceChannel = initOrgBo.SourceChannel

	org.Name = initOrgBo.OrgName
	org.LogoUrl = initOrgBo.OrgLogo
	org.Address = initOrgBo.OrgProvince + initOrgBo.OrgCity
	org.IsAuthenticated = isAuth
	//插入org 信息
	dbErr := store.Mysql.TransInsert(tx, org)
	if dbErr != nil {
		logger.ErrorF("组织初始化，添加组织时异常: %s", strs.ObjectToString(dbErr))
		return 0, errs.MysqlOperateError
	}
	return orgId, nil
}

//组织外部信息初始化
func GeneralOrgOutInfoInit(initOrgBo bo.InitOrgBo, orgId int64, isDingTalk bool, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
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
func GeneralOrgConfigInfoInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
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
		logger.ErrorF("组织初始化，添加组织配置信息时异常: %s", strs.ObjectToString(dbErr))
		return 0, errs.MysqlOperateError
	}

	return orgConfigId, nil
}
