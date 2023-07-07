package open_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/service/open_service"
)

type deptOpen int

var DeptOpen deptOpen

// @Summary 获取部门信息列表
// @Description 获取部门信息列表
// @Tags 部门（OpenApi调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} []open_resp.DeptInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/dept/list [post]
func (deptOpen) GetDeptList(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetDeptList(operator.OrgId)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 获取部门信息列表，根据条件
// @Description 获取部门信息列表,根据条件
// @Tags 部门（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.DeptQueryReq true "入参"
// @Success 200 {object} []open_resp.DeptInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/dept/list-by-query-cond [post]
func (deptOpen) GetDeptListByQueryCond(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}

	var reqParam open_req.DeptQueryReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetDeptListByQueryCond(operator.OrgId, reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 根据部门ID列表获取部门信息列表
// @Description 根据部门ID列表获取部门信息列表
// @Tags 部门（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} []open_resp.DeptInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/dept/list-by-ids [post]
func (deptOpen) GetDeptListByIds(c *gin.Context) {
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
	respVO, err := open_service.GetDeptListByIds(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 根据上级ID列表获取部门信息列表
// @Description 根据上级ID列表获取部门信息列表
// @Tags 部门（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body open_req.IdsReq true "入参"
// @Success 200 {object} []open_resp.DeptInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/dept/children [post]
func (deptOpen) GetDeptChildrenList(c *gin.Context) {
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
	respVO, err := open_service.GetDeptChildrenList(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 获取部门信息列表,包含成员
// @Description 获取部门信息列表,包含成员
// @Tags 部门（OpenApi调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} []open_resp.DeptMemberListResp
// @Failure 400
// @Router /open/usercenter/api/v1/dept/list-have-member [post]
func (deptOpen) GetDeptHaveMemberList(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}

	members, err := open_service.GetDeptHaveMemberList(operator.OrgId)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, members)
}
