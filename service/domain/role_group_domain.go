package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

var DefaultRoleGroupName = "默认"

// CreateRoleGroup 创建角色组
func CreateRoleGroup(orgId, operatorUid int64, name string, remark string) (int64, error) {
	roleGroup := po.PpmRolRoleGroup{
		Id:      snowflake.Id(),
		OrgId:   orgId,
		Name:    name,
		Remark:  remark,
		Creator: operatorUid,
		Updator: operatorUid,
	}
	dbErr := store.Mysql.Insert(&roleGroup)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return roleGroup.Id, nil
}

// UpdateRoleGroup 修改角色组
func UpdateRoleGroup(orgId, operatorUid int64, groupId int64, version int, name string) (int, error) {
	//加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableRoleGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       groupId,
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

// DeleteRoleGroup 删除角色组
func DeleteRoleGroup(orgId, operatorUid int64, groupId int64) (int, error) {
	// 加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableRoleGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       groupId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorUid,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// GetRoleGroup 根据ID获取角色组
func GetRoleGroup(orgId, groupId int64) (*bo.RoleGroupBo, error) {
	if orgId == 0 || groupId == 0 {
		return nil, errs.OrgRoleGroupNotExist
	}
	var roleGroup po.PpmRolRoleGroup
	dbErr := store.Mysql.SelectOneByCond(consts.TableRoleGroup, db.Cond{
		consts.TcId:       groupId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &roleGroup)

	if dbErr != nil {
		return nil, dbErr
	}

	var roleGroupBo bo.RoleGroupBo
	_ = copyer.Copy(roleGroup, &roleGroupBo)
	return &roleGroupBo, nil
}

// GetRoleGroupByName 根据Name获取角色组
func GetRoleGroupByName(orgId int64, groupName string) (*bo.RoleGroupBo, error) {
	if orgId == 0 || groupName == "" {
		return nil, errs.OrgRoleGroupNotExist
	}
	var roleGroup po.PpmRolRoleGroup
	dbErr := store.Mysql.SelectOneByCond(consts.TableRoleGroup, db.Cond{
		consts.TcName:     groupName,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &roleGroup)
	if dbErr != nil {
		return nil, dbErr
	}
	var roleGroupBo bo.RoleGroupBo
	_ = copyer.Copy(roleGroup, &roleGroupBo)
	return &roleGroupBo, nil
}

// GetRoleGroupListByOrg
func GetRoleGroupListByOrg(orgId int64) ([]bo.RoleGroupBo, error) {

	var groups []po.PpmRolRoleGroup
	dbErr := store.Mysql.SelectAllByCond(consts.TableRoleGroup, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, &groups)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	var result []bo.RoleGroupBo
	_ = copyer.Copy(groups, &result)
	return result, nil
}

func GetRoleGroupListIds(orgId int64, ids []int64) ([]bo.RoleGroupBo, error) {
	if len(ids) == 0 {
		return []bo.RoleGroupBo{}, nil
	}
	var groups []po.PpmRolRoleGroup
	dbErr := store.Mysql.SelectAllByCond(consts.TableRoleGroup, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(ids),
	}, &groups)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	var result []bo.RoleGroupBo
	_ = copyer.Copy(groups, &result)
	return result, nil
}

func GetRoleCounts(orgId int64) (map[int64]uint64, error) {
	connect, err := store.Mysql.GetConnect()
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	roleUserCountList := make([]bo.RoleUserCount, 0)
	err = connect.Select(db.Raw("count(user_id) as count"), "role_id").From(consts.TableRoleUser).Where(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}).GroupBy("role_id").All(&roleUserCountList)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	roleUserCountMap := map[int64]uint64{}
	for _, roleUserCount := range roleUserCountList {
		roleUserCountMap[roleUserCount.RoleID] = roleUserCount.Count
	}
	return roleUserCountMap, nil
}
