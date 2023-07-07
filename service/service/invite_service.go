package service

import (
	"strconv"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/rand"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// GetInviteCode 获取邀请码
func GetInviteCode(orgId int64, operatorUid int64, sourcePlatform string) (*resp.GetInviteCodeResp, errs.SystemErrorInfo) {
	////用户角色权限校验
	//authErr := AuthOrgRole(orgId, operatorUid, consts2.RoleOperationPathOrgUser, consts2.RoleOperationInvite)
	//if authErr != nil {
	//	logger.Error(authErr)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	//}
	inviteInfo := bo.InviteInfoBo{
		InviterId:      operatorUid,
		OrgId:          orgId,
		SourcePlatform: sourcePlatform,
	}

	inviteCode := rand.RandomInviteCode(uuid.NewUuid() + strconv.FormatInt(operatorUid, 10) + sourcePlatform)
	err := domain.SetUserInviteCodeInfo(inviteCode, inviteInfo)
	if err != nil {
		logger.Info(strs.ObjectToString(err))
		return nil, err
	}
	return &resp.GetInviteCodeResp{InviteCode: inviteCode, Expire: consts.CacheUserInviteCodeExpire}, nil
}

func GetInviteInfo(inviteCode string) (*resp.GetInviteInfoResp, errs.SystemErrorInfo) {
	inviteInfo, err := domain.GetUserInviteCodeInfo(inviteCode)
	if err != nil {
		logger.Info(strs.ObjectToString(err))
		return nil, err
	}

	orgBaseInfo, err := domain.GetBaseOrgInfo("", inviteInfo.OrgId)
	if err != nil {
		logger.Info(strs.ObjectToString(err))
		return nil, err
	}
	userBaseInfo, err := domain.GetBaseUserInfo("", inviteInfo.OrgId, inviteInfo.InviterId)
	if err != nil {
		logger.Info(strs.ObjectToString(err))
		return nil, err
	}

	return &resp.GetInviteInfoResp{
		OrgID:       orgBaseInfo.OrgId,
		OrgName:     orgBaseInfo.OrgName,
		InviterID:   userBaseInfo.UserId,
		InviterName: userBaseInfo.Name,
	}, nil
}

//  RemoveInviteMember 删除成员邀请列表
func RemoveInviteMember(orgId, operatorUid int64, param req.RemoveInviteUserReq) errs.SystemErrorInfo {
	if len(param.Ids) == 0 && param.IsAll != 1 {
		return nil
	}
	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
		if param.IsAll != 1 {
			cond[consts.TcId] = db.In(param.Ids)
		}
		_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserInvite, cond, mysql.Upd{
			consts.TcUpdator:  operatorUid,
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}
