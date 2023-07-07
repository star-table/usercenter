package inner_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type roleInner int

var RoleInner roleInner

// @Summary 获取角色信息列表
// @Description 根据角色ID列表获取角色信息列表
// @Tags 角色（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.RoleListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.RoleInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/role/getListByIds [post]
func (roleInner) GetRoleListByIds(c *gin.Context) {

	var reqParam inner_req.RoleListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetRoleListByIds(reqParam.OrgId, reqParam.Ids, consts.AppStatusEnable, consts.AppIsNoDelete)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)

}

// @Summary 获取角色信息列表 （不分状态和是否删除）
// @Description 根据角色ID列表获取角色信息列表（不分状态和是否删除）
// @Tags 角色（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.RoleListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.RoleInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/role/getAllListByIds [post]
func (roleInner) GetAllRoleListByIds(c *gin.Context) {

	var reqParam inner_req.RoleListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetRoleListByIds(reqParam.OrgId, reqParam.Ids, 0, 0)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)

}

// @Summary 获取角色用户（角色对应用户数组）
// @Description 获取角色用户（角色对应用户数组）
// @Tags 角色（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetRoleUserIdsReq true "入参"
// @Success 200 {object} []inner_resp.GetRoleUserIdsResp
// @Failure 400
// @Router /usercenter/inner/api/v1/role/getUserIds [post]
func (roleInner) GetRoleUserIds(c *gin.Context) {
	var reqParam inner_req.GetRoleUserIdsReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetRoleUserIds(reqParam.OrgId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}
