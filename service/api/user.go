package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type user int

var User user = 1

// @Security token
// @Summary 新建成员
// @Description 新建成员接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.CreateOrgMemberReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /usercenter/api/v1/user/create [post]
func (b user) CreateOrgMember(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.CreateOrgMemberReq
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
	if !permission.HasInManageGroup() {
		Fail(c, errs.ForbiddenAccess)
		return
	}
	if len(reqParam.RoleIds) > 0 {
		if !permission.HasManageRoles(reqParam.RoleIds) {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	}
	if len(reqParam.DeptAndPositions) > 0 {
		deptIds := make([]int64, 0)
		for _, v := range reqParam.DeptAndPositions {
			deptIds = append(deptIds, v.DepartmentId)
		}
		if !permission.HasManageDepts(deptIds) {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	}

	id, err := service.CreateOrgMember(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, id)
}

// @Security token
// @Summary 更新成员
// @Description 更新成员接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.UpdateOrgMemberReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/update [post]
func (b user) UpdateOrgMemberInfo(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.UpdateOrgMemberReq
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
	// 自己可以更改信息，有权限的管理员可以更改
	if !permission.HasInManageGroup() && !(reqParam.UserId == operator.UserId || permission.HasOpForPolaris(consts.OperationOrgUserModifyStatus)) {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.UpdateOrgMemberInfo(operator.OrgId, operator.UserId, reqParam, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 成员列表
// @Description 成员列表接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.UserListReq true "入参"
// @Success 200 {object} resp.UserListResp
// @Failure 400
// @Router /usercenter/api/v1/user/userList [post]
func (b user) GetOrgMemberList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.UserListReq
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

	res, err := service.GetOrgMemberList(operator.OrgId, reqParam, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 根据ID获取成员详情
// @Description 根据ID获取成员详情接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.IdReq true "入参"
// @Success 200 {object} resp.OrgMemberInfoReq
// @Failure 400
// @Router /usercenter/api/v1/user/member-info [post]
func (b user) GetOrgMemberInfoById(c *gin.Context) {
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

	res, err := service.GetOrgMemberInfoById(operator.OrgId, reqParam.Id)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 邀请时搜索成员
// @Description 邀请时搜索成员接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @param input body req.SearchUserReq true "入参"
// @Success 200 {object} resp.SearchUserResp
// @Failure 400
// @Router /usercenter/api/v1/user/searchUser [post]
func (b user) SearchUser(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var userReq req.SearchUserReq
	err = ParseBody(c, &userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.SearchUser(operator.OrgId, userReq.Email)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 邀请成员/再次邀请
// @Description 邀请成员(再次邀请)接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @param input body req.InviteUserReq true "入参"
// @Success 200 {object} resp.InviteUserResp
// @Failure 400
// @Router /usercenter/api/v1/user/inviteUser [post]
func (b user) InviteUser(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var inviteReq req.InviteUserReq
	err = ParseBody(c, &inviteReq)
	if err != nil {
		Fail(c, err)
		return
	}

	inviteResp, err := service.InviteUser(operator.OrgId, operator.UserId, inviteReq)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, inviteResp)
}

// @Security token
// @Summary 未接受邀请成员列表
// @Description 未接受邀请成员列表接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @param input body req.InviteUserListReq true "入参"
// @Success 200 {object} resp.InviteUserListResp
// @Failure 400
// @Router /usercenter/api/v1/user/inviteUserList [post]
func (b user) InviteUserList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var userReq req.InviteUserListReq
	err = ParseBody(c, &userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.InviteUserList(operator.OrgId, userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 删除成员邀请列表
// @Description 删除成员邀请列表接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @param input body req.RemoveInviteUserReq true "入参"
// @Success 200
// @Failure 400
// @Router /usercenter/api/v1/user/removeInviteUser [post]
func (b user) RemoveInviteUser(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var userReq req.RemoveInviteUserReq
	err = ParseBody(c, &userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	err = service.RemoveInviteMember(operator.OrgId, operator.UserId, userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, nil)
}

// @Security token
// @Summary 导出通讯录
// @Description 导出通讯录接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.ExportAddressListReq true "入参"
// @Success 200 {object} resp.ExportAddressListResp
// @Failure 400
// @Router /usercenter/api/v1/user/exportAddressList [post]
func (b user) ExportAddressList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var exportReq req.ExportAddressListReq
	err = ParseBody(c, &exportReq)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasInManageGroup() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	res, err := service.ExportOrgMemberList(operator.OrgId, operator.UserId, exportReq, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 获取邀请码
// @Description 获取邀请码接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.GetInviteCodeResp
// @Failure 400
// @Router /usercenter/api/v1/user/getInviteCode [get]
func (b user) GetInviteCode(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.GetInviteCode(operator.OrgId, operator.UserId, "")
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 成员tab统计
// @Description 成员tab统计接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.UserStatResp
// @Failure 400
// @Router /usercenter/api/v1/user/userStat [post]
func (b user) OegMemberStat(c *gin.Context) {
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
	res, err := service.OrgMemberStat(operator.OrgId, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 用户组织列表
// @Description 用户组织列表接口
// @Tags 用户
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.UserOrganizationListResp
// @Failure 400
// @Router /usercenter/api/v1/user/userOrgList [post]
func (b user) UserOrgList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.GetUserParticipateOrgList(operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 移除成员
// @Description 移除成员接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.RemoveOrgMemberReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/removeUser [post]
func (b user) RemoveOrgMember(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	param := &req.RemoveOrgMemberReq{}
	err = ParseBody(c, &param)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasInManageGroup() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.RemoveOrgMember(operator.OrgId, operator.UserId, *param, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 获取当前成员信息
// @Description 获取当前成员信息接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.PersonalInfo
// @Failure 400
// @Router /usercenter/api/v1/user/personalInfo [post]
func (b user) PersonalInfo(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.PersonalInfo(operator.OrgId, operator.UserId, operator.SourceChannel)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 更新当前成员信息
// @Description 更新当前成员信息接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.UpdateUserInfoReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/updatePersonalInfo [post]
func (b user) UpdatePersonalInfo(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	var userReq req.UpdateUserInfoReq
	err = ParseBody(c, &userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.UpdateCurrentUserInfo(operator.OrgId, operator.UserId, userReq)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 启用/禁用成员
// @Description 启用/禁用成员接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.UpdateOrgMemberStatusReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/updateUserStatus [post]
func (b user) ChangeOrgMemberStatus(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	param := &req.UpdateOrgMemberStatusReq{}
	err = ParseBody(c, &param)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if !permission.HasInManageGroup() {
		Fail(c, errs.ForbiddenAccess)
		return
	}

	ok, err := service.ChangeOrgMemberStatus(operator.OrgId, operator.UserId, *param, permission)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 获取邀请信息
// @Description 获取邀请信息接口
// @Tags 成员邀请
// @accept application/json
// @Produce application/json
// @param input body req.GetInviteInfoReq true "入参"
// @Success 200 {ParseBody} resp.GetInviteInfoResp
// @Failure 400
// @Router /usercenter/api/v1/user/getInviteInfo [post]
func (b user) GetInviteInfo(c *gin.Context) {
	var reqParam req.GetInviteInfoReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.GetInviteInfo(reqParam.InviteCode)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, res)
}

// @Security token
// @Summary 成员管理组权限
// @Description 成员管理组权限
// @Tags 成员权限
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.UserManageAuthResp
// @Failure 400
// @Router /usercenter/api/v1/user/get-user-manage-auth [get]
func (b user) GetUserManageAuth(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	manageAuth, resErr := service.GetUserManageAuth(operator.OrgId, operator.UserId)
	if resErr != nil {
		Fail(c, resErr)
		return
	}

	Suc(c, manageAuth)
}
