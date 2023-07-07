package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/vo"
	"github.com/star-table/usercenter/service/service"
)

type sendAuthCode int

var SendAuthCode sendAuthCode = 1

// @Security token
// @Summary 发送各种验证码(目前支持手机号以及邮箱)
// @Description 发送各种验证码(目前支持手机号以及邮箱)接口
// @Tags 验证码
// @accept application/json
// @Produce application/json
// @param input body req.SendAuthCodeReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/sendAuthCode [post]
func (b sendAuthCode) SendAuthCode(c *gin.Context) {
	var reqParam vo.SendAuthCodeReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	if !domain.CheckAuthTypeValid(reqParam.AuthType) {
		err = service.VerifyCaptcha(reqParam.CaptchaID, reqParam.CaptchaPassword, reqParam.Address)
		if err != nil {
			Fail(c, err)
			return
		}
	}

	ok, err := service.SendAuthCode(reqParam)

	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary **校验**手机短信验证码
// @Description 校验手机短信验证码接口
// @Tags 验证码
// @accept application/json
// @Produce application/json
// @param input body req.AuthSmsCodeReq true "入参"
// @Success 200 {object} resp.AuthSmsCodeResp
// @Failure 400
// @Router /usercenter/api/v1/user/auth-sms-code [post]
func (b sendAuthCode) AuthSmsCode(c *gin.Context) {
	var reqParam req.AuthSmsCodeReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	if !domain.CheckAuthTypeValid(reqParam.AuthType) {
		Fail(c, errs.SMSLoginCodeInvalid)
		return
	}
	resp, err := service.AuthSmsCode(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, resp)
}
