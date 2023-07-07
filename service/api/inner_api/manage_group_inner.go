package inner_api

import (
	"strconv"

	"github.com/star-table/usercenter/core/errs"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"github.com/star-table/usercenter/service/service"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type manageGroupInner int

var ManageGroupInner manageGroupInner

// @Summary 新增应用包至人员所在管理组
// @Description 新增应用包至人员所在管理组
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.AddPkgReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/add-pkg [post]
func (manageGroupInner) AddAppPkgToManageGroup(c *gin.Context) {
	var reqParam inner_req.AddPkgReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.AddAppPkgToManageGroup(reqParam.OrgId, reqParam.UserId, reqParam.PkgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 从人员所在管理组删除应用包
// @Description 从人员所在管理组删除应用包
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.DeletePkgReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/delete-pkg [post]
func (manageGroupInner) DeleteAppPkgFromManageGroup(c *gin.Context) {
	var reqParam inner_req.DeletePkgReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.DeleteAppPkgFromManageGroup(reqParam.OrgId, reqParam.UserId, reqParam.PkgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 新增应用至人员所在管理组
// @Description 新增应用至改人员所在管理组
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.AddAppReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/add-app [post]
func (manageGroupInner) AddAppToManageGroup(c *gin.Context) {
	var reqParam inner_req.AddAppReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.AddAppToManageGroup(reqParam.OrgId, reqParam.UserId, reqParam.AppId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 从人员所在管理组删除应用
// @Description 从人员所在管理组删除应用
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.DeleteAppReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/delete-app [post]
func (manageGroupInner) DeleteAppFromManageGroup(c *gin.Context) {
	var reqParam inner_req.DeleteAppReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.DeleteAppFromManageGroup(reqParam.OrgId, reqParam.UserId, reqParam.AppId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取管理组管理的应用集合
// @Description 获取管理组管理的应用集合
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetManagerReq true "入参"
// @Success 200 {object} []inner_resp.GetManagerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/getManageInfo [post]
func (manageGroupInner) GetManager(c *gin.Context) {
	var reqParam inner_req.GetManagerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	respVo, err := inner_service.GetManager(reqParam.OrgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVo)
}

// GetManageGroupTree
// @Summary 管理组树状
// @Description 管理组树状接口，内部调用
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.ManageGroupTreeResp
// @Failure 400
// @Router /usercenter/inner/api/v1/adminGroup/tree [post]
func (manageGroupInner) GetManageGroupTree(c *gin.Context) {
	var reqParam inner_req.GetManageGroupTreeReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	tree, err := service.GetManageGroupTree(reqParam.OrgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, tree)
}

// UpdateContents
// @Summary 修改管理组包含内容
// @Description 修改管理组包含内容接口，内部调用
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body req.UpdateManageGroupContents true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/adminGroup/updateContents [put]
func (manageGroupInner) UpdateContents(c *gin.Context) {
	var input inner_req.UpdateManageGroupContents
	err := api.ParseBody(c, &input)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	if input.Key == "usage_ids" {
		input.Key = "opt_auth"
	}

	ok, err := service.UpdateManageGroupContents(input.OrgId, input.OperatorUserId, true, input.Id, req.UpdateManageGroupContents{
		Id:     input.Id,
		Values: input.Values,
		Key:    input.Key,
	})
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, ok)
}

// GetManageGroupDetail
// @Summary 管理组详情
// @Description 管理组详情接口，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param id query int64 true "管理组ID"
// @Success 200 {object} resp.ManageGroupDetailResp
// @Failure 400
// @Router /usercenter/inner/api/v1/adminGroup/detail [get]
func (manageGroupInner) GetManageGroupDetail(c *gin.Context) {
	groupId := int64(0)
	orgId := int64(0)
	userId := int64(0)
	if idStr, ok := c.GetQuery("id"); ok {
		groupId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("userId"); ok {
		userId, _ = strconv.ParseInt(idStr, 10, 64)
	}

	tree, err := service.GetManageGroupDetail(orgId, userId, groupId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, tree)
}

// ManageGroupInit
// @Summary 初始化管理组
// @Description 初始化管理组，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @param userId query int64 true "用户id"
// @param sourceFrom query string false "应用来源。默认可以不传，如果是极星，则传 `polaris`。"
// @Success 200 {object} inner_resp.ManagerGroupInitResp
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/init [post]
func (manageGroupInner) ManageGroupInit(c *gin.Context) {
	orgId := int64(0)
	userId := int64(0)
	sourceFrom := ""
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("userId"); ok {
		userId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if sourceStr, ok := c.GetQuery("sourceFrom"); ok {
		sourceFrom = sourceStr
	}

	id, err := inner_service.ManageGroupInit(orgId, userId, sourceFrom)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, inner_resp.ManagerGroupInitResp{
		SysGroupID: id,
	})
}

// AddUserToSysManageGroup
// @Summary 增加人员到系统管理组
// @Description 增加人员到系统管理组，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @param userId query int64 true "用户id"
// @param input body inner_req.AddUserToSysManageGroupReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/add-user-to-sys-manage-group [post]
func (manageGroupInner) AddUserToSysManageGroup(c *gin.Context) {
	orgId := int64(0)
	userId := int64(0)
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("userId"); ok {
		userId, _ = strconv.ParseInt(idStr, 10, 64)
	}

	var input inner_req.AddUserToSysManageGroupReq
	err := api.ParseBody(c, &input)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	addErr := inner_service.AddUserToSysManageGroup(orgId, userId, input.UserIds)
	if addErr != nil {
		api.InnerFail(c, addErr)
		return
	}
	api.InnerSuc(c, true)
}

// DeleteOneUserFromOrg 将一个用户从该组织的所有管理组中移除
// @Summary 将一个用户从该组织的所有管理组中移除
// @Description 将一个用户从该组织的所有管理组中移除，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @param targetUserId query int64 true "要移除的用户id"
// @param operateUid query int64 true "操作人id，没有则传 0"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/adminGroup/deleteOneUserFromOrg [post]
func (manageGroupInner) DeleteOneUserFromOrg(c *gin.Context) {
	orgId := int64(0)
	operateUid := int64(0)
	targetUserId := int64(0)
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("operateUid"); ok {
		operateUid, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("targetUserId"); ok {
		targetUserId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	err := inner_service.DeleteOneUserFromOrg(orgId, operateUid, targetUserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, true)
}

// ReplaceSuperAdmin 更换组织超管
// @Summary 更换组织超管
// @Description 更换组织超管，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @param targetUserId query int64 true "要移除的用户id"
// @param operateUid query int64 true "操作人id，没有则传 0"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/replace-super-admin [post]
func (manageGroupInner) ReplaceSuperAdmin(c *gin.Context) {
	orgId := int64(0)
	operateUid := int64(0)
	targetUserId := int64(0)
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("operateUid"); ok {
		operateUid, _ = strconv.ParseInt(idStr, 10, 64)
	}
	if idStr, ok := c.GetQuery("targetUserId"); ok {
		targetUserId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	err := inner_service.ReplaceSuperAdmin(orgId, operateUid, targetUserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, true)
}

// GetOrgSuperAdminIds 获取组织的超管（id 列表）
// @Summary 获取组织的超管（id 列表）
// @Description 获取组织的超管（id 列表），内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @Success 200 {object} []int64
// @Failure 400
// @Router /usercenter/inner/api/v1/manage-group/get-super-admin-ids [get]
func (manageGroupInner) GetOrgSuperAdminIds(c *gin.Context) {
	orgId := int64(0)
	if idStr, ok := c.GetQuery("orgId"); ok {
		orgId, _ = strconv.ParseInt(idStr, 10, 64)
	}
	adminUidArr, err := inner_service.GetSuperAdminIds(orgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, adminUidArr)
}

// AddNewMenuToRole 向组织角色中追加新的权限项
// @Summary 向组织角色中追加新的权限项
// @Description 向组织角色中追加新的权限项，内部调用！
// @Tags 管理组（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.AddNewMenuToRoleReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/adminGroup/addNewMenuToRole [post]
func (manageGroupInner) AddNewMenuToRole(c *gin.Context) {
	var input inner_req.AddNewMenuToRoleReq
	err := api.ParseBody(c, &input)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	businessErr := inner_service.AddNewMenuToRole(&input)
	if businessErr != nil {
		api.InnerFail(c, businessErr)
		return
	}

	api.InnerSuc(c, true)
}

// 获取普通管理员可以管理的app
func (manageGroupInner) GetCommAdminMangeApps(c *gin.Context) {
	var reqParam inner_req.GetCommAdminMangeAppsReq
	if err := c.BindJSON(&reqParam); err != nil {
		api.InnerFail(c, errs.ReqParamsValidateError)
		return
	}

	resp, resErr := inner_service.GetCommAdminMangeApps(reqParam.OrgId)
	if resErr != nil {
		api.InnerFail(c, resErr)
		return
	}
	api.InnerSuc(c, resp)
}
