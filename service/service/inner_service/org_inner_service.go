package inner_service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req/inner_req"
)

func GetOrgInfo(orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	info, err := domain.GetBaseOrgInfo("", orgId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	config, configErr := domain.GetOrgConfig(orgId)
	if configErr != nil {
		logger.Error(configErr)
		return nil, configErr
	}

	info.PayLevel = config.PayLevel

	return info, nil
}

func AddOutCollaborator(req inner_req.AddOrgOutCollaborator) errs.SystemErrorInfo {
	return domain.AddOrgOutCollaborator(req.OrgID, req.UserID)
}

// CheckAndSetSuperAdmin 检查组织拥有者是否是超管，如果不是，则设置为超管。
func CheckAndSetSuperAdmin(req inner_req.CheckAndSetSuperAdminReq) errs.SystemErrorInfo {
	org, err := domain.GetBaseOrgInfo("", req.OrgID)
	if err != nil {
		logger.Error(err)
		return err
	}
	sysGroup, dbErr := domain.GetSysManageGroup(req.OrgID)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	userIdArr := make([]int64, 0)
	oriErr := json.FromJson(sysGroup.UserIds, &userIdArr)
	if oriErr != nil {
		logger.Error(oriErr)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
	}
	if isOk, _ := slice.Contain(userIdArr, org.OrgOwnerId); org.OrgOwnerId > 0 && !isOk {
		_, dbErr = domain.AppendUserIntoAdminGroup(sysGroup.Id, 0, org.OrgOwnerId, nil)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
	}
	return nil
}
