package inner_service

import (
	"fmt"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/jsonx"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
)

// GetUserListByIds 内部调用
func GetUserListByIds(orgId int64, userIds []int64, checkStatus, status, isDelete int) ([]inner_resp.UserInfoInnerResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	if len(userIds) == 0 {
		return []inner_resp.UserInfoInnerResp{}, nil
	}
	//去重
	userIds = slice.SliceUniqueInt64(userIds)
	userInfos, dbErr := domain.GetOrgMemberBaseInfoList(orgId, userIds, checkStatus, status, isDelete)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	users := make([]inner_resp.UserInfoInnerResp, 0)
	for _, userInfo := range userInfos {
		users = append(users, inner_resp.UserInfoInnerResp{
			Id:          userInfo.UserId,
			LoginName:   userInfo.LoginName,
			Name:        userInfo.Name,
			NamePy:      userInfo.NamePinyin,
			Avatar:      userInfo.Avatar,
			Email:       userInfo.Email,
			PhoneRegion: userInfo.MobileRegion,
			PhoneNumber: userInfo.Mobile,
			Status:      userInfo.Status,
			Creator:     userInfo.Creator,
			CreateTime:  userInfo.CreateTime,
			Updator:     userInfo.Updator,
			UpdateTime:  userInfo.UpdateTime,
			IsDelete:    userInfo.IsDelete,
			Type:        userInfo.Type,
		})
	}

	// 查询部门和角色
	deptAndRoleMap, err := GetUserDeptRoleBindListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		user := &users[i]
		if deptAndRole, ok := deptAndRoleMap[user.Id]; ok {
			user.UserBindDeptAndRoleResp = deptAndRole
		} else {
			user.DeptList = []inner_resp.UserDeptBindData{}
			user.RoleList = []inner_resp.UserRoleBindData{}
		}
	}
	return users, nil
}

// GetUserDeptRoleBindListByUsers 根据成员列表，查询绑定的部门和角色列表
func GetUserDeptRoleBindListByUsers(orgId int64, userIds []int64) (map[int64]inner_resp.UserBindDeptAndRoleResp, errs.SystemErrorInfo) {
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
	bindMap := make(map[int64]inner_resp.UserBindDeptAndRoleResp)
	for _, uid := range userIds {
		bind := inner_resp.UserBindDeptAndRoleResp{}
		if deptList, ok := deptBindMap[uid]; ok {
			bind.DeptList = deptList
		} else {
			bind.DeptList = []inner_resp.UserDeptBindData{}
		}

		if roleList, ok := roleBindMap[uid]; ok {
			bind.RoleList = roleList
		} else {
			bind.RoleList = []inner_resp.UserRoleBindData{}
		}
		bindMap[uid] = bind
	}
	return bindMap, nil
}

