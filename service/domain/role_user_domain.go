package domain

import (
	"fmt"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// BindRoleUsers 新增角色用户关系
func BindRoleUsers(orgId, operator, roleId int64, userIds []int64) errs.SystemErrorInfo {
	if len(userIds) == 0 {
		return errs.RoleUsersIsEmpty
	}
	// 获取锁
	{
		newUUID := uuid.NewUuid()
		lockKey := fmt.Sprintf("%s:%d", consts.ModifyUserRoleLock, orgId)
		suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, newUUID)
		if redisErr != nil {
			logger.Error(redisErr)
			return errs.TryDistributedLockError
		}
		if suc {
			defer func() {
				if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, newUUID); redisErr != nil {
					logger.Error(redisErr)
				}
			}()
		} else {
			return errs.RoleUserRefModifyBusy
		}
	}

	_, dbErr := GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return errs.RoleNotExist
		}
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	// 验证用户信息
	userIds = slice.SliceUniqueInt64(userIds)
	userInfos, dbErr := GetEnableOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	if len(userIds) != len(userInfos) {
		return errs.OrgMemberNotExistOrDisable
	}
	// 过滤已经存在的
	userRoleBindList, dbErr := GetUserRoleBindListByRole(orgId, roleId)

	refMap := make(map[int64]bo.UserRoleBo)
	for _, v := range userRoleBindList {
		refMap[v.UserId] = v
	}

	if dbErr != nil {
		return errs.MysqlOperateError
	}
	for i := 0; i < len(userIds); i++ {
		if _, ok := refMap[userIds[i]]; ok {
			userIds = append(userIds[:i], userIds[i+1:]...)
			i--
		}
	}

	if len(userIds) == 0 {
		return nil
	}

	addList := make([]interface{}, 0)
	for _, userId := range userIds {
		roleUserRef := po.PpmRolRoleUser{
			Id:       snowflake.Id(),
			OrgId:    orgId,
			RoleId:   roleId,
			UserId:   userId,
			Creator:  operator,
			Updator:  operator,
			IsDelete: consts.AppIsNoDelete,
		}
		addList = append(addList, roleUserRef)
	}
	dbErr = store.Mysql.BatchInsert(&po.PpmRolRoleUser{}, addList)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	return nil

}

