package domain

import (
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"upper.io/db.v3"
)

// UpdateUserPassword 修改密码
func UpdateUserPassword(operatorUid int64, userId int64, password string, salt string) (int, error) {
	logger.InfoF("[修改密码]operatorUid: %d, userId: %d", operatorUid, userId)
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcPassword:        password,
		consts.TcPasswordSalt:    salt,
		consts.TcLastEditPwdTime: time.Now(),
		consts.TcUpdator:         operatorUid,
		consts.TcVersion:         db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// BindMobileOrEmail 换绑邮箱或者手机号
func BindMobileOrEmail(userId int64, addressType int, loginName string) errs.SystemErrorInfo {
	var lockKey string
	if addressType == consts.ContactAddressTypeMobile {
		lockKey = consts.UserBindMobileLock + loginName
	} else if addressType == consts.ContactAddressTypeEmail {
		lockKey = consts.UserBindEmailLock + loginName
	}
	uid := uuid.NewUuid()
	suc, redisErr := store.Redis.TryGetDistributedLock(lockKey, uid)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.BindLoginNameFail
	}
	if suc {
		defer func() {
			if _, redisErr := store.Redis.ReleaseDistributedLock(lockKey, uid); redisErr != nil {
				logger.Error(redisErr)
			}
		}()
	} else {
		return errs.BindLoginNameFail
	}

	//判断是否被其他账号绑定过
	_, dbErr := GetUserByLoginNameOrMobileOrEmail("", loginName, loginName)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			return errs.MysqlOperateError
		}
	} else {
		if addressType == consts.ContactAddressTypeEmail {
			return errs.EmailAlreadyBindByOtherAccountError
		} else {
			return errs.MobileAlreadyBindOtherAccountError
		}
	}

	upd := mysql.Upd{
		consts.TcUpdator: userId,
	}
	if addressType == consts.ContactAddressTypeEmail {
		upd[consts.TcEmail] = loginName
	} else if addressType == consts.ContactAddressTypeMobile {
		upd[consts.TcMobile] = loginName
	}

	_, dbErr = store.Mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.UnBindLoginNameFail
	}
	return nil
}
