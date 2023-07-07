package domain

import (
	"strconv"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

func UpdateUserConfig(orgId, operatorId int64, userConfigBo bo.UserConfigBo) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if util.IsBool(userConfigBo.DailyReportMessageStatus) {
		upd[consts.TcDailyReportMessageStatus] = userConfigBo.DailyReportMessageStatus
	}
	if util.IsBool(userConfigBo.OwnerRangeStatus) {
		upd[consts.TcOwnerRangeStatus] = userConfigBo.OwnerRangeStatus
	}
	if util.IsBool(userConfigBo.ParticipantRangeStatus) {
		upd[consts.TcParticipantRangeStatus] = userConfigBo.ParticipantRangeStatus
	}
	if util.IsBool(userConfigBo.AttentionRangeStatus) {
		upd[consts.TcAttentionRangeStatus] = userConfigBo.AttentionRangeStatus
	}
	if util.IsBool(userConfigBo.CreateRangeStatus) {
		upd[consts.TcCreateRangeStatus] = userConfigBo.CreateRangeStatus
	}
	if util.IsBool(userConfigBo.RemindMessageStatus) {
		upd[consts.TcRemindMessageStatus] = userConfigBo.RemindMessageStatus
	}
	if util.IsBool(userConfigBo.CommentAtMessageStatus) {
		upd[consts.TcCommentAtMessageStatus] = userConfigBo.CommentAtMessageStatus
	}
	if util.IsBool(userConfigBo.ModifyMessageStatus) {
		upd[consts.TcModifyMessageStatus] = userConfigBo.ModifyMessageStatus
	}
	if util.IsBool(userConfigBo.RelationMessageStatus) {
		upd[consts.TcRelationMessageStatus] = userConfigBo.RelationMessageStatus
	}
	if util.IsBool(userConfigBo.DailyProjectReportMessageStatus) {
		upd[consts.TcDailyProjectReportMessageStatus] = userConfigBo.DailyProjectReportMessageStatus
	}

	//更新人必填
	upd[consts.TcUpdator] = operatorId

	_, err := store.Mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   operatorId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)

	if err != nil {
		//配置更新失败
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateUserPcConfig(orgId, operatorId int64, userConfigBo bo.UserConfigBo) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if util.IsBool(userConfigBo.PcNoticeOpenStatus) {
		upd[consts.TcPcNoticeOpenStatus] = userConfigBo.PcNoticeOpenStatus
	}
	if util.IsBool(userConfigBo.PcIssueRemindMessageStatus) {
		upd[consts.TcPcIssueRemindMessageStatus] = userConfigBo.PcIssueRemindMessageStatus
	}
	if util.IsBool(userConfigBo.PcOrgMessageStatus) {
		upd[consts.TcPcOrgMessageStatus] = userConfigBo.PcOrgMessageStatus
	}
	if util.IsBool(userConfigBo.PcProjectMessageStatus) {
		upd[consts.TcPcProjectMessageStatus] = userConfigBo.PcProjectMessageStatus
	}
	if util.IsBool(userConfigBo.PcCommentAtMessageStatus) {
		upd[consts.TcPcCommentAtMessageStatus] = userConfigBo.PcCommentAtMessageStatus
	}
	_, err := store.Mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   operatorId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if err != nil {
		//配置更新失败
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func InsertUserConfig(orgId, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	userConfig := &bo.UserConfigBo{}
	var err errs.SystemErrorInfo
	//如果不存在就插入
	userIdStr := strconv.FormatInt(userId, 10)
	lockKey := consts.AddUserConfigLock + userIdStr
	suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, userIdStr)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if suc {
		defer store.Redis.ReleaseDistributedLock(lockKey, userIdStr)

		userConfig, err = GetUserConfigInfo(orgId, userId)
		if err != nil {
			err := insertUserConfig(orgId, userId)
			if err != nil {
				return nil, err
			}

		}
	} else {
		userConfig, err = GetUserConfigInfo(orgId, userId)
		if err != nil {
			logger.Error(err)
			return nil, errs.UserConfigNotExist
		}
	}

	userConfigBo := &bo.UserConfigBo{}
	copyErr := copyer.Copy(userConfig, userConfigBo)
	if err != nil {
		logger.Error(copyErr)
		return nil, errs.ObjectCopyError
	}
	return userConfigBo, nil
}

//用户配置不存在插入用户配置
func insertUserConfig(orgId, userId int64) errs.SystemErrorInfo {

	userConfig := &po.PpmOrgUserConfig{}
	userConfigId := snowflake.Id()
	userConfig.Id = userConfigId
	userConfig.OrgId = orgId
	userConfig.UserId = userId
	userConfig.Creator = userId
	userConfig.Updator = userId
	userConfig.IsDelete = consts.AppIsNoDelete
	err2 := store.Mysql.Insert(userConfig)
	if err2 != nil {
		logger.Error(err2)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return nil
}
