package open_service

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/model/resp/open_resp"
	"upper.io/db.v3"
)

func _copyMemberBoToMemberResp(memberBaseInfo *bo.OrgMemberBaseInfoBo, member *open_resp.OrgMemberBaseResp) {
	member.UserID = memberBaseInfo.UserId
	member.LoginName = memberBaseInfo.LoginName
	member.Name = memberBaseInfo.Name
	member.NamePy = memberBaseInfo.NamePinyin
	member.Avatar = memberBaseInfo.Avatar
	member.Email = memberBaseInfo.Email
	member.PhoneRegion = memberBaseInfo.MobileRegion
	member.PhoneNumber = memberBaseInfo.Mobile
	member.EmpNo = memberBaseInfo.EmpNo
	member.WeiboIds = memberBaseInfo.WeiboIds
	member.Status = memberBaseInfo.Status
	member.OrgOwner = memberBaseInfo.OrgOwner
	member.Creator = memberBaseInfo.Creator
	member.CreateTime = memberBaseInfo.CreateTime
	member.Updator = memberBaseInfo.Updator
	member.UpdateTime = memberBaseInfo.UpdateTime
}

// GetOrgMemberList 获取组织内成员列表
func GetOrgMemberList(orgId int64) ([]open_resp.OrgMemberBaseResp, errs.SystemErrorInfo) {
	memberBaseInfoList, dbErr := domain.GetOrgMemberBaseInfoListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	users := make([]open_resp.OrgMemberBaseResp, 0)
	userIds := make([]int64, 0)
	for _, memberBaseInfo := range memberBaseInfoList {
		member := open_resp.OrgMemberBaseResp{}
		_copyMemberBoToMemberResp(&memberBaseInfo, &member)
		users = append(users, member)
		userIds = append(userIds, member.UserID)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		user := &users[i]
		if deptAndRole, ok := deptAndRoleMap[user.UserID]; ok {
			user.UserBindDeptAndRoleResp = deptAndRole
		} else {
			user.DeptList = []open_resp.UserDeptBindData{}
			user.RoleList = []open_resp.UserRoleBindData{}
		}
	}
	return users, nil
}

// GetOrgMemberListByQueryCond 获取组织内成员列表
func GetOrgMemberListByQueryCond(orgId int64, reqParam open_req.MemberQueryReq) ([]open_resp.OrgMemberBaseResp, errs.SystemErrorInfo) {
	memberBaseInfoList, dbErr := domain.GetOrgMemberBaseInfoListByQueryCond(orgId, reqParam)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	users := make([]open_resp.OrgMemberBaseResp, 0)
	userIds := make([]int64, 0)
	for _, memberBaseInfo := range memberBaseInfoList {
		member := open_resp.OrgMemberBaseResp{}
		_copyMemberBoToMemberResp(&memberBaseInfo, &member)
		users = append(users, member)
		userIds = append(userIds, member.UserID)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		user := &users[i]
		if deptAndRole, ok := deptAndRoleMap[user.UserID]; ok {
			user.UserBindDeptAndRoleResp = deptAndRole
		} else {
			user.DeptList = []open_resp.UserDeptBindData{}
			user.RoleList = []open_resp.UserRoleBindData{}
		}
	}
	return users, nil
}

// GetOrgMemberListByUserIds 获取组织内成员列表 根据成员ID列表查询
func GetOrgMemberListByUserIds(orgId int64, userIds []int64) ([]open_resp.OrgMemberBaseResp, errs.SystemErrorInfo) {
	if len(userIds) == 0 {
		return []open_resp.OrgMemberBaseResp{}, nil
	}
	//去重
	userIds = slice.SliceUniqueInt64(userIds)
	memberBaseInfoList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	users := make([]open_resp.OrgMemberBaseResp, 0)
	userIds = make([]int64, 0)
	for _, memberBaseInfo := range memberBaseInfoList {
		member := open_resp.OrgMemberBaseResp{}
		_copyMemberBoToMemberResp(&memberBaseInfo, &member)
		users = append(users, member)
		userIds = append(userIds, member.UserID)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		user := &users[i]
		if deptAndRole, ok := deptAndRoleMap[user.UserID]; ok {
			user.UserBindDeptAndRoleResp = deptAndRole
		} else {
			user.DeptList = []open_resp.UserDeptBindData{}
			user.RoleList = []open_resp.UserRoleBindData{}
		}
	}
	return users, nil
}

