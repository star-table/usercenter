package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/star-table/usercenter/pkg/util/json"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/md5"
	"github.com/star-table/usercenter/pkg/util/pinyin"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// ChangeOrgMemberStatus 修改组织成员状态, 1可用,2禁用,3离职
func ChangeOrgMemberStatus(orgId int64, operatorUid int64, memberIds []int64, status int, tx sqlbuilder.Tx) (int, error) {
	count, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcUserId:      db.In(memberIds),
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcStatus:      db.NotEq(status),
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
	}, mysql.Upd{
		consts.TcStatus:           status,
		consts.TcStatusChangerId:  operatorUid,
		consts.TcStatusChangeTime: time.Now(),
		consts.TcUpdator:          operatorUid,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	if count > 0 {
		// 非启用则删除管理组中绑定的用户关系
		if status != consts.AppStatusEnable {
			for _, uid := range memberIds {
				_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableManageGroup, db.Cond{
					consts.TcOrgId:    orgId,
					consts.TcIsDelete: consts.AppIsNoDelete,
					db.Raw("json_search(`user_ids`, 'one', ?)", uid): db.IsNotNull(),
				}, mysql.Upd{
					consts.TcUserIds: db.Raw("json_remove(`user_ids`,JSON_UNQUOTE(json_search(`user_ids`, 'one', ?)))", uid),
					consts.TcUpdator: operatorUid,
				})
				if dbErr != nil {
					logger.Error(dbErr)
					return 0, dbErr
				}
			}
		}
	}

	return int(count), nil
}

