package domain

import (
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/core/types"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/date"
	"github.com/star-table/usercenter/pkg/util/md5"
	"github.com/star-table/usercenter/pkg/util/pinyin"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

func GetOutUserName(orgId, userId int64, sourceChannel string) (string, errs.SystemErrorInfo) {
	//用户外部信息
	var userOutInfo po.PpmOrgUserOutInfo
	dbErr := store.Mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        userId,
	}, &userOutInfo)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return "", errs.UserOutInfoNotExist
		} else {
			logger.Error(dbErr)
			return "", errs.MysqlOperateError
		}
	}

	return userOutInfo.Name, nil
}

// GetUserPoById 获取用户信息
func GetUserPoById(userId int64) (*po.PpmOrgUser, error) {
	var user po.PpmOrgUser
	dbErr := store.Mysql.SelectOneByCond(user.TableName(), db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &user)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return &user, nil
}

// GetUserOrgPo 获取用户组织信息
func GetUserOrgPo(orgId, userId int64) (*po.PpmOrgUserOrganization, error) {
	var userOrg po.PpmOrgUserOrganization
	dbErr := store.Mysql.SelectOneByCond(userOrg.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userOrg)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return &userOrg, nil
}

// GetOrgUserInfo 获取组织用户信息
func GetOrgUserInfo(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, errs.SystemErrorInfo) {

	baseOrgInfo, err := GetBaseOrgInfo(sourceChannel, orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.OrgNotInitError
	}

	baseUserInfo, err := GetBaseUserInfo(sourceChannel, orgId, userId)
	if err != nil {
		logger.Error(err)
		return nil, errs.UserNotInitError
	}

	userPo, dbErr := GetUserPoById(userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	userOrgPo, dbErr := GetUserOrgPo(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.UserOrgNotRelation
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	ownerInfo := &bo.UserInfoBo{
		ID:                 userPo.Id,
		EmplID:             &baseUserInfo.OutUserId,
		OrgID:              orgId,
		OrgName:            baseOrgInfo.OrgName,
		OrgOwner:           baseOrgInfo.OrgOwnerId,
		Name:               userPo.Name,
		NamePinyin:         userPo.NamePinyin,
		LoginName:          userPo.LoginName,
		LoginNameEditCount: userPo.LoginNameEditCount,
		Email:              userPo.Email,
		MobileRegion:       userPo.MobileRegion,
		Mobile:             userPo.Mobile,
		Birthday:           userPo.Birthday,
		Password:           userPo.Password,
		PasswordSalt:       userPo.PasswordSalt,
		Avatar:             userPo.Avatar,
		SourceChannel:      userPo.SourceChannel,
		Language:           userPo.Language,
		Motto:              userPo.Motto,
		LastLoginIP:        userPo.LastLoginIp,
		LastLoginTime:      userPo.LastLoginTime,
		LoginFailCount:     userPo.LoginFailCount,
		LastEditPwdTime:    userPo.LastEditPwdTime,
		RemindBindPhone:    userPo.RemindBindPhone,
		CreateTime:         userPo.CreateTime,
		UpdateTime:         userPo.UpdateTime,
		Sex:                userPo.Sex,
		Rimanente:          10,     //这里先默认写死
		Level:              1,      //这里先默认写死
		LevelName:          "试用级别", //这里先默认写死
		UserType:           userOrgPo.Type,
	}
	//获取第三方名称
	if ownerInfo.SourceChannel != "" {
		outName, err := GetOutUserName(orgId, userId, ownerInfo.SourceChannel)
		if err != nil {
			return nil, err
		}
		ownerInfo.ThirdName = outName
	}

	return ownerInfo, nil
}

//如果err不等于空，说明用户未注册
func GetUserInfoByMobile(phoneNumber string) (*po.PpmOrgUser, error) {
	var userPo po.PpmOrgUser
	dbErr := store.Mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcMobile:   phoneNumber,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userPo)
	if dbErr != nil {
		return nil, dbErr
	}
	return &userPo, nil
}

func GetUserInfoByEmail(email string) (*po.PpmOrgUser, error) {
	var userPo po.PpmOrgUser
	dbErr := store.Mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcEmail:    email,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userPo)
	if dbErr != nil {
		return nil, dbErr
	}
	return &userPo, nil
}

