package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/vo"
	"github.com/star-table/usercenter/service/model/vo/orgvo"
	"github.com/star-table/usercenter/service/service"
)

type org int

var Org org = 1

// @Security token
// @Summary 创建按组织
// @Description 创建组织接口
// @Tags 组织
// @accept application/json
// @Produce application/json
// @param input body req.CreateOrgReq true "入参"
// @Success 200 {object} resp.CreateOrgRespVoData
// @Failure 400
// @Router /usercenter/api/v1/user/createOrg [post]
func (b org) CreateOrg(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	token, err := getToken(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var orgReq vo.CreateOrgReq
	err = ParseBody(c, &orgReq)
	if err != nil {
		Fail(c, err)
		return
	}
	sourceChannel := ""
	sourcePlatform := ""
	if orgReq.SourceChannel != nil {
		sourceChannel = *orgReq.SourceChannel
	}
	if orgReq.SourcePlatform != nil {
		sourcePlatform = *orgReq.SourcePlatform
	}
	orgId, err := service.CreateOrg(orgvo.CreateOrgReqVo{
		Data: orgvo.CreateOrgReqVoData{
			CreatorId:    operator.UserId,
			CreateOrgReq: orgReq,
			UserToken:    token,
		},
		OrgId:  operator.OrgId,
		UserId: operator.UserId,
	}, sourceChannel, sourcePlatform)

	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, resp.CreateOrgRespVoData{OrgId: orgId})
}

// @Security token
// @Summary 获取组织设置
// @Description 获取组织设置接口
// @Tags 组织
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.OrgConfig
// @Failure 400
// @Router /usercenter/api/v1/user/getOrgConfig [get]
func (b org) GetOrgConfig(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	configReq, err := service.GetOrgConfig(operator.OrgId)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, configReq)
}

// @Security token
// @Summary 组织切换
// @Description 登录的用户进行组织切换
// @Tags 用户
// @accept application/json
// @Produce application/json
// @param input body vo.SwitchUserOrganizationReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/switchUserOrganization [post]
func (b org) SwitchUserOrganization(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var switchReq vo.SwitchUserOrganizationReq
	err = ParseBody(c, &switchReq)
	if err != nil {
		Fail(c, err)
		return
	}
	token := c.GetHeader(consts.AppHeaderTokenName)
	ok, err := service.SwitchUserOrganization(switchReq.OrgID, operator.UserId, token)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 将外部协作人转为内部成员
// @Description 将外部协作人转为内部成员
// @Tags 组织
// @accept application/json
// @Produce application/json
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/org/switch-to-inner-member [post]
func (b org) SwitchToInnerMember(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var switchReq req.SwitchToInnerMemberReq
	err = ParseBody(c, &switchReq)
	if err != nil {
		Fail(c, err)
		return
	}

	permission, err := service.GetOrgUserPerContext(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}
	if permission.CheckOrgSourceChannelIsPolaris() {
		if !permission.HasOpForPolaris(consts.OperationOrgUserModifyStatus) {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	} else {
		if !permission.HasAllPermission() {
			Fail(c, errs.ForbiddenAccess)
			return
		}
	}

	ok, err := service.SwitchToInnerMember(operator.OrgId, operator.UserId, switchReq)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 开启OpenApi
// @Description 开启OpenApi，并生成apiKey
// @Tags 组织
// @accept application/json
// @Produce application/json
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/org/generate-api-key [post]
func (b org) GenerateApiKey(c *gin.Context) {
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
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}
	ok, err := service.GenerateApiKey(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 重置ApiKey
// @Description 重置ApiKey
// @Tags 组织
// @accept application/json
// @Produce application/json
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/org/reset-api-key [post]
func (b org) ResetApiKey(c *gin.Context) {
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
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}
	ok, err := service.ResetApiKey(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 关闭OpenApi
// @Description 关闭OpenApi，并删除apiKey
// @Tags 组织
// @accept application/json
// @Produce application/json
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/org/close-api-key [delete]
func (b org) CloseOpenApi(c *gin.Context) {
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
	if !permission.HasAllPermission() {
		Fail(c, errs.ForbiddenAccess)
		return
	}
	ok, err := service.RemoveApiKey(operator.OrgId, operator.UserId)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}
