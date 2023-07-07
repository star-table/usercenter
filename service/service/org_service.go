package service

import (
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/pinyin"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/vo"
	"github.com/star-table/usercenter/service/model/vo/orgvo"
	"upper.io/db.v3"
)

func GetBaseOrgInfoByOutOrgId(sourceChannel string, outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfoByOutOrgId(sourceChannel, outOrgId)
}

// CreateOrg 创建组织
func CreateOrg(reqParam orgvo.CreateOrgReqVo, sourceChannel, sourcePlatform string) (int64, errs.SystemErrorInfo) {
	logger.InfoF("[创建组织] -> reqParam: %s", json.ToJsonIgnoreError(reqParam))
	creatorId := reqParam.Data.CreatorId
	createReqInfo := reqParam.Data.CreateOrgReq
	if createReqInfo.CreatorName != nil {
		if !format.VerifyNicknameFormat(strings.Trim(*(createReqInfo.CreatorName), " ")) {
			return 0, errs.BuildSystemErrorInfo(errs.NicknameLenError)
		}
	}
	user, dbErr := domain.GetUserPoById(creatorId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return 0, errs.UserNotExist
		}
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}
	orgId, err := domain.CreateOrg(bo.CreateOrgBo{
		OrgName:    createReqInfo.OrgName,
		IndustryId: createReqInfo.IndustryId,
		Scale:      createReqInfo.Scale,
	}, creatorId, sourceChannel, sourcePlatform, "")
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	// 绑定用户和组织关联
	_, dbErr = domain.BindUserOrgRelation(orgId, creatorId, true, false, false, 0)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	upd := mysql.Upd{}
	// 更新用户的orgId
	if user.OrgId == 0 {
		upd[consts.TcOrgId] = orgId
	}
	// 更新用户的名称
	if createReqInfo.CreatorName != nil {
		creatorName := strings.Trim(*(createReqInfo.CreatorName), " ")

		upd[consts.TcName] = creatorName
		upd[consts.TcNamePinyin] = pinyin.ConvertToPinyin(creatorName)
	}

	// 更新用户的信息
	if len(upd) > 0 {
		upd[consts.TcUpdator] = creatorId
		upd[consts.TcVersion] = db.Raw("version + 1")
		dbErr := store.Mysql.UpdateSmart(consts.TableUser, creatorId, upd)
		if dbErr != nil {
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
	}

	//刷新用户缓存
	userToken := reqParam.Data.UserToken
	err = UpdateCacheUserInfoOrgId(userToken, orgId)
	if err != nil {
		logger.Error(err)
	}

	err = domain.ClearBaseUserInfo(orgId, creatorId)
	if err != nil {
		logger.Error(err)
	}

	return orgId, nil
}