// GetUserAuthorityByUserId 获取用户权限信息。兼容无码授权模式。
func GetUserAuthorityByUserId(orgId, userId int64) (*inner_resp.UserAuthorityInnerResp, errs.SystemErrorInfo) {
	baseOrgInfo, errSys := domain.GetBaseOrgInfo("", orgId)
	if errSys != nil {
		logger.Error(errSys)
		return nil, errSys
	}
	baseUserInfo, errSys := domain.GetBaseUserInfo("", orgId, userId)
	if errSys != nil {
		logger.Error(errSys)
		return nil, errSys
	}
	if baseUserInfo.OrgUserStatus != consts.AppStatusEnable {
		if baseOrgInfo.OutOrgId != "" {
			return nil, errs.OrgUserInvalidErr
		}
		return nil, errs.OrgUserUnabledErr
	}
	if baseUserInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
		return nil, errs.OrgUserCheckStatusUnabledErr
	}
	if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted {
		return nil, errs.OrgUserDeletedErr
	}

	// 管理组
	var isOwner = baseOrgInfo.OrgOwnerId == userId
	isSysAdmin := false
	isSubAdmin := false
	hasDeptOptAuth := isOwner
	hasRoleOptAuth := isOwner
	hasAppPackageOptAuth := isOwner
	manageDepts := make([]int64, 0)
	manageRoles := make([]int64, 0)
	manageAppPackages := make([]int64, 0)
	manageApps := make([]int64, 0)
	optAuth := make([]string, 0)

	groups, dbErr := domain.GetManageGroupListByUsers(orgId, []int64{userId})
	if dbErr != nil {
		// 不在管理组则忽视
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}
	if groups == nil {
		groups = make([]po.LcPerManageGroup, 0)
	}
	// 如果用户没有角色，则默认“团队成员”角色
	if len(groups) < 1 {
		normalUserGroup, dbErr := domain.GetOrgDefaultAdminGroupForPolaris(orgId)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
		if normalUserGroup != nil {
			groups = append(groups, *normalUserGroup)
		}
	}

	for _, group := range groups {
		if group.LangCode == consts.ManageGroupSys {
			isSysAdmin = true
			hasDeptOptAuth = true
			hasRoleOptAuth = true
			hasAppPackageOptAuth = true
		} else {
			if !isSubAdmin {
				isSubAdmin, _ = slice.Contain([]string{consts.ManageGroupSub, consts.ManageGroupSubNormalAdmin}, group.LangCode)
			}
			manageDeptsTemp := make([]int64, 0)
			manageRolesTemp := make([]int64, 0)
			manageAppPackagesTemp := make([]int64, 0)
			manageAppsTemp := make([]int64, 0)
			optAuthTemp := make([]interface{}, 0)

			_ = jsonx.FromJson(group.DeptIds, &manageDeptsTemp)
			_ = jsonx.FromJson(group.RoleIds, &manageRolesTemp)
			_ = jsonx.FromJson(group.AppPackageIds, &manageAppPackagesTemp)
			_ = jsonx.FromJson(group.AppIds, &manageAppsTemp)
			_ = jsonx.FromJson(group.OptAuth, &optAuthTemp)

			for _, v := range optAuthTemp {
				if v == consts.DeptOpt {
					hasDeptOptAuth = true
				} else if v == consts.RoleOpt {
					hasRoleOptAuth = true
				} else if v == consts.AppPackageOpt {
					hasAppPackageOptAuth = true
				}
			}
			for _, optAuthCode := range optAuthTemp {
				optAuthCodeStr := fmt.Sprintf("%v", optAuthCode)
				//if optAuthCodeStr != consts.MenuPermissionTrend {
				optAuth = append(optAuth, optAuthCodeStr)
				//}
			}
			manageDepts = append(manageDepts, manageDeptsTemp...)
			manageRoles = append(manageRoles, manageRolesTemp...)
			manageAppPackages = append(manageAppPackages, manageAppPackagesTemp...)
			manageApps = append(manageApps, manageAppsTemp...)
		}
	}
	optAuth = domain.TransferOperationArr(optAuth)

	// 查询子部门全部
	if len(manageDepts) > 0 {
		manageDeptIds, dbErr := domain.GetDeptAndChildrenIds(orgId, manageDepts)
		if dbErr != nil {
			return nil, errs.MysqlOperateError
		}
		manageDepts = manageDeptIds
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
	if len(refDeptIds) > 0 {
		parentIds, err := domain.GetDeptParentIds(orgId, refDeptIds)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		refDeptIds = append(refDeptIds, parentIds...)
	}

	userInfo := &inner_resp.UserAuthorityInnerResp{
		OrgId:                baseOrgInfo.OrgId,
		UserId:               userId,
		IsOrgOwner:           isOwner,
		IsSysAdmin:           isSysAdmin,
		IsSubAdmin:           isSubAdmin,
		OrgSourceChannel:     baseOrgInfo.SourceChannel,
		RefRoleIds:           refRoleIds,
		RefDeptIds:           refDeptIds,
		HasDeptOptAuth:       hasDeptOptAuth,
		HasRoleOptAuth:       hasRoleOptAuth,
		HasAppPackageOptAuth: hasAppPackageOptAuth,
		ManageDepts:          manageDepts,
		ManageRoles:          manageRoles,
		ManageAppPackages:    manageAppPackages,
		ManageApps:           manageApps,
		OptAuth:              optAuth,
		//IsOutCollaborator:    userBaseInfo.Type == 2,
	}
	return userInfo, nil
}

