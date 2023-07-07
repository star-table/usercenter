package service

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// BindUserRoles 绑定用户角色 全量
func BindUserRoles(orgId, operator, userId int64, roleIds []int64, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	if userId == 0 || len(roleIds) == 0 {
		return false, nil
	}

	userInfos, dbErr := domain.GetEnableOrgMemberBaseInfoListByUsers(orgId, []int64{userId})
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	if len(userInfos) == 0 {
		return false, errs.OrgMemberNotExistOrDisable
	}

	err := _bindUserRoles(orgId, operator, userId, roleIds, perContext)
	if err != nil {
		return false, err
	}
	return true, nil
}

// _bindUserRoles 绑定用户的角色 全量
func _bindUserRoles(orgId, operator, userId int64, roleIds []int64, perContext *inner_resp.OrgUserPerContext) errs.SystemErrorInfo {
	if len(roleIds) > 0 {
		//去重
		roleIds = slice.SliceUniqueInt64(roleIds)

		roles, dbErr := domain.GetRoleListByIds(orgId, roleIds)
		if dbErr != nil {
			return errs.MysqlOperateError
		}
		if len(roleIds) != len(roles) {
			return errs.RoleNotExist
		}
	}

	userRoleBindList, dbErr := domain.GetUserRoleBindListByUser(orgId, userId)
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	oldBindMap := make(map[int64]bo.UserRoleBo)
	for _, bindInfo := range userRoleBindList {
		oldBindMap[bindInfo.RoleId] = bindInfo
	}

	addBindList := make([]interface{}, 0)
	delRoleIdList := make([]int64, 0)
	for _, rId := range roleIds {
		if _, ok := oldBindMap[rId]; !ok {
			// 非自己管理的
			if !perContext.HasManageRole(rId) {
				return errs.ForbiddenAccess
			}
			roleUser := po.PpmRolRoleUser{
				Id:       snowflake.Id(),
				OrgId:    orgId,
				RoleId:   rId,
				UserId:   userId,
				Creator:  operator,
				Updator:  operator,
				IsDelete: consts.AppIsNoDelete,
			}
			addBindList = append(addBindList, roleUser)
		} else {
			// 存在的则无需操作
			delete(oldBindMap, rId)
		}
	}

	// 把剩余的删除掉
	for rId, _ := range oldBindMap {
		// 非自己管理的不进行删除
		if perContext.HasManageRole(rId) {
			delRoleIdList = append(delRoleIdList, rId)
		}
	}

	if len(addBindList) == 0 && len(delRoleIdList) == 0 {
		return nil
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 删除的
		if len(delRoleIdList) > 0 {
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableRoleUser, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcRoleId:   db.In(delRoleIdList),
				consts.TcUserId:   userId,
				consts.TcIsDelete: consts.AppIsNoDelete,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  operator,
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}
		// 新增的
		if len(addBindList) > 0 {
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmRolRoleUser{}, addBindList)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}
		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// BindRoleUsers 绑定角色和用户关联
func BindRoleUsers(orgId, operator int64, roleId int64, userIds []int64) (bool, errs.SystemErrorInfo) {
	if roleId == 0 || len(userIds) == 0 {
		return false, nil
	}
	//去重
	userIds = slice.SliceUniqueInt64(userIds)

	_, dbErr := domain.GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.RoleNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	err := domain.BindRoleUsers(orgId, operator, roleId, userIds)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnbindRoleUsers 移除角色和用户关联
func UnbindRoleUsers(orgId, operator int64, roleId int64, userIds []int64) (bool, errs.SystemErrorInfo) {
	if roleId == 0 || len(userIds) == 0 {
		return false, nil
	}
	//去重
	userIds = slice.SliceUniqueInt64(userIds)

	_, dbErr := domain.GetRole(orgId, roleId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.RoleNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	err := domain.UnbindRoleOfUsers(orgId, operator, roleId, userIds)
	if err != nil {
		return false, err
	}
	return true, nil
}
