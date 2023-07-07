package service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
)

func GetBaseOrgInfo(sourceChannel string, orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfo(sourceChannel, orgId)
}

func GetDingTalkBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetDingTalkBaseUserInfoByEmpId(orgId, empId)
}

func GetBaseUserInfoByEmpId(sourceChannel string, orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoByEmpId(sourceChannel, orgId, empId)
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	return domain.GetUserConfigInfo(orgId, userId)
}

func GetBaseUserInfo(sourceChannel string, orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfo(sourceChannel, orgId, userId)
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetDingTalkBaseUserInfo(orgId, userId)
}

func GetBaseUserInfoBatch(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoBatch(sourceChannel, orgId, userIds)
}
