package inner_service

import (
	"strings"

	"github.com/spf13/cast"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
)

func UniversalGetUserBaseInfoByIds(orgId int64, ids []string) ([]*inner_resp.UserBaseInfo, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}

	var users []*inner_resp.UserBaseInfo
	if len(ids) == 0 {
		return users, nil
	}

	ids = slice.SliceUniqueString(ids)

	var userIds []int64
	var deptIds []int64
	for _, id := range ids {
		if strings.HasPrefix(id, consts.MemberTypeUser) {
			userIds = append(userIds, cast.ToInt64(strings.TrimLeft(id, consts.MemberTypeUser)))
		} else if strings.HasPrefix(id, consts.MemberTypeDept) {
			deptIds = append(deptIds, cast.ToInt64(strings.TrimLeft(id, consts.MemberTypeDept)))
		} else {
			userIds = append(userIds, cast.ToInt64(id))
		}
	}

	// dept
	if len(deptIds) > 0 {
		// 获取部门树
		_, deptMap, err := domain.GetDeptTree(orgId)
		if err != nil {
			logger.ErrorF("[UniversalGetUserBaseInfoByIds] GetDeptTree failed, orgId: %v, err: %v", orgId, err)
			return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
		}

		// 获取部门以及子部门Ids
		allDeptIds := domain.GetAllDeptIdsWithChildren(deptMap, deptIds)

		// 获取所有用户信息
		logger.InfoF("[UniversalGetUserBaseInfoByIds] GetUserIdsByDeptIds, orgId: %v, deptIds: %v", orgId, allDeptIds)
		allDeptUserIds, err := domain.GetUserIdsByDeptIds(orgId, allDeptIds)
		if err != nil {
			logger.ErrorF("[UniversalGetUserBaseInfoByIds] GetUserIdsByDeptIds failed, orgId: %v, deptIds: %v, err: %v", orgId, allDeptIds, err)
			return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
		}
		logger.InfoF("[UniversalGetUserBaseInfoByIds] GetUserIdsByDeptIds, orgId: %v, deptIds: %v, result: %v", orgId, allDeptIds, allDeptUserIds)

		userIds = append(userIds, allDeptUserIds...)
	}

	// 去重
	userIds = slice.SliceUniqueInt64(userIds)

	logger.InfoF("[UniversalGetUserBaseInfoByIds] GetUserBaseInfos, orgId: %v, userIds: %v", orgId, userIds)
	if len(userIds) > 0 {
		userInfos, dbErr := domain.GetUserBaseInfos(orgId, userIds)
		if dbErr != nil {
			return nil, errs.MysqlOperateError
		}
		for _, userInfo := range userInfos {
			user := &inner_resp.UserBaseInfo{
				Id:            userInfo.Id,
				Name:          userInfo.Name,
				Avatar:        userInfo.Avatar,
				Email:         userInfo.Email,
				SourceChannel: userInfo.SourceChannel,
				OpenId:        userInfo.OpenId,
			}
			if userInfo.Mobile2 != nil {
				user.Mobile = *userInfo.Mobile2
			} else {
				user.Mobile = userInfo.Mobile1
			}
			users = append(users, user)
		}
	}
	return users, nil
}
