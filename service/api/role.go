package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/service"
)

type role int

var Role role = 1

// @Security token
// @Summary 创建角色
// @Description 创建角色接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.CreateRoleReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/role/create [post]
func (role) Create(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var roleReq req.CreateRoleReq
	err = ParseBody(c, &roleReq)
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

	id, err := service.CreateRole(operator.OrgId, operator.UserId, roleReq)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}

// @Security token
// @Summary 修改角色
// @Description 修改角色接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.UpdateRoleReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/role/update [put]
func (role) Update(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var roleReq req.UpdateRoleReq
	err = ParseBody(c, &roleReq)
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

	ok, err := service.UpdateRole(operator.OrgId, operator.UserId, roleReq.Id, roleReq)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 删除角色
// @Description 删除角色接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.IdReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/role/delete [delete]
func (role) Delete(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var roleReq req.IdReq
	err = ParseBody(c, &roleReq)
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

	ok, err := service.DeleteRole(operator.OrgId, operator.UserId, roleReq.Id)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 移动角色
// @Description 移动角色接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.MoveRoleReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/role/move [put]
func (role) Move(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var roleReq req.MoveRoleReq

	err = ParseBody(c, &roleReq)
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

	ok, err := service.MoveRole(operator.OrgId, operator.UserId, roleReq.RoleGroupId, roleReq.Id)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 角色组和角色列表
// @Description 角色组和角色列表接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.RoleListResp
// @Failure 400
// @Router /usercenter/api/v1/role/list [get]
func (role) GetOrgRoleList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	listResp, err := service.GetOrgRoleListResp(operator.OrgId, permission)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, listResp)
}

// @Security token
// @Summary 角色查询
// @Description 角色查询
// @Tags 选人控件
// @accept application/json
// @Produce application/json
// @param input body req.RoleFilterReq true "请求结构体"
// @Success 200 {object} resp.RoleFilterResp
// @Failure 400
// @Router /usercenter/api/v1/role/filter [post]
func (role) Filter(c *gin.Context, args struct {
	Body req.RoleFilterReq `param:"body"`
}) (*resp.RoleFilterResp, error) {
	operator, err := GetOperator(c)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return service.ContactRoleFilter(operator.OrgId, args.Body)
}
