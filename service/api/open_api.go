package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/service"
)

type openApi int

var OpenApi openApi

// @Security token
// @Summary 校验ApiKey
// @Description 校验ApiKey接口
// @Tags 组织
// @accept application/x-www-form-urlencoded
// @Produce application/json
// @param apiKey query string true "apiKey"
// @Success 200 {object} resp.CacheUserInfoBo
// @Failure 400
// @Router /usercenter/inner/api/v1/auth/api-key-auth [post]
func (b openApi) ApiKeyAuth(c *gin.Context) {
	apiKey := c.Query("apiKey")
	if apiKey == "" {
		Fail(c, errs.ApiKeyAuthErr)
		return
	}
	operator, err := service.ApiKeyAuth(apiKey)
	if err != nil {
		Fail(c, err)
		return
	}

	res := &resp.CacheUserInfoBo{}
	_ = copyer.Copy(operator, res)
	Suc(c, res)
}