// GetUserParticipateOrgList 用户绑定过的组织列表
func GetUserParticipateOrgList(userId int64) (*resp.UserOrganizationListResp, errs.SystemErrorInfo) {

	// 获取组织成员列表
	userOrgBindList, dbErr := domain.GetOrgMemberListByUser(userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	orgIds := make([]int64, 0)
	enableOrgIdMap := make(map[int64]bool)
	for _, organization := range userOrgBindList {
		orgIds = append(orgIds, organization.OrgId)
		if organization.Status == consts.AppStatusEnable {
			enableOrgIdMap[organization.OrgId] = true
		}
	}

	// 获取组织信息
	orgIds = slice.SliceUniqueInt64(orgIds)
	orgList, dbErr := domain.GetOrgListByIds(orgIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	functionInfo, dbErr := domain.GetFunctionConfig(orgIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	functionMap := map[int64][]string{}
	for _, config := range functionInfo {
		if ok, _ := slice.Contain(consts.SupportFunctionCode, config.FunctionCode); ok {
			functionMap[config.OrgId] = append(functionMap[config.OrgId], config.FunctionCode)
		}
	}

	resultList := make([]resp.UserOrganization, 0)

	for _, organization := range orgList {
		var userOrg resp.UserOrganization
		_ = copyer.Copy(organization, &userOrg)
		if ok, _ := enableOrgIdMap[userOrg.ID]; ok {
			userOrg.OrgIsEnabled = consts.AppStatusEnable
		} else {
			userOrg.OrgIsEnabled = consts.AppStatusDisabled
		}
		if functions, ok := functionMap[userOrg.ID]; ok {
			userOrg.Functions = functions
		} else {
			userOrg.Functions = []string{}
		}
		resultList = append(resultList, userOrg)
	}

	return &resp.UserOrganizationListResp{
		List: resultList,
	}, nil
}

// SwitchUserOrganization 切换默认组织
func SwitchUserOrganization(orgId, userId int64, token string) (bool, errs.SystemErrorInfo) {
	//监测可用性
	baseUserInfo, err := GetBaseUserInfo("", orgId, userId)
	if err != nil {
		return false, err
	}

	err = BaseUserInfoOrgStatusCheck(*baseUserInfo)
	if err != nil {
		return false, err
	}

	//更改用户缓存的orgId
	err = UpdateCacheUserInfoOrgId(token, orgId)
	if err != nil {
		return false, err
	}

	//修改用户默认组织
	dbErr := domain.UpdateUserDefaultOrg(userId, orgId)
	if dbErr != nil {
		logger.ErrorF("[修改默认组织] -> 失败,  dbErr: %s", dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

func SwitchToInnerMember(orgId, userId int64, reqVo req.SwitchToInnerMemberReq) (bool, errs.SystemErrorInfo) {
	// todo 判断权限
	_, err := domain.UpdateOrgUserType(orgId, userId, reqVo.UserIds, consts.OrgUserTypeInner)
	if err != nil {
		logger.Error(err)
		return false, errs.MysqlOperateError
	}
	return true, nil
}

//获取组织信息
func OrganizationInfo(req orgvo.OrganizationInfoReqVo) (*vo.OrganizationInfoResp, errs.SystemErrorInfo) {
	orgInfo, dbErr := domain.GetOrgById(req.OrgId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	ownerInfo, err := domain.GetBaseUserInfo(orgInfo.SourceChannel, req.OrgId, orgInfo.Owner)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	ownerResp := vo.UserIDInfo{
		UserID:     orgInfo.Owner,
		Name:       ownerInfo.Name,
		Avatar:     ownerInfo.Avatar,
		EmplID:     ownerInfo.OutOrgUserId,
		IsDeleted:  ownerInfo.OrgUserIsDelete == consts.AppIsDeleted,
		IsDisabled: ownerInfo.OrgUserStatus == consts.AppStatusDisabled,
	}

	infoResp := vo.OrganizationInfoResp{
		OrgID:      orgInfo.Id,
		OrgName:    orgInfo.Name,
		Code:       orgInfo.Code,
		WebSite:    orgInfo.WebSite,
		IndustryID: orgInfo.IndustryId,
		Scale:      orgInfo.Scale,
		CountryID:  orgInfo.CountryId,
		ProvinceID: orgInfo.ProvinceId,
		CityID:     orgInfo.CityId,
		Address:    orgInfo.Address,
		LogoURL:    orgInfo.LogoUrl,
		Owner:      orgInfo.Owner,
		OwnerInfo:  &ownerResp,
	}

	return &infoResp, nil
}

//对于自己创建的组织，暂时不支持转让
//
//对于加入的企业只有查看全，无操作权
//
//暂时只做基本设置
func UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) (int64, errs.SystemErrorInfo) {

	input := req.Input
	// Owns转让的成员需要判断是否在这个组织里面 暂定
	orgInfo, dbErr := domain.GetOrgById(input.OrgID)

	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return 0, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}
	//不是所有者不可以更改信息
	if orgInfo.Owner != req.UserId {
		return 0, errs.OrgOwnTransferError
	}

	//更改Own的接口拆开来,  这一期暂时也不做
	updateOrgBo, err := assemblyOrganization(input, req.UserId, orgInfo)

	if err != nil {
		return 0, err
	}

	err = domain.UpdateOrg(*updateOrgBo)

	if err != nil {
		return 0, err
	}
	return input.OrgID, nil
}

func assemblyOwner(userId int64, input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "owner") && input.Owner != nil {
		orgId := input.OrgID
		orgInfo, dbErr := domain.GetOrgById(orgId)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				return errs.OrgNotExist
			}
			logger.Error(dbErr)
			return errs.MysqlOperateError
		}

		if orgInfo.Owner == *input.Owner {
			//组织所有者没有变动
			return nil
		}
		//不是所有者不可以更改信息
		if orgInfo.Owner != userId {
			return errs.OrgOwnTransferError
		}
		//查看新用户是否存在
		_, userErr := domain.GetBaseUserInfo(orgInfo.SourceChannel, orgId, *input.Owner)
		if userErr != nil {
			logger.Error(userErr)
			return userErr
		}

		//修改组织拥有者
		(*upd)[consts.TcOwner] = *input.Owner
	}

	return nil
}

func assemblyOrganization(input vo.UpdateOrganizationSettingsReq, userId int64, orgOrganization *po.PpmOrgOrganization) (*bo.UpdateOrganizationBo, errs.SystemErrorInfo) {
	//公用初始化
	orgBo := bo.OrganizationBo{Id: input.OrgID}

	upd := &mysql.Upd{}
	//名字
	err := assemblyOrgName(input, upd)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	//网址
	err = assemblyCode(input, upd)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	//行业
	assemblyIndustryID(input, upd)
	//组织规模
	assemblyScaleID(input, upd)
	//所在国家
	assemblyCountryID(input, upd)
	// 所在省份
	assemblyProvince(input, upd)
	// 所在城市
	assemblyCity(input, upd)
	// 组织地址
	err = assemblyAddress(input, upd)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// 组织logo地址
	err = assemblyLogoUrl(input, upd)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//组织所有者
	err = assemblyOwner(userId, input, upd)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//sourceChannel
	orgBo.SourceChannel = orgOrganization.SourceChannel

	return &bo.UpdateOrganizationBo{
		Bo:                     orgBo,
		OrganizationUpdateCond: *upd,
	}, nil
}

func assemblyOrgName(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "orgName") {
		orgName := strings.TrimSpace(input.OrgName)
		//orgNameLen := len(orgName)
		//if orgNameLen == 0 || orgNameLen > 256 {
		//	return errs.BuildSystemErrorInfo(errs.OrgNameLenError)
		//}
		isOrgNameRight := format.VerifyOrgNameFormat(orgName)
		if !isOrgNameRight {
			return errs.OrgNameLenError
		}

		(*upd)[consts.TcName] = input.OrgName
	}
	return nil
}

//网址
func assemblyCode(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "code") {
		if input.Code != nil {
			orgCode := *input.Code
			orgCode = strings.TrimSpace(orgCode)

			//判断当前组织有没有设置过code
			organizationBo, dbErr := domain.GetOrgById(input.OrgID)
			if dbErr != nil {
				if dbErr == db.ErrNoMoreRows {
					return errs.OrgNotExist
				}
				logger.Error(dbErr)
				return errs.MysqlOperateError
			}
			if organizationBo.Code != consts.BlankString {
				return errs.BuildSystemErrorInfo(errs.OrgCodeAlreadySetError)
			}

			//orgCodeLen := strs.Len(orgCode)
			////判断长度
			//if orgCodeLen > sconsts.OrgCodeLength || orgCodeLen < 1 {
			//	return errs.BuildSystemErrorInfo(errs.OrgCodeLenError)
			//}
			isOrgCodeRight := format.VerifyOrgCodeFormat(orgCode)
			if !isOrgCodeRight {
				return errs.OrgCodeLenError
			}

			_, dbErr = domain.GetOrgByCode(orgCode)
			//查不到才能更改
			if dbErr != nil {
				if dbErr == db.ErrNoMoreRows {
					(*upd)[consts.TcCode] = orgCode
				}
				logger.Error(dbErr)
				return errs.MysqlOperateError
			} else {
				return errs.OrgCodeExistError
			}
		}
	}

	return nil
}