// CheckLoginNameAndPhoneAndEmail 检查用户名/手机号/邮箱是否已注册
func CheckLoginNameAndPhoneAndEmail(loginName, mobile, email string) errs.SystemErrorInfo {
	user, dbErr := GetUserByLoginNameOrMobileOrEmail(loginName, mobile, email)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil
		}
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}
	if loginName == user.LoginName || user.LoginName != "" {
		return errs.AccountRegisterConflict
	}
	if mobile == user.Mobile || user.Mobile != "" {
		return errs.MobileAlreadyRegister
	}
	if email == user.Email || user.Email != "" {
		return errs.EmailAlreadyRegister
	}
	return nil
}

// GetUserByLoginNameOrMobileOrEmail 根据账号/手机号/邮箱获取账号信息
func GetUserByLoginNameOrMobileOrEmail(loginName, mobile, email string) (*po.PpmOrgUser, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	condList := make([]db.Compound, 0)
	loginNameList := make([]string, 0)
	if mobile != "" {
		condList = append(condList, db.Cond{
			consts.TcMobile: mobile,
		})
		loginNameList = append(loginNameList, mobile)
	}
	if email != "" {
		condList = append(condList, db.Cond{
			consts.TcEmail: email,
		})
		loginNameList = append(loginNameList, email)
	}
	if loginName != "" {
		loginNameList = append(loginNameList, loginName)
	}
	if len(loginNameList) > 0 {
		condList = append(condList, db.Cond{
			consts.TcLoginName: db.In(loginNameList),
		})
	}

	var userPo po.PpmOrgUser
	dbErr = conn.Collection(consts.TableUser).Find(db.And(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, db.Or(condList...))).One(&userPo)

	if dbErr != nil {
		return nil, dbErr
	}

	return &userPo, nil
}

// RegisterUser 注册账号
func RegisterUser(reqParam bo.UserSMSRegisterInfo) (int64, errs.SystemErrorInfo) {

	uid := uuid.NewUuid()
	// 对手机号加锁
	if reqParam.PhoneNumber != "" {
		lockKey := consts.UserBindMobileLock + reqParam.PhoneNumber
		suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, uid)
		if redisErr != nil {
			logger.Error(redisErr)
			return 0, errs.UserRegisterError
		}
		if suc {
			defer func() {
				if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, uid); redisErr != nil {
					logger.Error(redisErr)
				}
			}()
		} else {
			logger.Error("注册失败")
			return 0, errs.UserRegisterError
		}
	}
	// 对邮箱加锁
	if reqParam.Email != "" {
		lockKey := consts.UserBindEmailLock + reqParam.Email
		suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, uid)
		if redisErr != nil {
			logger.Error(redisErr)
			return 0, errs.UserRegisterError
		}
		if suc {
			defer func() {
				if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, uid); redisErr != nil {
					logger.Error(redisErr)
				}
			}()
		} else {
			logger.Error("注册失败")
			return 0, errs.UserRegisterError
		}
	}

	// 注册新用户

	userId := snowflake.Id()
	userPo := &po.PpmOrgUser{
		Id:                 userId,
		OrgId:              0,
		Name:               reqParam.Name,
		NamePinyin:         pinyin.ConvertToPinyin(reqParam.Name),
		Avatar:             "",
		LoginNameEditCount: 0,
		Email:              reqParam.Email,
		Mobile:             reqParam.PhoneNumber,
		SourceChannel:      reqParam.SourceChannel,
		SourcePlatform:     reqParam.SourcePlatform,
		LastEditPwdTime:    time.Now(),
		Creator:            userId,
		Updator:            userId,
	}

	if reqParam.Pwd != "" {
		userPo.PasswordSalt = md5.Md5(uuid.NewUuid())
		userPo.Password = util.PwdEncrypt(userPo.LoginName+reqParam.Pwd, userPo.PasswordSalt)
	}

	//插入用户
	dbErr := store.Mysql.Insert(userPo)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.UserRegisterError
	}

	//即使注册时插入失败，查看时也会做二次check并插入
	err := insertUserConfig(userPo.OrgId, userPo.Id)
	if err != nil {
		logger.Error(err)
	}

	return userId, nil
}

