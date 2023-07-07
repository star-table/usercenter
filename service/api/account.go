package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/service"
)

type account int

var Account account = 1

// @Security token
// @Summary 注册
// @Description 注册接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.UserRegisterReq true "注册信息"
// @Success 200 {object} string
// @Failure 400
// @Router /usercenter/api/v1/user/register [post]
func (b account) Register(c *gin.Context) {
	var reqParam req.UserRegisterReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	res, err := service.RegisterUser(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, res.Token)
}

// @Security token
// @Summary 设置本人密码
// @Description 设置本人密码接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param password body req.PasswordReq true "密码"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/setPassword [post]
func (b account) SetPassword(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var password req.PasswordReq
	err = ParseBody(c, &password)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.SetPassword(req.SetPasswordReq{
		UserId: operator.UserId,
		Input:  password,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 重置密码
// @Description 重置密码接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param password body req.SetPasswordReq true "密码"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/reset-password [post]
func (b account) ResetPassword(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	var reqParam req.SetPasswordReq
	err = ParseBody(c, &reqParam)
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

	ok, err := service.ResetPassword(operator.OrgId, operator.UserId, reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 找回密码
// @Description 找回密码接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.RetrievePasswordReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/retrievePassword [post]
func (b account) RetrievePassword(c *gin.Context) {
	var reqParam req.RetrievePasswordReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.RetrievePassword(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 修改本人密码
// @Description 修改本人密码接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.OldNewPasswordReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/update-password [post]
func (b account) UpdatePassword(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	var passwordReq req.OldNewPasswordReq
	err = ParseBody(c, &passwordReq)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.UpdatePassword(req.UpdatePasswordReq{
		OrgId:  operator.OrgId,
		UserId: operator.UserId,
		Input:  passwordReq,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Summary 密码过期修改密码(无需认证)
// @Description 密码过期修改密码，进入该接口的必定是设置过密码的
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.UpdatePasswordByLoginNameReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/update-password-by-username [post]
func (b account) UpdatePasswordByUsername(c *gin.Context) {
	var reqParam req.UpdatePasswordByLoginNameReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.UpdatePasswordByUsername(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}

// @Security token
// @Summary 验证原邮箱/手机号
// @Description 验证原邮箱/手机号接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.UnbindLoginNameReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/verifyOldName [post]
func (b account) VerifyOldName(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	param := &req.UnbindLoginNameReq{}
	err = ParseBody(c, &param)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.VerifyOldPhoneOrEmail(operator.OrgId, operator.UserId, *param)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 绑定邮箱/手机号
// @Description 绑定邮箱/手机号接口
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.BindLoginNameReq true "入参"
// @Success 200
// @Failure 400
// @Router /usercenter/api/v1/user/bindLoginName [post]
func (b account) BindLoginName(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}

	param := &req.BindLoginNameReq{}
	err = ParseBody(c, &param)
	if err != nil {
		Fail(c, err)
		return
	}

	resErr := service.BindLoginName(operator.OrgId, operator.UserId, *param)
	if resErr != nil {
		Fail(c, resErr)
		return
	}

	Suc(c, nil)
}

// @Security token
// @Summary 解绑登录名
// @Description 解绑登录名，如解绑手机号、邮箱等
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.UnbindLoginNameReq true "入参"
// @Success 200 {object} resp.UnbindLoginNameResp
// @Failure 400
// @Router /usercenter/api/v1/user/unbindLoginName [post]
func (b account) UnbindLoginName(c *gin.Context) {
	operator, err := GetOperator(c)
	if err != nil {
		Fail(c, err)
		return
	}
	param := &req.UnbindLoginNameReq{}
	err = ParseBody(c, &param)
	if err != nil {
		Fail(c, err)
		return
	}
	resp, resErr := service.UnbindLoginName(operator.OrgId, operator.UserId, *param)
	if resErr != nil {
		Fail(c, resErr)
		return
	}

	Suc(c, resp)
}