//组织行业Id
func assemblyIndustryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "industryId") {

		if input.IndustryID != nil {
			(*upd)[consts.TcIndustryId] = *input.IndustryID
		} else {
			(*upd)[consts.TcIndustryId] = 0
		}
	}
}

//组织规模
func assemblyScaleID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "scale") {

		if input.Scale != nil {
			(*upd)[consts.TcScale] = *input.Scale
		} else {
			(*upd)[consts.TcScale] = 0
		}
	}
}

//所在国家
func assemblyCountryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "countryId") {

		if input.CountryID != nil {
			(*upd)[consts.TcCountryId] = *input.CountryID
		} else {
			(*upd)[consts.TcCountryId] = 0
		}
	}
}

//省份
func assemblyProvince(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "provinceId") {

		if input.ProvinceID != nil {
			(*upd)[consts.TcProvinceId] = *input.ProvinceID
		} else {
			(*upd)[consts.TcProvinceId] = 0
		}
	}
}

//城市
func assemblyCity(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "cityId") {

		if input.CityID != nil {
			(*upd)[consts.TcCityId] = *input.CityID
		} else {
			(*upd)[consts.TcCityId] = 0
		}
	}
}

//地址
func assemblyAddress(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "address") {

		if input.Address != nil {
			//len := strs.Len(*input.Address)
			//if len > 256 {
			//	return errs.BuildSystemErrorInfo(errs.OrgAddressLenError)
			//}
			isAdressRight := format.VerifyOrgAdressFormat(*input.Address)
			if !isAdressRight {
				return errs.OrgAddressLenError
			}

			(*upd)[consts.TcAddress] = *input.Address
		}
	}
	return nil
}

