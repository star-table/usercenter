package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/jsonx"
	"github.com/star-table/usercenter/service/model/bo"

	//"github.com/micro/go-micro/v2/errors"
	"io/ioutil"
)

const (
	// OK ok
	OK int32 = 0

	// RequestErr request error
	RequestErr int32 = -400

	// ServerErr server error
	ServerErr int32 = -500

	contextErrCode = "context/err/code"
)

type res struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Fail(c *gin.Context, err errs.SystemErrorInfo) {
	errCode := RequestErr
	errMsg := "request err"

	if err != nil {
		errCode = int32(err.Code())
		errMsg = err.Message()
	}
	responseBody := jsonx.ToJsonIgnoreError(res{
		Code:    errCode,
		Message: errMsg,
	})
	logger.InfoF("Response FAIL------------------body:%s, ", responseBody)
	c.Set(contextErrCode, errCode)
	c.Data(200, "application/json;charset=utf8", []byte(responseBody))
}

func Suc(c *gin.Context, data interface{}) {
	code := OK
	c.Set(contextErrCode, code)

	responseBody := jsonx.ToJsonIgnoreError(res{
		Code: code,
		Data: data,
	})
	logger.InfoF("Response SUCCESS ------------------body:%s, ", responseBody)
	c.Data(200, "application/json;charset=utf8", []byte(responseBody))
}

// ParseBody 解析body信息
func ParseBody(c *gin.Context, data interface{}) errs.SystemErrorInfo {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		return errs.SystemError
	}
	_ = json.FromJson(string(body), data)
	logger.InfoF("ParseBody------------------body:%s, parse:%s", string(body), json.ToJsonIgnoreError(data))
	return nil
}

// GetOperator 获取操作人信息
func GetOperator(c *gin.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	orgIdInterface, ok := c.Get("orgId")
	if !ok {
		return nil, errs.TokenAuthError
	}
	userIdInterface, ok := c.Get("userId")
	if !ok {
		return nil, errs.TokenAuthError
	}

	orgId := orgIdInterface.(int64)
	userId := userIdInterface.(int64)

	logger.InfoF("==[GetOperator]== orgId:%v, userId:%v", orgId, userId)

	return &bo.CacheUserInfoBo{
		UserId: userId,
		OrgId:  orgId,
	}, nil
}

// AuthHandler 获取认证信息
func AuthHandler(c *gin.Context) {
	orgIdStr := c.GetHeader(consts.IdentityOrgHeader)
	userIdStr := c.GetHeader(consts.IdentityUserHeader)
	if orgIdStr == "" || userIdStr == "" {
		c.Abort()
		Fail(c, errs.TokenAuthError)
		return
	}
	orgId, err := strconv.ParseInt(orgIdStr, 10, 64)
	if err != nil {
		c.Abort()
		Fail(c, errs.TokenAuthError)
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.Abort()
		Fail(c, errs.TokenAuthError)
		return
	}
	c.Set("orgId", orgId)
	c.Set("userId", userId)
}

func getToken(c *gin.Context) (string, errs.SystemErrorInfo) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		return "", errs.TokenAuthError
	}
	return accessToken, nil
}

func InnerFail(c *gin.Context, err errs.SystemErrorInfo) {
	errCode := RequestErr
	errMsg := "request err"

	if err != nil {
		errCode = int32(err.Code())
		errMsg = err.Message()
	}
	responseBody := json.ToJsonIgnoreError(res{
		Code:    errCode,
		Message: errMsg,
	})
	logger.InfoF("Response FAIL ------------------body:%s, ", responseBody)
	c.Set(contextErrCode, errCode)
	c.Data(200, "application/json;charset=utf8", []byte(responseBody))
}

func InnerSuc(c *gin.Context, data interface{}) {
	code := OK
	c.Set(contextErrCode, code)
	responseBody := json.ToJsonIgnoreError(res{
		Code: code,
		Data: data,
	})
	logger.InfoF("Response SUCCESS ------------------body:%s, ", responseBody)
	c.Data(200, "application/json;charset=utf8", []byte(responseBody))
}
