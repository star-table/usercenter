/**
2 * @Author: Nico
3 * @Date: 2020/1/31 11:20
4 */
package service

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/random"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/vo"
	"upper.io/db.v3"
)

// RegisterUser 注册账号
func RegisterUser(input req.UserRegisterReq) (*vo.UserRegisterResp, errs.SystemErrorInfo) {
	authCode := input.AuthCode

	if input.Name != "" {
		input.Name = strings.TrimSpace(input.Name)
		//检测姓名是否合法
		if !format.VerifyNicknameFormat(input.Name) {
			return nil, errs.NicknameLenError
		}
	}
	if input.PhoneNumber != "" {
		input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)
	}

	if input.Email != "" {
		input.Email = strings.TrimSpace(input.Email)
	}

	if input.Password != "" {
		input.Password = strings.TrimSpace(input.Password)
		if !format.VerifyAccountPwdFormat(input.Password) {
			return nil, errs.PwdLengthError
		}
	}

	if authCode == "" {
		return nil, errs.AuthCodeIsNull
	}
	//验证码是否正确
	err := domain.AuthCodeVerify(consts.AuthCodeTypeRegister, consts.ContactAddressTypeMobile, input.PhoneNumber, authCode)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//检测手机号或邮箱是否存在
	err = domain.CheckLoginNameAndPhoneAndEmail("", input.PhoneNumber, input.Email)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	param := bo.UserSMSRegisterInfo{
		PhoneNumber:    input.PhoneNumber,
		SourceChannel:  input.SourceChannel,
		SourcePlatform: input.SourcePlatform,
		Name:           input.Name,
		Email:          input.Email,
		Pwd:            input.Password,
	}
	// 邀请码
	if input.InviteCode != "" {
		param.InviteCode = input.InviteCode
	}
	//注册
	userId, err := domain.RegisterUser(param)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	userPo, dbErr := domain.GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.UserRegisterError
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	//邀请hook
	trulyOrgId, hookErr := UserAlreadyRegisterHandleHook(input.InviteCode, userPo.Id, userPo.Email)

	if hookErr != nil {
		logger.Error(hookErr)
		return nil, hookErr
	}
	orgId := trulyOrgId
	if orgId == 0 {
		orgId = userPo.OrgId
	}

	token := random.Token()
	//自动登录
	cacheUserBo := &bo.CacheUserInfoBo{
		UserId:        userPo.Id,
		SourceChannel: input.SourceChannel,
		OrgId:         orgId,
	}
	redisErr := store.Redis.SetEx(sconsts.CacheUserToken+token, json.ToJsonIgnoreError(cacheUserBo), consts.CacheUserTokenExpire)
	if redisErr != nil {
		logger.Error(redisErr)
	}

	return &vo.UserRegisterResp{
		Token: token,
	}, nil
}

// 检查账号是否存在，检查范围是全局
func CheckLoginNameExist(loginName string) (bool, errs.SystemErrorInfo) {
	err := domain.CheckLoginNameAndPhoneAndEmail(loginName, loginName, loginName)
	if err != nil {
		return true, err
	}
	return false, nil
}
