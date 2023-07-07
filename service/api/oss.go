package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type oss int

var Oss oss

// @Security token
// @Summary 获取文件上传策略信息
// @Description 获取文件上传策略信息
// @Tags OSS
// @accept application/json
// @Produce application/json
// @param input body req.GetOssPostPolicyReq true "入参"
// @Success 200 {object} resp.GetOssPostPolicyResp
// @Failure 400
// @Router /usercenter/api/v1/oss/getOssPostPolicy [post]
func (oss) GetOssPostPolicy(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var policyReq req.GetOssPostPolicyReq
	err = ParseBody(c, &policyReq)
	if err != nil {
		Fail(c, err)
		return
	}
	res, err := service.GetOssPostPolicy(operator.OrgId, operator.UserId, policyReq)
	Suc(c, res)
}
