package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type department int

var Department department = 1

// @Security token
// @Summary 部门列表
// @Description 部门列表接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.DepartmentListReq true "入参"
// @Success 200 {object} resp.DepartmentList
// @Failure 400
// @Router /usercenter/api/v1/user/departmentList [post]
func (b department) DepartmentList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.DepartmentListReq
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

	res, err := service.DepartmentList(operator.OrgId, reqParam, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 创建部门
// @Description 创建部门接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.CreateDepartmentReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/user/createDepartment [post]
func (b department) CreateDepartment(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.CreateDepartmentReq
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

	if !permission.HasManageDept(reqParam.ParentID) && !permission.HasOpForPolaris(consts.OperationOrgDeptCreate) {
		Fail(c, errs.PolarisForbiddenAccess)
		return
	}

	id, err := service.CreateDepartment(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, id)
}

// @Security token
// @Summary 更新部门
// @Description 更新部门接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.UpdateDepartmentReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/updateDepartment [post]
func (b department) UpdateDepartment(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.UpdateDepartmentReq
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
	if !permission.HasManageDept(reqParam.DepartmentId) && !permission.HasOpForPolaris(consts.OperationOrgDeptModify) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdateDepartment(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 删除部门
// @Description 删除部门接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.DeleteDepartmentReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/deleteDepartment [post]
func (b department) DeleteDepartment(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.DeleteDepartmentReq

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
	if !permission.HasManageDept(reqParam.DepartmentId) && !permission.HasOpForPolaris(consts.OperationOrgDeptDelete) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.DeleteDepartment(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 变更用户部门/职级
// @Description 变更用户部门/职级接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.ChangeUserDeptAndPositionReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/change-dept-and-position [post]
func (b department) ChangeUserDeptAndPosition(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.ChangeUserDeptAndPositionReq
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

	ok, err := service.ChangeUserDeptAndPosition(operator.OrgId, operator.UserId, reqParam, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 给用户分配部门
// @Description 给用户分配部门
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.AllocateUserDeptReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/allocate-dept [post]
func (b department) AllocateUserDept(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.AllocateUserDeptReq
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
	// 当前用户是否具有权限项
	if !permission.HasOpForPolaris(consts.OperationOrgUserModifyUserDept) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.AllocateUserDept(operator.OrgId, operator.UserId, reqParam, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 切换用户管理组
// @Description 切换用户管理组接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param input body req.ChangeUserAdminGroupReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/change-user-admin-group [post]
func (b department) ChangeUserAdminGroup(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.ChangeUserAdminGroupReq
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
	if reqParam.SourceFrom == "polaris" {
		if !permission.HasOpForPolaris(consts.OperationOrgUserModifyUserAdminGroup) {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	} else {
		if !permission.HasAllPermission() {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	}

	if reqParam.DstAdminGroupIds != nil && len(reqParam.DstAdminGroupIds) > 0 {
		reqParam.DstAdminGroupIds = slice.SliceUniqueInt64(reqParam.DstAdminGroupIds)
		_, err := service.ChangeUserAdminGroupBatch(operator.OrgId, operator.UserId, reqParam, permission)
		if err != nil {
			Fail(c, err)
			return
		}
	} else {
		Fail(c, errs.UserMustHasAdminGroup)
		return
	}

	Suc(c, true)
}

// @Security token
// @Summary 设置部门负责人
// @Description 设置部门负责人接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.UpdateDeptLeaderReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/setUserDepartmentLevel [post]
func (b department) UpdateDeptLeader(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.UpdateDeptLeaderReq
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
	if !permission.HasManageDept(reqParam.DepartmentId) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdateDeptLeader(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 改变部门排序
// @Description 改变部门排序接口
// @Tags 部门
// @accept application/json
// @Produce application/json
// @param input body req.ChangeDeptSortReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/dept/change-dept-sort [post]
func (b department) ChangeDeptSort(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.ChangeDeptSortReq
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
	// 只有超管可以修改
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.ChangeDeptSort(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// DeptRemoveUsers 将多个用户移出某个部门接口
// @Security token
// @Summary
// @Description 将多个用户移出某个部门接口
// @Tags 用户
// @accept application/json
// @Produce application/json
// @param input body req.DeptRemoveUsersReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/dept/remove-users [post]
func (b department) DeptRemoveUsers(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.DeptRemoveUsersReq
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
	// 管理员或者拥有”修改用户部门“的权限才可以修改
	if !permission.HasOpForPolaris(consts.OperationOrgUserModifyUserDept) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.DeptRemoveUsers(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}
