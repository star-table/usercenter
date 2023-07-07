package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/network"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/vo"
	"github.com/star-table/usercenter/service/service"
)

type authenticate int

var Authenticate authenticate = 1

// @Security token
// @Summary 登录
// @Description 登录接口
// @Tags 认证
// @accept application/json
// @Produce application/json
// @param input body req.UserLoginReq true "登录信息"
// @Success 200 {object} resp.UserLoginResp
// @Failure 400
// @Router /usercenter/api/v1/user/login [post]
func (b authenticate) Login(c *gin.Context) {
	var reqParam vo.UserLoginReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := service.Login(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	//更新上次登录时间
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorF("捕获到的错误：%s", r)
			}
		}()
		msg := "登陆成功"
		if err != nil {
			msg = err.Error()
		}
		domain.UserLoginHook(res.OrgID, res.UserID, reqParam.SourceChannel, c.ClientIP(), network.GetUserAgent(c.Request.Header), msg)
	}()
	Suc(c, res)
}

// @Security token
// @Summary 校验token
// @Description 校验token接口
// @Tags 认证
// @accept application/x-www-form-urlencoded
// @Produce application/json
// @param token query string true "token"
// @Success 200 {object} resp.CacheUserInfoBo
// @Failure 400
// @Router /usercenter/inner/api/v1/user/auth [post]
func (b authenticate) Auth(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		Fail(c, errs.TokenAuthError)
		return
	}
	operator, err := service.AuthToken(token, false)
	if err != nil {
		Fail(c, err)
		return
	}

	res := &resp.CacheUserInfoBo{}
	_ = copyer.Copy(operator, res)
	Suc(c, res)
}

// @Security token
// @Summary 校验token并验证状态
// @Description 校验token并验证状态接口
// @Tags 认证
// @accept application/x-www-form-urlencoded
// @Produce application/json
// @param token query string true "token"
// @Success 200 {object} resp.CacheUserInfoBo
// @Failure 400
// @Router /usercenter/inner/api/v1/user/auth-check-status [post]
func (b authenticate) AuthCheckStatus(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		Fail(c, errs.TokenAuthError)
		return
	}
	operator, err := service.AuthToken(token, true)
	if err != nil {
		Fail(c, err)
		return
	}

	res := &resp.CacheUserInfoBo{}
	_ = copyer.Copy(operator, res)
	Suc(c, res)
}

// @Security token
// @Summary 退出登录
// @Description 退出登录接口
// @Tags 认证
// @accept application/json
// @Produce application/json
// @param input body req.UserQuitReq true "入参"
// @Success 200 {object} bool
// @Failure 400
// @Router /usercenter/api/v1/user/quit [post]
func (b authenticate) Logout(c *gin.Context) {
	var reqParam req.UserQuitReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	ok, err := service.Logout(reqParam)
	if err != nil {
		Fail(c, err)
		return
	}

	Suc(c, ok)
}

// @Security token
// @Summary 校验账号是否存在
// @Description 校验账号是否存在
// @Tags 账号
// @accept application/json
// @Produce application/json
// @param input body req.CheckLoginNameExistReq true "用户名是否存在：true 存在，false 不存在"
// @Success 200 {object} string
// @Failure 400
// @Router /usercenter/api/v1/user/check-authenticate-name-exist [post]
func (b authenticate) CheckLoginNameExist(c *gin.Context) {
	var reqParam req.CheckLoginNameExistReq
	err := ParseBody(c, &reqParam)
	if err != nil {
		Fail(c, err)
		return
	}
	ok, err := service.CheckLoginNameExist(reqParam.Name)
	if err != nil {
		Fail(c, err)
		return
	}
	Suc(c, ok)
}
