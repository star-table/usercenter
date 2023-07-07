package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"github.com/star-table/usercenter/service/service/remote"
	"github.com/tealeg/xlsx"
	"upper.io/db.v3"

	"upper.io/db.v3/lib/sqlbuilder"
)

// CreateOrgMember 创建组织成员
func CreateOrgMember(orgId, operatorUid int64, reqParam req.CreateOrgMemberReq) (int64, errs.SystemErrorInfo) {

	// 检测账号是否合法
	if reqParam.LoginName != "" {
		reqParam.LoginName = strings.TrimSpace(reqParam.LoginName)
		if !format.VerifyAccountFormat(reqParam.LoginName) {
			return 0, errs.UsernameLenError
		}
	}

	// 检测手机号
	if reqParam.PhoneNumber != "" {
		reqParam.PhoneNumber = strings.TrimSpace(reqParam.PhoneNumber)
		reqParam.PhoneRegion = strings.TrimSpace(reqParam.PhoneRegion)
		if !format.VerifyMobileFormat(reqParam.PhoneNumber) || !format.VerifyMobileRegionFormat(reqParam.PhoneRegion) {
			return 0, errs.PhoneNumberFormatError
		}
	}

	// 检测邮箱是否合法
	if reqParam.Email != "" {
		reqParam.Email = strings.TrimSpace(reqParam.Email)
		if !format.VerifyEmailFormat(reqParam.Email) {
			return 0, errs.EmailFormatErr
		}
	}

	// 检测昵称是否合法
	if reqParam.Name == "" {
		reqParam.Name = reqParam.LoginName
	}
	if reqParam.Name == "" {
		reqParam.Name = reqParam.PhoneNumber
	}
	if reqParam.Name == "" {
		reqParam.Name = reqParam.Email
	}
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	if !format.VerifyNicknameFormat(reqParam.Name) {
		return 0, errs.NicknameLenError
	}

	// 检测密码
	if reqParam.Password != "" {
		reqParam.Password = strings.TrimSpace(reqParam.Password)
		if !format.VerifyAccountPwdFormat(reqParam.Password) {
			return 0, errs.PwdLengthError
		}
	}

	// 查询部门是否有效
	if len(reqParam.DeptAndPositions) > 0 {
		deptIds := make([]int64, 0)
		for _, v := range reqParam.DeptAndPositions {
			deptIds = append(deptIds, v.DepartmentId)
		}
		deptList, dbErr := domain.GetDeptListByDeptIds(orgId, deptIds)
		if dbErr != nil {
			return 0, errs.MysqlOperateError
		}
		if len(deptList) != len(deptIds) {
			return 0, errs.DepartmentNotExist
		}
	}
	// 查询角色是否有效
	if len(reqParam.RoleIds) > 0 {
		reqParam.RoleIds = slice.SliceUniqueInt64(reqParam.RoleIds)
		roles, dbErr := domain.GetRoleListByIds(orgId, reqParam.RoleIds)
		if dbErr != nil {
			return 0, errs.MysqlOperateError
		}

		if len(roles) != len(reqParam.RoleIds) {
			return 0, errs.RoleNotExist
		}
	}

	// 添加
	userId, err := domain.CreateOrgUser(orgId, operatorUid, reqParam)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	// 启用异步线程去同步给网盘
	if orgId == 202505060582035456 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.ErrorF("[同步新增成员给网盘]失败 ->：%s", r)
				}
			}()
			member, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, userId)
			if dbErr != nil {
				logger.ErrorF("[同步新增成员给网盘]失败 -> 成员信息查询失败：%s", dbErr)
				return
			}
			// 执行同步操作
			_ = remote.SynAddUserToCloudreve(*member)
		}()
	}

	return userId, nil
}

