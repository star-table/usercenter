package service

import (
	"strings"
	"time"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/date"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/pinyin"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/temp"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func PersonalInfo(orgId, userId int64, sourceChannel string) (*resp.PersonalInfo, errs.SystemErrorInfo) {
	userInfoBo, err := domain.GetOrgUserInfo(orgId, userId, sourceChannel)
	if err != nil {
		logger.Error(err)
		return nil, errs.GetUserInfoError
	}

	var personalInfo resp.PersonalInfo
	_ = copyer.Copy(userInfoBo, &personalInfo)
	// 是否设置过密码
	if userInfoBo.Password != "" {
		personalInfo.PasswordSet = 1
	}
	// 是否为外部协作人
	personalInfo.IsOutCollaborator = userInfoBo.UserType == 2
	// 组织拥有者
	personalInfo.IsOrgOwner = userInfoBo.OrgOwner == userId
	// 管理组权限
	group, dbErr := domain.GetManageGroupListByUser(orgId, userId)
	if dbErr != nil {
		// 不在管理组则忽视
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}
	if group != nil {
		if group.LangCode == consts.ManageGroupSys {
			personalInfo.IsSysAdmin = true
		} else {
			personalInfo.IsSubAdmin = true
		}
	}

	bos, dbErr := domain.GetFunctionConfig([]int64{orgId})
	if dbErr != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := make([]string, 0)
	for _, bo := range bos {
		if ok, _ := slice.Contain(consts.SupportFunctionCode, bo.FunctionCode); ok {
			res = append(res, bo.FunctionCode)
		}
	}
	personalInfo.Functions = res

	//查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptPositionDataList := make([]resp.UserDeptPositionData, 0)
	for _, info := range userDeptBindInfoList {
		deptPositionDataList = append(deptPositionDataList, resp.UserDeptPositionData{
			DepartmentId:  info.DepartmentId,
			IsLeader:      info.IsLeader,
			DeparmentName: info.DepartmentName,
			PositionId:    info.OrgPositionId,
			PositionName:  info.PositionName,
			PositionLevel: info.PositionLevel,
		})
	}

	// 查询角色信息
	userRoleBindInfoList, dbErr := domain.GetUserRoleBindListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userRoleDataList := make([]resp.UserRoleData, 0)
	for _, user := range userRoleBindInfoList {
		userRoleDataList = append(userRoleDataList, resp.UserRoleData{
			RoleId:   user.RoleId,
			RoleName: user.RoleName,
		})
	}

	personalInfo.DepartmentList = deptPositionDataList
	personalInfo.RoleList = userRoleDataList

	return &personalInfo, nil
}

// UpdateCurrentUserInfo 更新当前用户信息
func UpdateCurrentUserInfo(orgId, userId int64, input req.UpdateUserInfoReq) (bool, errs.SystemErrorInfo) {

	upd := &mysql.Upd{}
	//头像
	assemblyAvatar(input, upd)
	//姓名
	err := assemblyName(input, upd)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	//出生日期
	assemblyBirthday(input, upd)
	//性别
	err = assemblySex(input, upd)

	if err != nil {
		logger.Error(err)
		return false, err
	}

	dbErr := store.Mysql.UpdateSmart(consts.TableUser, userId, *upd)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	err = domain.ClearBaseUserInfo(orgId, userId)
	if err != nil {
		logger.Error(err)
	}

	return true, nil
}

func assemblySex(input req.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "sex") {

		if input.Sex != nil {

			if *input.Sex != consts.Male && *input.Sex != consts.Female {
				return errs.UserSexFail
			}
			(*upd)[consts.TcSex] = *input.Sex
		}
	}
	return nil
}

func assemblyBirthday(input req.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "birthday") {

		if input.Birthday != nil {
			birthday := time.Time(*input.Birthday)
			(*upd)[consts.TcBirthday] = birthday
		}
	}
}

//组装个人头像信息
func assemblyAvatar(input req.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "avatar") {

		if input.Avatar != nil {
			(*upd)[consts.TcAvatar] = *input.Avatar
		}
	}
}

//组装名字
func assemblyName(input req.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "name") {

		if input.Name != nil {
			name := strings.Trim(*input.Name, " ")
			isNameRight := format.VerifyNicknameFormat(name)
			if !isNameRight {
				return errs.NicknameLenError
			}
			(*upd)[consts.TcName] = name
			(*upd)[consts.TcNamePinyin] = pinyin.ConvertToPinyin(name)
		}
	}
	return nil
}

func SearchUser(orgId int64, email string) (*resp.SearchUserResp, errs.SystemErrorInfo) {
	if email == "" {
		return nil, errs.EmailFormatErr
	}
	//查看邮箱是否已注册
	var info po.PpmOrgUser
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	dbErr = conn.Select(db.Raw("u.*")).From("ppm_org_user u", "ppm_org_user_organization o").Where(db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"u." + consts.TcEmail:    email,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}).One(&info)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			//如果没有的话就去判断是否邀请过
			isInvited, dbErr := store.Mysql.IsExistByCond(consts.TableUserInvite, db.Cond{
				consts.TcIsDelete:       consts.AppIsNoDelete,
				consts.TcEmail:          email,
				consts.TcIsRegister:     2,
				consts.TcOrgId:          orgId,
				consts.TcLastInviteTime: db.Gte(date.Format(time.Now().AddDate(0, 0, -1))),
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return nil, errs.MysqlOperateError
			}
			if !isInvited {
				return &resp.SearchUserResp{
					Status:   1, //可邀请
					UserInfo: nil,
				}, nil
			} else {
				return &resp.SearchUserResp{
					Status:   2, //已邀请
					UserInfo: nil,
				}, nil
			}
		} else {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
	}

	return &resp.SearchUserResp{
		Status: 3, //已注册
		UserInfo: &resp.OrgMemberInfoReq{
			UserID: info.Id,
			Name:   info.Name,
			NamePy: info.NamePinyin,
			Avatar: info.Avatar,
		},
	}, nil
}

