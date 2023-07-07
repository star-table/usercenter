package open_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/service/open_service"
)

type memberOpen int

var MemberOpen memberOpen

// @Summary 组织成员信息列表
// @Description 获取组织成员列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} []open_resp.OrgMemberBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/list [post]
func (memberOpen) GetOrgMemberList(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	memberList, err := open_service.GetOrgMemberList(operator.OrgId)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, memberList)

}

// @Summary 组织成员信息列表，根据条件
// @Description 获取组织成员列表，根据条件
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body  open_req.MemberQueryReq true "入参"
// @Success 200 {object} []open_resp.OrgMemberBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/list-by-query-cond [post]
func (memberOpen) GetOrgMemberListByQueryCond(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.MemberQueryReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	memberList, err := open_service.GetOrgMemberListByQueryCond(operator.OrgId, reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, memberList)

}

// @Summary 获取组织成员信息列表
// @Description 根据成员ID列表，获取组织成员信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} []open_resp.OrgMemberBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/list-by-ids [post]
func (memberOpen) GetOrgMemberListByUserIds(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	memberList, err := open_service.GetOrgMemberListByUserIds(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, memberList)

}

// @Summary 获取组织成员信息
// @Description 根据成员ID，获取组织成员信息
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdReq true "入参"
// @Success 200 {object} open_resp.OrgMemberBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/by-id [post]
func (memberOpen) GetOrgMemberByUserId(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	member, err := open_service.GetOrgMemberByUserId(operator.OrgId, reqParam.Id)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, member)
}

// @Summary 获取组织成员信息列表
// @Description 获取组织成员信息列表，排除指定成员ID以后
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} open_resp.OrgMemberBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/exclude-ids [post]
func (memberOpen) GetOrgMemberListByExcludeIds(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	member, err := open_service.GetOrgMemberListByExcludeIds(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, member)
}

// @Summary 获取部门下成员信息列表
// @Description 获取部门下成员信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdReq true "入参"
// @Success 200 {object} []open_resp.UserDeptBindResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/by-dept [post]
func (memberOpen) GetOrgMemberListByDept(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	members, err := open_service.GetOrgMemberListByDept(operator.OrgId, reqParam.Id)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, members)
}

// @Summary 成员绑定的部门信息列表
// @Description 根据成员ID，获取绑定部门信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdReq true "入参"
// @Success 200 {object} []open_resp.UserDeptBindResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/dept-list-by-user [post]
func (memberOpen) GetUserDeptBindListByUser(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetUserDeptBindListByUser(operator.OrgId, reqParam.Id)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 成员绑定的部门信息列表
// @Description 根据成员ID列表，获取绑定部门信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} []open_resp.UserDeptBindResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/dept-list-by-users [post]
func (memberOpen) GetUserDeptBindListByUsers(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetUserDeptBindListByUsers(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 成员绑定的角色信息列表
// @Description 根据成员ID，获取绑定角色信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdReq true "入参"
// @Success 200 {object} []open_resp.UserRoleBindResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/role-list-by-user [post]
func (memberOpen) GetUserRoleBindListByUser(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetUserRoleBindListByUser(operator.OrgId, reqParam.Id)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 成员绑定的角色信息列表
// @Description 根据成员ID列表，获取绑定角色信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} []open_resp.UserRoleBindResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/role-list-by-users [post]
func (memberOpen) GetUserRoleBindListByUsers(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetUserRoleBindListByUsers(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 成员绑定的部门和角色信息列表
// @Description 根据成员ID列表，获取绑定部门和角色信息列表
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} map[int64]open_resp.UserBindDeptAndRoleResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/dept-role-list-by-users [post]
func (memberOpen) GetUserBindDeptAndRoleListByUsers(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetUserDeptRoleBindListByUsers(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)
}

// @Summary 获取平级或者上级的个数
// @Description 获取平级或者上级的个数 根据同一部门下的职级等级对比
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.CompareUser true "入参"
// @Success 200 {object} int 个数
// @Failure 400
// @Router /open/usercenter/api/v1/member/same-superior-count [post]
func (memberOpen) GetSameOrSuperiorCount(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.CompareUser
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	count, err := open_service.GetSameOrSuperiorCount(operator.OrgId, reqParam.UserId, reqParam.SuperiorId)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, count)
}

// @Summary 获取平级或者下级的成员ID列表
// @Description 获取平级或者下级的成员ID列表 根据同一部门下的职级等级对比
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.SameSubordinateUsers true "入参"
// @Success 200 {object} []int64 成员ID列表
// @Failure 400
// @Router /open/usercenter/api/v1/member/get-same-subordinate-users [post]
func (memberOpen) GetSameOrSubordinateMembers(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.SameSubordinateUsers
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	ids, err := open_service.GetSameOrSubordinateMembers(operator.OrgId, reqParam.Id, reqParam.DeptIds)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, ids)

}

// @Summary 人员权限基本信息
// @Description 人员权限基本信息
// @Tags 成员（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdReq true "入参"
// @Success 200 {object} open_resp.OrgUserAuthBaseResp
// @Failure 400
// @Router /open/usercenter/api/v1/member/user-auth-base-info [post]
func (memberOpen) GetUserAuthBaseInfo(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam open_req.IdReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	authBaseInfo, err := open_service.GetUserAuthBaseInfo(operator.OrgId, reqParam.Id)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, authBaseInfo)

}
