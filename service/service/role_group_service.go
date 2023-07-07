package service

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
)

// CreateRoleGroup 创建角色组
func CreateRoleGroup(orgId, operatorUid int64, groupReq req.RoleGroupReq) (int64, errs.SystemErrorInfo) {
	if orgId == 0 {
		return 0, errs.OrgNotExist
	}
	if !format.VerifyRoleGroupNameFormat(groupReq.Name) {
		return 0, errs.RoleGroupNameLenErr
	}

	// 插入角色组
	id, dbErr := domain.CreateRoleGroup(orgId, operatorUid, groupReq.Name, consts.BlankString)
	if dbErr != nil {
		return 0, errs.MysqlOperateError
	}
	return id, nil
}

// UpdateRoleGroup 修改角色组信息
func UpdateRoleGroup(orgId, operatorUid int64, groupId int64, reqParam req.RoleGroupReq) (bool, errs.SystemErrorInfo) {

	groupBo, dbErr := domain.GetRoleGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgRoleGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	if !format.VerifyRoleGroupNameFormat(reqParam.Name) {
		return false, errs.RoleGroupNameLenErr
	}

	count, dbErr := domain.UpdateRoleGroup(orgId, operatorUid, groupId, groupBo.Version, reqParam.Name)
	if dbErr != nil {
		return false, errs.OrgRoleGroupNotExist
	}
	return count > 0, nil
}

// DeleteRoleGroup 删除角色组
func DeleteRoleGroup(orgId, operatorUid int64, groupId int64) (bool, errs.SystemErrorInfo) {
	_, dbErr := domain.GetRoleGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgRoleGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 是否存在有效角色
	roles, dbErr := domain.GetRoleListByGroup(orgId, groupId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	if len(roles) > 0 {
		return false, errs.RoleGroupHaveRoleErr
	}

	// 删除
	count, dbErr := domain.DeleteRoleGroup(orgId, operatorUid, groupId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

// GetGroupList 获取组织角色组列表
func GetGroupList(orgId int64) ([]resp.RoleGroup, errs.SystemErrorInfo) {
	groups, dbErr := domain.GetRoleGroupListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	var list []resp.RoleGroup
	_ = copyer.Copy(groups, &list)

	return list, nil
}