// BindUserOrgRelation 添加用户和组织关联
func BindUserOrgRelation(orgId, userId int64, inUsed bool, inCheck bool, inDisabled bool, inviteId int64) (int64, error) {
	bindId := snowflake.Id()

	useStatus := consts.AppStatusDisabled
	if inUsed {
		useStatus = consts.AppStatusEnable
	}

	checkStatus := consts.AppCheckStatusSuccess
	status := consts.AppStatusEnable
	if inCheck {
		checkStatus = consts.AppCheckStatusWait
		status = consts.AppStatusDisabled
	}

	if inDisabled {
		status = consts.AppStatusDisabled
	}

	userOrgPo := po.PpmOrgUserOrganization{
		Id:          bindId,
		OrgId:       orgId,
		UserId:      userId,
		InviteId:    inviteId,
		CheckStatus: checkStatus,
		UseStatus:   useStatus,
		Status:      status,
		Creator:     userId,
		Updator:     userId,
	}

	dbErr := store.Mysql.Insert(&userOrgPo)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return bindId, nil
}

// UpdateUserDefaultOrg 修改用户的默认组织
func UpdateUserDefaultOrg(userId, orgId int64) error {
	dbErr := store.Mysql.UpdateSmart(consts.TableUser, userId, mysql.Upd{
		consts.TcOrgId: orgId,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	return nil
}

// GetUserListByIds 获取用户信息
func GetUserListByIds(userIds []int64) ([]po.PpmOrgUser, error) {
	var pos []po.PpmOrgUser
	dbErr := store.Mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(userIds),
	}, &pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return pos, nil
}

// GetUsersByCond 根据条件获取用户
func GetUsersByCond(cond db.Cond) ([]po.PpmOrgUser, error) {
	var pos []po.PpmOrgUser
	dbErr := store.Mysql.SelectAllByCond(consts.TableUser, cond, &pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return pos, nil
}

// GetUserIdsOfOrgByCond 根据条件获取某个组织的用户 id
func GetUserIdsOfOrgByCond(cond db.Cond, page, size int) ([]po.PpmOrgUserOrganization, error) {
	var pos []po.PpmOrgUserOrganization
	_, dbErr := store.Mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, cond, nil, page, size, "user_id asc", &pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return pos, nil
}

// UserLoginHook 记录用户登陆行为
// 不抛出错误，对业务无影响
func UserLoginHook(orgId int64, userId int64, sourceChannel string, clientId string, userAgent string, msg string) {
	logger.InfoF("[记录用户登陆行为] ->  orgId: %d, userId: %d, sourceChannel:%s", orgId, userId, sourceChannel)

	// 更新用户最后登录时间
	dbErr := store.Mysql.UpdateSmart(consts.TableUser, userId, mysql.Upd{
		consts.TcLastLoginTime: date.FormatTime(types.NowTime()),
	})
	if dbErr != nil {
		logger.ErrorF("[更新用户最后登陆时间] -> 失败， dbErr: %s", dbErr)
	}
	if orgId != 0 {
		_, dbErr := UpdateOrgUserUseStatus(orgId, userId)
		if dbErr != nil {
			logger.ErrorF("[更新用户使用状态] -> 失败，dbErr: %s", dbErr)
		}
	}
	_, dbErr = AddLoginRecord(orgId, userId, sourceChannel, clientId, userAgent, msg)
	if dbErr != nil {
		logger.ErrorF("[记录用户登陆日志] -> 失败，dbErr: %s", dbErr)
	}
}
