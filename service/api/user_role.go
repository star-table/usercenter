package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type userRole int

var UserRole userRole

// @Security token
// @Summary 用户分配角色
// @Description 分配角色接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.AssignRoleRqe true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/userRole/assignRoles [post]
func (userRole) AssignRoles(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var refReq req.AssignRoleRqe
	err = ParseBody(c, &refReq)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.BindUserRoles(operator.OrgId, operator.UserId, refReq.UserId, refReq.RoleIds, permission)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 角色添加用户
// @Description 修改角色组接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.UpdateRoleUserRefRqe true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/userRole/create [post]
func (userRole) Create(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var refRqe req.UpdateRoleUserRefRqe
	err = ParseBody(c, &refRqe)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	if !permission.HasManageRole(refRqe.RoleId) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.BindRoleUsers(operator.OrgId, operator.UserId, refRqe.RoleId, refRqe.UserIds)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 角色删除用户
// @Description 删除角色组接口
// @Tags 角色
// @accept application/json
// @Produce application/json
// @param input body req.UpdateRoleUserRefRqe true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/userRole/delete [delete]
func (userRole) Delete(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var refRqe req.UpdateRoleUserRefRqe
	err = ParseBody(c, &refRqe)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	if !permission.HasManageRole(refRqe.RoleId) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UnbindRoleUsers(operator.OrgId, operator.UserId, refRqe.RoleId, refRqe.UserIds)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}
