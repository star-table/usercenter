package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"upper.io/db.v3"
)

/**
职级数据交互
@author WangShiChang
@version v1.0
@date 2020-10-21
*/

/**
CreatePosition 创建职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func CreatePosition(orgId, operatorUid int64, reqParam req.CreatePositionReq) (int64, error) {
	positionPo := po.LcOrgPosition{
		Id:            snowflake.Id(),
		OrgPositionId: snowflake.Id(),
		OrgId:         orgId,
		Name:          reqParam.Name,
		PositionLevel: reqParam.PositionLevel,
		Remark:        reqParam.Remark,
		Status:        consts.AppStatusEnable,
		Creator:       operatorUid,
		Updator:       operatorUid,
	}
	dbErr := store.Mysql.Insert(&positionPo)
	if dbErr != nil {
		logger.ErrorF("[创建职级] -> db操作异常: %s", dbErr)
		return 0, dbErr
	}
	return positionPo.Id, nil
}

/**
InitPosition 初始化默认职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func InitPosition(orgId, operatorUid int64) error {
	logger.InfoF("[初始化职级] -> orgId: %d", orgId)
	positionPos := make([]interface{}, 0)
	for orgPositionId, m := range consts.DefaultPositions {
		positionPos = append(positionPos, po.LcOrgPosition{
			Id:            snowflake.Id(),
			OrgPositionId: orgPositionId,
			OrgId:         orgId,
			Name:          m[consts.TcName].(string),
			PositionLevel: m[consts.TcPositionLevel].(int),
			Status:        consts.AppStatusEnable,
			Creator:       operatorUid,
			Updator:       operatorUid,
		})
	}
	dbErr := store.Mysql.BatchInsert(&po.LcOrgPosition{}, positionPos)
	if dbErr != nil {
		logger.ErrorF("[初始化职级] -> db操作异常: %s", dbErr)
		return dbErr
	}
	logger.InfoF("[初始化职级]-> 成功  orgId: %d", orgId)

	return nil
}

/**
ModifyPositionInfo 修改职级信息
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func ModifyPositionInfo(orgId, operatorUid int64, orgPositionId int64, reqParam req.ModifyPositionInfoReq) (int64, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TablePosition, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcOrgPositionId: orgPositionId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcName:          reqParam.Name,
		consts.TcPositionLevel: reqParam.PositionLevel,
		consts.TcRemark:        reqParam.Remark,
		consts.TcUpdator:       operatorUid,
	})
	if dbErr != nil {
		logger.ErrorF("[修改职级信息] -> db操作异常: %s", dbErr)
		return 0, dbErr
	}
	return count, nil
}

/**
UpdatePositionStatus 修改职级状态
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func UpdatePositionStatus(orgId, operatorUid int64, orgPositionId int64, status int) (int64, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TablePosition, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcOrgPositionId: orgPositionId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcStatus:  status,
		consts.TcUpdator: operatorUid,
	})
	if dbErr != nil {
		logger.ErrorF("[修改职级状态] -> db操作异常: %s", dbErr)
		return 0, dbErr
	}
	return count, nil
}

/**
RemovePosition 删除职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func RemovePosition(orgId, operatorUid int64, orgPositionId int64) (int64, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TablePosition, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcOrgPositionId: orgPositionId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcStatus:   consts.AppStatusDisabled, // 状态也设置为关闭
		consts.TcUpdator:  operatorUid,
	})
	if dbErr != nil {
		logger.ErrorF("[删除职级] -> db操作异常: %s", dbErr)
		return 0, dbErr
	}
	return count, nil
}

/**
GetPositionById 根据ID获取职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionById(orgId, orgPositionId int64) (*po.LcOrgPosition, error) {
	var positionBo po.LcOrgPosition

	dbErr := store.Mysql.SelectOneByCond(consts.TablePosition, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcOrgPositionId: orgPositionId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, &positionBo)
	if dbErr != nil {
		return nil, dbErr
	}
	return &positionBo, nil
}

/**
GetPositionList 获取职级列表，按照Level正序
status == 0 忽视条件
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionList(orgId int64, status int) ([]po.LcOrgPosition, error) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if status != 0 {
		cond[consts.TcStatus] = status
	}
	var positionBos []po.LcOrgPosition
	_, dbErr := store.Mysql.SelectAllByCondWithPageAndOrder(consts.TablePosition,
		cond,
		nil,
		0,
		0,
		"position_level asc, create_time desc",
		&positionBos)

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return positionBos, nil
}

/**
GetPositionPageList 获取职级分页列表，按照Level正序
status == 0 忽视条件
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionPageList(orgId int64, status int, page, size int) ([]po.LcOrgPosition, int64, error) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if status != 0 {
		cond[consts.TcStatus] = status
	}
	var positionBos []po.LcOrgPosition
	total, dbErr := store.Mysql.SelectAllByCondWithPageAndOrder(consts.TablePosition,
		cond,
		nil,
		page,
		size,
		"position_level asc, create_time desc",
		&positionBos)

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, 0, dbErr
	}
	return positionBos, int64(total), nil
}

/**
GetPositionHaveUserList 获取职级包含的人员
status == 0 忽视条件
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionHaveUserList(orgId, orgPositionId int64) ([]po.UserPosition, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	var userPositions []po.UserPosition
	positionAlias := "ps"
	userDeptAlias := "ud"
	dbErr = conn.Select(
		"ps.id as id",
		"ps.org_position_id as org_position_id",
		"ps.name as position_name",
		"ps.position_level as position_level",
		"ud.user_id as user_id",
	).From(consts.TablePosition + " as " + positionAlias).
		Join(consts.TableUserDepartment + " as " + userDeptAlias).
		On(db.Raw(" ud.org_position_id = ps.org_position_id and ud.org_id = ? and ud.is_delete = ?", orgId, consts.AppIsNoDelete)).
		Where(db.Cond{
			positionAlias + "." + consts.TcOrgId:         orgId,
			positionAlias + "." + consts.TcOrgPositionId: orgPositionId,
			positionAlias + "." + consts.TcIsDelete:      consts.AppIsNoDelete,
		}).All(&userPositions)

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return userPositions, nil
}
