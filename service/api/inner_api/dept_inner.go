package inner_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type deptInner int

var DeptInner deptInner

// @Summary 获取部门信息列表（可用）
// @Description 根据部门ID列表获取部门信息列表
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.DeptListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getListByIds [post]
func (deptInner) GetDeptListByIds(c *gin.Context) {
	var reqParam inner_req.DeptListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetDeptListByIds(reqParam.OrgId, reqParam.Ids, consts.AppStatusEnable, consts.AppIsNoDelete)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 获取所有部门信息列表
// @Description 根据部门ID列表获取部门信息列表
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.DeptListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getAllListByIds [post]
func (deptInner) GetAllDeptListByIds(c *gin.Context) {

	var reqParam inner_req.DeptListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := inner_service.GetDeptListByIds(reqParam.OrgId, reqParam.Ids, 0, 0)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)

}

// @Summary 统计部门下的人数
// @Description 统计部门下的人数
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetUserCountByDeptIdsReq true "入参"
// @Success 200 {object} []inner_resp.GetUserCountByDeptIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getUserCountByDeptIds [post]
func (deptInner) GetUserCountByDeptIds(c *gin.Context) {
	var reqParam inner_req.GetUserCountByDeptIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserCountByDeptIds(reqParam.OrgId, reqParam.DeptIds)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取人员所在的部门id集合
// @Description 获取人员所在的部门id集合
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetUserDeptIdsReq true "入参"
// @Success 200 {object} []inner_resp.GetUserDeptIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getUserDeptIds [post]
func (deptInner) GetUserDeptIds(c *gin.Context) {
	var reqParam inner_req.GetUserDeptIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserDeptIds(reqParam.OrgId, reqParam.UserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取人员所在的部门id集合（批量）
// @Description 获取人员所在的部门id集合（批量）
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetUserDeptIdsBatchReq true "入参"
// @Success 200 {object} inner_resp.GetUserDeptIdsBatchResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getUserDeptIdsBatch [post]
func (deptInner) GetUserDeptIdsBatch(c *gin.Context) {
	var reqParam inner_req.GetUserDeptIdsBatchReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserDeptIdsBatch(reqParam.OrgId, reqParam.UserIds)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取部门下的所有人员id
// @Description 获取部门下的所有人员id
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetUserIdsByDeptIdsReq true "入参"
// @Success 200 {object} []inner_resp.GetUserIdsByDeptIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getUserIdsByDeptIds [post]
func (deptInner) GetUserIdsByDeptIds(c *gin.Context) {
	var reqParam inner_req.GetUserIdsByDeptIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserIdsByDeptIds(reqParam.OrgId, reqParam.DeptIds)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取部门下的所有leader
// @Description 获取部门下的所有leader
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetLeadersByDeptIdsReq true "入参"
// @Success 200 {object} inner_resp.GetLeadersByDeptIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getLeadersByDeptIds [post]
func (deptInner) GetLeadersByDeptIds(c *gin.Context) {
	var reqParam inner_req.GetLeadersByDeptIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetLeadersByDeptIds(reqParam.OrgId, reqParam.DeptIds)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 通过部门id获取完整部门路径
// @Description 通过部门id获取完整部门路径
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetFullDeptByIdsReq true "入参"
// @Success 200 {object} inner_resp.GetFullDeptByIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/get-full-dept-by-ids [post]
func (deptInner) GetDeptFullNameByIds(c *gin.Context) {
	var reqParam inner_req.GetFullDeptByIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetFullDeptByIds(reqParam.OrgId, reqParam.DeptIds)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取部门用户（部门对应用户数组）
// @Description 获取部门用户（部门对应用户数组）
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetDeptUserIdsReq true "入参"
// @Success 200 {object} []inner_resp.GetDeptUserIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/getUserIds [post]
func (deptInner) GetDeptUserIds(c *gin.Context) {
	var reqParam inner_req.GetDeptUserIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetDeptUserIds(reqParam.OrgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取部门列表
// @Description 获取部门列表
// @Tags 部门（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetDeptListReq true "入参"
// @Success 200 {object} []inner_resp.DeptInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/dept/get-dept-list [post]
func (deptInner) GetDeptList(c *gin.Context) {
	var reqParam inner_req.GetDeptListReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetDeptList(reqParam.OrgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// 获取人员所在的部门id集合（包含父部门）
func (deptInner) GetUserDeptIdsWithParentId(c *gin.Context) {
	var reqParam inner_req.GetUserDeptIdsWithParentIdReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	resp, err := inner_service.GetUserDeptIdsWithParentId(reqParam.OrgId, reqParam.UserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, resp)
}
