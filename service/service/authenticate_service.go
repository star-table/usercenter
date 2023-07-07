package service

import (
	"strings"
	"time"

	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/random"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/vo"
	"upper.io/db.v3"
)

// Login 登陆
func Login(reqParam vo.UserLoginReq) (*vo.UserLoginResp, errs.SystemErrorInfo) {

	if reqParam.Name != "" {
		reqParam.Name = strings.TrimSpace(reqParam.Name)
		if !format.VerifyNicknameFormat(reqParam.Name) {
			return nil, errs.NicknameLenError
		}
	}

	loginType := reqParam.LoginType
	var userId int64
	var err errs.SystemErrorInfo
	if loginType == 1 || loginType == 3 {
		// 验证码登陆
		userId, err = captchaLogin(reqParam)
	} else if loginType == 2 {
		userId, err = userPwdLogin(reqParam)
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if userId == 0 {
		//暂时不支持的登录类型
		return nil, errs.UnSupportLoginType
	}

	userPo, dbErr := domain.GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.AccountNotRegister
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	// 邀请hook
	trulyOrgId, hookErr := UserAlreadyRegisterHandleHook(reqParam.InviteCode, userPo.Id, userPo.Email)
	if hookErr != nil {
		logger.Error(hookErr)
		return nil, hookErr
	}

	orgId := trulyOrgId
	if orgId == 0 {
		orgId = userPo.OrgId
	}

	token := random.Token()
	res := &vo.UserLoginResp{
		Token:       token,
		UserID:      userPo.Id,
		OrgID:       orgId,
		Name:        userPo.Name,
		Avatar:      userPo.Avatar,
		NeedInitOrg: orgId == 0,
	}

	if orgId != 0 {
		orgInfo, dbErr := domain.GetOrgById(orgId)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				// todo 不抛出异常？
			}
			//todo 不抛出异常？
			logger.Error(err)
		} else {
			res.OrgCode = orgInfo.Code
			res.OrgName = orgInfo.Name
		}
	}

	cacheUserBo := &bo.CacheUserInfoBo{
		UserId:        userPo.Id,
		SourceChannel: userPo.SourceChannel,
		OrgId:         orgId,
	}
	redisErr := store.Redis.SetEx(sconsts.CacheUserToken+res.Token, json.ToJsonIgnoreError(cacheUserBo), consts.CacheUserTokenExpire)
	if redisErr != nil {
		logger.Error(redisErr)
	}

	return res, nil
}

// captchaLogin 验证码登陆
func captchaLogin(req vo.UserLoginReq) (int64, errs.SystemErrorInfo) {

	authCode := req.AuthCode
	loginType := req.LoginType

	name := ""
	inviteCode := ""
	if req.Name != "" {
		name = req.Name
	}
	if req.InviteCode != "" {
		inviteCode = req.InviteCode
	}
	if authCode == "" {
		logger.Error("验证码不能为空")
		return 0, errs.AuthCodeIsNull
	}
	addressType := consts.ContactAddressTypeEmail
	if loginType == consts.LoginTypeSMSCode {
		addressType = consts.ContactAddressTypeMobile
	}
	// 校验验证码
	err := domain.AuthCodeVerify(consts.AuthCodeTypeLogin, addressType, req.LoginName, authCode)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	var loginUserId int64

	if loginType == consts.LoginTypeSMSCode {
		// 手机验证码登录
		// 做登录和自动注册逻辑
		userPo, dbErr := domain.GetUserInfoByMobile(req.LoginName)
		if dbErr != nil {
			// 其他异常抛出
			if dbErr != db.ErrNoMoreRows {
				logger.Error(dbErr)
				return 0, errs.MysqlOperateError
			}
			//除邀请外登录不做注册逻辑
			if inviteCode == "" {
				return 0, errs.AccountNotRegister
			}
			// 开始注册逻辑
			logger.InfoF("用户%s未注册，开始注册....", req.LoginName)
			userId, err := domain.RegisterUser(bo.UserSMSRegisterInfo{
				PhoneNumber:    req.LoginName,
				SourceChannel:  req.SourceChannel,
				SourcePlatform: req.SourcePlatform,
				Name:           name,
				InviteCode:     inviteCode,
			})
			if err != nil {
				logger.Error(err)
				return 0, err
			}
			loginUserId = userId
		} else {
			loginUserId = userPo.Id
		}
	} else if loginType == consts.LoginTypeMail {
		userPo, dbErr := domain.GetUserInfoByEmail(req.LoginName)
		if dbErr != nil {
			// 其他异常抛出
			if dbErr != db.ErrNoMoreRows {
				logger.Error(dbErr)
				return 0, errs.MysqlOperateError
			}
			if inviteCode == "" {
				return 0, errs.AccountNotRegister
			}
			logger.InfoF("用户%s未注册，开始注册....", req.LoginName)
			userId, err := domain.RegisterUser(bo.UserSMSRegisterInfo{
				Email:          req.LoginName,
				SourceChannel:  req.SourceChannel,
				SourcePlatform: req.SourcePlatform,
				Name:           name,
				InviteCode:     inviteCode,
			})
			if err != nil {
				logger.Error(err)
				return 0, err
			}
			loginUserId = userId
		} else {
			loginUserId = userPo.Id
		}
	}

	return loginUserId, nil
}

