package consts

import (
	"math/rand"
	"time"
)

const CacheKeyPrefix = "polaris:"
const CacheKeyOfSys = "sys:"
const CacheKeyOfOrg = "org_{{." + CacheKeyOrgIdConstName + "}}:"
const CacheKeyOfOutOrg = "outorg_{{." + CacheKeyOutOrgIdConstName + "}}:"
const CacheKeyOfUser = "user_{{." + CacheKeyUserIdConstName + "}}:"
const CacheKeyOfOutUser = "outuser_{{." + CacheKeyOutUserIdConstName + "}}:"
const CacheKeyOfProject = "project_{{." + CacheKeyProjectIdConstName + "}}:"
const CacheKeyOfProcess = "process_{{." + CacheKeyProcessIdConstName + "}}:"
const CacheKeyOfRole = "role_{{." + CacheKeyRoleIdConstName + "}}:"
const CacheKeyOfSourceChannel = "source_channel_{{." + CacheKeySourceChannelConstName + "}}:"
const CacheKeyOfPhone = "phone_{{." + CacheKeyPhoneConstName + "}}:"
const CacheKeyOfAuthType = "authType_{{." + CacheKeyAuthTypeConstName + "}}:"
const CacheKeyOfAddressType = "addressType_{{." + CacheKeyAddressTypeConstName + "}}:"
const CacheKeyOfLoginName = "login_name_{{." + CacheKeyLoginNameConstName + "}}:"
const CacheKeyOfRoleGroup = "group_{{." + CacheKeyRoleGroupConstName + "}}"
const CacheKeyOfChangeBind = "user_id_{{." + CacheKeyUserIdConstName + "}}"

//服务名
const (
	AppsvcApplicationName      = "appsvc:"
	IdsvcApplicationName       = "idsvc:"
	MsgsvcApplicationName      = "msgsvc:"
	CallsvcApplicationName     = "callsvc:"
	OrgsvcApplicationName      = "orgsvc:"
	ProcesssvcApplicationName  = "processsvc:"
	ProjectsvcApplicationName  = "projectsvc:"
	ResourcesvcApplicationName = "resourcesvc:"
	RolesvcApplicationName     = "rolesvc:"
	RrendssvcApplicationName   = "trendssvc:"
	SchedulesvcApplicationName = "scheduletsvc:"
	CommonsvcApplicationName   = "commonsvc"
)

//失效时间
const (
	//用户Token失效时间: 15天
	CacheUserTokenExpire = 60 * 60 * 24 * 15
	//通用失效时间: 3小时
	CacheBaseExpire = int64(60 * 60 * 3)
	//用户信息缓存: 1小时
	CacheBaseUserInfoExpire = int64(60 * 60 * 1)
)

func GetCacheBaseExpire() int64 {
	rand.Seed(time.Now().Unix())
	return CacheBaseExpire + rand.Int63n(30)
}
func GetCacheBaseUserInfoExpire() int64 {
	rand.Seed(time.Now().Unix())
	return CacheBaseUserInfoExpire + rand.Int63n(30)
}

const (
	CacheKeyOrgIdConstName         = "orgId"
	CacheKeyUserIdConstName        = "userId"
	CacheKeyOutOrgIdConstName      = "outOrgId"
	CacheKeyOutUserIdConstName     = "outUserId"
	CacheKeyProjectIdConstName     = "projectId"
	CacheKeyFolderNameConstName    = "folderName"
	CacheKeyIssueIdConstName       = "issueId"
	CacheKeyObjectCodeConstName    = "objectCode"
	CacheKeyProcessIdConstName     = "processId"
	CacheKeyRoleIdConstName        = "roleId"
	CacheKeySourceChannelConstName = "sourceChannel"
	CacheKeyYearConstName          = "year"
	CacheKeyMonthConstName         = "month"
	CacheKeyDayConstName           = "day"
	CacheKeyPhoneConstName         = "phone"
	CacheKeyAuthTypeConstName      = "authType"
	CacheKeyAddressTypeConstName   = "addressType"
	CacheKeyLoginNameConstName     = "loginName"
	CacheKeyRoleGroupConstName     = "roleGroup"
)

