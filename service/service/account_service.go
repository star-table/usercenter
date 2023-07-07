package service

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/md5"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"github.com/star-table/usercenter/service/service/inner_service"
	"upper.io/db.v3"
)

// SetPassword 设置本人密码
func SetPassword(reqParam req.SetPasswordReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[设置本人密码] ->   userId: %d", reqParam.UserId)

	targetPassword := strings.TrimSpace(reqParam.Input.Password)

	if !format.VerifyAccountPwdFormat(targetPassword) {
		return false, errs.PwdLengthError
	}

	user, dbErr := domain.GetUserPoById(reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.UserNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	if strings.TrimSpace(user.Password) != consts.BlankString {
		return false, errs.PwdAlreadySettingsErr
	}

	salt := md5.Md5(uuid.NewUuid())
	pwd := util.PwdEncrypt(user.LoginName+targetPassword, salt)
	_, dbErr = domain.UpdateUserPassword(reqParam.UserId, reqParam.UserId, pwd, salt)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// ResetPassword 重置密码,仅供管理员使用
func ResetPassword(orgId, operatorUid int64, reqParam req.SetPasswordReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[重置密码] ->  operatorUid:%d  targetUid: %d", operatorUid, reqParam.UserId)
	targetPassword := strings.TrimSpace(reqParam.Input.Password)

	if !format.VerifyAccountPwdFormat(targetPassword) {
		return false, errs.PwdLengthError
	}

	// 获取组织成员
	member, err := domain.GetOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExist
		}
		logger.Error(err)
		return false, errs.MysqlOperateError
	}

	salt := md5.Md5(uuid.NewUuid())
	pwd := util.PwdEncrypt(member.LoginName+targetPassword, salt)
	_, dbErr := domain.UpdateUserPassword(operatorUid, reqParam.UserId, pwd, salt)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// UpdatePassword 修改密码
func UpdatePassword(reqParam req.UpdatePasswordReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改本人密码] ->   userId: %d", reqParam.UserId)

	input := reqParam.Input

	targetPassword := strings.TrimSpace(input.NewPassword)
	if !format.VerifyAccountPwdFormat(targetPassword) {
		return false, errs.PwdLengthError
	}

	user, dbErr := domain.GetUserPoById(reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.UserNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 无密码不可修改
	if strings.TrimSpace(user.Password) == consts.BlankString {
		return false, errs.PasswordNotSetError
	}

	if !util.PwdMatch(user.LoginName+input.CurrentPassword, user.PasswordSalt, user.Password) {
		return false, errs.PasswordNotMatchError
	}

	salt := md5.Md5(uuid.NewUuid())
	newPassword := util.PwdEncrypt(user.LoginName+targetPassword, salt)
	_, dbErr = domain.UpdateUserPassword(reqParam.UserId, reqParam.UserId, newPassword, salt)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// UpdatePasswordByUsername 密码过期，修改密码
func UpdatePasswordByUsername(reqParam req.UpdatePasswordByLoginNameReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[通过用户名修改密码] ->   Username: %s", reqParam.Username)

	input := reqParam.Input

	targetPassword := strings.TrimSpace(input.NewPassword)
	if !format.VerifyAccountPwdFormat(targetPassword) {
		return false, errs.PwdLengthError
	}

	user, dbErr := domain.GetUserByLoginNameOrMobileOrEmail(reqParam.Username, reqParam.Username, reqParam.Username)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.AccountNotRegister
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 无密码不可修改 进入该接口的必定是设置过密码的
	if strings.TrimSpace(user.Password) == consts.BlankString {
		return false, errs.PasswordNotSetError
	}

	if !util.PwdMatch(user.LoginName+input.CurrentPassword, user.PasswordSalt, user.Password) {
		if !CheckIsFusePolarisAccount(input.CurrentPassword, user) {
			return false, errs.PasswordNotMatchError
		}
	}

	salt := md5.Md5(uuid.NewUuid())
	newPassword := util.PwdEncrypt(user.LoginName+targetPassword, salt)
	_, dbErr = domain.UpdateUserPassword(user.Id, user.Id, newPassword, salt)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// RetrievePassword 找回密码
func RetrievePassword(reqParam req.RetrievePasswordReq) (bool, errs.SystemErrorInfo) {
	authCode := reqParam.AuthCode
	loginName := reqParam.Username

	// 密码验证格式
	if reqParam.NewPassword == "" {
		return false, errs.PasswordEmptyError
	}
	reqParam.NewPassword = strings.TrimSpace(reqParam.NewPassword)
	if !format.VerifyAccountPwdFormat(reqParam.NewPassword) {
		return false, errs.PwdLengthError
	}

	// 查找账号信息
	user, dbErr := domain.GetUserByLoginNameOrMobileOrEmail("", loginName, loginName)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.NotBindAccountError
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	userId := user.Id

	addressType := consts.ContactAddressTypeMobile
	if format.VerifyEmailFormat(loginName) {
		addressType = consts.ContactAddressTypeEmail
	}

	err := domain.AuthCodeVerify(consts.AuthCodeTypeRetrievePwd, addressType, loginName, authCode)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	salt := md5.Md5(uuid.NewUuid())
	newPassword := util.PwdEncrypt(user.LoginName+reqParam.NewPassword, salt)
	_, dbErr = domain.UpdateUserPassword(userId, userId, newPassword, salt)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// 解绑登录名
func UnbindLoginName(orgId, userId int64, req req.UnbindLoginNameReq) (*resp.UnbindLoginNameResp, errs.SystemErrorInfo) {
	addressType := req.AddressType
	authCode := req.AuthCode
	userPo, dbErr := domain.GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	loginName := ""
	if addressType == consts.ContactAddressTypeEmail {
		if strings.TrimSpace(userPo.Email) == consts.BlankString {
			return nil, errs.EmailNotBindError
		}
		loginName = userPo.Email
	} else if addressType == consts.ContactAddressTypeMobile {
		if strings.TrimSpace(userPo.Mobile) == consts.BlankString {
			return nil, errs.MobileNotBindError
		}
		loginName = userPo.Mobile
	} else {
		return nil, errs.NotSupportedContactAddressType
	}

	if strings.TrimSpace(userPo.Email) == consts.BlankString || strings.TrimSpace(userPo.Mobile) == consts.BlankString {
		return nil, errs.HaveNoContract
	}
	err := domain.AuthCodeVerify(consts.AuthCodeTypeUnBind, addressType, loginName, authCode)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// 产品需求：防止解绑完，但还未来得及重新绑定时造成的问题
	// 用户只有重新绑定，才能解绑，然后绑定。否则无法单独进行解绑，因此去除这一步的解绑操作。改为：在更新时，才会被替换掉就的登录名
	//err = domain.UnbindUserName(userId, addressType)
	//if err != nil {
	//	logger.Error(err)
	//	return err
	//}
	// 生成新的解绑验证码，返回给前端，更新新的登录名时，会进行校验
	code, err := EncodeChangeBindVerifyCode(orgId, userId, addressType, authCode)
	if err != nil {
		return nil, err
	}

	return &resp.UnbindLoginNameResp{
		ChangeBindCode: code,
	}, nil
}

func EncodeChangeBindVerifyCode(orgId, userId int64, addressType int, authCode string) (string, errs.SystemErrorInfo) {
	cacheValue := md5.Md5(authCode)
	// 从 redis 中获取缓存值，进行校验
	key, keyErr := util.ParseCacheKey(consts.ChangeBindCode, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyUserIdConstName:      userId,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return "", errs.TemplateRenderError
	}
	redisErr := store.Redis.SetEx(key, cacheValue, int64(60*5))
	if redisErr != nil {
		logger.Error(redisErr)
		return "", errs.RedisOperateError
	}
	return cacheValue, nil
}

func CheckChangeBindVerifyCode(orgId, userId int64, addressType int, checkValue string) (bool, errs.SystemErrorInfo) {
	// 从 redis 中获取缓存值，进行校验
	key, keyErr := util.ParseCacheKey(consts.ChangeBindCode, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyUserIdConstName:      userId,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return false, errs.TemplateRenderError
	}
	value, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return false, errs.RedisOperateError
	}

	return value == checkValue, nil
}

func BindLoginName(orgId, userId int64, input req.BindLoginNameReq) errs.SystemErrorInfo {
	loginName := input.Address
	addressType := input.AddressType
	authCode := input.AuthCode

	user, dbErr := domain.GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return errs.UserNotExist
		}
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	// 绑定时，有2种情况：1新用户第一次绑定，此时是空的，可以正常绑定；2用户换绑联系方式。
	// 第2种情况，需要额外验证 ChangeBindCode，ChangeBindCode 是验证原联系方式正确解绑的结果 code，如果该 code 不正确，则不允许绑定
	if addressType == consts.ContactAddressTypeEmail {
		if input.ChangeBindCode == consts.BlankString {
			if strings.TrimSpace(user.Email) != consts.BlankString {
				return errs.EmailAlreadyBindError
			}
		} else {
			if ok, err := CheckChangeBindVerifyCode(orgId, userId, addressType, input.ChangeBindCode); !ok {
				if err != nil {
					logger.Error(err)
					return err
				}
				logger.Error("换绑时，换绑验证code错误。")
				return errs.UserRebindCodeError
			}
		}
	} else if addressType == consts.ContactAddressTypeMobile {
		if input.ChangeBindCode == consts.BlankString {
			if strings.TrimSpace(user.Mobile) != consts.BlankString {
				return errs.MobileAlreadyBindError
			}
		} else {
			if ok, err := CheckChangeBindVerifyCode(orgId, userId, addressType, input.ChangeBindCode); !ok {
				if err != nil {
					logger.Error(err)
					return err
				}
				logger.Error("换绑时，换绑验证code错误。")
				return errs.UserRebindCodeError
			}
		}
	} else {
		return errs.NotSupportedContactAddressType
	}

	err := domain.AuthCodeVerify(consts.AuthCodeTypeBind, addressType, loginName, authCode)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = domain.BindMobileOrEmail(userId, addressType, loginName)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// VerifyOldPhoneOrEmail 解绑验证
func VerifyOldPhoneOrEmail(orgId, userId int64, input req.UnbindLoginNameReq) (bool, errs.SystemErrorInfo) {
	addressType := input.AddressType
	authCode := input.AuthCode

	user, dbErr := domain.GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.UserNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	loginName := ""
	if addressType == consts.ContactAddressTypeEmail {
		if strings.TrimSpace(user.Email) == consts.BlankString {
			return false, errs.EmailNotBindError
		}
		loginName = user.Email
	} else if addressType == consts.ContactAddressTypeMobile {
		if strings.TrimSpace(user.Mobile) == consts.BlankString {
			return false, errs.MobileNotBindError
		}
		loginName = user.Mobile
	} else {
		return false, errs.NotSupportedContactAddressType
	}

	err := domain.AuthCodeVerify(consts.AuthCodeTypeUnBind, addressType, loginName, authCode)
	if err != nil {
		return false, err
	}

	//将行为记录到缓存，在有效期内可用
	err = domain.SetChangeLoginNameSign(orgId, userId, input.AddressType)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AuthToken 对token进行认证
func AuthToken(accessToken string, check bool) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	if accessToken == "" {
		return nil, errs.TokenAuthError
	}

	// 模板预览需求-token判断 认证处理
	if accessToken == consts.PreviewTplToken {
		return &bo.CacheUserInfoBo{
			OutUserId:     "",
			SourceChannel: "",
			UserId:        consts.PreviewTplUserId,
			CorpId:        "",
			OrgId:         consts.PreviewTplOrgId,
		}, nil
	} else if accessToken == consts.PreviewTplTokenForWrite {
		return &bo.CacheUserInfoBo{
			OutUserId:     "",
			SourceChannel: "",
			UserId:        consts.PreviewOrWriteTplUserId,
			CorpId:        "",
			OrgId:         consts.PreviewTplOrgId,
		}, nil
	}

	loginUserJson, redisErr := store.Redis.Get(sconsts.CacheUserToken + accessToken)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if loginUserJson == "" {
		return nil, errs.TokenAuthError
	}

	loginUser := &bo.CacheUserInfoBo{}
	_ = json.FromJson(loginUserJson, loginUser)
	if check {
		// 未初始化组织，应当提示并引导用户创建组织
		if loginUser.OrgId == 0 {
			return nil, errs.OrgNotInitError
		}
		//校验用户
		//baseUserInfo, err := domain.GetBaseUserInfo("", loginUser.OrgId, loginUser.UserId)
		//if err != nil {
		//	logger.Error(err)
		//	return nil, err
		//}
		//err = BaseUserInfoOrgStatusCheck(*baseUserInfo)
		//if err != nil {
		//	logger.Error(err)
		//	return nil, err
		//}
	}
	return loginUser, nil
}

func GetOrgUserPerContext(orgId, userId int64) (*inner_resp.OrgUserPerContext, errs.SystemErrorInfo) {
	auth, err := inner_service.GetUserAuthorityByUserId(orgId, userId)
	if err != nil {
		return nil, err
	}
	// 查询可以管理的人员
	ctx := inner_resp.NewOrgUserPermissionContext(*auth)
	logger.InfoF("用户权限上下文 %s", json.ToJsonIgnoreError(auth))
	return ctx, nil
}

func GetUserManageAuth(orgId, userId int64) (*resp.UserManageAuthResp, errs.SystemErrorInfo) {
	auth, err := inner_service.GetUserAuthorityByUserId(orgId, userId)
	if err != nil {
		return nil, err
	}

	return &resp.UserManageAuthResp{
		IsOrgOwner: auth.IsOrgOwner,
		IsSysAdmin: auth.IsSysAdmin,
		IsSubAdmin: auth.IsSubAdmin,
	}, nil
}