// RemoveOrgMember 删除组织成员
func RemoveOrgMember(orgId int64, operatorUid int64, memberIds []int64, tx sqlbuilder.Tx) (int, error) {
	count, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcUserId:      db.In(memberIds),
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcIsDelete:    consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcAuditorId: operatorUid,
		consts.TcUpdator:   operatorUid,
		consts.TcIsDelete:  consts.AppIsDeleted,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// GetEnableOrgMemberBaseInfoListByUsers 批量获取组织成员 根据成员ID列表（审批通过，并且启用的）
func GetEnableOrgMemberBaseInfoListByUsers(orgId int64, userIds []int64) ([]bo.OrgMemberBaseInfoBo, error) {
	if len(userIds) == 0 {
		return []bo.OrgMemberBaseInfoBo{}, nil
	}
	return _getOrgMemberBaseInfoList(orgId, userIds, nil, consts.AppCheckStatusSuccess, consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetEnableOrgMemberBaseInfoList 批量获取组织成员（审批通过，并且启用的）
func GetEnableOrgMemberBaseInfoList(orgId int64) ([]bo.OrgMemberBaseInfoBo, error) {
	return _getOrgMemberBaseInfoList(orgId, nil, nil, consts.AppCheckStatusSuccess, consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetOrgMemberBaseInfoListByUsers 批量获取组织成员 根据成员ID列表（审批通过的）
func GetOrgMemberBaseInfoListByUsers(orgId int64, userIds []int64) ([]bo.OrgMemberBaseInfoBo, error) {
	if len(userIds) == 0 {
		return []bo.OrgMemberBaseInfoBo{}, nil
	}
	return _getOrgMemberBaseInfoList(orgId, userIds, nil, consts.AppCheckStatusSuccess, 0, consts.AppIsNoDelete)
}

// GetOrgMemberBaseInfoListByOrg 批量获取组织成员（审批通过的）
func GetOrgMemberBaseInfoListByOrg(orgId int64) ([]bo.OrgMemberBaseInfoBo, error) {
	return _getOrgMemberBaseInfoList(orgId, nil, nil, consts.AppCheckStatusSuccess, 0, consts.AppIsNoDelete)
}

// GetOrgMemberBaseInfoListByExcludeIds 批量获取组织成员 排除指定成员（审批通过）
func GetOrgMemberBaseInfoListByExcludeIds(orgId int64, excludeIds []int64) ([]bo.OrgMemberBaseInfoBo, error) {
	return _getOrgMemberBaseInfoList(orgId, nil, excludeIds, consts.AppCheckStatusSuccess, consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetOrgMemberBaseInfoList 批量获取组织成员，各状态自己填
func GetOrgMemberBaseInfoList(orgId int64, userIds []int64, checkStatus int, status int, isDelete int) ([]bo.OrgMemberBaseInfoBo, error) {
	return _getOrgMemberBaseInfoList(orgId, userIds, nil, checkStatus, status, isDelete)
}

func GetSysUserIds(allIds []int64) ([]int64, []int64) {
	var sysIds []int64
	var userIds []int64
	for _, id := range allIds {
		if id <= 100 {
			sysIds = append(sysIds, id)
		} else {
			userIds = append(userIds, id)
		}
	}
	return sysIds, userIds
}

func GetUserBaseInfos(orgId int64, userIds []int64) ([]*bo.UserBaseInfoBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	sysIds, uIds := GetSysUserIds(userIds)
	if len(sysIds) == 0 && len(uIds) == 0 {
		return nil, nil
	}

	var cmds []db.Compound
	if len(sysIds) > 0 {
		cmds = append(cmds, db.Cond{
			"userOrg." + consts.TcOrgId:  999,
			"userOrg." + consts.TcUserId: db.In(sysIds),
		})
	}
	if len(uIds) > 0 {
		cond := db.Cond{
			"userOrg." + consts.TcOrgId:       orgId,
			"userOrg." + consts.TcUserId:      db.In(uIds),
			"userOrg." + consts.TcIsDelete:    consts.AppIsNoDelete,
			"userOrg." + consts.TcCheckStatus: consts.AppCheckStatusSuccess,
			"userOrg." + consts.TcStatus:      consts.AppStatusEnable,
		}
		cmds = append(cmds, cond)
	}

	var baseUserInfos []*bo.UserBaseInfoBo
	query := conn.Select(
		"user.id as id",
		"user.name as name",
		"user.avatar as avatar",
		"user.email as email",
		"user.mobile as mobile1",
		"gUser.mobile as mobile2",
		"org.source_channel as source_channel",
		"userOut.out_user_id as open_id",
	).From(consts.TableUserOrganization + " userOrg").
		Join(consts.TableUser + " user").
		On(db.Cond{
			"user." + consts.TcId:       db.Raw("userOrg.user_id"),
			"user." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Join(consts.TableUserOutInfo + " userOut").
		On(db.Cond{
			"userOut." + consts.TcOrgId:    db.Raw("userOrg.org_id"),
			"userOut." + consts.TcUserId:   db.Raw("userOrg.user_id"),
			"userOut." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Join(consts.TableOrganization + " org").
		On(db.Cond{
			"org." + consts.TcId:       db.Raw("userOrg.org_id"),
			"org." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		LeftJoin(consts.TableGlobalUserRelation + " gUserRel").
		On(db.Cond{
			"gUserRel." + consts.TcUserId:   db.Raw("userOrg.user_id"),
			"gUserRel." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		LeftJoin(consts.TableGlobalUser + " gUser").
		On(db.Cond{
			"gUser." + consts.TcId:       db.Raw("gUserRel.global_user_id"),
			"gUser." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Where(db.Or(cmds...))
	logger.InfoF("[GetUserBaseInfos] query: %v, args: %v", query.String(), json.ToJsonIgnoreError(query.Arguments()))
	dbErr = query.All(&baseUserInfos)
	if dbErr != nil {
		logger.ErrorF("[GetUserBaseInfos] err: %v", dbErr)
		return nil, dbErr
	}

	return baseUserInfos, nil
}

// _getOrgMemberBaseInfoList 批量获取组织成员  禁止其他包使用
// len(userIds)==0 忽视条件
// len(excludeUserIds)==0 忽视条件
// checkStatus==0 忽视条件
// status==0 忽视条件
func _getOrgMemberBaseInfoList(orgId int64, userIds []int64, excludeIds []int64, checkStatus int, status int, isDelete int) ([]bo.OrgMemberBaseInfoBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	sysIds, uIds := GetSysUserIds(userIds)

	var cmds []db.Compound
	// 支持下获取系统用户的信息
	if len(sysIds) > 0 {
		cmds = append(cmds, db.Cond{
			"userOrg." + consts.TcOrgId:  999,
			"userOrg." + consts.TcUserId: db.In(sysIds),
		})
	}
	if len(uIds) > 0 {
		cond := db.Cond{
			"userOrg." + consts.TcOrgId: orgId,
		}
		if isDelete > 0 {
			cond["userOrg."+consts.TcIsDelete] = isDelete
		}
		if len(uIds) > 0 {
			cond["userOrg."+consts.TcUserId] = db.In(uIds)
		}
		if len(excludeIds) > 0 {
			cond["userOrg."+consts.TcUserId+" "] = db.NotIn(excludeIds)
		}
		if checkStatus == consts.AppCheckStatusSuccess || checkStatus == consts.AppCheckStatusFail || checkStatus == consts.AppCheckStatusWait {
			cond["userOrg."+consts.TcCheckStatus] = checkStatus
		}
		if status == consts.AppStatusEnable || status == consts.AppStatusDisabled || status == consts.OrgUserStatusResigned {
			cond["userOrg."+consts.TcStatus] = status
		}
		cmds = append(cmds, cond)
	}

	var baseUserInfos []bo.OrgMemberBaseInfoBo
	dbErr = conn.Select(
		"org.id as org_id",
		"org.name as org_name",
		"org.owner as org_owner",
		"org.creator as org_creator",
		"user.id as user_id",
		"user.login_name",
		"user.name",
		"user.name_pinyin",
		"user.sex",
		"user.birthday",
		"user.email",
		"user.mobile_region",
		"user.mobile",
		"user.avatar as avatar",
		"user.language",
		"userOrg.check_status",
		"userOrg.status as status",
		"userOrg.status_change_time",
		"userOrg.use_status",
		"userOrg.auditor_id",
		"userOrg.audit_time",
		"userOrg.emp_no",
		"userOrg.weibo_ids",
		"userOrg.creator",
		"userOrg.create_time",
		"userOrg.updator",
		"userOrg.update_time",
		"userOrg.is_delete as is_delete",
		"userOrg.type as type",
	).From(consts.TableUserOrganization + " userOrg").
		Join(consts.TableUser + " user").
		On(db.Cond{
			"user." + consts.TcId: db.Raw("userOrg.user_id"),
			//"user." + consts.TcStatus:   consts.AppStatusEnable,
			"user." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Join(consts.TableOrganization + " org").
		On(db.Cond{
			"org." + consts.TcId:       db.Raw("userOrg.org_id"),
			"org." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Where(db.Or(cmds...)).All(&baseUserInfos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return baseUserInfos, nil
}

// GetOrgMemberBaseInfoListByQueryCond 批量获取组织成员  涉及模糊查询 慎用!!
func GetOrgMemberBaseInfoListByQueryCond(orgId int64, queryReq open_req.MemberQueryReq) ([]bo.OrgMemberBaseInfoBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	cond := db.Cond{
		"userOrg." + consts.TcOrgId:    orgId,
		"userOrg." + consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if len(queryReq.UserIds) > 0 {
		cond["userOrg."+consts.TcUserId] = db.In(queryReq.UserIds)
	}
	if queryReq.Nickname != "" {
		cond["user."+consts.TcName] = db.Like("%" + queryReq.Nickname + "%")
	}
	if queryReq.LoginName != "" {
		cond["user."+consts.TcLoginName] = db.Like("%" + queryReq.LoginName + "%")
	}
	if queryReq.Mobile != "" {
		cond["user."+consts.TcMobile] = db.Like("%" + queryReq.Mobile + "%")
	}
	if queryReq.Email != "" {
		cond["user."+consts.TcEmail] = db.Like("%" + queryReq.Email + "%")
	}
	if queryReq.Status == consts.AppStatusEnable || queryReq.Status == consts.AppStatusDisabled || queryReq.Status == consts.OrgUserStatusResigned {
		cond["userOrg."+consts.TcStatus] = queryReq.Status
	}
	var baseUserInfos []bo.OrgMemberBaseInfoBo
	dbErr = conn.Select(
		"org.id as org_id",
		"org.name as org_name",
		"org.owner as org_owner",
		"org.creator as org_creator",
		"user.id as user_id",
		"user.login_name",
		"user.name",
		"user.name_pinyin",
		"user.sex",
		"user.birthday",
		"user.email",
		"user.mobile_region",
		"user.mobile",
		"user.avatar",
		"user.language",
		"userOrg.check_status",
		"userOrg.status as status",
		"userOrg.status_change_time",
		"userOrg.use_status",
		"userOrg.auditor_id",
		"userOrg.audit_time",
		"userOrg.emp_no",
		"userOrg.weibo_ids",
		"userOrg.creator",
		"userOrg.create_time",
		"userOrg.updator",
		"userOrg.update_time",
		"userOrg.is_delete as is_delete",
	).From(consts.TableUserOrganization + " userOrg").
		Join(consts.TableUser + " user").
		On(db.Cond{
			"user." + consts.TcId:       db.Raw("userOrg.user_id"),
			"user." + consts.TcStatus:   consts.AppStatusEnable,
			"user." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Join(consts.TableOrganization + " org").
		On(db.Cond{
			"org." + consts.TcId:       db.Raw("userOrg.org_id"),
			"org." + consts.TcIsDelete: consts.AppIsNoDelete,
		}).
		Where(cond).All(&baseUserInfos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return baseUserInfos, nil
}

// GetEnableOrgMemberBaseInfoByUser 获取组织成员 根据成员ID （审批通过，并且启用的）
func GetEnableOrgMemberBaseInfoByUser(orgId int64, userId int64) (*bo.OrgMemberBaseInfoBo, error) {
	return getOrgMemberBaseInfoByUser(orgId, userId, consts.AppCheckStatusSuccess, consts.AppStatusEnable)
}

// GetOrgMemberBaseInfoByUser 获取组织成员 根据成员ID （审批通过的）
func GetOrgMemberBaseInfoByUser(orgId int64, userId int64) (*bo.OrgMemberBaseInfoBo, error) {
	return getOrgMemberBaseInfoByUser(orgId, userId, consts.AppCheckStatusSuccess, 0)
}

// getOrgMemberBaseInfoByUser 获取组织成员 根据成员ID
// checkStatus=0忽视条件
// status=0忽视条件
func getOrgMemberBaseInfoByUser(orgId int64, userId int64, checkStatus int, status int) (*bo.OrgMemberBaseInfoBo, error) {
	members, dbErr := _getOrgMemberBaseInfoList(orgId, []int64{userId}, nil, checkStatus, status, consts.AppIsNoDelete)
	if dbErr != nil {
		return nil, dbErr
	}
	if len(members) == 0 {
		return nil, db.ErrNoMoreRows
	}
	return &members[0], nil
}

// CreateOrgUser 创建组织用户
func CreateOrgUser(orgId, operatorUid int64, reqParam req.CreateOrgMemberReq) (int64, errs.SystemErrorInfo) {
	uid := uuid.NewUuid()
	if reqParam.LoginName != "" {
		// 对账号加锁
		lockKey := consts.UserBindLoginNameLock + reqParam.LoginName
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
	// 对工号加锁
	if reqParam.EmpNo != "" {
		lockKey := consts.UserBindEmpNoLock + reqParam.EmpNo
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

	// 检测账号是否存在
	err := CheckLoginNameAndPhoneAndEmail(reqParam.LoginName, reqParam.PhoneNumber, reqParam.Email)
	if err != nil {
		return 0, err
	}

	// 检测工号是否存在
	if reqParam.EmpNo != "" {
		reqParam.EmpNo = strings.TrimSpace(reqParam.EmpNo)
		exist, dbErr := store.Mysql.IsExistByCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcEmpNo:       reqParam.EmpNo,
			consts.TcEmpNo + " ": db.NotEq(""), // 排除空
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
		if exist {
			return 0, errs.OrgEmpNoConflictErr
		}
	}

	userPo := &po.PpmOrgUser{
		Id:                 snowflake.Id(),
		OrgId:              orgId,
		Name:               reqParam.Name,
		NamePinyin:         pinyin.ConvertToPinyin(reqParam.Name),
		Avatar:             reqParam.Avatar,
		LoginName:          reqParam.LoginName,
		LoginNameEditCount: 0,
		Email:              reqParam.Email,
		MobileRegion:       reqParam.PhoneRegion,
		Mobile:             reqParam.PhoneNumber,
		LastEditPwdTime:    time.Now(),
		Creator:            operatorUid,
		Updator:            operatorUid,
	}

	// 可以不设置密码
	if reqParam.Password != "" {
		userPo.PasswordSalt = md5.Md5(uuid.NewUuid())
		userPo.Password = util.PwdEncrypt(reqParam.LoginName+reqParam.Password, userPo.PasswordSalt)
	}

	// 插入用户
	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		//用户表
		dbErr := store.Mysql.TransInsert(tx, userPo)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		//用户配置表
		userConfig := &po.PpmOrgUserConfig{}
		userConfigId := snowflake.Id()
		userConfig.Id = userConfigId
		userConfig.OrgId = orgId
		userConfig.UserId = userPo.Id
		userConfig.Creator = operatorUid
		userConfig.Updator = operatorUid
		dbErr = store.Mysql.TransInsert(tx, userConfig)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		if reqParam.Status == 0 {
			// 默认启用
			reqParam.Status = consts.AppStatusEnable
		}
		if reqParam.WeiboIds == nil {
			reqParam.WeiboIds = []string{}
		}
		//用户组织表
		userOrgPo := &po.PpmOrgUserOrganization{
			Id:          snowflake.Id(),
			OrgId:       orgId,
			UserId:      userPo.Id,
			CheckStatus: consts.AppCheckStatusSuccess,
			UseStatus:   consts.AppStatusDisabled,
			Status:      reqParam.Status,
			EmpNo:       reqParam.EmpNo,
			WeiboIds:    strings.Join(reqParam.WeiboIds, ","),
			Creator:     operatorUid,
		}

		dbErr = store.Mysql.TransInsert(tx, userOrgPo)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		// 成员部门表
		if len(reqParam.DeptAndPositions) > 0 {
			var userDepartment []interface{}
			for _, deptAndPosition := range reqParam.DeptAndPositions {
				if deptAndPosition.PositionId == 0 {
					deptAndPosition.PositionId = consts.PositionMemberId
				}
				userDepartment = append(userDepartment, po.PpmOrgUserDepartment{
					Id:            snowflake.Id(),
					OrgId:         orgId,
					UserId:        userPo.Id,
					DepartmentId:  deptAndPosition.DepartmentId,
					OrgPositionId: deptAndPosition.PositionId,
					Creator:       operatorUid,
				})
			}

			dbErr = store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, userDepartment)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		//用户角色表
		if len(reqParam.RoleIds) > 0 {
			var userRole []interface{}
			for _, id := range reqParam.RoleIds {
				userRole = append(userRole, po.PpmRolRoleUser{
					Id:      snowflake.Id(),
					OrgId:   orgId,
					RoleId:  id,
					UserId:  userPo.Id,
					Creator: operatorUid,
					Updator: operatorUid,
				})
			}

			dbErr = store.Mysql.TransBatchInsert(tx, &po.PpmRolRoleUser{}, userRole)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	return userPo.Id, nil
}

// UpdateOrgUserUseStatus 修改组织用户的使用状态
func UpdateOrgUserUseStatus(orgId int64, userId int64) (int, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcUserId:    userId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcUseStatus: consts.AppStatusDisabled,
	}, mysql.Upd{
		consts.TcUseStatus: consts.AppStatusEnable,
		consts.TcVersion:   db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// UpdateOrgUserType 修改组织用户的类型
func UpdateOrgUserType(orgId, operatorId int64, userIds []int64, t int) (int, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcType:    t,
		consts.TcUpdator: operatorId,
		consts.TcVersion: db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// GetOrgMemberListByUser 根据用户获取组织成员列表
func GetOrgMemberListByUser(userId int64) ([]po.PpmOrgUserOrganization, error) {
	cond := db.Cond{
		consts.TcUserId:      userId,
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
	}

	var orgMembers []po.PpmOrgUserOrganization
	dbErr := store.Mysql.SelectAllByCond(consts.TableUserOrganization, cond, &orgMembers)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return orgMembers, nil
}

// GetNewestOrgMemberByOrgAndUser 用来获取最新的组织成员信息
func GetNewestOrgMemberByOrgAndUser(orgId, userId int64) (*po.PpmOrgUserOrganization, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	var orgMember po.PpmOrgUserOrganization
	dbErr = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
	}).OrderBy("id desc").Limit(1).One(&orgMember)
	if dbErr != nil {
		return nil, dbErr
	}
	return &orgMember, nil
}

// GetOrgUserCount 获取组织总人数
func GetOrgUserCount(orgId int64) (uint64, error) {
	total, err := store.Mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	return total, nil
}

func GetOrgUserSimpleInfo(orgId int64, needDelete bool) ([]inner_resp.SimpleInfo, errs.SystemErrorInfo) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	condition := db.Cond{
		"o." + consts.TcOrgId: orgId,
		"u." + consts.TcId:    db.Raw("o." + consts.TcUserId),
	}
	if !needDelete {
		condition["o."+consts.TcIsDelete] = consts.AppIsNoDelete
		condition["u."+consts.TcIsDelete] = consts.AppIsNoDelete
	}

	var userInfo []bo.OrgUserSimpleInfoBo
	err := conn.Select("u.id", "u.name", "u.avatar", "o.status").From("ppm_org_user_organization as o", "ppm_org_user as u").
		Where(condition).All(&userInfo)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := make([]inner_resp.SimpleInfo, 0)
	for _, user := range userInfo {
		res = append(res, inner_resp.SimpleInfo{
			Id:     user.Id,
			Name:   user.Name,
			Avatar: user.Avatar,
			Status: user.Status,
		})
	}

	return res, nil
}

func GetRepeatUserInfo(orgId int64) ([]inner_resp.RepeatMemberInfo, errs.SystemErrorInfo) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	var userInfo []po.PpmOrgUser
	err := conn.Select("u.id", "u.name").From("ppm_org_user_organization as o", "ppm_org_user as u").Where(db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcId:       db.Raw("o." + consts.TcUserId),
	}).All(&userInfo)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	//整合姓名id
	nameMap := make(map[string][]int64, 0)
	for _, user := range userInfo {
		nameMap[user.Name] = append(nameMap[user.Name], user.Id)
	}

	//找出重复的人员
	userIds := make([]int64, 0)
	for _, int64s := range nameMap {
		if len(int64s) <= 1 {
			continue
		}

		for _, i2 := range int64s {
			userIds = append(userIds, i2)
		}
	}

	if len(userIds) == 0 {
		return []inner_resp.RepeatMemberInfo{}, nil
	}

	//查询部门/职级信息
	userDeptBindInfoList, dbErr := GetUserDeptBindInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	userDepartmentMap := map[int64][]string{}
	for _, bindBo := range userDeptBindInfoList {
		userDepartmentMap[bindBo.UserId] = append(userDepartmentMap[bindBo.UserId], fmt.Sprintf("%s/%s", bindBo.DepartmentName, bindBo.PositionName))
	}

	res := make([]inner_resp.RepeatMemberInfo, 0)
	for name, s := range nameMap {
		if len(s) <= 1 {
			continue
		}
		for _, i2 := range s {
			temp := inner_resp.RepeatMemberInfo{
				Id:         i2,
				Name:       name,
				Department: []string{},
			}
			if depts, ok := userDepartmentMap[i2]; ok {
				temp.Department = depts
			}
			res = append(res, temp)
		}
	}

	return res, nil
}