func UserAlreadyRegisterHandleHook(inviteCode string, userId int64, email string) (int64, errs.SystemErrorInfo) {
	var orgId int64
	//判断邀请逻辑
	if inviteCode != "" {
		inviteCodeInfo, err := domain.GetUserInviteCodeInfo(inviteCode)
		if err != nil {
			logger.Error(err)
			return 0, err
		}
		orgId = inviteCodeInfo.OrgId
		inviterId := inviteCodeInfo.InviterId

		//修改用户默认组织
		dbErr := domain.UpdateUserDefaultOrg(userId, orgId)
		if dbErr != nil {
			logger.Error(dbErr)
		}

		err = AddOrgMember(orgId, userId, inviterId, false, false)
		if err != nil {
			logger.Error(err)
			return 0, err
		}
	}
	//判断是否在邀请表，在的话就变更
	_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUserInvite, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcEmail:    email,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsRegister: 1,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}
	return orgId, nil
}

// userPwdLogin 账号密码登录
func userPwdLogin(req vo.UserLoginReq) (int64, errs.SystemErrorInfo) {
	if req.Password == "" {
		return 0, errs.PasswordEmptyError
	}
	// 查找账号
	user, dbErr := domain.GetUserByLoginNameOrMobileOrEmail(req.LoginName, req.LoginName, req.LoginName)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return 0, errs.AccountNotRegister
		}
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	if user.Password == "" {
		return 0, errs.AccountNotSetPwdErr
	}

	// 匹配密码
	if !util.PwdMatch(user.LoginName+req.Password, user.PasswordSalt, user.Password) {
		if !CheckIsFusePolarisAccount(req.Password, user) {
			return 0, errs.PwdLoginUsrOrPwdNotMatch
		}
	}

	// 是否超过三个月未修改密码
	if time.Now().After(user.LastEditPwdTime.AddDate(0, 3, 0)) {
		return 0, errs.PasswordEditTimeError
	}
	return user.Id, nil
}

// CheckIsFusePolarisAccount 校验是否是合法的极星账户。极星的用户需要在无码系统也能登录。
func CheckIsFusePolarisAccount(inputPwdStr string, userPo *po.PpmOrgUser) bool {
	shouldPwdEncrypted := util.PwdEncryptForFusePolaris(inputPwdStr, userPo.PasswordSalt)
	return userPo.Password == shouldPwdEncrypted
}

// Logout 用户退出
func Logout(outReq req.UserQuitReq) (bool, errs.SystemErrorInfo) {
	err := domain.ClearUserCacheInfo(outReq.Token)
	if err != nil {
		return false, err
	}
	return true, nil
}
