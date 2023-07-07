package inner_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type universalInner int

var UniversalInner universalInner

func (universalInner) GetUserBaseInfoByIds(c *gin.Context) {
	var req inner_req.GetBaseInfoByIdsReq
	err := api.ParseBody(c, &req)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.UniversalGetUserBaseInfoByIds(req.OrgId, req.Ids)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)
}
