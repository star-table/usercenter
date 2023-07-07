package service

import (
	"strings"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/pkg/util/copyer"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// CreateManageGroup 创建管理组
func CreateManageGroup(orgId, operator int64, reqParam req.CreateManageGroup) (int64, errs.SystemErrorInfo) {
	if orgId == 0 {
		return 0, errs.OrgNotExist
	}

	// 验证名称
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	if !format.VerifyManageGroupNameFormat(reqParam.Name) {
		return 0, errs.ManageGroupNameLenErr
	}
	// 不允许和系统管理组重名
	if ok, _ := slice.Contain(consts.DefaultManageGroupName, reqParam.Name); ok {
		return 0, errs.DefaultManageGroupErr
	}

	var id int64
	var dbErr error
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		id, dbErr = domain.CreateManageGroup(orgId, operator, reqParam, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}
	return id, nil
}

// UpdateManageGroup 修改管理组
func UpdateManageGroup(orgId, operator, groupId int64, reqParam req.UpdateManageGroup) (bool, errs.SystemErrorInfo) {
	if orgId == 0 {
		return false, errs.OrgNotExist
	}
	manageGroup, dbErr := domain.GetManageGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.ManageGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	// 不允许修改系统管理组
	if consts.IsDefaultManageGroup(manageGroup.LangCode) {
		return false, errs.DefaultManageGroupCantModify
	}
	// 验证名称
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	if !format.VerifyManageGroupNameFormat(reqParam.Name) {
		return false, errs.ManageGroupNameLenErr
	}
	// 不允许和系统管理组重名
	if ok, _ := slice.Contain(consts.DefaultManageGroupName, reqParam.Name); ok {
		return false, errs.DefaultManageGroupErr
	}

	count, dbErr := domain.UpdateManageGroup(orgId, operator, groupId, reqParam)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return count > 0, nil
}

// DeleteManageGroup 删除管理组
func DeleteManageGroup(orgId, operator, groupId int64) (bool, errs.SystemErrorInfo) {
	if orgId == 0 {
		return false, errs.OrgNotExist
	}

	manageGroup, dbErr := domain.GetManageGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.ManageGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 不允许删除系统管理组
	if consts.IsDefaultManageGroup(manageGroup.LangCode) {
		return false, errs.CannotRemoveDefaultManageGroupErr
	}

	count, dbErr := domain.DeleteManageGroup(orgId, operator, groupId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

// UpdateManageGroupContents 更新管理组内容信息
func UpdateManageGroupContents(orgId, operator int64, isOrgOwner bool, groupId int64, input req.UpdateManageGroupContents) (bool, errs.SystemErrorInfo) {
	if orgId == 0 {
		return false, errs.OrgNotExist
	}

	manageGroup, dbErr := domain.GetManageGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.ManageGroupNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	input.Values = slice.SliceUniqueInt64(input.Values)

	if input.Key == consts.TcAppIds {
		appIds := []int64{}
		copyErr := copyer.Copy(input.ValueIf, &appIds)
		if copyErr != nil {
			logger.Error(copyErr)
		}
		input.Values = appIds
	}

	// 系统管理组除了用户不允许修改其他选项
	if manageGroup.LangCode == consts.ManageGroupSys {
		if input.Key != consts.TcUserIds {
			input.Values = []int64{}
		} else {
			if len(input.Values) == 0 {
				return false, errs.SysAdminGroupMustHasMember
			}
			switch input.SourceFrom {
			case "polaris":
				// 更换超管时，校验是否已通过手机验证码的验证（根据请求携带的 token）。
				adminUserIds := make([]int64, 0)
				json.FromJson(manageGroup.UserIds, &adminUserIds)
				if !CheckIsPrivateDeploy() {
					_, err := CheckTokenForChangeSuperAdmin(&input, adminUserIds)
					if err != nil {
						return false, err
					}
				}
			default:
				// 非系统拥有者不可把自己从系统管理组移除
				ok, _ := slice.Contain(input.Values, operator)
				if !isOrgOwner && !ok {
					return false, errs.CannotDeleteSelf
				}
			}
		}
	} else {
		// 限制一名子管理员。极星方，允许设置多个管理组成员。
		if input.SourceFrom != "polaris" {
			if input.Key == consts.TcUserIds && len(input.Values) > 1 {
				return false, errs.ManageGroupMemberCountLimitErr
			}
		}
	}

	err := verifyContentsRef(orgId, groupId, input.Key, input.Values)
	if err != nil {
		return false, err
	}
	// 加入更新人
	count, dbErr := domain.UpdateManageGroupContents(orgId, operator, groupId, manageGroup.LangCode, input)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return count > 0, nil
}

// 检查是否是私有化部署
func CheckIsPrivateDeploy() bool {
	logger.InfoF("application config:%v", conf.Cfg.Application)
	runMode := conf.Cfg.Application.RunMode
	if runMode == 3 || runMode == 4 {
		return true
	}
	return false
}

func CheckTokenForChangeSuperAdmin(input *req.UpdateManageGroupContents, userIds []int64) (bool, errs.SystemErrorInfo) {
	if len(userIds) < 1 {
		return false, errs.NoAdminInAdminGroup
	}
	adminUserId := userIds[0]
	// 更换系统管理组的 user_ids 时，需要校验验证码是否正确
	// 产品：更换超管时，需要通过手机验证
	// 获取用户手机号
	adminUserInfo, err := domain.GetUserPoById(adminUserId)
	if err != nil {
		logger.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tmpToken, err := domain.GetSMSLoginCode(consts.AuthCodeTypeChangeSuperAdmin, consts.ContactAddressTypeMobile, adminUserInfo.Mobile)
	if err != nil {
		logger.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.CaptchaError, err)
	}
	if tmpToken != input.AuthToken {
		return false, errs.CaptchaError
	}

	return true, nil
}

// verifyContentsRef 查询内容是否有效
func verifyContentsRef(orgId int64, groupId int64, key string, values []int64) errs.SystemErrorInfo {
	switch key {
	case consts.TcUserIds:
		if len(values) > 0 {
			// 不可重复加入
			//groups, dbErr := domain.GetManageGroupListByUsers(orgId, values)
			//if dbErr != nil {
			//	logger.Error(dbErr)
			//	return errs.MysqlOperateError
			//}
			//if len(groups) > 1 || (len(groups) == 1 && groups[0].Id != groupId) {
			//	return errs.ManageGroupMemberConflict
			//}
			orgMemberList, dbErr := domain.GetEnableOrgMemberBaseInfoListByUsers(orgId, values)
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.MysqlOperateError
			}
			if len(values) != len(orgMemberList) {
				return errs.OrgMemberNotExistOrDisable
			}
		}
	case consts.TcAppPackageIds:
		if len(values) > 0 {
			packageList, err := domain.GetAppPackageList(orgId)
			if err != nil {
				logger.Error(err)
				return err
			}
			appPackageMap := make(map[int64]resp.AppPackageData)
			for _, pkg := range packageList {
				appPackageMap[pkg.Id] = pkg
			}
			for _, v := range values {
				if _, ok := appPackageMap[v]; !ok {
					return errs.AppPackageNotExist
				}
			}
		}
	case consts.TcAppIds:
		if len(values) > 0 {
			if values[0] != -1 {
				appList, err := domain.GetAppList(orgId)
				if err != nil {
					logger.Error(err)
					return err
				}
				appMap := make(map[int64]resp.AppData)
				for _, pkg := range appList {
					appMap[pkg.Id] = pkg
				}
				for _, v := range values {
					if _, ok := appMap[v]; !ok {
						return errs.AppNotExist
					}
				}
			}
		}
	case consts.TcDeptIds:
		if len(values) > 0 {
			departmentBos, dbErr := domain.GetDeptListByDeptIds(orgId, values)
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.MysqlOperateError
			}
			if len(departmentBos) != len(values) {
				return errs.DepartmentNotExist
			}
		}
	case consts.TcRoleIds:
		if len(values) > 0 {
			roles, dbErr := domain.GetRoleListByIds(orgId, values)
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.MysqlOperateError
			}
			if len(roles) != len(values) {
				return errs.RoleNotExist
			}
		}
	case consts.TcOptAuth, consts.TcUsageIds:
		{
		}
	default:
		return errs.ManageGroupOptionsNotExist
	}
	return nil
}

