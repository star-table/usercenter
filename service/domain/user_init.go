package domain

import (
	"strconv"
	"strings"

	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/pinyin"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

func initDepartmentUser(userNativeId int64, userDetail *sdk.GetUserDetailResp, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	//暂不考虑后期增加了部门导致的部门不存在
	//1.获取外部部门与内部部门的id对应关系
	deptRelationList, err := GetOutDeptAndInnerDept(orgId, tx)
	if err != nil {
		return err
	}
	IsLeaderInDepts := dealLeaderString(userDetail.IsLeaderInDepts)

	//该用户为负责人的所有部门
	leaderDepts := make([]string, 0)
	for k, v := range IsLeaderInDepts {
		if v == "true" {
			leaderDepts = append(leaderDepts, k)
		}
	}
	departmentUserList := make([]interface{}, 0)
	for _, v := range userDetail.Department {
		deptId, ok := deptRelationList[strconv.FormatInt(v, 10)]
		if !ok {
			continue
		}

		id := snowflake.Id()
		isLeader := 2 //默认不是部门负责人
		if bol, _ := slice.Contain(leaderDepts, strconv.FormatInt(v, 10)); bol {
			isLeader = 1
		}
		departmentUserList = append(departmentUserList, po.PpmOrgUserDepartment{
			Id:           id,
			OrgId:        orgId,
			UserId:       userNativeId,
			DepartmentId: deptId,
			IsLeader:     isLeader,
		})
	}
	if len(deptRelationList) == 0 {
		return nil
	}
	dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUserList)
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

func dealLeaderString(str string) map[string]string {
	res := map[string]string{}
	str = strings.ReplaceAll(str, " ", "")
	str = str[1 : len(str)-1]
	strArr := strings.Split(str, ",")
	if len(strArr) == 0 {
		return res
	}
	for _, v := range strArr {
		kv := strings.Split(v, ":")
		if len(kv) >= 2 {
			res[kv[0]] = kv[1]
		}
	}

	return res
}

//组装用户信息
func assemblyUserInfo(userOutInfo *po.PpmOrgUserOutInfo, user *po.PpmOrgUser) (int64, int64, int64, errs.SystemErrorInfo) {

	userOutInfoId := snowflake.Id()
	userNativeId := snowflake.Id()
	userConfigId := snowflake.Id()

	return userOutInfoId, userNativeId, userConfigId, nil
}

//恢复注销用户
func restoreCancellateUser(userOutInfo *po.PpmOrgUserOutInfo, userDetail *sdk.GetUserDetailResp, tx sqlbuilder.Tx, err errs.SystemErrorInfo) (int64, errs.SystemErrorInfo) {
	logger.Info("user update ")
	//用户之前被注销掉，恢复状态
	if userOutInfo.IsDelete == consts.AppIsDeleted {
		userOutInfo.IsDelete = consts.AppIsNoDelete
		userOutInfo.SourceChannel = consts.AppSourceChannelDingTalk
		AssemblyUserOutInfo(userOutInfo, *userDetail)

		dbErr := store.Mysql.TransUpdate(tx, userOutInfo)
		if dbErr != nil {
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
	}
	return userOutInfo.Id, nil
}

func assemblyRegisteredInfo(registered bool, orgId, userConfigId int64, userNativeId *int64, user *po.PpmOrgUser, userDetail *sdk.GetUserDetailResp,
	userConfig *po.PpmOrgUserConfig, tx sqlbuilder.Tx, registeredUserOutInfo *po.PpmOrgUserOutInfo) errs.SystemErrorInfo {
	if !registered {
		user.Id = *userNativeId
		user.OrgId = orgId
		user.SourceChannel = consts.AppSourceChannelDingTalk
		AssemblyUser(user, *userDetail)

		dbErr := store.Mysql.TransInsert(tx, user)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.MysqlOperateError
		}
		logger.Info("初始化用户成功")

		userConfig.Id = userConfigId
		userConfig.OrgId = orgId
		userConfig.UserId = *userNativeId

		dbErr = store.Mysql.TransInsert(tx, userConfig)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.MysqlOperateError
		}
		logger.Info("初始化用户配置成功")
	} else {
		*userNativeId = registeredUserOutInfo.Id
	}

	return nil
}

func AssemblyUserOutInfo(userOutInfo *po.PpmOrgUserOutInfo, userDetailResp sdk.GetUserDetailResp) {
	isActive := 0
	if userDetailResp.Active {
		isActive = 1
	}

	userOutInfo.Name = userDetailResp.Name
	userOutInfo.Avatar = userDetailResp.Avatar
	userOutInfo.IsActive = isActive
	userOutInfo.JobNumber = userDetailResp.JobNumber
}

func AssemblyUser(user *po.PpmOrgUser, userDetailResp sdk.GetUserDetailResp) {
	user.Name = userDetailResp.Name
	user.NamePinyin = pinyin.ConvertToPinyin(user.Name)
	user.LoginName = user.Name
	if user.Avatar == "" {
		user.Avatar = userDetailResp.Avatar
	}
}
