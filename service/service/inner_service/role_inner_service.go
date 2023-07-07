package inner_service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
)

/**
角色 (内部调用)
*/
// GetRoleList 角色列表 内部调用
func GetRoleList(orgId int64) ([]inner_resp.RoleInfoInnerResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	//去重
	roleBos, dbErr := domain.GetRoleListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	roles := make([]inner_resp.RoleInfoInnerResp, 0)
	for _, info := range roleBos {
		var role inner_resp.RoleInfoInnerResp
		_ = copyer.Copy(info, &role)
		roles = append(roles, role)
	}
	return roles, nil
}

// GetRoleListByIds 根据id列表获取角色列表 内部调用
func GetRoleListByIds(orgId int64, roleIds []int64, status, isDelete int) ([]inner_resp.RoleInfoInnerResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	if len(roleIds) == 0 {
		return []inner_resp.RoleInfoInnerResp{}, nil
	}
	//去重
	roleIds = slice.SliceUniqueInt64(roleIds)
	roleBos, dbErr := domain.GetRoleListByIdsAndStatus(orgId, roleIds, status, isDelete)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	roles := make([]inner_resp.RoleInfoInnerResp, 0)
	for _, info := range roleBos {
		var role inner_resp.RoleInfoInnerResp
		_ = copyer.Copy(info, &role)
		roles = append(roles, role)
	}
	return roles, nil
}

func _copyBindBoToUserRoleData(infoBo *bo.UserRoleBo, infoResp *inner_resp.UserRoleBindData) {
	infoResp.RoleId = infoBo.RoleId
	infoResp.RoleName = infoBo.RoleName
}

// GetUserRoleBindDataListByUsers 根据成员列表，查询角色列表
func GetUserRoleBindDataListByUsers(orgId int64, userIds []int64) (map[int64][]inner_resp.UserRoleBindData, errs.SystemErrorInfo) {
	// 查询绑定角色信息
	userRoleBindList, dbErr := domain.GetUserRoleBindListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respMap := make(map[int64][]inner_resp.UserRoleBindData, 0)
	for _, infoBo := range userRoleBindList {
		infoResp := inner_resp.UserRoleBindData{}
		_copyBindBoToUserRoleData(&infoBo, &infoResp)
		respMap[infoBo.UserId] = append(respMap[infoBo.UserId], infoResp)
	}
	return respMap, nil
}

func GetRoleUserIds(orgId int64) (*inner_resp.GetRoleUserIdsResp, errs.SystemErrorInfo) {
	data, err := domain.GetRoleUserIds(orgId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &inner_resp.GetRoleUserIdsResp{Data: data}, nil
}