// UpdateOrgMemberInfo 修改组织成员信息
func UpdateOrgMemberInfo(orgId, operatorUid int64, reqParam req.UpdateOrgMemberReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	// 获取组织成员信息，非启用状态也可以获取到
	memberBaseInfo, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	manageGroup, dbErr := domain.GetManageGroupListByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
	}
	isSysManager := false
	if manageGroup != nil && manageGroup.LangCode == consts.ManageGroupSys {
		isSysManager = true
	}

	upd := mysql.Upd{}
	orgUserUpd := mysql.Upd{}
	// 状态修改 判断权限
	if reqParam.Status != nil && *reqParam.Status != memberBaseInfo.Status {
		if reqParam.UserId == operatorUid {
			return false, errs.CannotChangeSelfStatus
		}
		if reqParam.UserId == memberBaseInfo.OrgOwner {
			return false, errs.CannotChangeOwnerStatus
		}
		// 子管理员不可以修改系统管理员
		if !perContext.HasAllPermission() && isSysManager {
			return false, errs.CannotChangeAdminStatus
		}
		orgUserUpd[consts.TcStatus] = *reqParam.Status
	}

	if (reqParam.PhoneNumber != nil && *reqParam.PhoneNumber != memberBaseInfo.Mobile) ||
		(reqParam.Email != nil && *reqParam.Email != memberBaseInfo.Email) ||
		(reqParam.Name != nil && *reqParam.Name != memberBaseInfo.Name) ||
		(reqParam.EmpNo != nil && *reqParam.EmpNo != memberBaseInfo.EmpNo) ||
		strings.Join(reqParam.WeiboIds, ",") != memberBaseInfo.WeiboIds ||
		(reqParam.Avatar != nil && *reqParam.Avatar != memberBaseInfo.Avatar) {
		// 非拥有者不可修改组织拥有者信息
		if !perContext.IsOrgOwner() && reqParam.UserId == memberBaseInfo.OrgOwner {
			return false, errs.CannotEditOrgOwnerMainField
		}
		// 子管理员不可以修改系统管理员
		if !perContext.HasAllPermission() && isSysManager {
			return false, errs.CannotEditAdminMainField
		}
	}

	if reqParam.Name != nil && *reqParam.Name != memberBaseInfo.Name {
		// 昵称
		*reqParam.Name = strings.TrimSpace(*reqParam.Name)
		isNameRight := format.VerifyNicknameFormat(*reqParam.Name)
		if !isNameRight {
			return false, errs.NicknameLenError
		}
		upd[consts.TcName] = *reqParam.Name
	}

	// 邮箱
	if reqParam.Email != nil && *reqParam.Email != memberBaseInfo.Email {
		*reqParam.Email = strings.TrimSpace(*reqParam.Email)
		if !format.VerifyEmailFormat(*reqParam.Email) {
			return false, errs.EmailFormatErr
		}
		exist, dbErr := store.Mysql.IsExistByCond(consts.TableUser, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcId:          db.NotEq(reqParam.UserId),
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcEmail:       *reqParam.Email,
			consts.TcEmail + " ": db.NotEq(""), // 排除空
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}

		if exist {
			return false, errs.EmailAlreadyBindByOtherAccountError
		}
		upd[consts.TcEmail] = reqParam.Email
	}

	// 手机号
	if (reqParam.PhoneNumber != nil && *reqParam.PhoneNumber != memberBaseInfo.Mobile) ||
		(reqParam.PhoneRegion != nil && *reqParam.PhoneRegion != memberBaseInfo.MobileRegion) {
		*reqParam.PhoneNumber = strings.TrimSpace(*reqParam.PhoneNumber)
		*reqParam.PhoneRegion = strings.TrimSpace(*reqParam.PhoneRegion)
		if !format.VerifyMobileFormat(*reqParam.PhoneNumber) || !format.VerifyMobileRegionFormat(*reqParam.PhoneRegion) {
			return false, errs.PhoneNumberFormatError
		}
		exist, dbErr := store.Mysql.IsExistByCond(consts.TableUser, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcId:           db.NotEq(reqParam.UserId),
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcMobile:       *reqParam.PhoneNumber,
			consts.TcMobile + " ": db.NotEq(""), // 排除空
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}

		if exist {
			return false, errs.MobileAlreadyBindOtherAccountError
		}
		upd[consts.TcMobile] = reqParam.PhoneNumber
		upd[consts.TcMobileRegion] = reqParam.PhoneRegion
	}

	// 头像
	if reqParam.Avatar != nil && *reqParam.Avatar != memberBaseInfo.Avatar {
		upd[consts.TcAvatar] = reqParam.Avatar
	}

	// 工号
	if reqParam.EmpNo != nil && *reqParam.EmpNo != memberBaseInfo.EmpNo {
		*reqParam.EmpNo = strings.TrimSpace(*reqParam.EmpNo)
		exist, dbErr := store.Mysql.IsExistByCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcId:          db.NotEq(reqParam.UserId),
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcEmpNo:       *reqParam.EmpNo,
			consts.TcEmpNo + " ": db.NotEq(""), // 排除空
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
		if exist {
			return false, errs.OrgEmpNoConflictErr
		}
		orgUserUpd[consts.TcEmpNo] = *reqParam.EmpNo
	}
	// 更新管理组
	if reqParam.AdminGroup != nil && len(reqParam.AdminGroup) > 0 {
		// 校验是否可以修改用户管理组 ModifyUserAdminGroup
		if !perContext.HasOpForPolaris(consts.OperationOrgUserModifyUserAdminGroup) {
			return false, errs.NoOperationPermissions
		}
		_, err := ChangeUserAdminGroupBatch(orgId, operatorUid, req.ChangeUserAdminGroupReq{
			UserId:           reqParam.UserId,
			DstAdminGroupIds: reqParam.AdminGroup,
		}, perContext)
		if err != nil {
			logger.Error(err)
			return false, err
		}
	}

	// 微博
	if reqParam.WeiboIds != nil && strings.Join(reqParam.WeiboIds, ",") != memberBaseInfo.WeiboIds {
		orgUserUpd[consts.TcWeiboIds] = strings.Join(reqParam.WeiboIds, ",")
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(upd) > 0 {
			upd[consts.TcUpdator] = operatorUid
			_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
				consts.TcId: reqParam.UserId,
			}, upd)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		// 修改组织成员表
		if len(orgUserUpd) > 0 {
			orgUserUpd[consts.TcUpdator] = operatorUid
			orgUserUpd[consts.TcStatusChangeTime] = time.Now()
			_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
				consts.TcUserId: reqParam.UserId,
			}, orgUserUpd)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		// 状态修改
		if reqParam.Status != nil && *reqParam.Status != memberBaseInfo.Status {
			if ok, _ := slice.Contain([]int{consts.AppStatusEnable, consts.AppStatusDisabled, consts.OrgUserStatusResigned}, reqParam.Status); ok {
				_, dbErr := domain.ChangeOrgMemberStatus(orgId, operatorUid, []int64{reqParam.UserId}, *reqParam.Status, tx)
				if dbErr != nil {
					logger.Error(dbErr)
					return dbErr
				}
			}
		}

		// 启用异步线程去同步给网盘
		if orgId == 202505060582035456 &&
			(reqParam.Name != nil && reqParam.Status != nil) &&
			(*reqParam.Name != memberBaseInfo.Name || *reqParam.Status != memberBaseInfo.Status) {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.ErrorF("[同步修改成员给网盘]失败 ->：%s", r)
					}
				}()
				member, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
				if dbErr != nil {
					logger.ErrorF("[同步修改成员给网盘]失败 -> 成员信息查询失败：%s", dbErr)
					return
				}
				// 执行同步操作
				_ = remote.SynEditUserToCloudreve([]bo.OrgMemberBaseInfoBo{*member})
			}()
		}
		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 更新用户的部门信息
	if reqParam.DeptAndPositions != nil {
		err := _changeUserDeptAndPosition(orgId, operatorUid, reqParam.UserId, reqParam.DeptAndPositions, perContext)
		if err != nil {
			return false, err
		}
	}
	if reqParam.DepartmentIds != nil && len(reqParam.DepartmentIds) > 0 {

	}
	// 更改用户的角色信息
	if reqParam.RoleIds != nil {
		err := _bindUserRoles(orgId, operatorUid, reqParam.UserId, reqParam.RoleIds, perContext)
		if err != nil {
			return false, err
		}
	}
	// 清除user缓存
	cacheErr := domain.ClearBaseUserInfo(orgId, reqParam.UserId)
	if cacheErr != nil {
		logger.Error(cacheErr)
	}

	return true, nil
}