func InviteUserList(orgId int64, listReq req.InviteUserListReq) (*resp.InviteUserListResp, errs.SystemErrorInfo) {
	var userInvites []po.PpmOrgUserInvite
	total, dnErr := store.Mysql.SelectAllByCondWithPageAndOrder(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	}, nil, listReq.Page, listReq.Size, "create_time desc", &userInvites)
	if dnErr != nil {
		logger.Error(dnErr)
		return nil, errs.MysqlOperateError
	}

	var resInfo []resp.InviteUserInfo
	for _, invite := range userInvites {
		resInfo = append(resInfo, resp.InviteUserInfo{
			Id:              invite.Id,
			Name:            invite.Name,
			Email:           invite.Email,
			InviteTime:      invite.LastInviteTime,
			IsInvitedRecent: invite.LastInviteTime.Before(time.Now().AddDate(0, 0, -1)),
		})
	}

	return &resp.InviteUserListResp{
		Total: int64(total),
		List:  resInfo,
	}, nil
}

// InviteUser
func InviteUser(orgId, userId int64, param req.InviteUserReq) (*resp.InviteUserResp, errs.SystemErrorInfo) {
	result := &resp.InviteUserResp{}
	if len(param.Data) == 0 {
		return result, nil
	}

	validEmail := make([]string, 0)
	for _, datum := range param.Data {
		if !format.VerifyEmailFormat(datum.Email) {
			result.InvalidEmail = append(result.InvalidEmail, datum.Email)
		} else {
			validEmail = append(validEmail, datum.Email)
		}
	}
	if len(validEmail) == 0 {
		return result, nil
	}

	validEmail = slice.SliceUniqueString(validEmail)
	emailNameMap := map[string]string{}
	for _, datum := range param.Data {
		emailNameMap[datum.Email] = datum.Name
	}

	//查看是否是用户
	var userInfos []po.PpmOrgUser
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	dbErr = conn.Select(db.Raw("u.*")).From("ppm_org_user u", "ppm_org_user_organization o").Where(db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"u." + consts.TcEmail:    db.In(validEmail),
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}).All(&userInfos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	if len(userInfos) > 0 {
		for _, user := range userInfos {
			result.IsUserEmail = append(result.IsUserEmail, user.Email)
		}
	}

	notUserEmail := make([]string, 0)
	for _, s := range validEmail {
		if ok, _ := slice.Contain(result.IsUserEmail, s); !ok {
			notUserEmail = append(notUserEmail, s)
		}
	}
	if len(notUserEmail) == 0 {
		return result, nil
	}

	//查看是否已邀请
	var inviteInfo []po.PpmOrgUserInvite
	dbErr = store.Mysql.SelectAllByCond(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcEmail:      db.In(notUserEmail),
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	}, &inviteInfo)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	needInviteAgain := make([]int64, 0)
	if len(inviteInfo) > 0 {
		for _, invite := range inviteInfo {
			if invite.LastInviteTime.Before(time.Now().AddDate(0, 0, -1)) {
				//已邀请（需要再次邀请,更新数据库）
				needInviteAgain = append(needInviteAgain, invite.Id)
			} else {
				result.InvitedEmail = append(result.InvitedEmail, invite.Email)
			}
		}
	}
	for _, s := range notUserEmail {
		if ok, _ := slice.Contain(result.InvitedEmail, s); !ok {
			result.SuccessEmail = append(result.SuccessEmail, s)
		}
	}

	if len(result.SuccessEmail) == 0 {
		return result, nil
	}
	codeResp, err := GetInviteCode(userId, orgId, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = SendMailRelaxed(result.SuccessEmail, consts.MailTemplateSubjectInvite, temp.RenderIgnoreError(consts.MailTemplateContentInvite, map[string]string{
		consts.SMSParamsNameInviteHref: conf.Cfg.Application.Domain + "/user/entry?inviteCode=" + codeResp.InviteCode,
		consts.SMSParamsNameInviteUrl:  conf.Cfg.Application.Domain + "/user/entry?inviteCode=" + codeResp.InviteCode,
	}))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		now := time.Now()
		if len(needInviteAgain) > 0 {
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserInvite, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcId:    db.In(needInviteAgain),
			}, mysql.Upd{
				consts.TcUpdator:        userId,
				consts.TcLastInviteTime: now,
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		userInvites := make([]interface{}, 0)
		for _, s := range result.SuccessEmail {
			var name string
			if _, ok := emailNameMap[s]; ok {
				name = emailNameMap[s]
			}
			userInvites = append(userInvites, po.PpmOrgUserInvite{
				Id:             snowflake.Id(),
				OrgId:          orgId,
				Name:           name,
				Email:          s,
				InviteUserId:   userId,
				LastInviteTime: now,
				Creator:        userId,
				Updator:        userId,
			})
		}

		dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserInvite{}, userInvites)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	return result, nil
}
