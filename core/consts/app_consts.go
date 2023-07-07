package consts

import "time"

//Head Name
const (
	AppHeaderTokenName = "PM-TOKEN"
	AppHeaderOrgName   = "PM-ORG"
	AppHeaderProName   = "PM-PRO"
	AppHeaderEnvName   = "PM-ENV"
	AppHeaderPlatName  = "PM-PLAT"
	AppHeaderVerName   = "PM-VER"
)

//TraceName
const TraceServiceName = "PolarisTrace"
const JaegerContextTraceKey = "JaegerContextTraceKey"
const JaegerContextSpanKey = "JaegerContextSpanKey"

// linux操作系统
const LinuxGOOS = "linux"

//app服务的项目名
const AppApplicationName = "app"

//api版本（增加了飞书群聊功能）
const ApiVersionFsChat = "1.4.3"

//来源
const (
	AppSourceChannelDingTalk = "ding"
	AppSourceChannelFeiShu   = "fs"
	AppSourceChannelInvite   = "invite"

	AppSourceChannelDingTalkDefaultLang = "zh_CN"

	AppSourcePlatformLarkXYJH2019 = "lark-xyjh2019"

	AppSourceChannelWeb  = "web"
	AppSourceChannelLark = "fs"
)

//默认空时间
const BlankTime = "1970-01-01 00:00:00"
const BlankDate = "1970-01-01"
const BlankElasticityTime = "1970-01-02 00:00:00"

var BlankTimeObject, _ = time.Parse(AppTimeFormat, BlankTime)

//默认空字符串
const BlankString = ""

//是否被删除
const (
	AppIsDeleted  = 1
	AppIsNoDelete = 2
)

//审核状态,1待审核,2审核通过,3审核不过
const (
	AppCheckStatusWait    = 1
	AppCheckStatusSuccess = 2
	AppCheckStatusFail    = 3
)

const (
	OrgUserTypeInner = 1 //内部成员
	OrgUserTypeOut   = 2 //外部成员
)

//是否隐藏
const (
	AppIsHiding    = 1
	AppIsNotHiding = 2
)

//待联系状态
const (
	//待联系
	ContactStatusWait = 1
	//已联系
	ContactStatusCompleted = 2
	//联系失败
	ContactStatusFail = 3
)

const (
	AppUserIsInUse     = 1
	AppUserIsNotInUser = 2
)

//是否流程初始化状态
const (
	//是
	AppIsInitStatus = 1
	//否
	AppIsNotInitStatus = 2
)

//是否可用
const (
	AppStatusEnable   = 1
	AppStatusDisabled = 2
)

//人员状态
const (
	OrgUserStatusResigned = 3 // 离职
)

//是否展示
const (
	AppShowEnable   = 1
	AppShowDisabled = 2
)

//是否默认
const (
	APPIsDefault    = 1
	AppIsNotDefault = 2
)

//是否提醒
const (
	AppIsRemind    = 1
	AppIsNotRemind = 2
)

const (
	AppIsFilling    = 1
	AppIsNotFilling = 2
)

//全局日期格式
const AppDateFormat = "2006-01-02"
const AppTimeFormat = "2006-01-02 15:04:05"
const AppTimeFormatYYYYMMDDHHmm = "2006-01-02 15:04"
const AppTimeFormatYYYYMMDDHHmmTimezone = "2006-01-02 15:04 -0700"
const AppSystemTimeFormat = "2006-01-02T15:04:05Z"
const AppSystemTimeFormat8 = "2006-01-02T15:04:05+08:00"

const (
	// SAAS运行模式
	AppRunmodeSaas = 1
	// 单机部署
	AppRunmodeSingle = 2
	// 私有化部署
	AppRunmodePrivate = 3
	// 私有化单库部署
	AppRunmodePrivateSingleDb = 4
)

//初始化时的一些常量定义
const (
	InitDefaultTeamName     = "默认团队"
	InitDefaultTeamNickname = "默认团队昵称"
)

// context key
//const (
//	TraceIdKey     = "PM-TRACE-ID"
//	//TraceIdKey = "_traceId"
//	HttpContextKey = "_httpContext"
//)

// 默认对象id步长
const (
	DefaultObjectIdMaxId = int64(1000)
	DefaultObjectIdStep  = 500
)

// 系统缓存模式
const (
	CacheModeRedis  = "Redis"
	CacheModeInside = "Inside"
)

// 系统消息队列模式
const (
	MQModeRocketMQ = "RocketMQ"
	MQModeDB       = "DB"
	MQModeKafka    = "Kafka"
)

// 发送消息状态
const (
	SendMQStatusSuccess = 1
	SendMQStatusFail    = 2
)

