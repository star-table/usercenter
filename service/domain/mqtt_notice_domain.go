package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/mqtt"
	"github.com/star-table/usercenter/service/model/bo"
)

func PushAddOrgMemberNotice(orgId, depId int64, memberIds []int64, operatorId int64) {
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)

	baseUserInfos, err := GetBaseUserInfoBatch("", orgId, memberIds)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, baseUserInfo := range baseUserInfos {
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId: baseUserInfo.UserId,
			ObjectValue: bo.BaseUserInfoExtBo{
				BaseUserInfoBo: baseUserInfo,
				DepartmentId:   depId,
			},
		})
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId:         orgId,
		Action:        consts.MQTTDataRefreshActionAdd,
		Type:          consts.MQTTDataRefreshTypeMember,
		OperationId:   operatorId,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil {
		logger.Error(err)
	}
}

func PushRemoveOrgMemberNotice(orgId int64, memberIds []int64, operatorId int64) {
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)

	baseUserInfos, err := GetBaseUserInfoBatch("", orgId, memberIds)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, baseUserInfo := range baseUserInfos {
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId:    baseUserInfo.UserId,
			ObjectValue: baseUserInfo,
		})
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId:         orgId,
		Action:        consts.MQTTDataRefreshActionDel,
		Type:          consts.MQTTDataRefreshTypeMember,
		OperationId:   operatorId,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil {
		logger.Error(err)
	}
}