//系统缓存
const (
	//DingTalk Suite Ticket
	CacheDingTalkSuiteTicket = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelDingTalk + ":suite_ticket"
	//飞书 AppTicket
	CacheFeiShuAppTicket = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelFeiShu + ":app_ticket"
	//飞书 AppAccessToken
	CacheFeiShuAppAccessToken = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelFeiShu + ":app_access_token:"
	//飞书 TenantAccessToken
	CacheFeiShuTenantAccessToken = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelFeiShu + ":tenant_access_token:"
	//飞书 授权范围
	CacheFeiShuScope = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelFeiShu + ":scope:"

	//mqtt root key
	CacheMQTTRootKey = CacheKeyPrefix + CacheKeyOfSys + "mqtt:root_key"

	//飞书 卡片通知回调消息refresh-token, 网络抖动等极端情况下，会出现卡片点击失败但是业务方已经处理过 action 的现象，所以业务方接口存在被重复调用的风险。X-Refresh-Token 只有在卡片点击事件成功被响应并通知到客户端的时候才会刷新，如果业务方的接口非幂等，可以通过缓存并验证该字段防止接口被重复调用。
	CacheFeiShuCardCallBackRefreshToken = CacheKeyPrefix + CacheKeyOfSys + AppSourceChannelFeiShu + ":card_call_back:refresh_token:"

	//用户token
	//CacheUserToken = CacheKeyPrefix + CacheKeyOfSys + "user:token:"
	////对象id缓存key前缀
	//CacheObjectIdPreKey = CacheKeyPrefix + CacheKeyOfSys + "object_id:"
	// 角色操作列表
	//CacheRoleOperationList = CacheKeyPrefix + CacheKeyOfSys + "role_operation_list"
	////部门对应关系
	//CacheDeptRelation = CacheKeyPrefix + CacheKeyOfSys + "dept_relation_list"

	//用户配置缓存
	CacheUserConfig = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + "config"
	//用户基础信息缓存key
	CacheBaseUserInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + "baseinfo"
	//用户外部信息缓存key
	CacheBaseUserOutInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + "outinfo"

	//组织基础信息
	CacheBaseOrgInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + "baseinfo2022"
	//CacheBaseOrgInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + "baseinfo"
	//组织外部信息
	CacheBaseOrgOutInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + "outinfo"
	//获取外部组织id关联的内部组织id
	CacheOutOrgIdRelationId = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOutOrg + CacheKeyOfSourceChannel + "org_id_info"

	//获取外部用户id关联的内部用户id
	CacheOutUserIdRelationId = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfSourceChannel + CacheKeyOfOutUser + "user_id"
	//部门对应关系
	CacheDeptRelation = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + "dept_relation_list"
	//用户token
	CacheUserToken = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + "user:token:"
	//fs用户refresh_token和user_access_token
	CacheFsUserAccessToken = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + "fstoken"

	//用户邀请code, 拼接 inviteCode
	CacheUserInviteCode = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + "invite_code:"
	//用户邀请code有效时间为: 24小时
	CacheUserInviteCodeExpire = 60 * 60 * 24

	// 短信验证码相关 + 手机号, 验证失败五次间隔调整为五分钟，验证失败50次冻结一天
	// 短信发送时间间隔: 一分钟，五分钟
	CacheSmsSendLoginCodeFreezeTime = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + CacheKeyOfAuthType + CacheKeyOfAddressType + CacheKeyOfPhone + "sms_auth_code:freeze_time"
	// 短信验证码: 五分钟
	CacheSmsLoginCode = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + CacheKeyOfAuthType + CacheKeyOfAddressType + CacheKeyOfPhone + "sms_auth_code"
	// 号码白名单
	CachePhoneNumberWhiteList = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + "sms_white_list"
	// 登录短信验证失败次数
	CacheSmsLoginCodeVerifyFailTimes = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + CacheKeyOfAuthType + CacheKeyOfAddressType + CacheKeyOfPhone + "sms_auth_code:verify_times"
	//登录图形验证码：一分钟
	CacheLoginGraphCode = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSys + CacheKeyOfLoginName + "graph_auth_code"
	//换绑账号行为记录：5分钟
	ChangeLoginNameSign = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + CacheKeyOfAddressType + "change_login_name_sign"
	// 邮箱、手机号换绑时验证码校验的存储
	ChangeBindCode = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + CacheKeyOfAddressType + "change_bind_code"
	//钉钉扫码登录：5分钟
	LoginByDingCode = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfSourceChannel + CacheKeyOfOutUser + "login_by_ding_code"
)
