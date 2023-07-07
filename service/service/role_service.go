package service

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
)

// CreateRole 创建角色
func CreateRole(orgId, operatorUid int64, roleReq req.CreateRoleReq) (int64, errs.SystemErrorInfo) {

	if !format.VerifyRoleNameFormat(roleReq.Name) {
		return 0, errs.RoleNameLenErr
	}
	if ok, _ := slice.Contain(domain.DefaultRoleName, roleReq.Name); ok {
		return 0, errs.DefaultRoleNameErr
	}
	// 获取角色组
	if roleReq.RoleGroupId == 0 {
		// 这里不放在同一个事务，没有必要，下次添加的时候，就无需自动生成
		groupId, dbErr := domain.CreateRoleGroup(orgId, operatorUid, domain.DefaultRoleGroupName, consts.BlankString)
		if dbErr != nil {
			return 0, errs.MysqlOperateError
		}
		roleReq.RoleGroupId = groupId
	}

	//插入角色
	id, dbErr := domain.CreateRole(orgId, operatorUid, roleReq.RoleGroupId, roleReq.Name)
	if dbErr != nil {
		return 0, errs.MysqlOperateError
	}
	return id, nil
}

// UpdateRole 修改角色
func UpdateRole(orgId, operatorUid int64, roleId int64, roleReq req.UpdateRoleReq) (bool, errs.SystemErrorInfo) {
	if ok, _ := slice.Contain(domain.DefaultRoleName, roleReq.Name); ok {
		return false, errs.DefaultRoleNameErr
	}
	if !format.VerifyRoleNameFormat(roleReq.Name) {
		return false, errs.RoleNameLenErr
	}

	roleBo, dbErr := domain.GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.RoleNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	count, dbErr := domain.UpdateRole(orgId, operatorUid, roleId, roleReq.Name, roleBo.Version)

	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	return count > 0, nil
}

// MoveRole 移动角色
func MoveRole(orgId, operatorUid int64, groupId, roleId int64) (bool, errs.SystemErrorInfo) {
	roleBo, dbErr := domain.GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.RoleNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	// 获取角色组
	_, dbErr = domain.GetRoleGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgRoleGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	count, dbErr := domain.MoveRole(orgId, operatorUid, groupId, roleId, roleBo.Version)

	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	return count > 0, nil
}

// DeleteRole 删除角色
func DeleteRole(orgId, operatorUid int64, roleId int64) (bool, errs.SystemErrorInfo) {
	_, dbErr := domain.GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.RoleNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	count, dbErr := domain.DeleteRole(orgId, operatorUid, roleId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

// GetOrgRoleListResp 通讯录-获取该组织下角色组和角色列表
func GetOrgRoleListResp(orgId int64, perContext *inner_resp.OrgUserPerContext) (*resp.RoleListResp, errs.SystemErrorInfo) {
	listResp := &resp.RoleListResp{
		RoleGroups: []resp.RoleGroupInfo{},
		Roles:      []resp.RoleInfo{},
	}
	groupBos, dbErr := domain.GetRoleGroupListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	if len(groupBos) == 0 {
		return listResp, nil
	}

	roles, dbErr := domain.GetRoleListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	for _, groupBo := range groupBos {
		var groupInfo resp.RoleGroupInfo
		_ = copyer.Copy(groupBo, &groupInfo)
		listResp.RoleGroups = append(listResp.RoleGroups, groupInfo)
	}
	for _, role := range roles {
		var roleInfo resp.RoleInfo
		_ = copyer.Copy(role, &roleInfo)
		// 是否有权限修改成员
		roleInfo.EditableMember = perContext.HasManageRole(role.Id)
		listResp.Roles = append(listResp.Roles, roleInfo)
	}
	return listResp, nil
}

func ContactRoleFilter(orgID int64, filterReq req.RoleFilterReq) (*resp.RoleFilterResp, error) {
	query := strings.TrimSpace(filterReq.Query)
	roleList, err := domain.GetRoleListByOrgWithFilter(orgID, query)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	roleGroupIDs := make([]int64, 0)
	for _, role := range roleList {
		roleGroupIDs = append(roleGroupIDs, role.RoleGroupId)
	}
	roleGroupIDs = slice.SliceUniqueInt64(roleGroupIDs)
	roleGroupList, err := domain.GetRoleGroupListIds(orgID, roleGroupIDs)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	roleGroupMap := map[int64]bo.RoleGroupBo{}
	for _, roleGroup := range roleGroupList {
		roleGroupMap[roleGroup.Id] = roleGroup
	}
	roleUserCountMap, err := domain.GetRoleCounts(orgID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	roleNodeList := make([]resp.RoleNode, 0)
	for _, role := range roleList {
		roleNode := resp.RoleNode{
			ID:        role.Id,
			Name:      role.Name,
			GroupID:   role.RoleGroupId,
			UserCount: roleUserCountMap[role.Id],
		}
		roleGroup, ok := roleGroupMap[role.RoleGroupId]
		if ok {
			roleNode.GroupName = roleGroup.Name
		}
		roleNodeList = append(roleNodeList, roleNode)
	}
	return &resp.RoleFilterResp{
		RoleList: roleNodeList,
	}, nil
}