// ChangeOrgMemberStatus 更新成员状态
func ChangeOrgMemberStatus(orgId, operatorUid int64, reqParam req.UpdateOrgMemberStatusReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	// 验证修改权限
	err := checkChangeStatusPermission(orgId, operatorUid, reqParam.MemberIds, perContext)
	if err != nil {
		return false, err
	}
	count := 0
	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		count1, dbErr := domain.ChangeOrgMemberStatus(orgId, operatorUid, reqParam.MemberIds, reqParam.Status, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		count = count1
		return nil
	})
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	//最后将用户信息缓存清掉
	cacheErr := domain.ClearBaseUserInfoBatch(orgId, reqParam.MemberIds)
	if cacheErr != nil {
		logger.Error(cacheErr)
	}

	// 启用异步线程去同步给网盘
	if orgId == 202505060582035456 && count > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.ErrorF("[同步修改成员给网盘]失败 ->：%s", r)
				}
			}()
			members, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, reqParam.MemberIds)
			if dbErr != nil {
				logger.ErrorF("[同步修改成员给网盘]失败 -> 成员信息查询失败：%s", dbErr)
				return
			}
			// 执行同步操作
			_ = remote.SynEditUserToCloudreve(members)
		}()
	}
	return count > 0, nil
}

// RemoveOrgMember 移除组织成员
func RemoveOrgMember(orgId, operatorUid int64, reqParam req.RemoveOrgMemberReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	// 验证修改权限
	err := checkChangeStatusPermission(orgId, operatorUid, reqParam.MemberIds, perContext)
	if err != nil {
		return false, errs.CannotRemoveUser
	}
	//如果有权限，移除成员
	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		count, dbErr := domain.RemoveOrgMember(orgId, operatorUid, reqParam.MemberIds, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		// 如果更新数量与预期不符，认为动作失败
		if count != len(reqParam.MemberIds) {
			return errs.UpdateMemberStatusFail
		}
		// 将用户从组织移除之后 - 将该用户从部门移除
		dbErr = domain.UnBoundDepartmentUser(orgId, reqParam.MemberIds, operatorUid, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		// 删除角色关系
		dbErr = domain.TransUnbindUserHaveRoles(orgId, operatorUid, reqParam.MemberIds, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		// 删除管理组中绑定的用户关系
		for _, uid := range reqParam.MemberIds {
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableManageGroup, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIsDelete: consts.AppIsNoDelete,
				db.Raw("json_search(`user_ids`, 'one', ?)", uid): db.IsNotNull(),
			}, mysql.Upd{
				consts.TcUserIds: db.Raw("json_remove(`user_ids`,JSON_UNQUOTE(json_search(`user_ids`, 'one', ?)))", uid),
				consts.TcUpdator: operatorUid,
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		return nil
	})
	if dbErr != nil {
		if dbErr == errs.UpdateMemberStatusFail {
			return false, errs.UpdateMemberStatusFail
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	//最后将用户信息缓存清掉
	cacheErr := domain.ClearBaseUserInfoBatch(orgId, reqParam.MemberIds)
	if cacheErr != nil {
		logger.Error(cacheErr)
	}

	// 启用异步线程去同步给网盘
	if orgId == 202505060582035456 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.ErrorF("[同步删除成员给网盘]失败 ->：%s", r)
				}
			}()
			// 执行同步操作
			_ = remote.SynRemoveUserToCloudreve(reqParam.MemberIds)
		}()
	}
	return true, nil
}