// GetOrgMemberByUserId 获取组织内成员信息
func GetOrgMemberByUserId(orgId int64, userId int64) (*open_resp.OrgMemberBaseResp, errs.SystemErrorInfo) {
	memberBaseInfo, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, nil
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	member := open_resp.OrgMemberBaseResp{}
	_copyMemberBoToMemberResp(memberBaseInfo, &member)

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, []int64{member.UserID})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if deptAndRole, ok := deptAndRoleMap[member.UserID]; ok {
		member.UserBindDeptAndRoleResp = deptAndRole
	} else {
		member.DeptList = []open_resp.UserDeptBindData{}
		member.RoleList = []open_resp.UserRoleBindData{}
	}

	return &member, nil
}

// GetOrgMemberListByExcludeIds 获取组织内成员列表 排除指定成员ID
func GetOrgMemberListByExcludeIds(orgId int64, excludeIds []int64) ([]open_resp.OrgMemberBaseResp, errs.SystemErrorInfo) {
	//去重
	excludeIds = slice.SliceUniqueInt64(excludeIds)
	memberBaseInfoList, dbErr := domain.GetOrgMemberBaseInfoListByExcludeIds(orgId, excludeIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	users := make([]open_resp.OrgMemberBaseResp, 0)
	userIds := make([]int64, 0)

	for _, memberBaseInfo := range memberBaseInfoList {
		member := open_resp.OrgMemberBaseResp{}
		_copyMemberBoToMemberResp(&memberBaseInfo, &member)
		users = append(users, member)
		userIds = append(userIds, member.UserID)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		user := &users[i]
		if deptAndRole, ok := deptAndRoleMap[user.UserID]; ok {
			user.UserBindDeptAndRoleResp = deptAndRole
		} else {
			user.DeptList = []open_resp.UserDeptBindData{}
			user.RoleList = []open_resp.UserRoleBindData{}
		}
	}
	return users, nil
}

// GetOrgMemberListByDept 获取部门下成员列表
func GetOrgMemberListByDept(orgId int64, deptId int64) ([]open_resp.UserDeptBindResp, errs.SystemErrorInfo) {
	deptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByDept(orgId, deptId)

	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userIds := make([]int64, 0)
	for _, info := range deptBindInfoList {
		userIds = append(userIds, info.UserId)
	}

	// 获取成员信息
	memberMap := make(map[int64]bo.OrgMemberBaseInfoBo)
	memberList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	for _, member := range memberList {
		memberMap[member.UserId] = member
	}
	deptBindRespList := make([]open_resp.UserDeptBindResp, 0)
	for _, info := range deptBindInfoList {
		if member, ok := memberMap[info.UserId]; ok {
			infoResp := open_resp.UserDeptBindResp{}
			_copyBindBoToUserDeptData(&info, &infoResp.UserDeptBindData)
			infoResp.UserId = member.UserId
			infoResp.Nickname = member.Name
			infoResp.Status = member.Status
			deptBindRespList = append(deptBindRespList, infoResp)
		}
	}
	return deptBindRespList, nil
}

// GetDeptHaveMemberList 获取部门列表（附带成员信息）
func GetDeptHaveMemberList(orgId int64) ([]open_resp.DeptMemberListResp, errs.SystemErrorInfo) {
	// 获取机构列表
	deptList, err := GetDeptList(orgId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// 获取成员信息
	memberMap := make(map[int64]bo.OrgMemberBaseInfoBo)
	memberList, dbErr := domain.GetOrgMemberBaseInfoListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	for _, member := range memberList {
		memberMap[member.UserId] = member
	}

	// 获取部门和用户绑定关系
	deptUserMap := make(map[int64][]int64)
	deptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	for _, deptUser := range deptBindInfoList {
		deptUserMap[deptUser.DepartmentId] = append(deptUserMap[deptUser.DepartmentId], deptUser.UserId)
	}

	deptListResp := make([]open_resp.DeptMemberListResp, 0)
	// 构建
	for _, dept := range deptList {
		deptResp := open_resp.DeptMemberListResp{
			DeptInfoResp: dept,
			UserList:     []open_resp.SimpleMemberInfo{},
		}
		if userIds, ok := deptUserMap[dept.Id]; ok {
			for _, uid := range userIds {
				if member, ok := memberMap[uid]; ok {
					deptResp.UserList = append(deptResp.UserList, open_resp.SimpleMemberInfo{
						UserId:   member.UserId,
						Nickname: member.Name,
						Status:   member.Status,
					})
				}
			}
		}
		deptListResp = append(deptListResp, deptResp)
	}
	return deptListResp, nil
}

// GetUserDeptRoleBindListByUsers 根据成员列表，查询绑定的部门和角色列表
func GetUserDeptRoleBindListByUsers(orgId int64, userIds []int64) (map[int64]open_resp.UserBindDeptAndRoleResp, errs.SystemErrorInfo) {
	userIds = slice.SliceUniqueInt64(userIds)

	// 查询绑定角色信息
	deptBindMap, err := GetUserDeptBindDataListByUsers(orgId, userIds)
	if err != nil {
		return nil, err
	}
	roleBindMap, err := GetUserRoleBindDataListByUsers(orgId, userIds)
	if err != nil {
		return nil, err
	}
	userIds = []int64{}
	for k, _ := range deptBindMap {
		userIds = append(userIds, k)
	}
	for k, _ := range roleBindMap {
		userIds = append(userIds, k)
	}
	userIds = slice.SliceUniqueInt64(userIds)
	bindMap := make(map[int64]open_resp.UserBindDeptAndRoleResp)
	for _, uid := range userIds {
		bind := open_resp.UserBindDeptAndRoleResp{}
		if deptList, ok := deptBindMap[uid]; ok {
			bind.DeptList = deptList
		} else {
			bind.DeptList = []open_resp.UserDeptBindData{}
		}

		if roleList, ok := roleBindMap[uid]; ok {
			bind.RoleList = roleList
		} else {
			bind.RoleList = []open_resp.UserRoleBindData{}
		}
		bindMap[uid] = bind
	}
	return bindMap, nil
}

// GetUserAuthBaseInfo 获取组织内成员列表
func GetUserAuthBaseInfo(orgId int64, userId int64) (*open_resp.OrgUserAuthBaseResp, errs.SystemErrorInfo) {
	userBaseInfo, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	// 管理组
	var isOwner = userBaseInfo.OrgOwner == userId
	isSysAdmin := false
	isSubAdmin := false
	manageGroup, dbErr := domain.GetManageGroupListByUser(orgId, userId)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}
	if manageGroup != nil {
		if manageGroup.LangCode == consts.ManageGroupSys {
			isSysAdmin = true
		} else {
			isSubAdmin = true
		}
	}

	// 角色信息
	userRoleBindList, dbErr := domain.GetUserRoleBindListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	refRoleIds := make([]int64, len(userRoleBindList))
	for i, role := range userRoleBindList {
		refRoleIds[i] = role.RoleId
	}
	// 部门信息
	depts, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	refDeptIds := make([]int64, len(depts))
	for i, dept := range depts {
		refDeptIds[i] = dept.DepartmentId
	}
	userInfo := &open_resp.OrgUserAuthBaseResp{
		UserID:     userBaseInfo.UserId,
		OrgOwner:   userBaseInfo.OrgOwner,
		IsOrgOwner: isOwner,
		IsSysAdmin: isSysAdmin,
		IsSubAdmin: isSubAdmin,
		RoleIds:    refRoleIds,
		DeptIds:    refDeptIds,
	}
	return userInfo, nil
}
