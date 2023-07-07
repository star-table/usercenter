package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//特殊角色名称
var DefaultRoleName = []string{
	"系统管理员",
	"子管理员",
}

// CreateRole 创建角色
func CreateRole(orgId, operatorUid int64, groupId int64, name string) (int64, error) {
	rolePo := po.PpmRolRole{
		Id:          snowflake.Id(),
		OrgId:       orgId,
		RoleGroupId: groupId,
		Name:        name,
		Creator:     operatorUid,
		Updator:     operatorUid,
	}
	dbErr := store.Mysql.Insert(&rolePo)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return rolePo.Id, nil
}

// UpdateRole 更新角色
func UpdateRole(orgId, operatorUid int64, roleId int64, name string, version int) (int, error) {
	//加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableRole, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       roleId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcVersion:  version,
	}, mysql.Upd{
		consts.TcName:    name,
		consts.TcUpdator: operatorUid,
		consts.TcVersion: version + 1,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// DeleteRole 删除角色 同时删除关联角色的用户
func DeleteRole(orgId, operatorUid int64, roleId int64) (int, error) {
	count := int64(0)
	var dbErr error = nil
	// 查询用户角色绑定关系
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 删除角色
		count, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableRole, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcId:       roleId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		// 删除用户角色绑定关系
		_, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableRoleUser, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcRoleId:   roleId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		// 删除管理组中绑定的角色关系
		_, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableManageGroup, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			db.Raw("json_search(role_ids, 'one', ?)", roleId): db.IsNotNull(),
		}, mysql.Upd{
			consts.TcRoleIds: db.Raw("json_remove(role_ids,JSON_UNQUOTE(json_search(role_ids, 'one', ?)))", roleId),
			consts.TcUpdator: operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return int(count), nil
}

// MoveRole 移动角色
func MoveRole(orgId, operatorUid int64, groupId, roleId int64, version int) (int, error) {
	// 加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableRole, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       roleId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcVersion:  version,
	}, mysql.Upd{
		consts.TcRoleGroupId: groupId,
		consts.TcUpdator:     operatorUid,
		consts.TcVersion:     version + 1,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// GetRole 获取角色
func GetRole(orgId, roleId int64) (*po.PpmRolRole, error) {
	var role po.PpmRolRole
	dbErr := store.Mysql.SelectOneByCond(consts.TableRole, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcId:       roleId,
	}, &role)
	if dbErr != nil {
		return nil, dbErr
	}

	return &role, nil
}

// GetRoleListByIds 获取角色列表
// len(roleIds) == 0 则返回空切片
func GetRoleListByIds(orgId int64, roleIds []int64) ([]po.PpmRolRole, error) {
	if len(roleIds) == 0 {
		return []po.PpmRolRole{}, nil
	}

	return _getRoleList(orgId, roleIds, 0, consts.AppIsNoDelete)
}

// GetRoleListByIdsAndStatus 获取角色列表
// len(roleIds) == 0 返回空切片
// status <= 0 忽视条件
// isDelete <= 0 忽视条件
func GetRoleListByIdsAndStatus(orgId int64, roleIds []int64, status, isDelete int) ([]po.PpmRolRole, error) {
	if len(roleIds) == 0 {
		return []po.PpmRolRole{}, nil
	}
	return _getRoleList(orgId, roleIds, status, isDelete)
}

// _getRoleList 获取角色列表
// len(roleIds) == 0 忽视条件
// status <= 0 忽视条件
// isDelete <= 0 忽视条件
func _getRoleList(orgId int64, roleIds []int64, status, isDelete int) ([]po.PpmRolRole, error) {

	//拼装条件
	cond := db.Cond{
		consts.TcOrgId: orgId,
	}
	if len(roleIds) > 0 {
		cond[consts.TcId] = db.In(roleIds)
	}
	//if status > 0{
	//	cond[consts.TcStatus] = status
	//}
	if isDelete > 0 {
		cond[consts.TcIsDelete] = isDelete
	}

	var roles []po.PpmRolRole
	dbErr := store.Mysql.SelectAllByCond(consts.TableRole, cond, &roles)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return roles, nil
}

// GetRoleListByOrg 获取角色列表
func GetRoleListByOrg(orgId int64) ([]po.PpmRolRole, error) {

	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}

	var roles []po.PpmRolRole
	dbErr := store.Mysql.SelectAllByCond(consts.TableRole, cond, &roles)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return roles, nil
}

// GetRoleListByOrgWithFilter 获取角色列表，名称模糊匹配
func GetRoleListByOrgWithFilter(orgId int64, query string) ([]po.PpmRolRole, error) {
	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}
	if query != "" {
		cond[consts.TcName] = db.Like("%" + query + "%")
	}

	var roles []po.PpmRolRole
	dbErr := store.Mysql.SelectAllByCond(consts.TableRole, cond, &roles)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return roles, nil
}

// GetRoleListByGroup 根据角色组获取角色
func GetRoleListByGroup(orgId, groupId int64) ([]po.PpmRolRole, error) {
	var roles []po.PpmRolRole
	dbErr := store.Mysql.SelectAllByCond(consts.TableRole, db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcRoleGroupId: groupId,
		consts.TcIsDelete:    consts.AppIsNoDelete,
	}, &roles)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return roles, nil
}

func GetRoleSimpleInfo(orgId int64) ([]inner_resp.SimpleInfo, errs.SystemErrorInfo) {
	info, err := GetRoleListByOrg(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := make([]inner_resp.SimpleInfo, 0)
	for _, role := range info {
		res = append(res, inner_resp.SimpleInfo{
			Id:   role.Id,
			Name: role.Name,
		})
	}

	return res, nil
}

func GetRepeatRoleInfo(orgId int64) ([]inner_resp.RepeatMemberInfo, errs.SystemErrorInfo) {
	var roleInfo []po.PpmRolRole
	err := store.Mysql.SelectAllByCond(consts.TableRole, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcName:     db.In(db.Raw("select name from ppm_rol_role where org_id = ? and is_delete = 2 group by name having count(1) > 1", orgId)),
	}, &roleInfo)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	if len(roleInfo) == 0 {
		return []inner_resp.RepeatMemberInfo{}, nil
	}

	groupIds := make([]int64, 0)
	roleMap := make(map[string][]po.PpmRolRole, 0)
	for _, role := range roleInfo {
		if role.RoleGroupId != int64(0) {
			groupIds = append(groupIds, role.RoleGroupId)
		}
		roleMap[role.Name] = append(roleMap[role.Name], role)
	}
	groupIds = slice.SliceUniqueInt64(groupIds)
	parentMap := make(map[int64]string, 0)
	if len(groupIds) > 0 {
		var parentInfo []po.PpmRolRoleGroup
		parentErr := store.Mysql.SelectAllByCond(consts.TableRoleGroup, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcId:       db.In(groupIds),
		}, &parentInfo)
		if parentErr != nil {
			logger.Error(parentErr)
			return nil, errs.MysqlOperateError
		}

		for _, group := range parentInfo {
			parentMap[group.Id] = group.Name
		}
	}

	res := make([]inner_resp.RepeatMemberInfo, 0)
	for _, list := range roleMap {
		for _, role := range list {
			temp := inner_resp.RepeatMemberInfo{
				Id:         role.Id,
				Name:       role.Name,
				Department: []string{},
			}
			if parentName, ok := parentMap[role.RoleGroupId]; ok {
				temp.Department = []string{parentName}
			}

			res = append(res, temp)
		}
	}

	return res, nil
}