// GetManageGroupTree 获取管理组树
func GetManageGroupTree(orgId int64) (*resp.ManageGroupTreeResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	manageGroupList, dbErr := domain.GetManageGroupListByOrg(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	if len(manageGroupList) == 0 {
		return nil, nil
	}
	tree := resp.ManageGroupTreeResp{
		GeneralGroups: []resp.SimpleManageGroupInfo{},
	}
	for _, groupBo := range manageGroupList {
		optAuthArr := make([]string, 0)
		json.FromJson(groupBo.OptAuth, &optAuthArr)
		groupInfo := resp.SimpleManageGroupInfo{
			Id:       groupBo.Id,
			OrgId:    groupBo.OrgId,
			Name:     groupBo.Name,
			LangCode: groupBo.LangCode,
			OptAuth:  optAuthArr,
		}
		if groupInfo.LangCode == consts.ManageGroupSys {
			tree.SysGroup = &groupInfo
		} else {
			tree.GeneralGroups = append(tree.GeneralGroups, groupInfo)
		}
	}
	return &tree, nil
}

// GetManageGroupDetail 获取管理组详情
func GetManageGroupDetail(orgId, operatorUid int64, groupId int64) (*resp.ManageGroupDetailResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}

	manageGroup, dbErr := domain.GetManageGroup(orgId, groupId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.ManageGroupNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	isSysGroup := manageGroup.LangCode == consts.ManageGroupSys
	groupType := 2
	if isSysGroup {
		groupType = 1
	}
	detailResp := resp.ManageGroupDetailResp{
		AdminGroup: &resp.ManageGroupInfo{
			Id:         manageGroup.Id,
			OrgId:      manageGroup.OrgId,
			Name:       manageGroup.Name,
			Type:       groupType,
			LangCode:   manageGroup.LangCode,
			Creator:    manageGroup.Creator,
			CreateTime: manageGroup.CreateTime,
		},
		HeadData: &resp.ManageGroupHead{
			Users:       []resp.ManageGroupHeadDataUser{},
			AppPackages: []resp.ManageGroupHeadData{},
			Apps:        []resp.ManageGroupHeadData{},
			Depts:       []resp.ManageGroupHeadData{},
			Roles:       []resp.ManageGroupHeadData{},
		},
	}
	info := detailResp.AdminGroup
	_ = json.FromJson(manageGroup.UserIds, &info.UserIds)
	_ = json.FromJson(manageGroup.AppPackageIds, &info.AppPackageIds)
	_ = json.FromJson(manageGroup.AppIds, &info.AppIds)
	_ = json.FromJson(manageGroup.OptAuth, &info.OptAuth)
	// 为了兼容前端使用旧字段 usageIds。
	_ = json.FromJson(manageGroup.OptAuth, &info.UsageIds)
	_ = json.FromJson(manageGroup.DeptIds, &info.DeptIds)
	_ = json.FromJson(manageGroup.RoleIds, &info.RoleIds)

	// 人员信息
	if len(info.UserIds) > 0 {
		userInfos, dbErr := domain.GetEnableOrgMemberBaseInfoListByUsers(orgId, info.UserIds)
		if dbErr != nil {
			return nil, errs.MysqlOperateError
		}
		userMap := make(map[int64]bo.OrgMemberBaseInfoBo)
		for _, v := range userInfos {
			userMap[v.UserId] = v
		}
		userIds := make([]int64, 0)
		for _, id := range info.UserIds {
			if user, ok := userMap[id]; ok {
				phone := user.Mobile
				if !isSysGroup {
					// 只有管理员，才会返回完整的手机号。
					phone = strs.TransferPhone(user.Mobile, 1)
				}
				detailResp.HeadData.Users = append(detailResp.HeadData.Users, resp.ManageGroupHeadDataUser{
					Id:          user.UserId,
					Name:        user.Name,
					Avatar:      user.Avatar,
					PhoneNumber: phone,
				})
				userIds = append(userIds, id)
			}
		}
		info.UserIds = userIds
	}
	// 系统分组只获取人员
	if manageGroup.LangCode == consts.ManageGroupSys {
		return &detailResp, nil
	}

	// 查询应用包信息
	if len(info.AppPackageIds) > 0 {
		appPkgInfos, err := domain.GetAppPackageList(orgId)
		if err != nil {
			return nil, err
		}
		appPkgMap := make(map[int64]resp.AppPackageData)
		for _, v := range appPkgInfos {
			appPkgMap[v.Id] = v
		}
		appPackageIds := make([]int64, 0)
		for i := 0; i < len(info.AppPackageIds); i++ {
			id := info.AppPackageIds[i]
			if app, ok := appPkgMap[id]; ok {
				detailResp.HeadData.AppPackages = append(detailResp.HeadData.AppPackages, resp.ManageGroupHeadData{
					Id:   app.Id,
					Name: app.Name,
				})
				appPackageIds = append(appPackageIds, id)
			} else {
				if groupType == 2 {
					logger.InfoF("[删除不存在的pkg] ->  orgId:%d  pkgId:%d", orgId, id)
					_, dbErr = domain.RemoveContentsFromSubManageGroup(orgId, operatorUid, consts.TcAppPackageIds, id)
					if dbErr != nil {
						logger.ErrorF("[删除不存在的pkg] -> 失败: %s", dbErr)
					}
				}
			}
		}
		info.AppPackageIds = appPackageIds
	}

	// 查询应用信息
	if len(info.AppIds) > 0 {
		if info.AppIds[0] != -1 {
			appInfoList, err := domain.GetAppList(orgId)
			if err != nil {
				return nil, err
			}
			appMap := make(map[int64]resp.AppData)
			for _, v := range appInfoList {
				appMap[v.Id] = v
			}
			appIds := make([]int64, 0)
			for i := 0; i < len(info.AppIds); i++ {
				id := info.AppIds[i]
				if app, ok := appMap[id]; ok {
					detailResp.HeadData.Apps = append(detailResp.HeadData.Apps, resp.ManageGroupHeadData{
						Id:   app.Id,
						Name: app.Name,
					})
					appIds = append(appIds, id)
				} else {
					if groupType == 2 {
						logger.InfoF("[删除不存在的app] ->  orgId:%d  pkgId:%d", orgId, id)
						_, dbErr = domain.RemoveContentsFromSubManageGroup(orgId, operatorUid, consts.TcAppPackageIds, id)
						if dbErr != nil {
							logger.ErrorF("[删除不存在的app] -> 失败: %s", dbErr)
						}
					}
				}
			}
			info.AppIds = appIds
		} else {
			info.AppIds = []int64{-1}
		}
	}

	// 查询部门
	if len(info.DeptIds) > 0 {
		deptList, dbErr := domain.GetDeptListByDeptIds(orgId, info.DeptIds)
		if dbErr != nil {
			return nil, errs.MysqlOperateError
		}
		deptMap := make(map[int64]bo.OrgDeptBo)
		for _, v := range deptList {
			deptMap[v.Id] = v
		}
		deptIds := make([]int64, 0)
		for i := 0; i < len(info.DeptIds); i++ {
			id := info.DeptIds[i]
			if dept, ok := deptMap[id]; ok {
				detailResp.HeadData.Depts = append(detailResp.HeadData.Depts, resp.ManageGroupHeadData{
					Id:   dept.Id,
					Name: dept.Name,
				})
				deptIds = append(deptIds, id)
			}
		}
		info.DeptIds = deptIds
	}

	// 查询角色是否有效
	if len(info.RoleIds) > 0 {
		roles, dbErr := domain.GetRoleListByIds(orgId, info.RoleIds)
		if dbErr != nil {
			return nil, errs.MysqlOperateError
		}
		roleMap := make(map[int64]po.PpmRolRole)
		for _, v := range roles {
			roleMap[v.Id] = v
		}
		roleIds := make([]int64, 0)
		for i := 0; i < len(info.RoleIds); i++ {
			id := info.RoleIds[i]
			if role, ok := roleMap[id]; ok {
				detailResp.HeadData.Roles = append(detailResp.HeadData.Roles, resp.ManageGroupHeadData{
					Id:   role.Id,
					Name: role.Name,
				})
				roleIds = append(roleIds, id)
			}
		}
		info.RoleIds = roleIds
	}

	return &detailResp, nil
}

// GetManageGroupOperationConfig 获取管理组权限项
func GetManageGroupOperationConfig(orgId, operatorUid int64) (*resp.GetOperationConfigResp, errs.SystemErrorInfo) {
	optAuthArr := make([]resp.GetOperationConfigRespItem, 0)
	errJson := json.FromJson(consts.OptAuthList, &optAuthArr)
	if errJson != nil {
		logger.ErrorF("[GetManageGroupOperationConfig] json err:%v", errJson)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
	}
	extraInfo := make(map[string]interface{}, 0)
	errJson = json.FromJson(consts.OptAuthExtraInfo, &extraInfo)
	if errJson != nil {
		logger.ErrorF("[GetManageGroupOperationConfig] json err:%v", errJson)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
	}

	return &resp.GetOperationConfigResp{
		OptAuthList:      optAuthArr,
		OptAuthExtraInfo: extraInfo,
	}, nil
}
