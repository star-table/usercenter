package remote

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req/remote"
)

const (
	CloudreveServerName = "cloudreve"
)

// SynAddUserToCloudreve 同步新增用户给网盘
func SynAddUserToCloudreve(user bo.OrgMemberBaseInfoBo) errs.SystemErrorInfo {
	userStr := json.ToJsonIgnoreError(remote.UserList{
		UserList: []remote.UserData{
			{
				Id:          user.UserId,
				LoginName:   user.LoginName,
				Name:        user.Name,
				NamePy:      user.NamePinyin,
				Avatar:      user.Avatar,
				Email:       user.Email,
				PhoneRegion: user.MobileRegion,
				PhoneNumber: user.Mobile,
				Status:      user.Status,
				Creator:     user.Creator,
				CreateTime:  user.CreateTime,
				Updator:     user.Updator,
				UpdateTime:  user.UpdateTime,
			},
		},
	})
	bodyStr, code, err := nacos.DoPost(CloudreveServerName, "cloudreve/api/v3/syncData/userAdd", map[string]interface{}{}, userStr)
	if err != nil {
		logger.ErrorF("[同步新增用户至云盘]失败 -> code:%d err:%s", code, err.Error())
		return errs.SynAddUserToCloudreveErr
	}
	logger.InfoF("[同步新增用户至云盘]成功 -> bodyStr:%s", bodyStr)
	return nil
}

// SynEditUserToCloudreve 同步修改用户给网盘
func SynEditUserToCloudreve(users []bo.OrgMemberBaseInfoBo) errs.SystemErrorInfo {
	userList := make([]remote.UserData, 0)
	for _, user := range users {
		userList = append(userList, remote.UserData{
			Id:          user.UserId,
			LoginName:   user.LoginName,
			Name:        user.Name,
			NamePy:      user.NamePinyin,
			Avatar:      user.Avatar,
			Email:       user.Email,
			PhoneRegion: user.MobileRegion,
			PhoneNumber: user.Mobile,
			Status:      user.Status,
			Creator:     user.Creator,
			CreateTime:  user.CreateTime,
			Updator:     user.Updator,
			UpdateTime:  user.UpdateTime,
		})
	}
	userStr := json.ToJsonIgnoreError(remote.UserList{
		UserList: userList,
	})

	bodyStr, code, err := nacos.DoPost(CloudreveServerName, "cloudreve/api/v3/syncData/userUpdate", map[string]interface{}{}, userStr)
	if err != nil {
		logger.ErrorF("[同步修改用户至云盘]失败 -> code:%d err:%s", code, err.Error())
		return errs.SynEditUserToCloudreveErr
	}
	logger.InfoF("[同步修改用户至成员]成功 -> bodyStr:%s", bodyStr)
	return nil
}

// SynRemoveUserToCloudreve 同步删除用户给网盘
func SynRemoveUserToCloudreve(userIds []int64) errs.SystemErrorInfo {
	userIdStr := json.ToJsonIgnoreError(map[string][]int64{"ids": userIds})
	bodyStr, code, err := nacos.DoPost(CloudreveServerName, "cloudreve/api/v3/syncData/userDelete", map[string]interface{}{}, userIdStr)
	if err != nil {
		logger.ErrorF("[同步删除用户至云盘]失败 -> code:%d err:%s", code, err.Error())
		return errs.SynRemoveUserToCloudreveErr
	}
	logger.InfoF("[同步删除至云盘]成功 -> bodyStr:%s", bodyStr)
	return nil
}
