package inner_service

import (
	"fmt"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/str"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// AddAppPkgToManageGroup
func AddAppPkgToManageGroup(orgId int64, operatorUid int64, pkgId int64) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[新增应用包至人员所在管理组] -> orgId:%d, operatorUid:%d, pkgId:%d", orgId, operatorUid, pkgId)
	manageGroup, dbErr := domain.GetManageGroupListByUser(orgId, operatorUid)
	if dbErr != nil {
		// 不在管理组则忽视
		if dbErr == db.ErrNoMoreRows {
			return true, nil
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	// 在系统管理组中则无需添加
	if manageGroup.LangCode == consts.ManageGroupSys {
		return true, nil
	}

	_, dbErr = domain.AppendContentsToManageSubGroup(orgId, operatorUid, manageGroup.Id, consts.TcAppPackageIds, pkgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return true, nil

}

// DeleteAppPkgFromManageGroup
func DeleteAppPkgFromManageGroup(orgId int64, operatorUid int64, pkgId int64) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[从管理组删除应用包] -> orgId:%d, operatorUid:%d, pkgId:%d", orgId, operatorUid, pkgId)
	_, dbErr := domain.RemoveContentsFromSubManageGroup(orgId, operatorUid, consts.TcAppPackageIds, pkgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return true, nil

}

// AddAppToManageGroup
func AddAppToManageGroup(orgId int64, operatorUid int64, appId int64) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[新增应用至人员所在管理组] -> orgId:%d, operatorUid:%d, appId:%d", orgId, operatorUid, appId)
	manageGroup, dbErr := domain.GetManageGroupListByUser(orgId, operatorUid)
	if dbErr != nil {
		// 不在管理组则忽视
		if dbErr == db.ErrNoMoreRows {
			return true, nil
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	// 在系统管理组中则无需添加
	if manageGroup.LangCode == consts.ManageGroupSys {
		return true, nil
	}

	_, dbErr = domain.AppendContentsToManageSubGroup(orgId, operatorUid, manageGroup.Id, consts.TcAppIds, appId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return true, nil

}

// DeleteAppFromManageGroup
func DeleteAppFromManageGroup(orgId int64, operatorUid int64, appId int64) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[从管理组删除应用] -> orgId:%d, operatorUid:%d, appId:%d", orgId, operatorUid, appId)
	_, dbErr := domain.RemoveContentsFromSubManageGroup(orgId, operatorUid, consts.TcAppIds, appId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return true, nil

}

func GetManager(orgId int64) (*inner_resp.GetManagerResp, errs.SystemErrorInfo) {
	logger.InfoF("[查询组织管理员管理的应用] -> orgId:%d", orgId)
	list, err := domain.GetManageGroupListByOrg(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	userManagerMap := map[int64]string{}
	userAppIdsMap := map[int64][]int64{}

	roleManagerMap := map[int64]string{}
	roleAppIdsMap := map[int64][]int64{}

	deptManagerMap := map[int64]string{}
	deptAppIdsMap := map[int64][]int64{}

	for _, group := range list {
		var userIds []int64
		json.FromJsonIgnoreError(group.UserIds, &userIds)

		var roleIds []int64
		json.FromJsonIgnoreError(group.RoleIds, &roleIds)

		var deptIds []int64
		json.FromJsonIgnoreError(group.DeptIds, &deptIds)

		var appIds []int64
		json.FromJsonIgnoreError(group.AppIds, &appIds)

		for _, id := range userIds {
			langCode := userManagerMap[id]
			if group.LangCode == consts.ManageGroupSys {
				userManagerMap[id] = group.LangCode
			}
			if langCode != consts.ManageGroupSys && langCode != consts.ManageGroupSub && langCode != consts.ManageGroupSubNormalAdmin {
				userManagerMap[id] = group.LangCode
			}
			userAppIdsMap[id] = append(userAppIdsMap[id], appIds...)
		}

		for _, id := range roleIds {
			langCode := roleManagerMap[id]
			if group.LangCode == consts.ManageGroupSys {
				roleManagerMap[id] = group.LangCode
			}
			if langCode != consts.ManageGroupSys && langCode != consts.ManageGroupSub && langCode != consts.ManageGroupSubNormalAdmin {
				roleManagerMap[id] = group.LangCode
			}
			roleAppIdsMap[id] = append(roleAppIdsMap[id], appIds...)
		}

		for _, id := range deptIds {
			langCode := deptManagerMap[id]
			if group.LangCode == consts.ManageGroupSys {
				deptManagerMap[id] = group.LangCode
			}
			if langCode != consts.ManageGroupSys && langCode != consts.ManageGroupSub && langCode != consts.ManageGroupSubNormalAdmin {
				deptManagerMap[id] = group.LangCode
			}
			deptAppIdsMap[id] = append(deptAppIdsMap[id], appIds...)
		}
	}

	res := make([]inner_resp.GetManagerData, 0)
	for id, s := range userManagerMap {
		temp := inner_resp.GetManagerData{
			MemberType: consts.MemberTypeUser,
			MemberId:   id,
			LangCode:   s,
			AppIds:     []int64{},
		}
		if appIds, ok := userAppIdsMap[id]; ok {
			temp.AppIds = appIds
		}
		res = append(res, temp)
	}

	for id, s := range deptManagerMap {
		temp := inner_resp.GetManagerData{
			MemberType: "D_",
			MemberId:   id,
			LangCode:   s,
			AppIds:     []int64{},
		}
		if appIds, ok := deptAppIdsMap[id]; ok {
			temp.AppIds = appIds
		}
		res = append(res, temp)
	}

	for id, s := range roleManagerMap {
		temp := inner_resp.GetManagerData{
			MemberType: "R_",
			MemberId:   id,
			LangCode:   s,
			AppIds:     []int64{},
		}
		if appIds, ok := roleAppIdsMap[id]; ok {
			temp.AppIds = appIds
		}
		res = append(res, temp)
	}
	for i, re := range res {
		if re.LangCode == consts.ManageGroupSys {
			res[i].IsSysAdmin = true
		}
		if re.LangCode == consts.ManageGroupSub || re.LangCode == consts.ManageGroupSubNormalAdmin {
			res[i].IsSubAdmin = true
		}
	}

	return &inner_resp.GetManagerResp{Data: res}, nil
}

func ManageGroupInit(orgId, creatorId int64, sourceFrom string) (int64, errs.SystemErrorInfo) {
	var id int64
	var err error
	switch sourceFrom {
	case "polaris": // 极星的初始化特殊一些，需要创建三个管理组：超管、组织管理员、组织成员
		id, err = domain.ManageGroupInitForPolaris(orgId, creatorId)
	default:
		id, err = domain.ManageGroupInitDefault(orgId, creatorId)
	}
	if err != nil {
		logger.Error(err)
		return id, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return id, nil
}

//增加人员到系统管理组
func AddUserToSysManageGroup(orgId int64, operator int64, addUserIds []int64) errs.SystemErrorInfo {
	if len(addUserIds) == 0 {
		return nil
	}
	group, err := domain.GetSysManageGroup(orgId)
	if err != nil {
		if err == errs.ManageGroupNotExist {
			//没有初始化就初始化系统管理组
			_, initErr := ManageGroupInit(orgId, operator, "")
			if initErr != nil {
				logger.Error(initErr)
				return initErr
			}
			//再次获取(再获取不到报错好了)
			group, err = domain.GetSysManageGroup(orgId)
			if err != nil {
				logger.Error(err)
				return errs.MysqlOperateError
			}
		} else {
			logger.Error(err)
			return errs.MysqlOperateError
		}
	}

	var userIds []int64
	jsonErr := json.FromJson(group.UserIds, &userIds)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return errs.JSONConvertError
	}
	newUserIds := append(userIds, addUserIds...)
	//去重
	newUserIds = slice.SliceUniqueInt64(newUserIds)

	// 不可重复加入
	//groups, dbErr := domain.GetManageGroupListByUsers(orgId, newUserIds)
	//if dbErr != nil {
	//	logger.Error(dbErr)
	//	return errs.MysqlOperateError
	//}
	//if len(groups) > 1 || (len(groups) == 1 && groups[0].Id != group.Id) {
	//	return errs.ManageGroupMemberConflict
	//}
	//orgMemberList, dbErr := domain.GetEnableOrgMemberBaseInfoListByUsers(orgId, newUserIds)
	//if dbErr != nil {
	//	logger.Error(dbErr)
	//	return errs.MysqlOperateError
	//}
	//if len(newUserIds) != len(orgMemberList) {
	//	return errs.OrgMemberNotExistOrDisable
	//}
	// 加入更新人
	_, dbErr := domain.UpdateManageGroupContents(orgId, operator, group.Id, group.LangCode, req.UpdateManageGroupContents{
		Values: newUserIds,
		Key:    consts.TcUserIds,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// DeleteOneUserFromOrg 将一个用户从该组织的**所有管理组**中移除
func DeleteOneUserFromOrg(orgId, operateUid, targetUserId int64) errs.SystemErrorInfo {
	// 1.先查询存在该用户的所有管理组
	groups, dbErr := domain.GetManageGroupListByUsers(orgId, []int64{targetUserId})
	if dbErr != nil {
		logger.ErrorF("[DeleteOneUserFromOrg] GetManageGroupListByUsers err: %v, orgId: %d", dbErr, orgId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	// 2.逐个进行移除
	if len(groups) > 0 {
		// 普通查询放在事务外
		effectiveUids := make([]int64, 0)
		effectiveUids, dbErr = domain.GetOneNormalAdminForUpgrade(orgId)
		if dbErr != nil || len(effectiveUids) < 1 {
			// 尝试将一个普通成员设置为超管
			effectiveUid, dbErr := domain.GetOneNormalUserForSuperAdmin(orgId)
			logger.InfoF("[DeleteOneUserFromOrg] Set normalAdmin into superAdmin but normalAdmin notExist. orgId: %d, userId:%v, targetId:%v err: %v", orgId, effectiveUid, targetUserId, dbErr)
			if dbErr != nil {
				logger.ErrorF("[DeleteOneUserFromOrg] GetOneNormalUserForSuperAdmin err: %v", dbErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
			} else {
				effectiveUids = append(effectiveUids, effectiveUid)
			}
		}
		dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
			for _, group := range groups {
				tmpUserIdArr := domain.TransferUserIdsFromUserIdJson(group.UserIds)
				searchedIndex := domain.GetSearchedIndexArr(tmpUserIdArr, targetUserId)
				if searchedIndex != -1 {
					effectiveUid := int64(0)
					// 对于超管，如果仅有的超管离职了，则将管理员设为超管
					if group.LangCode == consts.ManageGroupSys {
						logger.InfoF("[DeleteOneUserFromOrg] get effectiveUid. tmpUserIdArr: %d, userId:%v, err: %v", tmpUserIdArr, effectiveUids, dbErr)
						if len(tmpUserIdArr) == 1 && tmpUserIdArr[0] == targetUserId {
							// 重新设置超管
							for _, oneUserId := range effectiveUids {
								if oneUserId != tmpUserIdArr[0] {
									effectiveUid = oneUserId
								}
							}
						}
					}
					// 删除旧管理组中特定的用户
					// eg: update lc_per_manage_group set user_ids = JSON_REMOVE(`user_ids`, '$[4]') WHERE org_id=1449 and id=1396673623212785665
					_, dbErr := domain.RemoveUserFromAdminGroup(group.Id, operateUid, searchedIndex, tx)
					if dbErr != nil {
						logger.ErrorF("[DeleteOneUserFromOrg] RemoveUserFromAdminGroup err: %v", dbErr)
						return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
					}
					// 如果需要将某个用户设为超管
					if effectiveUid > 0 {
						_, dbErr = domain.AppendUserIntoAdminGroup(group.Id, operateUid, effectiveUid, tx)
						if dbErr != nil {
							logger.ErrorF("[DeleteOneUserFromOrg] AppendUserIntoAdminGroup err: %v", dbErr)
							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
						}
					}
				}
			}
			return nil
		})
		if dbErr != nil {
			logger.ErrorF("[DeleteOneUserFromOrg] transX err: %v", dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
	}

	return nil
}

// ReplaceSuperAdmin 更换组织超管
// 由于产品设定的是只有一个超管，因此直接将目标用户作为唯一的超管，无需删除旧超管
func ReplaceSuperAdmin(orgId, operateUid, targetUserId int64) errs.SystemErrorInfo {
	logger.InfoF("[更换组织超管] -> orgId:%d, operatorUid:%d, targetUserId:%d", orgId, operateUid, targetUserId)
	sysGroup, dbErr := domain.GetSysManageGroup(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	if targetUserId < 1 {
		return errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "更新后的用户 id 不合法。")
	}

	newUserIdArr := []int64{targetUserId}
	_, dbErr = domain.UpdateManageGroupContents(orgId, operateUid, sysGroup.Id, sysGroup.LangCode, req.UpdateManageGroupContents{
		Values: newUserIdArr,
		Key:    consts.TcUserIds,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// GetSuperAdminIds 获取组织超管 ids
func GetSuperAdminIds(orgId int64) ([]int64, errs.SystemErrorInfo) {
	logger.InfoF("[获取组织超管id] -> orgId:%d", orgId)
	sysGroup, dbErr := domain.GetSysManageGroup(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	userIdsJson := sysGroup.UserIds
	if userIdsJson == "" {
		userIdsJson = "[]"
	}
	uidArr := domain.TransferUserIdsFromUserIdJson(userIdsJson)

	return uidArr, nil
}

// AddNewMenuToRole 向角色中追加权限项。目前是跑脚本用到。
func AddNewMenuToRole(input *inner_req.AddNewMenuToRoleReq) errs.SystemErrorInfo {
	logger.InfoF("[AddNewMenuToRole] input: %s", json.ToJsonIgnoreError(input))
	var oriErr error
	if len(input.OrgIds) < 1 {
		err := errs.ParamError
		logger.ErrorF("[AddNewMenuToRole] err: %v", err)
		return err
	}
	// 查询组织的非超管角色：普通管理员、团队成员、用户自定义角色
	cond := db.Cond{
		consts.TcOrgId: db.In(input.OrgIds),
		// 历史原因，这几个（`consts.ManageGroupSubNormalUser`, `consts.ManageGroupSubUserCustom`, `consts.ManageGroupSub`）是“团队成员”角色标识
		consts.TcLangCode: db.In([]string{consts.ManageGroupSubNormalAdmin, consts.ManageGroupSubNormalUser, consts.ManageGroupSubUserCustom, consts.ManageGroupSub}),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	roles := make([]po.LcPerManageGroup, 0)
	dbErr := store.Mysql.SelectAllByCondWithNumAndOrder(consts.TableManageGroup, cond, nil, 0, -1, "create_time asc", &roles)
	if dbErr != nil {
		logger.ErrorF("[AddNewMenuToRole] err: %v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	// 过滤掉已经存在回收站的角色
	roleIds := domain.FilterExistTrashMenuRoles(roles)
	if len(roleIds) < 1 {
		return nil
	}

	// update 语句不能带有 limit，所以换一种方式：Exec。
	conn, oriErr := store.Mysql.GetConnect()
	if oriErr != nil {
		logger.ErrorF("[AddNewMenuToRole] store.Mysql.GetConnect err: %v", oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	appendSqlStr := domain.GetAppendPrmStrForSql([]string{consts.MenuPermissionOrgTrash})
	_, oriErr = conn.Exec("UPDATE " + consts.TableManageGroup +
		" SET " + "opt_auth=" + fmt.Sprintf("JSON_ARRAY_APPEND(`opt_auth`, %s)", appendSqlStr) +
		fmt.Sprintf(" where id in (%s)", str.Int64Implode(roleIds, ",")) +
		fmt.Sprintf(" limit %d", len(roleIds)),
	)
	if oriErr != nil {
		logger.ErrorF("[AddNewMenuToRole] conn.Exec err: %v", oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	return nil
}