func checkChangeStatusPermission(orgId, operatorUid int64, userIds []int64, perContext *inner_resp.OrgUserPerContext) errs.SystemErrorInfo {
	// 获取组织成员信息，非启用状态也可以获取到
	memberBaseInfoList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	manageGroups, dbErr := domain.GetManageGroupListByUsers(orgId, userIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	isSysManager := false
	for _, group := range manageGroups {
		if group.LangCode == consts.ManageGroupSys {
			isSysManager = true
			break
		}
	}
	for _, member := range memberBaseInfoList {
		if member.UserId == operatorUid {
			return errs.CannotChangeSelfStatus
		}
		if member.UserId == member.OrgOwner {
			return errs.CannotChangeOwnerStatus
		}
		// 子管理员不可以修改系统管理员
		if !perContext.HasAllPermission() && isSysManager {
			return errs.CannotChangeAdminStatus
		}
	}
	return nil
}

// GetOrgMemberList 组织成员列表
func GetOrgMemberList(orgId int64, req req.UserListReq, perContext *inner_resp.OrgUserPerContext) (*resp.UserListResp, errs.SystemErrorInfo) {
	// 查询的列表中，超管放在第一位。所以先查询该组织的超管
	groups, dbErr := domain.GetAdminGroupsByLangCode(orgId, []string{consts.ManageGroupSys, consts.ManageGroupSubNormalAdmin})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	// 区分出超管、普通管理员角色
	adminUids := make([]int64, 0)
	normalAdminUids := make([]int64, 0)
	for _, item := range groups {
		tmpUids := domain.TransferUserIdsFromUserIdJson(item.UserIds)
		if item.LangCode == consts.ManageGroupSys {
			adminUids = append(adminUids, tmpUids...)
		} else {
			normalAdminUids = append(normalAdminUids, tmpUids...)
		}
	}

	union := &db.Union{}
	if req.SearchCode != nil && *req.SearchCode != "" {
		union = union.Or(db.Cond{
			//用户
			consts.TcUserId: db.In(db.Raw("select id from ppm_org_user where name like ? "+
				" union "+
				" select id from ppm_org_user where lower(`name_pinyin`) like ? "+
				" union "+
				" select id from ppm_org_user where lower(`login_name`) like ?"+
				" union "+
				" select id from ppm_org_user where email like ?",
				"%"+*req.SearchCode+"%", "%"+strings.ToLower(*req.SearchCode)+"%", "%"+strings.ToLower(*req.SearchCode)+"%", "%"+*req.SearchCode+"%",
			)),
		})
	}

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if req.IsAllocate != nil {
		raw := db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2", orgId)
		if *req.IsAllocate == 2 {
			cond[consts.TcUserId] = db.NotIn(raw)
		} else {
			cond[consts.TcUserId] = db.In(raw)
		}
	}

	if req.Status != nil {
		cond[consts.TcStatus] = *req.Status
		// 产品胡子龙：禁用成员中不包括审核不通过的成员
		//if *req.Status == consts.AppStatusDisabled {
		//	raw := db.Raw("select user_id from ppm_org_user_organization where org_id = ? and is_delete = 2 and check_status = ?", orgId, consts.AppCheckStatusFail)
		//	cond[consts.TcUserId] = db.NotIn(raw)
		//}
	}
	if req.CheckStatus != nil && len(req.CheckStatus) > 0 {
		cond[consts.TcCheckStatus] = db.In(req.Status)
	} else {
		// 因组织总人数统计中使用了 `check_status=2` 的条件，因此后续的查询，默认情况下需要带上 `check_status=2`
		cond[consts.TcCheckStatus] = consts.AppCheckStatusSuccess
	}

	// 部门条件
	if req.DepartmentId != nil {
		if req.AuthFilter && !perContext.HasManageAllDept() {
			cond[consts.TcUserId+" "] = db.In(db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2 and department_id = ? and department_id in ? ", orgId, *req.DepartmentId, perContext.GetManageDeptIds()))
		} else {
			// 查询部门下所有的子部门
			childDeptIds, err := domain.GetDeptChildrenIdsById(orgId, *req.DepartmentId)
			if err != nil {
				logger.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
			allDeptIds := append(childDeptIds, *req.DepartmentId)
			cond[consts.TcUserId+" "] = db.In(db.Raw("select  DISTINCT(user_id) from ppm_org_user_department where org_id = ? and is_delete = 2 and department_id in ?", orgId, allDeptIds))
		}
	}
	// 角色条件
	if req.RoleId != nil {
		if req.AuthFilter && !perContext.HasManageAllRole() {
			cond[consts.TcUserId+"  "] = db.In(db.Raw("select user_id from ppm_rol_role_user where org_id = ? and is_delete = 2 and role_id = ? and role_id in ?", orgId, *req.RoleId, perContext.GetManageRoleIds()))
		} else {
			cond[consts.TcUserId+"  "] = db.In(db.Raw("select user_id from ppm_rol_role_user where org_id = ? and is_delete = 2 and role_id = ?", orgId, *req.RoleId))
		}
	}
	// 职级条件
	if req.PositionId != nil {
		cond[consts.TcUserId+"   "] = db.In(db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2 and org_position_id = ?", orgId, *req.PositionId))
	}

	// 开启权限过滤
	union1 := &db.Union{}
	if req.AuthFilter {
		if !perContext.HasManageAllDept() {
			union1 = union1.Or(db.Cond{
				//用户
				consts.TcUserId: db.In(db.Raw("select d1.user_id from ppm_org_user_department d1 where d1.org_id = ? and d1.is_delete = 2  and d1.department_id in ? ", orgId, perContext.GetManageDeptIds())),
			})
		}
		if !perContext.HasManageAllRole() {
			union1 = union1.Or(db.Cond{
				consts.TcUserId: db.In(db.Raw("select r1.user_id from ppm_rol_role_user r1 where r1.org_id = ? and r1.is_delete = 2 and r1.role_id in ?", orgId, perContext.GetManageRoleIds())),
			})
		}
	}
	var order interface{}
	if req.Order != "" {
		order = db.Raw(req.Order)
	} else {
		// 默认，组织列表中，需要将超管、普通管理员放在前面
		adminUidStr := strs.Int64Implode(adminUids, ",")
		normalAdminUidStr := strs.Int64Implode(normalAdminUids, ",")
		if len(adminUids) > 0 && len(normalAdminUids) > 0 {
			order = db.Raw(fmt.Sprintf(" case when user_id in (%s) then 0 when user_id in (%s) then 1 else 2 end, create_time asc ", adminUidStr, normalAdminUidStr))
		} else if len(adminUids) > 0 && len(normalAdminUids) == 0 {
			order = db.Raw(fmt.Sprintf(" case when user_id in (%s) then 0 else 2 end, create_time asc ", adminUidStr))
		} else if len(adminUids) == 0 && len(normalAdminUids) > 0 {
			order = db.Raw(fmt.Sprintf(" case when user_id in (%s) then 0 else 2 end, create_time asc ", normalAdminUidStr))
		} else {
			order = db.Raw(fmt.Sprintf(" create_time asc "))
		}
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
	//查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUsersRmDupli(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userDepartmentMap := map[int64][]resp.UserDeptPositionData{}
	for _, info := range userDeptBindInfoList {
		userDepartmentMap[info.UserId] = append(userDepartmentMap[info.UserId], resp.UserDeptPositionData{
			DepartmentId:  info.DepartmentId,
			IsLeader:      info.IsLeader,
			DeparmentName: info.DepartmentName,
			PositionId:    info.OrgPositionId,
			PositionName:  info.PositionName,
			PositionLevel: info.PositionLevel,
		})
	}

	//查询角色信息
	userRoleBindInfoList, dbErr := domain.GetUserRoleBindListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userRoleMap := map[int64][]resp.UserRoleData{}
	for _, user := range userRoleBindInfoList {
		userRoleMap[user.UserId] = append(userRoleMap[user.UserId], resp.UserRoleData{
			RoleId:   user.RoleId,
			RoleName: user.RoleName,
		})
	}
	// 查询用户的管理组
	userAdminGroupInfoMap, err := domain.GetUserAdminGroupBindListByUsers(orgId, userIds)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	convertedAdminGroupMap := make(map[int64][]bo.ManageGroupInfoBo, 0)
	copyer.Copy(userAdminGroupInfoMap, &convertedAdminGroupMap)

	res := &resp.UserListResp{
		Total: int64(total),
		List:  []*resp.OrgMemberInfoReq{},
	}

	//查询组织创建者
	var orgInfo po.PpmOrgOrganization
	dbErr = store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcId: orgId,
	}, &orgInfo)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgNotExist
		} else {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}
	// 查询默认管理组
	defaultAdminGroup, dbErr := domain.GetOrgDefaultAdminGroupForPolaris(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	// 获取用户手机号码
	globalUsersMap, err2 := getMobileByUserIds(userIds)
	if err2 != nil {
		return nil, err2
	}

	for _, orgUser := range orgMemberList {
		if _, ok := userInfoMap[orgUser.UserId]; !ok {
			continue
		}
		mobile := ""
		hasPassword := false
		if globalUser, ok := globalUsersMap[orgUser.UserId]; ok {
			mobile = globalUser.Mobile
			if globalUser.Password != "" {
				hasPassword = true
			}
		}
		// 账号名登录的密码存在user表中
		if !hasPassword {
			if user, ok := userInfoMap[orgUser.UserId]; ok {
				if user.Password != "" {
					hasPassword = true
				}
			}
		}
		userInfo := userInfoMap[orgUser.UserId]
		loginName := userInfo.LoginName
		if mobile != "" {
			loginName = mobile
		}
		u := &resp.OrgMemberInfoReq{
			UserID:           userInfo.Id,
			LoginName:        loginName,
			Name:             userInfo.Name,
			NamePy:           userInfo.NamePinyin,
			Avatar:           userInfo.Avatar,
			Email:            userInfo.Email,
			PhoneRegion:      userInfo.MobileRegion,
			PhoneNumber:      mobile,
			CreateTime:       orgUser.CreateTime,
			StatusChangeTime: orgUser.StatusChangeTime,
			Status:           orgUser.Status,
			IsCreator:        false,
			EmpNo:            orgUser.EmpNo,
			WeiboIds:         strings.Split(orgUser.WeiboIds, ","),
			HasPassword:      hasPassword,
		}
		// 处理strings.Split函数分割后的切片长度最少为1的问题
		if len(u.WeiboIds) == 1 && u.WeiboIds[0] == "" {
			u.WeiboIds = make([]string, 0)
		}
		if _, ok := userDepartmentMap[orgUser.UserId]; ok {
			u.DepartmentList = userDepartmentMap[orgUser.UserId]
		} else {
			u.DepartmentList = []resp.UserDeptPositionData{}
		}
		if _, ok := userRoleMap[orgUser.UserId]; ok {
			u.RoleList = userRoleMap[orgUser.UserId]
		} else {
			u.RoleList = []resp.UserRoleData{}
		}
		if groups, ok := convertedAdminGroupMap[orgUser.UserId]; ok {
			u.AdminGroupList = groups
		} else {
			// 极星：返回一个默认角色
			if defaultAdminGroup != nil {
				groupBo := bo.ManageGroupInfoBo{}
				copyer.Copy(defaultAdminGroup, &groupBo)
				u.AdminGroupList = []bo.ManageGroupInfoBo{
					groupBo,
				}
			}
		}
		if orgUser.UserId == orgInfo.Creator {
			u.IsCreator = true
		}
		if orgUser.UserId == orgInfo.Owner {
			u.IsOwner = true
		}

		res.List = append(res.List, u)
	}

	return res, nil
}

func getMobileByUserIds(userIds []int64) (map[int64]*po.PpmOrgGlobalUser, errs.SystemErrorInfo) {
	if len(userIds) == 0 {
		return map[int64]*po.PpmOrgGlobalUser{}, nil
	}

	relations := make([]*po.PpmOrgGlobalUserRelation, 0, len(userIds))
	err := store.Mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(userIds),
	}, &relations)
	if err != nil {
		logger.ErrorF("[getMobileByUserIds] err:%v, userIds:%v", err, userIds)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	globalIdToUserIdMap := make(map[int64]int64)
	globalIds := make([]int64, 0, len(relations))
	for _, m := range relations {
		globalIdToUserIdMap[m.GlobalUserId] = m.UserId
		globalIds = append(globalIds, m.GlobalUserId)
	}

	globalUsers := make([]*po.PpmOrgGlobalUser, 0, len(globalIds))
	err = store.Mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(globalIds),
	}, &globalUsers)
	if err != nil {
		logger.ErrorF("[getMobileByUserIds] err:%v, globalIds:%v", err, globalIds)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	result := make(map[int64]*po.PpmOrgGlobalUser, len(globalUsers))

	for _, user := range globalUsers {
		result[globalIdToUserIdMap[user.Id]] = user
	}

	return result, nil
}

// OrgMemberStat 组织成员数
func OrgMemberStat(orgId int64, perContext *inner_resp.OrgUserPerContext) (*resp.UserStatResp, errs.SystemErrorInfo) {
	// 开启权限过滤
	union := &db.Union{}
	//if !perContext.HasManageAllDept() {
	//	union = union.Or(db.Cond{
	//		//用户
	//		consts.TcUserId: db.In(db.Raw("select d1.user_id from ppm_org_user_department d1 where d1.org_id = ? and d1.is_delete = 2  and d1.department_id in ? ", orgId, perContext.GetManageDeptIds())),
	//	})
	//}
	//if !perContext.HasManageAllRole() {
	//	union = union.Or(db.Cond{
	//		consts.TcUserId: db.In(db.Raw("select r1.user_id from ppm_rol_role_user r1 where r1.org_id = ? and r1.is_delete = 2 and r1.role_id in ?", orgId, perContext.GetManageRoleIds())),
	//	})
	//}

	allCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusEnable,
	}, union)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	//// 非系统管理员只查询全部
	//if !perContext.HasAllPermission() {
	//	return &resp.UserStatResp{
	//		AllCount: int64(allCount),
	//	}, nil
	//}

	forbiddenCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcStatus:      consts.AppStatusDisabled,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	resignedCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.OrgUserStatusResigned,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	unreceivedCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	unallocatedCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusEnable,
		consts.TcUserId:      db.NotIn(db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2", orgId)),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	//待审核：用户，未被删除，待审核。
	waitAuditCount, dbErr := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusWait,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	return &resp.UserStatResp{
		AllCount:         int64(allCount),
		UnallocatedCount: int64(unallocatedCount),
		UnreceivedCount:  int64(unreceivedCount),
		ForbiddenCount:   int64(forbiddenCount),
		ResignedCount:    int64(resignedCount),
		WaitAuditCount:   int64(waitAuditCount),
	}, nil
}

