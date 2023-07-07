package open_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service/open_service"
)

type roleOpen int

var RoleOpen roleOpen

// @Summary 获取角色信息列表
// @Description 获取角色信息列表
// @Tags 角色（OpenApi调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} []open_resp.RoleInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/role/list [post]
func (roleOpen) GetRoleList(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetRoleList(operator.OrgId)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 根据角色ID列表获取角色信息列表
// @Description 根据角色ID列表获取角色信息列表
// @Tags 角色（OpenApi调用）
// @accept application/json
// @Produce application/json
// @param input body req.IdsReq true "入参"
// @Success 200 {object} []open_resp.RoleInfoResp
// @Failure 400
// @Router /open/usercenter/api/v1/role/list-by-ids [post]
func (roleOpen) GetRoleListByIds(c *gin.Context) {
	operator, err := api.GetOperator(c)
	if err != nil {
		api.Fail(c, err)
		return
	}
	var reqParam req.IdsReq
	err = api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := open_service.GetRoleListByIds(operator.OrgId, reqParam.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}
