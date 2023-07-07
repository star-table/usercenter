package inner_api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type ossInner int

var OssInner ossInner

// @Security token
// @Summary 获取文件上传策略信息
// @Description 获取文件上传策略信息
// @Tags OSS（内部调用）
// @accept application/json
// @Produce application/json
// @param input body orgId true "入参"
// @Success 200 {object} resp.GetOssPostPolicyResp
// @Failure 400
// @Router /usercenter/inner/api/v1/oss/getOssPostPolicy [post]
func (ossInner) GetOssPostPolicy(c *gin.Context) {
	orgIdStr := c.Query("orgId")
	orgId, err := strconv.ParseInt(orgIdStr, 10, 64)
	if err != nil {
		api.InnerFail(c, errs.ParamError)
		return
	}
	defaultVal := int64(0)
	respVO, err1 := inner_service.GetOssPostPolicy(orgId, req.GetOssPostPolicyReq{
		PolicyType: 9, // 9:备忘录
		ProjectID:  &defaultVal,
		IssueID:    &defaultVal,
		FolderID:   &defaultVal,
	})
	if err1 != nil {
		api.InnerFail(c, err1)
		return
	}
	api.InnerSuc(c, respVO)
}
