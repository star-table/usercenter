package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/service"
)

type industry int

var Industry industry = 1

// @Security token
// @Summary 行业列表
// @Description 行业列表接口
// @Tags 行业列表
// @accept application/json
// @Produce application/json
// @Success 200 {object} resp.IndustryListResp
// @Failure 400
// @Router /usercenter/api/v1/user/industryList [get]
func (b industry) IndustryList(c *gin.Context) {
	res, err := service.IndustryList()

	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, res)
}
