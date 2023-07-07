package service

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
)

/**
职级业务处理
@author WangShiChang
*/

const (
	minPositionLevel = 1
	maxPositionLevel = 20
)

/**
CreatePosition 创建职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func CreatePosition(orgId, operatorUid int64, reqParam req.CreatePositionReq) (int64, errs.SystemErrorInfo) {
	logger.InfoF("[创建职级] -> 参数 reqParam: %s", json.ToJsonIgnoreError(reqParam))
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	if !format.VerifyPositionNameFormat(reqParam.Name) {
		return 0, errs.PositionNameLenErr
	}
	// 等级1到20
	if reqParam.PositionLevel < minPositionLevel || reqParam.PositionLevel > maxPositionLevel {
		return 0, errs.PositionLevelErr
	}

	id, dbErr := domain.CreatePosition(orgId, operatorUid, reqParam)
	if dbErr != nil {
		return 0, errs.MysqlOperateError
	}

	return id, nil
}

/**
ModifyPositionInfo 修改职级信息
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func ModifyPositionInfo(orgId, operatorUid int64, orgPositionId int64, reqParam req.ModifyPositionInfoReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改职级信息] -> 参数 orgPositionId: %d, reqParam: %s", orgPositionId, json.ToJsonIgnoreError(reqParam))
	// 默认的不可修改
	if orgPositionId == consts.PositionManagerId || orgPositionId == consts.PositionMemberId {
		return false, errs.DefaultPositionCannotEditErr
	}
	position, dbErr := domain.GetPositionById(orgId, orgPositionId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.PositionNotExistErr
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	if !format.VerifyPositionNameFormat(reqParam.Name) {
		return false, errs.PositionNameLenErr
	}
	// 等级1到20
	if reqParam.PositionLevel < minPositionLevel || reqParam.PositionLevel > maxPositionLevel {
		return false, errs.PositionLevelErr
	}

	count, dbErr := domain.ModifyPositionInfo(orgId, operatorUid, position.OrgPositionId, reqParam)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

/**
UpdatePositionStatus 修改职级状态
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func UpdatePositionStatus(orgId, operatorUid int64, orgPositionId int64, status int) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改职级状态] -> 参数 orgPositionId: %d, status: %d", orgPositionId, status)
	// 默认的不可修改
	if orgPositionId == consts.PositionManagerId || orgPositionId == consts.PositionMemberId {
		return false, errs.DefaultPositionCannotEditErr
	}
	position, dbErr := domain.GetPositionById(orgId, orgPositionId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.PositionNotExistErr
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 判断是否关联人员
	if status == consts.AppStatusDisabled {
		// 判断是否关联人员
		haveUsers, dbErr := domain.GetPositionHaveUserList(orgId, orgPositionId)
		if dbErr != nil {
			return false, errs.MysqlOperateError
		}
		if len(haveUsers) > 0 {
			return false, errs.PositionHaveUserExistErr
		}
	}

	count, dbErr := domain.UpdatePositionStatus(orgId, operatorUid, position.OrgPositionId, status)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

/**
DeletePosition 删除职级
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func DeletePosition(orgId, operatorUid int64, orgPositionId int64) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[删除职级] -> 参数 orgPositionId: %d", orgPositionId)
	// 默认的不可删除
	if orgPositionId == consts.PositionManagerId || orgPositionId == consts.PositionMemberId {
		return false, errs.DefaultPositionCannotDeleteErr
	}
	position, dbErr := domain.GetPositionById(orgId, orgPositionId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.PositionNotExistErr
		}
		return false, errs.MysqlOperateError
	}

	// 判断是否关联人员
	haveUsers, dbErr := domain.GetPositionHaveUserList(orgId, orgPositionId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	if len(haveUsers) > 0 {
		return false, errs.PositionHaveUserExistErr
	}

	count, dbErr := domain.RemovePosition(orgId, operatorUid, position.OrgPositionId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}

	return count > 0, nil
}

/**
GetPositionList 获取职级列表
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionList(orgId int64, reqParam req.SearchPositionListReq) ([]resp.PositionInfoResp, errs.SystemErrorInfo) {
	positions, dbErr := domain.GetPositionList(orgId, reqParam.Status)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	list := make([]resp.PositionInfoResp, 0)
	// 是否默认
	for _, position := range positions {
		item := resp.PositionInfoResp{
			Id:            position.Id,
			PositionId:    position.OrgPositionId,
			OrgId:         position.OrgId,
			Name:          position.Name,
			PositionLevel: position.PositionLevel,
			Remark:        position.Remark,
			Status:        position.Status,
			Creator:       position.Creator,
			CreateTime:    position.CreateTime,
			Updator:       position.Updator,
			UpdateTime:    position.UpdateTime,
			IsDelete:      position.IsDelete,
			IsDefault:     position.OrgPositionId == consts.PositionManagerId || position.OrgPositionId == consts.PositionMemberId,
		}
		list = append(list, item)
	}
	return list, nil
}

/**
GetPositionPageList 获取职级分页列表
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func GetPositionPageList(orgId int64, reqParam req.SearchPositionPageListReq) (*resp.PositionPageListResp, errs.SystemErrorInfo) {
	positions, total, dbErr := domain.GetPositionPageList(orgId, reqParam.Status, reqParam.Page, reqParam.Size)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	list := make([]resp.PositionInfoResp, 0) // 是否默认
	for _, position := range positions {
		item := resp.PositionInfoResp{
			Id:            position.Id,
			PositionId:    position.OrgPositionId,
			OrgId:         position.OrgId,
			Name:          position.Name,
			PositionLevel: position.PositionLevel,
			Remark:        position.Remark,
			Status:        position.Status,
			Creator:       position.Creator,
			CreateTime:    position.CreateTime,
			Updator:       position.Updator,
			UpdateTime:    position.UpdateTime,
			IsDelete:      position.IsDelete,
			IsDefault:     position.OrgPositionId == consts.PositionManagerId || position.OrgPositionId == consts.PositionMemberId,
		}
		list = append(list, item)
	}
	return &resp.PositionPageListResp{
		Total: total,
		List:  list,
	}, nil
}