// UnbindRoleOfUsers 移除用户角色关系
func UnbindRoleOfUsers(orgId, operator, roleId int64, userIds []int64) errs.SystemErrorInfo {
	if len(userIds) == 0 {
		return nil
	}
	// 获取锁
	{
		newUUID := uuid.NewUuid()
		lockKey := fmt.Sprintf("%s:%d", consts.ModifyUserRoleLock, orgId)
		suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, newUUID)
		if redisErr != nil {
			logger.Error(redisErr)
			return errs.TryDistributedLockError
		}
		if suc {
			defer func() {
				if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, newUUID); redisErr != nil {
					logger.Error(redisErr)
				}
			}()
		} else {
			return errs.RoleUserRefModifyBusy
		}
	}

	_, dbErr := GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return errs.RoleNotExist
		}
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	if len(userIds) == 0 {
		return nil
	}

	_, dbErr = store.Mysql.UpdateSmartWithCond(consts.TableRoleUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcRoleId:   roleId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operator,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// TransUnbindUserHaveRoles 移除用户的所以角色关系
func TransUnbindUserHaveRoles(orgId, operator int64, userIds []int64, tx sqlbuilder.Tx) error {
	_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableRoleUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operator,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	return nil
}

//  GetUserRoleBindListByUsers 查询用户绑定的角色信息
func GetUserRoleBindListByUsers(orgId int64, userIds []int64) ([]bo.UserRoleBo, error) {
	if len(userIds) == 0 {
		return []bo.UserRoleBo{}, nil
	}
	return _getUserRoleBindList(orgId, nil, userIds)
}

// GetUserRoleBindListByUser 查询用户绑定的角色信息
func GetUserRoleBindListByUser(orgId int64, userId int64) ([]bo.UserRoleBo, error) {
	return GetUserRoleBindListByUsers(orgId, []int64{userId})
}

// GetUserRoleBindListByRoles 查询用户绑定的角色信息
func GetUserRoleBindListByRoles(orgId int64, roleIds []int64) ([]bo.UserRoleBo, error) {
	return _getUserRoleBindList(orgId, roleIds, nil)
}

// GetUserRoleBindListByRole 查询用户绑定的角色信息
func GetUserRoleBindListByRole(orgId int64, roleId int64) ([]bo.UserRoleBo, error) {
	return GetUserRoleBindListByRoles(orgId, []int64{roleId})

}

// _getUserRoleBindList 查询用户角色绑定信息列表
func _getUserRoleBindList(orgId int64, roles, userIds []int64) ([]bo.UserRoleBo, error) {
	cond := db.Cond{
		"u." + consts.TcRoleId:   db.Raw("d.id"),
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"d." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcOrgId:    orgId,
		"d." + consts.TcOrgId:    orgId,
	}
	if len(roles) > 0 {
		cond["d."+consts.TcId] = db.In(roles)
	}
	if len(userIds) > 0 {
		cond["u."+consts.TcUserId] = db.In(userIds)
	}
	if len(userIds) == 0 {
		return []bo.UserRoleBo{}, nil
	}
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	var pos []bo.UserRoleBo
	dbErr = conn.Select(db.Raw("u.user_id,u.role_id,d.name as role_name")).From("ppm_rol_role_user u", "ppm_rol_role d").Where(cond).All(&pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return pos, nil
}

// GetUserAdminGroupBindListByUser 查询用户绑定的管理组
func GetUserAdminGroupBindListByUser(orgId int64, userId int64) (map[int64][]*bo.ManageGroupInfoBo, error) {
	return GetUserAdminGroupBindListByUsers(orgId, []int64{userId})
}

// GetUserAdminGroupBindListByUsers 查询用户绑定的管理组
func GetUserAdminGroupBindListByUsers(orgId int64, userIds []int64) (map[int64][]*bo.ManageGroupInfoBo, error) {
	result := make(map[int64][]*bo.ManageGroupInfoBo, 0)
	if len(userIds) == 0 {
		return result, nil
	}
	pos, err := GetManageGroupListByUsers(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return result, err
	}
	bos := make([]*bo.ManageGroupInfoBo, 0)
	copyer.Copy(pos, &bos)
	for _, boObj := range bos {
		tmpUserIdsStr := boObj.UserIds
		tmpUserIds := make([]int64, 0)
		json.FromJson(tmpUserIdsStr, &tmpUserIds)
		tmpUserIds = slice.SliceUniqueInt64(tmpUserIds)
		for _, userId := range tmpUserIds {
			if _, ok := result[userId]; ok {
				result[userId] = append(result[userId], boObj)
			} else {
				result[userId] = []*bo.ManageGroupInfoBo{boObj}
			}
		}
	}
	return result, nil
}

// GetUserIdsByRoles 查询角色和部门下的人员ID
func GetUserIdsByRoles(orgId int64, deptIds, roleIds []int64) ([]int64, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	var userIds []int64
	dbErr = conn.Select(db.Raw("Distinct(u.user_id) as user_ids")).From("ppm_rol_role_user u", "ppm_rol_role d").Where(db.Cond{
		"u." + consts.TcRoleId:   db.Raw("d.id"),
		"d." + consts.TcRoleId:   db.In(roleIds),
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"d." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcOrgId:    orgId,
		"d." + consts.TcOrgId:    orgId,
	}).All(&userIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	var userIds1 []int64
	dbErr = conn.Select(db.Raw("Distinct(u.user_id) as user_ids")).From("ppm_org_user_department u", "ppm_org_department d").Where(db.Cond{
		"u." + consts.TcDepartmentId: db.Raw("d.id"),
		"d." + consts.TcDepartmentId: db.In(deptIds),
		"u." + consts.TcIsDelete:     consts.AppIsNoDelete,
		"d." + consts.TcIsDelete:     consts.AppIsNoDelete,
		"u." + consts.TcOrgId:        orgId,
		"d." + consts.TcOrgId:        orgId,
	}).All(&userIds1)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	userIds = append(userIds, userIds1...)
	userIds = slice.SliceUniqueInt64(userIds)
	return userIds, nil
}

func GetRoleUserIds(orgId int64) (map[int64][]int64, errs.SystemErrorInfo) {
	var roleUsers []po.PpmRolRoleUser
	err := store.Mysql.SelectAllByCond(consts.TableRoleUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, &roleUsers)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := map[int64][]int64{}
	for _, user := range roleUsers {
		if ok, _ := slice.Contain(res[user.RoleId], user.UserId); !ok {
			res[user.RoleId] = append(res[user.RoleId], user.UserId)
		}
	}

	return res, nil
}
