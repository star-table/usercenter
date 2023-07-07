package inner_api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type orgInner int

var OrgInner orgInner

// @Summary 获取企业基础信息
// @Description 获取企业基础信息
// @Tags 企业（内部调用）
// @accept application/json
// @Produce application/json
// @param orgId query int64 true "组织id"
// @Success 200 {object} bo.BaseOrgInfoBo
// @Failure 400
// @Router /usercenter/inner/api/v1/org/info [get]
func (orgInner) GetOrgInfo(c *gin.Context) {
	orgIdStr := c.Query("orgId")
	orgId, err := strconv.ParseInt(orgIdStr, 10, 64)
	if err != nil {
		api.InnerFail(c, errs.ParamError)
		return
	}
	respVO, err1 := inner_service.GetOrgInfo(orgId)
	if err1 != nil {
		api.InnerFail(c, err1)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 新增外部协作人
// @Description 新增外部协作人
// @Tags 企业（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.AddOrgOutCollaborator true "入参"
// @Success 200 {object} bo.BaseOrgInfoBo
// @Failure 400
// @Router /usercenter/inner/api/v1/org/add-out-collaborator [post]
func (orgInner) AddOutCollaborator(c *gin.Context) {
	var reqParam inner_req.AddOrgOutCollaborator
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	err = inner_service.AddOutCollaborator(reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, "")
}

// @Summary 检查组织拥有者是否是超管，如果不是，则设置为超管。
// @Description 检查组织拥有者是否是超管，如果不是，则设置为超管。
// @Tags 企业（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.CheckAndSetSuperAdminReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/inner/api/v1/org/check-and-set-super-admin [post]
func (orgInner) CheckAndSetSuperAdmin(c *gin.Context) {
	var reqParam inner_req.CheckAndSetSuperAdminReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	err = inner_service.CheckAndSetSuperAdmin(reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, true)
}