// GetUserAuthorityByUserIdSimple 查询用户授权信息。用于极星的简化逻辑的授权信息查询
// 该接口来源于 `GetUserAuthorityByUserId` 函数。完整版可以搜索查看。
// 该接口主要用于极星的授权信息查询。目前提供的信息只有：`OptAuth`, `IsSysAdmin`, `IsSubAdmin`，如果需要其他，请按需要添加，并更新注释文档。
func GetUserAuthorityByUserIdSimple(orgId, userId int64) (*inner_resp.UserAuthorityInnerResp, errs.SystemErrorInfo) {
	userBaseInfo, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgMemberNotExistOrDisable
		}
		return nil, errs.MysqlOperateError
	}
	// 查询组织信息
	orgInfo, err := domain.GetBaseOrgInfo("", orgId)
	if err != nil {
		return nil, err
	}

	// 管理组
	var isOwner = userBaseInfo.OrgOwner == userId
	isSysAdmin := false
	isSubAdmin := false
	hasDeptOptAuth := isOwner
	hasRoleOptAuth := isOwner
	hasAppPackageOptAuth := isOwner
	optAuth := make([]string, 0)
	manageApps := make([]int64, 0)

	groups, dbErr := domain.GetManageGroupListByUsers(orgId, []int64{userId})
	if dbErr != nil {
		// 不在管理组则忽视
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}
	if groups == nil {
		groups = make([]po.LcPerManageGroup, 0)
	}
	// 如果用户没有角色，则默认“团队成员”角色
	if len(groups) < 1 {
		normalUserGroup, dbErr := domain.GetOrgDefaultAdminGroupForPolaris(orgId)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
		if normalUserGroup != nil {
			groups = append(groups, *normalUserGroup)
		}
	}

	for _, group := range groups {
		if group.LangCode == consts.ManageGroupSys {
			isSysAdmin = true
			hasDeptOptAuth = true
			hasRoleOptAuth = true
			hasAppPackageOptAuth = true
		} else {
			if !isSubAdmin {
				isSubAdmin, _ = slice.Contain([]string{consts.ManageGroupSub, consts.ManageGroupSubNormalAdmin}, group.LangCode)
			}
			optAuthTemp := make([]interface{}, 0)
			manageAppsTemp := make([]int64, 0)

			_ = jsonx.FromJson(group.OptAuth, &optAuthTemp)
			_ = jsonx.FromJson(group.AppIds, &manageAppsTemp)

			for _, v := range optAuthTemp {
				if v == consts.DeptOpt {
					hasDeptOptAuth = true
				} else if v == consts.RoleOpt {
					hasRoleOptAuth = true
				} else if v == consts.AppPackageOpt {
					hasAppPackageOptAuth = true
				}
			}
			manageApps = append(manageApps, manageAppsTemp...)
			for _, optAuthCode := range optAuthTemp {
				optAuthCodeStr := fmt.Sprintf("%v", optAuthCode)
				//if optAuthCodeStr != consts.MenuPermissionTrend {
				optAuth = append(optAuth, optAuthCodeStr)
				//}
			}
		}
	}
	optAuth = domain.TransferOperationArr(optAuth)
	userInfo := &inner_resp.UserAuthorityInnerResp{
		OrgId:                userBaseInfo.OrgId,
		UserId:               userBaseInfo.UserId,
		IsOrgOwner:           isOwner,
		IsSysAdmin:           isSysAdmin,
		IsSubAdmin:           isSubAdmin,
		OrgSourceChannel:     orgInfo.SourceChannel,
		HasDeptOptAuth:       hasDeptOptAuth,
		HasRoleOptAuth:       hasRoleOptAuth,
		HasAppPackageOptAuth: hasAppPackageOptAuth,
		OptAuth:              optAuth,
		ManageApps:           manageApps,
		IsOutCollaborator:    userBaseInfo.Type == 2,
	}

	return userInfo, nil
}

func GetMemberSimpleInfo(orgId int64, memberType int, needDelete bool) (*inner_resp.GetMemberSimpleInfoResp, errs.SystemErrorInfo) {
	var info []inner_resp.SimpleInfo
	var err errs.SystemErrorInfo
	switch memberType {
	case 1: //人员
		info, err = domain.GetOrgUserSimpleInfo(orgId, needDelete)
	case 2: //部门
		info, err = domain.GetDeptSimpleInfo(orgId)
	case 3: //角色
		info, err = domain.GetRoleSimpleInfo(orgId)
	default:
		return nil, errs.ParamError
	}

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if info == nil {
		return &inner_resp.GetMemberSimpleInfoResp{Data: []inner_resp.SimpleInfo{}}, nil

	}
	return &inner_resp.GetMemberSimpleInfoResp{Data: info}, nil
}

func GetRepeatMember(orgId int64) (*inner_resp.RepeatMemberResp, errs.SystemErrorInfo) {
	res := &inner_resp.RepeatMemberResp{
		User:       nil,
		Department: nil,
		Role:       nil,
	}

	users, usersErr := domain.GetRepeatUserInfo(orgId)
	if usersErr != nil {
		logger.Error(usersErr)
		return nil, usersErr
	}
	res.User = users

	depts, deptErr := domain.GetRepeatDeptInfo(orgId)
	if deptErr != nil {
		logger.Error(deptErr)
		return nil, deptErr
	}
	res.Department = depts

	roles, rolesErr := domain.GetRepeatRoleInfo(orgId)
	if rolesErr != nil {
		logger.Error(rolesErr)
		return nil, rolesErr
	}
	res.Role = roles

	return res, nil
}