// ExportOrgMemberList 导出组织成员
func ExportOrgMemberList(orgId, userId int64, input req.ExportAddressListReq, perContext *inner_resp.OrgUserPerContext) (string, errs.SystemErrorInfo) {
	if len(input.ExportField) == 0 {
		return "", errs.ExportFieldIsNull
	}

	info, err := GetOrgMemberList(orgId, req.UserListReq{
		SearchCode:   input.SearchCode,
		IsAllocate:   input.IsAllocate,
		Status:       input.Status,
		RoleId:       input.RoleId,
		DepartmentId: input.DepartmentId,
		Page:         0,
		Size:         0,
	}, perContext)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	relatePath := "/user" + "/org_" + strconv.FormatInt(orgId, 10)
	excelDir := conf.Cfg.Resource.RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		logger.Error(mkdirErr)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := "通讯录.xlsx"
	excelPath := excelDir + "/" + fileName
	url := conf.Cfg.Resource.LocalDomain + relatePath + "/" + fileName

	var file *xlsx.File
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, ioErr := file.AddSheet("Sheet1")
	if ioErr != nil {
		logger.Error(ioErr)
		return "", errs.BuildSystemErrorInfo(errs.SystemError, ioErr)
	}

	row = sheet.AddRow()

	if util.FieldInUpdate(input.ExportField, "name") {
		cell = row.AddCell()
		cell.Value = "成员"
	}
	if util.FieldInUpdate(input.ExportField, "loginName") {
		cell = row.AddCell()
		cell.Value = "用户名"
	}
	if util.FieldInUpdate(input.ExportField, "empNo") {
		cell = row.AddCell()
		cell.Value = "工号"
	}
	if util.FieldInUpdate(input.ExportField, "mobile") {
		cell = row.AddCell()
		cell.Value = "手机"
	}
	if util.FieldInUpdate(input.ExportField, "email") {
		cell = row.AddCell()
		cell.Value = "邮箱"
	}
	if util.FieldInUpdate(input.ExportField, "department") {
		cell = row.AddCell()
		cell.Value = "部门/职级"
	}
	if util.FieldInUpdate(input.ExportField, "isLeader") {
		cell = row.AddCell()
		cell.Value = "部门负责人"
	}
	if util.FieldInUpdate(input.ExportField, "role") {
		cell = row.AddCell()
		cell.Value = "角色"
	}
	if util.FieldInUpdate(input.ExportField, "statusChangeTime") {
		cell = row.AddCell()
		cell.Value = "冻结时间"
	}
	if util.FieldInUpdate(input.ExportField, "createTime") {
		cell = row.AddCell()
		cell.Value = "创建时间"
	}

	for _, userInfo := range info.List {
		row = sheet.AddRow()

		if util.FieldInUpdate(input.ExportField, "name") {
			cell = row.AddCell()
			cell.Value = userInfo.Name
		}

		if util.FieldInUpdate(input.ExportField, "loginName") {
			cell = row.AddCell()
			cell.Value = userInfo.LoginName
		}

		if util.FieldInUpdate(input.ExportField, "empNo") {
			cell = row.AddCell()
			cell.Value = userInfo.EmpNo
		}
		if util.FieldInUpdate(input.ExportField, "mobile") {
			cell = row.AddCell()
			cell.Value = userInfo.PhoneNumber
		}
		if util.FieldInUpdate(input.ExportField, "email") {
			cell = row.AddCell()
			cell.Value = userInfo.Email
		}
		if util.FieldInUpdate(input.ExportField, "department") {
			cell = row.AddCell()
			department := ""
			for _, data := range userInfo.DepartmentList {
				department += data.DeparmentName + "/" + data.PositionName + ","
			}
			if len(department) > 0 {
				department = department[0 : len(department)-1]
			}
			cell.Value = department
		}
		if util.FieldInUpdate(input.ExportField, "isLeader") {
			cell = row.AddCell()
			value := "否"
			if input.DepartmentId != nil && *input.DepartmentId != 0 {
				for _, data := range userInfo.DepartmentList {
					if data.DepartmentId == *input.DepartmentId && data.IsLeader == 1 {
						value = "是"
					}
				}
			}
			cell.Value = value
		}
		if util.FieldInUpdate(input.ExportField, "role") {
			cell = row.AddCell()
			role := ""
			for _, data := range userInfo.RoleList {
				role += data.RoleName + ","
			}
			if len(role) > 0 {
				role = role[0 : len(role)-1]
			}
			cell.Value = role
		}
		if util.FieldInUpdate(input.ExportField, "statusChangeTime") {
			cell = row.AddCell()
			cell.Value = userInfo.StatusChangeTime.Format(consts.AppTimeFormat)
		}
		if util.FieldInUpdate(input.ExportField, "createTime") {
			cell = row.AddCell()
			cell.Value = userInfo.CreateTime.Format(consts.AppTimeFormat)
		}
	}

	saveErr := file.Save(excelPath)
	if saveErr != nil {
		logger.Error(saveErr)
		return "", errs.SystemError
	}

	return url, nil
}

