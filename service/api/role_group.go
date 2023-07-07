package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type roleGroup int

var RoleGroup roleGroup

// @Security token
// @Summary 创建角色组
// @Description 创建角色组接口
// @Tags 角色组
// @accept application/json
// @Produce application/json
// @param input body req.RoleGroupReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/roleGroup/create [post]
func (roleGroup) Create(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var groupReq req.RoleGroupReq
	err = ParseBody(c, &groupReq)
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

	id, err := service.CreateRoleGroup(operator.OrgId, operator.UserId, groupReq)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}

// @Security token
// @Summary 修改角色组
// @Description 修改角色组接口
// @Tags 角色组
// @accept application/json
// @Produce application/json
// @param input body req.UpdateRoleGroupReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/roleGroup/update [put]
func (roleGroup) Update(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var groupReq req.UpdateRoleGroupReq
	err = ParseBody(c, &groupReq)
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

	ok, err := service.UpdateRoleGroup(operator.OrgId, operator.UserId, groupReq.Id, groupReq.RoleGroupReq)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 删除角色组
// @Description 删除角色组接口
// @Tags 角色组
// @accept application/json
// @Produce application/json
// @param input body req.IdReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/roleGroup/delete [delete]
func (roleGroup) Delete(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var groupReq req.IdReq
	err = ParseBody(c, &groupReq)
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
	ok, err := service.DeleteRoleGroup(operator.OrgId, operator.UserId, groupReq.Id)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 组织角色组列表
// @Description 获取组织角色组列表接口
// @Tags 角色组
// @accept application/json
// @Produce application/json
// @Success 200 {object} []resp.RoleGroup
// @Failure 400
// @Router /usercenter/api/v1/roleGroup/list [get]
func (roleGroup) GetGroupList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	groups, err := service.GetGroupList(operator.OrgId)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, groups)
}