// GetUsersCouldManage [极星] 获取该组织的超管用户以及普通管理员的所有用户。
// 通过权限的 langCode 获取对应的用户
func GetUsersCouldManage(orgId int64, appId int64) (*inner_resp.GetUsersCouldManageResp, errs.SystemErrorInfo) {
	userList := make([]bo.SimpleUserInfoBo, 0)
	groups, dbErr := domain.GetAdminGroupsByLangCode(orgId, []string{consts.ManageGroupSys, consts.ManageGroupSubNormalAdmin})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	adminUserIds := make([]int64, 0)
	for _, group := range groups {
		if group.LangCode == consts.ManageGroupSubNormalAdmin {
			// 如果是普通管理员，需要过滤掉 没有管理该app权限的user
			appIds := []int64{}
			json.FromJsonIgnoreError(group.AppIds, &appIds)
			if len(appIds) > 0 {
				if slice.IntContains(appIds, appId) || appIds[0] == -1 {
					// 将 userIds json 转换为 []int64
					userIds := domain.TransferUserIdsFromUserIdJson(group.UserIds)
					if userIds != nil && len(userIds) > 0 {
						adminUserIds = append(adminUserIds, userIds...)
					}
				}
			}
		} else if group.LangCode == consts.ManageGroupSys {
			userIds := domain.TransferUserIdsFromUserIdJson(group.UserIds)
			adminUserIds = append(adminUserIds, userIds...)
		}

	}
	userPos, dbErr := domain.GetUserListByIds(slice.SliceUniqueInt64(adminUserIds))
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	copyErr := copyer.Copy(userPos, &userList)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &inner_resp.GetUsersCouldManageResp{
		List: userList,
	}, nil
}

// GetMemberSimpleInfoList 组织成员简单信息列表
func GetMemberSimpleInfoList(req inner_req.MemberSimpleInfoListReq) (*inner_resp.MemberSimpleInfoListResp, errs.SystemErrorInfo) {

	union := &db.Union{}

	cond := db.Cond{
		consts.TcOrgId: req.OrgId,
		// consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId: db.In(req.UserIds),
	}

	// if len(req.UserIds) > 0 {
	// 	cond[consts.TcUserId] = db.In(req.UserIds)
	// }

	if req.Status != nil {
		cond[consts.TcStatus] = *req.Status
	}

	// 开启权限过滤
	union1 := &db.Union{}

	var order interface{}
	if req.Order != "" {
		order = db.Raw(req.Order)
	} else {
		order = db.Raw(" create_time asc ")
	}

	var orgMemberList []po.PpmOrgUserOrganization
	total, dbErr := store.Mysql.SelectAllByCondWithPageAndOrderUnion(consts.TableUserOrganization, cond, req.Page, req.Size, order, &orgMemberList, union, union1)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	userIds := make([]int64, 0)
	for _, orgMember := range orgMemberList {
		userIds = append(userIds, orgMember.UserId)
	}

	// 获取用户信息
	userInfoList, dbErr := domain.GetUserListByIds(userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userInfoMap := make(map[int64]po.PpmOrgUser)
	for _, info := range userInfoList {
		userInfoMap[info.Id] = info
	}

	res := &inner_resp.MemberSimpleInfoListResp{
		Total: int64(total),
		List:  []*inner_resp.MemberSimpleInfo{},
	}

	//查询组织创建者
	var orgInfo po.PpmOrgOrganization
	dbErr = store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcId: req.OrgId,
	}, &orgInfo)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgNotExist
		} else {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}

	for _, orgUser := range orgMemberList {
		if _, ok := userInfoMap[orgUser.UserId]; !ok {
			continue
		}
		userInfo := userInfoMap[orgUser.UserId]
		u := &inner_resp.MemberSimpleInfo{
			Id:       userInfo.Id,
			Name:     userInfo.Name,
			Status:   orgUser.Status,
			Type:     consts.MemberTypeUser,
			IsDelete: orgUser.IsDelete,
		}
		res.List = append(res.List, u)
	}

	return res, nil
}

// GetCommAdminMangeApps 普通管理员可以管理的该项目的app
func GetCommAdminMangeApps(orgId int64) ([]*inner_resp.GetCommAdminMangeAppsData, errs.SystemErrorInfo) {
	groups, err := domain.GetAdminGroupsByLangCode(orgId, []string{consts.ManageGroupSubNormalAdmin})
	if err != nil {
		logger.ErrorF("[GetCommAdminMangeApps] err:%v, orgId:%v", err, orgId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	userAppIdsMap := map[int64][]int64{}

	for _, group := range groups {
		var userIds []int64
		json.FromJsonIgnoreError(group.UserIds, &userIds)

		var appIds []int64
		json.FromJsonIgnoreError(group.AppIds, &appIds)

		for _, id := range userIds {
			userAppIdsMap[id] = append(userAppIdsMap[id], appIds...)
		}

	}
	manageApps := []*inner_resp.GetCommAdminMangeAppsData{}
	for i, v := range userAppIdsMap {
		manageApps = append(manageApps, &inner_resp.GetCommAdminMangeAppsData{
			UserId: i,
			AppIds: v,
		})
	}

	return manageApps, nil

}