// GetOrgMemberInfoById 组织成员信息
func GetOrgMemberInfoById(orgId int64, userId int64) (*resp.OrgMemberInfoReq, errs.SystemErrorInfo) {
	baseOrgInfo, err := domain.GetBaseOrgOutInfo(orgId)
	if err != nil {
		logger.ErrorF("[GetOrgMemberInfoById] GetBaseOrgOutInfo err:%v, orgId:%v", err, orgId)
		return nil, err
	}
	orgMemberInfo, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			if baseOrgInfo.OutOrgId != "" {
				return nil, errs.OrgUserInvalidErr
			}
			return nil, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	//查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptPositionDataList := make([]resp.UserDeptPositionData, 0)
	for _, info := range userDeptBindInfoList {
		deptPositionDataList = append(deptPositionDataList, resp.UserDeptPositionData{
			DepartmentId:  info.DepartmentId,
			IsLeader:      info.IsLeader,
			DeparmentName: info.DepartmentName,
			PositionId:    info.OrgPositionId,
			PositionName:  info.PositionName,
			PositionLevel: info.PositionLevel,
		})
	}

	// 查询角色信息
	userRoleBindInfoList, dbErr := domain.GetUserRoleBindListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userRoleDataList := make([]resp.UserRoleData, 0)
	for _, user := range userRoleBindInfoList {
		userRoleDataList = append(userRoleDataList, resp.UserRoleData{
			RoleId:   user.RoleId,
			RoleName: user.RoleName,
		})
	}

	memberInfo := resp.OrgMemberInfoReq{
		UserID:           orgMemberInfo.UserId,
		LoginName:        orgMemberInfo.LoginName,
		Name:             orgMemberInfo.Name,
		NamePy:           orgMemberInfo.NamePinyin,
		Avatar:           orgMemberInfo.Avatar,
		Email:            orgMemberInfo.Email,
		PhoneRegion:      orgMemberInfo.MobileRegion,
		PhoneNumber:      orgMemberInfo.Mobile,
		CreateTime:       orgMemberInfo.CreateTime,
		StatusChangeTime: orgMemberInfo.StatusChangeTime,
		Status:           orgMemberInfo.Status,
		IsCreator:        userId == orgMemberInfo.OrgCreator,
		IsOwner:          userId == orgMemberInfo.OrgOwner,
		EmpNo:            orgMemberInfo.EmpNo,
		WeiboIds:         strings.Split(orgMemberInfo.WeiboIds, ","),
		DepartmentList:   deptPositionDataList,
		RoleList:         userRoleDataList,
	}
	// 处理strings.Split函数分割后的切片长度最少为1的问题
	if len(memberInfo.WeiboIds) == 1 && memberInfo.WeiboIds[0] == "" {
		memberInfo.WeiboIds = make([]string, 0)
	}
	return &memberInfo, nil
}

