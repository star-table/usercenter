package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

/**
日志接口
*/
type records int

var Records records

// @Security token
// @Summary 登陆日志
// @Description 登陆日志接口
// @Tags 日志
// @accept application/json
// @Produce application/json
// @param input body req.LoginRecordListReq true "入参"
// @Success 200 {object} resp.LoginRecordListResp
// @Failure 400
// @Router /usercenter/api/v1/records/login-records [post]
func (records) GetLoginRecordsList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.LoginRecordListReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	id, err := service.GetLoginRecordsList(operator.OrgId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}

// @Security token
// @Summary 导出通讯录
// @Description 导出通讯录接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.ExportLoginRecordListReq true "入参"
// @Success 200 {object} resp.ExportAddressListResp
// @Failure 400
// @Router /usercenter/api/v1/records/export-login-records [post]
func (records) ExportLoginRecordsList(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.ExportLoginRecordListReq
	err = ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	id, err := service.ExportLoginRecordsList(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, id)
}