// 消息处理状态
const (
	//待处理
	MQStatusWait = 1
	//处理中
	MQStatusHandle = 2
	//处理成功
	MQStatusSuccess = 3
	//处理失败
	MQStatusFail = 4
)

//SMS签名和模板定义
const (
	// 由于产品调整，以及短信平台参数的设定，签名值改为了“极星”
	SMSSignNameBeiJiXing = "极星"
	//登录验证码，需要code参数
	SMSTemplateCodeLoginAuthCode = "SMS_175533634"
	//绑定验证码，需要code参数
	SMSTemplateCodeBindAuthCode = "SMS_181500004"
	//解绑验证码，需要code参数
	SMSTemplateCodeUnBindAuthCode = "SMS_181490006"
	//找回密码验证码，需要code参数
	SMSTemplateCodeRetrievePwdAuthCode = "SMS_181490005"
	//注册验证码，需要code参数
	SMSTemplateCodeRegisterAuthCode = "SMS_181500003"
	//重置密码验证码，需要code参数
	SMSTemplateCodeResetPwdAuthCode = "SMS_181500009"
	//更换组织超级管理员，需要code参数
	SMSTemplateCodeChangeSuperAdmin = "SMS_219749723"
)

//SMS参数名配置
const (
	//验证码code
	SMSParamsNameCode = "code"
	//验证动作
	SMSParamsNameAction = "action"
	//链接地址
	SMSParamsNameInviteUrl = "inviteUrl"
	//跳转地址
	SMSParamsNameInviteHref = "inviteHref"
)

//短信验证Action
const (
	SMSAuthCodeActionLogin       = "登录"
	SMSAuthCodeActionRegister    = "注册"
	SMSAuthCodeActionResetPwd    = "修改密码"
	SMSAuthCodeActionRetrievePwd = "找回密码"
	SMSAuthCodeActionBind        = "绑定账号"
	SMSAuthCodeActionUnBind      = "解绑账号"
)

//邮箱验证码模板
const (
	MailTemplateSubjectAuthCodeLogin       = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionLogin + "，请验证邮箱"
	MailTemplateSubjectAuthCodeRegister    = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionRegister + "，请验证邮箱"
	MailTemplateSubjectAuthCodeResetPwd    = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionResetPwd + "，请验证邮箱"
	MailTemplateSubjectAuthCodeRetrievePwd = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionRetrievePwd + "，请验证邮箱"
	MailTemplateSubjectAuthCodeBind        = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionBind + "，请验证邮箱"
	MailTemplateSubjectAuthCodeUnBind      = "欢迎使用极星协作，您正在进行" + SMSAuthCodeActionUnBind + "，请验证邮箱"
	MailTemplateContentAuthCode            = "您好，欢迎使用 极星，您正在进行邮箱{{.action}}验证，为了保护您的信息安全，我们来信进行邮箱验证，如果此操作不是由您发起的，请忽略此邮件。<h1>验证码：{{.code}}</h1>"
)

//邮箱邀请模板
const (
	MailTemplateSubjectInvite = "欢迎使用LessCode平台"
	MailTemplateContentInvite = "诚邀您加入LessCode平台<br/>请复制此链接后从浏览器中打开<b><a href = \"{{.inviteHref}}\">{{.inviteUrl}}</a></b><br/>为了防止您的敏感信息泄露，请勿将此邮件转发给他人"
)
const (
	//未读通知
	NoticeUnReadStatus = 1
	//已读通知
	NoticeReadStatus = 2
)

const (
	PermissionTypeSys = 1
	PermissionTypeOrg = 2
	PermissionTypePro = 3
)

//任务计划时间提醒类型
const (
	IssuePlanTimeRemindTypePlanEndTime   = 1
	IssuePlanTimeRemindTypePlanStartTime = 2
)

// context key
const (
	TraceIdKey    = "X-POLARIS-TRACE-ID"
	TraceIdLogKey = "pmTraceId"

	HttpContextKey = "_httpContext"
)

// 日志相关常量
const (
	LogAppKey       = "appName"
	LogTagKey       = "tag"
	LogMqMessageKey = "mqMessage"
)

const (
	// 模板预览需求-用于模板预览的组织的用户登录的 token
	PreviewTplToken         = "o999u2t3b6b7eaeda0f403998db1dc26be4dd1b"  // 该 token 用于预览模板
	PreviewTplTokenForWrite = "o999u10t344a04f0bd0b47e598c5012a52adfd9c" // 该模板用于创建项目、新建模板
	PreviewTplOrgId         = 999
	PreviewTplUserId        = 2
	PreviewOrWriteTplUserId = 10
)

const (
	MemberTypeUser = "U_"
	MemberTypeDept = "D_"
)
