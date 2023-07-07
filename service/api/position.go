package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

/**
职级管理接口
*/
type position int

var Position position

// @Security token
// @Summary 创建职级
// @Description 创建职级接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.CreatePositionReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/positions/create [post]
func (p position) CreatePosition(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.CreatePositionReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	id, err := service.CreatePosition(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}

// @Security token
// @Summary 修改职级信息
// @Description 修改职级信息接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.CreatePositionReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/positions/update-info [post]
func (p position) ModifyPositionInfo(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.ModifyPositionInfoReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.ModifyPositionInfo(operator.OrgId, operator.UserId, reqParam.PositionId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 修改职级状态
// @Description 修改职级状态接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.UpdatePositionStatusReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/positions/update-status [post]
func (p position) UpdatePositionStatus(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.UpdatePositionStatusReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdatePositionStatus(operator.OrgId, operator.UserId, reqParam.PositionId, reqParam.Status)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 修改职级状态
// @Description 修改职级状态接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.IdReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/positions/delete [delete]
func (p position) DeletePosition(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.IdReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.DeletePosition(operator.OrgId, operator.UserId, reqParam.Id)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 职级列表查询
// @Description 职级列表查询接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.SearchPositionListReq true "入参"
// @Success 200 {object} []resp.PositionInfoResp
// @Failure 400
// @Router /usercenter/api/v1/positions/list [post]
func (p position) GetPositionList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.SearchPositionListReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.GetPositionList(operator.OrgId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 职级分页列表查询
// @Description 职级分页列表查询接口
// @Tags 职级
// @accept application/json
// @Produce application/json
// @param input body req.SearchPositionPageListReq true "入参"
// @Success 200 {object} resp.PositionPageListResp
// @Failure 400
// @Router /usercenter/api/v1/positions/page [post]
func (p position) GetPositionPageList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.SearchPositionPageListReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.GetPositionPageList(operator.OrgId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}
