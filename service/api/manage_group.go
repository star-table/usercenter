package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type manageGroup int

var ManageGroup manageGroup

// @Security token
// @Summary 管理组新增
// @Description 管理组新增接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param .+ body req.CreateManageGroup true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/create [post]
func (manageGroup) CreateGroup(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var input req.CreateManageGroup
	err = ParseBody(c, &input)
	if err != nil {
		Fail(c, err)
		return
	}
	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupCreate) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	id, err := service.CreateManageGroup(operator.OrgId, operator.UserId, input)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}

// @Security token
// @Summary 管理组修改
// @Description 管理组新增接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param input body req.UpdateManageGroup true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/update [put]
func (a manageGroup) UpdateGroup(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var input req.UpdateManageGroup
	err = ParseBody(c, &input)
	if err != nil {
		Fail(c, err)
		return
	}
	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupModify) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdateManageGroup(operator.OrgId, operator.UserId, input.Id, input)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)

}

// @Security token
// @Summary 管理组删除
// @Description 管理组新增接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param input body req.DeleteManageGroup true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/delete [delete]
func (a manageGroup) DeleteGroup(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var input req.DeleteManageGroup
	err = ParseBody(c, &input)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupDelete) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.DeleteManageGroup(operator.OrgId, operator.UserId, input.Id)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 修改管理组包含内容
// @Description 修改管理组包含内容接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param input body req.UpdateManageGroupContents true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/updateContents [put]
func (a manageGroup) UpdateContents(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var input req.UpdateManageGroupContents
	err = ParseBody(c, &input)
	if err != nil {
		Fail(c, err)
		return
	}
	if input.Key == "usage_ids" {
		input.Key = "opt_auth"
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasAllPermission() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupModify) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdateManageGroupContents(operator.OrgId, operator.UserId, permission.IsOrgOwner(), input.Id, input)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 管理组树状
// @Description 管理组树状接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.ManageGroupTreeResp
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/tree [get]
func (a manageGroup) GetManageGroupTree(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	//permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	//if err != nil {
	//	Fail(c, err)
	//	return
	//}
	//
	//if !permission.HasInManageGroup() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupView) {
	//	// 不报错，返回空对象
	//	Suc(c, &resp.ManageGroupTreeResp{
	//		SysGroup:      &resp.SimpleManageGroupInfo{},
	//		GeneralGroups: make([]resp.SimpleManageGroupInfo, 0),
	//	})
	//	return
	//}

	tree, err := service.GetManageGroupTree(operator.OrgId)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, tree)
}

// @Security token
// @Summary 管理详情
// @Description 管理详情接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @param id query int64 true "管理组ID"
// @Success 200 {object} resp.ManageGroupDetailResp
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/detail [get]
func (a manageGroup) GetManageGroupDetail(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	groupId := int64(0)
	if idStr, ok := c.GetQuery("id"); ok {
		groupId, _ = strconv.ParseInt(idStr, 10, 64)
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasInManageGroup() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupView) {
		Fail(c, errs.ForbiddenView)
		return
	}

	tree, err := service.GetManageGroupDetail(operator.OrgId, operator.UserId, groupId)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, tree)
}

// @Security token
// @Summary 管理组权限项配置信息
// @Description 管理组权限项配置信息接口
// @Tags 管理员
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.GetOperationConfigResp
// @Failure 400
// @Router /usercenter/api/v1/adminGroup/operationConfigs [get]
func (a manageGroup) GetManageGroupOperationConfig(c *gin.Context) {
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

	if !permission.HasInManageGroup() && !permission.HasOpForPolaris(consts.OperationOrgAdminGroupView) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	resp, err := service.GetManageGroupOperationConfig(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, resp)
}