// AddOrgMember 增加组织成员，添加关联以及加入顶级部门
// inCheck：是否需要被审核
func AddOrgMember(orgId, userId int64, operatorId int64, inCheck bool, inDisabled bool) errs.SystemErrorInfo {
	orgMember, dbErr := domain.GetNewestOrgMemberByOrgAndUser(orgId, userId)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return errs.MysqlOperateError
		}
	}
	//关联不存在或者已删除，或者审核未通过，此时允许新增关联
	if orgMember == nil || orgMember.IsDelete == consts.AppIsDeleted || orgMember.CheckStatus == consts.AppCheckStatusFail {
		logger.ErrorF("用户%d和组织%d需要做关联", userId, orgId)
		//上锁
		lockKey := fmt.Sprintf("%s%d:%d", consts.UserAndOrgRelationLockKey, orgId, userId)
		lockUuid := uuid.NewUuid()

		suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, lockUuid)
		if redisErr != nil {
			logger.Error(redisErr)
			return errs.TryDistributedLockError
		}
		if suc {
			defer func() {
				if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, lockUuid); redisErr != nil {
					logger.Error(redisErr)
				}
			}()
			//二次check
			userOrgRelation, dbErr := domain.GetNewestOrgMemberByOrgAndUser(orgId, userId)
			if dbErr != nil {
				if dbErr != db.ErrNoMoreRows {
					logger.Error(dbErr)
					return errs.MysqlOperateError
				}
			}
			if userOrgRelation == nil || userOrgRelation.IsDelete == consts.AppIsDeleted || userOrgRelation.CheckStatus == consts.AppCheckStatusFail {
				//组织用户做关联
				logger.ErrorF("用户%d和组织%d开始关联", userId, orgId)
				_, dbErr = domain.BindUserOrgRelation(orgId, userId, false, inCheck, inDisabled, operatorId)
				//判断关联是否失败
				if dbErr != nil {
					logger.Error(dbErr)
					return errs.MysqlOperateError
				}
			}
		}
	}
	clearErr := domain.ClearBaseUserInfo(orgId, userId)
	if clearErr != nil {
		logger.Error(clearErr)
	}
	return nil
}
