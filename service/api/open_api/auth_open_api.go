package open_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/service/open_service"
)

type authOpen int

var AuthOpen authOpen

// @Security token
// @Summary 校验token并验证状态
// @Description 校验token并验证状态接口
// @Tags 认证（open）
// @accept application/json
// @Produce application/json
// @param token body open_req.CheckTokenReq true "token"
// @Success 200 {object} resp.CacheUserInfoBo
// @Failure 400
// @Router /open/usercenter/api/v1/auth/check [post]
func (authOpen) AuthCheckStatus(c *gin.Context) {
	var reqParam open_req.CheckTokenReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, errs.TokenAuthError)
		return
	}
	if reqParam.Token == "" {
		api.Fail(c, errs.TokenAuthError)
		return
	}
	tokenUser, err := open_service.AuthToken(reqParam.Token)
	if err != nil {
		api.Fail(c, err)
		return
	}

	res := &resp.CacheUserInfoBo{}
	_ = copyer.Copy(tokenUser, res)
	api.Suc(c, res)
}
