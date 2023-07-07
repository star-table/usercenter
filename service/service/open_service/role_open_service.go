package open_service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/resp/open_resp"
	"upper.io/db.v3"
)

func _copyBindBoToUserRoleData(infoBo *bo.UserRoleBo, infoResp *open_resp.UserRoleBindData) {
	infoResp.RoleId = infoBo.RoleId
	infoResp.RoleName = infoBo.RoleName
}

func _copyRoleBoToRoleInfoResp(infoBo *po.PpmRolRole, infoResp *open_resp.RoleInfoResp) {
	infoResp.Id = infoBo.Id
	infoResp.OrgId = infoBo.OrgId
	infoResp.RoleGroupId = infoBo.RoleGroupId
	infoResp.Name = infoBo.Name
}

// GetRoleList 角色列表
func GetRoleList(orgId int64) ([]open_resp.RoleInfoResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	//去重
	roleBos, dbErr := domain.GetRoleListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	roleList := make([]open_resp.RoleInfoResp, 0)
	for _, info := range roleBos {
		dept := open_resp.RoleInfoResp{}
		_copyRoleBoToRoleInfoResp(&info, &dept)
		roleList = append(roleList, dept)
	}
	return roleList, nil
}

// GetRoleListByIds 根据id列表获取角色列表
func GetRoleListByIds(orgId int64, roleIds []int64) ([]open_resp.RoleInfoResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	if len(roleIds) == 0 {
		return []open_resp.RoleInfoResp{}, nil
	}
	//去重
	roleIds = slice.SliceUniqueInt64(roleIds)
	roleBos, dbErr := domain.GetRoleListByIds(orgId, roleIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	roleList := make([]open_resp.RoleInfoResp, 0)
	for _, info := range roleBos {
		dept := open_resp.RoleInfoResp{}
		_copyRoleBoToRoleInfoResp(&info, &dept)
		roleList = append(roleList, dept)
	}
	return roleList, nil
}

// GetUserRoleBindListByUser 根据成员id查询角色列表
func GetUserRoleBindListByUser(orgId int64, userId int64) ([]open_resp.UserRoleBindResp, errs.SystemErrorInfo) {
	member, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return []open_resp.UserRoleBindResp{}, nil
		}
		return nil, errs.MysqlOperateError
	}
	// 查询绑定角色信息
	userRoleBindList, dbErr := domain.GetUserRoleBindListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respList := make([]open_resp.UserRoleBindResp, 0)
	for _, infoBo := range userRoleBindList {
		infoResp := open_resp.UserRoleBindResp{}
		_copyBindBoToUserRoleData(&infoBo, &infoResp.UserRoleBindData)
		infoResp.UserId = infoBo.UserId
		infoResp.Nickname = member.Name
		infoResp.Status = member.Status
		respList = append(respList, infoResp)
	}
	return respList, nil
}

// GetUserRoleBindListByUsers 根据成员列表，查询角色列表
func GetUserRoleBindListByUsers(orgId int64, userIds []int64) ([]open_resp.UserRoleBindResp, errs.SystemErrorInfo) {
	// 获取成员信息
	memberMap := make(map[int64]bo.OrgMemberBaseInfoBo)
	memberList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	for _, member := range memberList {
		memberMap[member.UserId] = member
	}
	// 查询绑定角色信息
	userRoleBindList, dbErr := domain.GetUserRoleBindListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respList := make([]open_resp.UserRoleBindResp, 0)
	for _, infoBo := range userRoleBindList {
		if member, ok := memberMap[infoBo.UserId]; ok {
			infoResp := open_resp.UserRoleBindResp{}
			_copyBindBoToUserRoleData(&infoBo, &infoResp.UserRoleBindData)
			infoResp.UserId = member.UserId
			infoResp.Nickname = member.Name
			infoResp.Status = member.Status
			respList = append(respList, infoResp)
		}
	}
	return respList, nil
}

// GetUserRoleBindDataListByUsers 根据成员列表，查询角色列表
func GetUserRoleBindDataListByUsers(orgId int64, userIds []int64) (map[int64][]open_resp.UserRoleBindData, errs.SystemErrorInfo) {
	// 查询绑定角色信息
	userRoleBindList, dbErr := domain.GetUserRoleBindListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respMap := make(map[int64][]open_resp.UserRoleBindData, 0)
	for _, infoBo := range userRoleBindList {
		infoResp := open_resp.UserRoleBindData{}
		_copyBindBoToUserRoleData(&infoBo, &infoResp)
		respMap[infoBo.UserId] = append(respMap[infoBo.UserId], infoResp)
	}
	return respMap, nil
}