//Logo
func assemblyLogoUrl(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "logoUrl") {
		if input.LogoURL != nil {
			logoLen := strs.Len(*input.LogoURL)
			if logoLen > 512 {
				return errs.OrgLogoLenError
			}

			(*upd)[consts.TcLogoUrl] = *input.LogoURL
		}
	}
	return nil
}

// GetOrgConfig 获取组织配置
func GetOrgConfig(orgId int64) (*resp.OrgConfig, errs.SystemErrorInfo) {
	return domain.GetOrgConfig(orgId)
}

// GenerateApiKey 生成OpenApi
func GenerateApiKey(orgId, operatorUid int64) (bool, errs.SystemErrorInfo) {
	org, dbErr := domain.GetOrgById(orgId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	if org.ApiKey != "" {
		return false, errs.ApiKeyIsOpened
	}

	dbErr = domain.GenerateAndSetApiKey(org.Id, operatorUid)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	return true, nil
}

// ResetApiKey 重置ApiKey
func ResetApiKey(orgId, operatorUid int64) (bool, errs.SystemErrorInfo) {
	org, dbErr := domain.GetOrgById(orgId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	if org.ApiKey == "" {
		return false, errs.ApiKeyIsClosed
	}
	dbErr = domain.GenerateAndSetApiKey(org.Id, operatorUid)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	return true, nil
}

// RemoveApiKey 删除ApiKey
func RemoveApiKey(orgId, operatorUid int64) (bool, errs.SystemErrorInfo) {
	org, dbErr := domain.GetOrgById(orgId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	dbErr = domain.RemoveApiKey(org.Id, operatorUid)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	return true, nil
}
