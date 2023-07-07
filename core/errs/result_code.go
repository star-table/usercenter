package errs

var (
	//成功
	OK = AddResultCodeInfo(0, "OK", "ResultCode.OK")

	//token错误
	RequestError = AddResultCodeInfo(400, "请求错误", "ResultCode.RequestError")

	//认证错误
	Unauthorized = AddResultCodeInfo(401, "认证错误", "ResultCode.Unauthorized")

	//无权操作
	ForbiddenAccess = AddResultCodeInfo(403, "无权操作", "ResultCode.ForbiddenAccess")

	// 暂无查看权限
	ForbiddenView = AddResultCodeInfo(413, "暂无查看权限", "ResultCode.ForbiddenView")

	PolarisForbiddenAccess = AddResultCodeInfo(493, "无权操作.", "ResultCode.ForbiddenAccess")

	//请求地址不存在
	PathNotFound = AddResultCodeInfo(404, "请求地址不存在", "ResultCode.PathNotFound")

	//不支持该方法
	MethodNotAllowed = AddResultCodeInfo(405, "不支持该方法", "ResultCode.MethodNotAllowed")

	//Token过期
	TokenExpires = AddResultCodeInfo(450, "登录失效", "ResultCode.TokenExpires")

	//请求参数错误
	ServerError = AddResultCodeInfo(500, "服务器错误", "ResultCode.ServerError")

	//过载保护,服务暂不可用
	ServiceUnavailable = AddResultCodeInfo(503, "过载保护,服务暂不可用", "ResultCode.ServiceUnavailable")

	//服务调用超时
	Deadline = AddResultCodeInfo(504, "服务调用超时", "ResultCode.Deadline")

	//超出限制
	LimitExceed = AddResultCodeInfo(509, "超出限制", "ResultCode.LimitExceed")

	//参数错误
	ParamError = AddResultCodeInfo(600, "参数错误", "ResultCode.ParamError")

	//文件过大
	FileTooLarge = AddResultCodeInfo(610, "文件过大", "ResultCode.FileTooLarge")

	//文件类型错误
	FileTypeError = AddResultCodeInfo(611, "文件类型错误", "ResultCode.FileTypeError")

	//文件或目录不存在
	FileNotExist = AddResultCodeInfo(612, "文件或目录不存在", "ResultCode.FileNotExist")

	//文件路径为空
	FilePathIsNull = AddResultCodeInfo(613, "文件路径为空", "ResultCode.FilePathIsNull")

	//读取文件失败
	FileReadFail = AddResultCodeInfo(614, "读取文件失败", "ResultCode.FileReadFail")

	//错误未定义
	ErrorUndefined = AddResultCodeInfo(996, "错误未定义", "ResultCode.ErrorUndefined")

	//业务失败
	BusinessFail = AddResultCodeInfo(997, "业务失败", "ResultCode.BusinessFail")

	//系统异常
	SystemError = AddResultCodeInfo(998, "系统异常", "ResultCode.SystemError")

	RocketMQProduceInitError   = AddResultCodeInfo(200100, "RocketMQ Produce 初始化异常", "ResultCode.RocketMQProduceInitError")
	RocketMQSendMsgError       = AddResultCodeInfo(200101, "RocketMQ SendMsg 失败", "ResultCode.RocketMQSendMsgError")
	RocketMQConsumerInitError  = AddResultCodeInfo(200102, "RocketMQ Consumer 初始化异常", "ResultCode.RocketMQConsumerInitError")
	RocketMQConsumerStartError = AddResultCodeInfo(200103, "RocketMQ Consumer 启动异常", "ResultCode.RocketMQConsumerStartError")
	RocketMQConsumerStopError  = AddResultCodeInfo(200104, "RocketMQ Consumer 停止异常", "ResultCode.RocketMQConsumerStopError")

	KafkaMqSendMsgError           = AddResultCodeInfo(200200, "Kafka发送消息失败", "ResultCode.KafkaMqSendMsgError")
	KafkaMqSendMsgCantBeNullError = AddResultCodeInfo(200201, "Kafka发送的消息不能为空", "ResultCode.KafkaMqSendMsgCantBeNullError")
	KafkaMqConsumeMsgError        = AddResultCodeInfo(200202, "Kafka消费消息失败", "ResultCode.KafkaMqConsumeMsgError")
	KafkaMqConsumeStartError      = AddResultCodeInfo(200203, "Kafka消费启动失败", "ResultCode.KafkaMqConsumeStartError")

	TryDistributedLockError = AddResultCodeInfo(200300, "获取分布式锁异常", "ResultCode.TryDistributedLockError")
	GetDistributedLockError = AddResultCodeInfo(200301, "获取分布式锁失败", "ResultCode.GetDistributedLockError")
	MysqlOperateError       = AddResultCodeInfo(200302, "db操作出现异常", "ResultCode.MysqlOperateError")
	RedisOperateError       = AddResultCodeInfo(200303, "redis操作出现异常", "ResultCode.RedisOperateError")

	DbMQSendMsgError         = AddResultCodeInfo(200401, "Db 保存 message queue 失败", "ResultCode.DbMQSendMsgError")
	DbMQCreateConsumerError  = AddResultCodeInfo(200402, "Db 创建 message queue consumer 失败", "ResultCode.DbMQCreateConsumerError")
	DbMQConsumerStartedError = AddResultCodeInfo(200403, "Db 创建 message queue consumer 已启动", "ResultCode.DbMQConsumerStartedError")

	SystemBusy = AddResultCodeInfo(200500, "系统繁忙，请稍后重试", "ResultCode.SystemBusy")

	JSONConvertError = AddResultCodeInfo(200601, "Json转换出现异常", "ResultCode.JSONConvertError")
	ObjectCopyError  = AddResultCodeInfo(200602, "对象copy出现异常", "ResultCode.ObjectCopyError")
	CacheProxyError  = AddResultCodeInfo(200603, "缓存代理出现异常", "ResultCode.CacheProxyError")
	ObjectTypeError  = AddResultCodeInfo(200604, "对象类型错误", "ResultCode.ObjectTypeError")

	ApplyIdError        = AddResultCodeInfo(200701, "ID申请异常", "ResultCode.ApplyIdError")
	ApplyIdCountTooMany = AddResultCodeInfo(200702, "申请id数量过多", "ResultCode.ApplyIdCountTooMany")
	TypeConvertError    = AddResultCodeInfo(200703, "类型转换出现异常", "ResultCode.TypeConvertError")
	UpdateFiledIsEmpty  = AddResultCodeInfo(200704, "未更新任何信息", "ResultCode.UpdateFiledIsEmpty")

	TokenAuthError       = AddResultCodeInfo(200801, "身份认证异常，请重新登录", "ResultCode.TokenAuthError")
	TokenNotExist        = AddResultCodeInfo(200802, "身份认证失败，请重新登录", "ResultCode.TokenNotExist")
	SuiteTicketError     = AddResultCodeInfo(200803, "获取SuiteTicket异常", "ResultCode.SuiteTicketError")
	GetContextError      = AddResultCodeInfo(200804, "获取请求上下文异常", "ResultCode.GetContextError")
	TemplateRenderError  = AddResultCodeInfo(200805, "模板解析失败", "ResultCode.TemplateRenderError")
	DecryptError         = AddResultCodeInfo(200806, "参数解密异常", "ResultCode.DecryptError")
	CaptchaError         = AddResultCodeInfo(200807, "验证码错误", "ResultCode.CaptchaError")
	DingCodeCacheInvalid = AddResultCodeInfo(200808, "扫码认证已失效，请重新扫码", "ResultCode.DingCodeCacheInvalid")

	MQTTKeyGenError        = AddResultCodeInfo(200901, "生成key发生异常", "ResultCode.MQTTKeyGenError")
	MQTTPublishError       = AddResultCodeInfo(200902, "MQTT推送消息发生异常", "ResultCode.MQTTPublishError")
	MQTTConnectError       = AddResultCodeInfo(200903, "MQTT连接时发生异常", "ResultCode.MQTTConnectError")
	MQTTMissingConfigError = AddResultCodeInfo(200904, "MQTT缺少配置", "ResultCode.MQTTMissingConfigError")

	OssConfigError     = AddResultCodeInfo(200934, "文件存储缺少配置", "ResultCode.OssConfigError")
	OssPolicyTypeError = AddResultCodeInfo(200935, "错误的策略类型", "ResultCode.OssPolicyTypeError")

	//>>>业务异常
	InitDbFail                     = AddResultCodeInfo(201000, "初始化db失败", "ResultCode.InitDbFail")
	ObjectRecordNotFoundError      = AddResultCodeInfo(201001, "对象记录不存在", "ResultCode.ObjectRecordNotFoundError")
	DingTalkUserInfoNotInitedError = AddResultCodeInfo(201002, "钉钉用户没有初始化", "ResultCode.DingTalkUserInfoNotInitedError")
	CacheUserInfoNotExistError     = AddResultCodeInfo(201004, "令牌对应的用户信息不存在", "ResultCode.CacheUserInfoNotExistError")

	PageSizeOverflowMaxSizeError = AddResultCodeInfo(201101, "请求页长超出最大页长限制", "ResultCode.PageSizeOverflowMaxSizeError")
	OutOfConditionError          = AddResultCodeInfo(201102, "请求条件超出限制", "ResultCode.OutOfConditionError")
	ConditionHandleError         = AddResultCodeInfo(201103, "条件处理异常", "ResultCode.ConditionHandleError")
	ReqParamsValidateError       = AddResultCodeInfo(201104, "请求参数校验异常", "ResultCode.ReqParamsValidateError")

	IssueCommentLenError = AddResultCodeInfo(201202, "评论不得为空且不能超过200字", "Result.IssueCommentLenError")
	IssueRemarkLenError  = AddResultCodeInfo(201203, "描述不能超过10000字", "Result.IssueRemarkLenError")
	IssueTitleError      = AddResultCodeInfo(201204, "任务标题包含非法字符或超出500个字符", "Result.IssueTitleError")

	OrgNotInitError           = AddResultCodeInfo(201300, "组织未初始化", "ResultCode.OrgNotInitError")
	UserConfigNotExist        = AddResultCodeInfo(201301, "用户配置不存在", "Result.UserConfigNotExist")
	OrgNotExist               = AddResultCodeInfo(201302, "组织不存在", "ResultCode.OrgNotExist")
	OrgInitError              = AddResultCodeInfo(201303, "组织初始化异常", "ResultCode.OrgInitError")
	OrgOwnTransferError       = AddResultCodeInfo(201304, "非组织创建者不能更改信息", "ResultCode.OrgOwnTransferError")
	OrgOutInfoNotExist        = AddResultCodeInfo(201305, "组织外部信息不存在", "ResultCode.OrgOutInfoNotExist")
	UserOutInfoNotExist       = AddResultCodeInfo(201306, "用户外部信息不存在", "ResultCode.UserOutInfoNotExist")
	UserOutInfoNotError       = AddResultCodeInfo(201307, "用户外部信息错误", "ResultCode.UserOutInfoNotError")
	OrgCodeAlreadySetError    = AddResultCodeInfo(201308, "组织网址不能二次修改", "ResultCode.OrgCodeAlreadySetError")
	OrgCodeLenError           = AddResultCodeInfo(201310, "组织网址后缀只能输入20个字符,包含数字和英文", "ResultCode.OrgWebSiteSettingLenError")
	OrgCodeExistError         = AddResultCodeInfo(201311, "组织网址后缀已被占用，请重新输入", "ResultCode.OrgCodeExistError")
	OrgAddressLenError        = AddResultCodeInfo(201313, "详情地址不得超过100字", "ResultCode.OrgAddressLenError")
	OrgLogoLenError           = AddResultCodeInfo(201314, "组织logo路径长度不能超过512字", "ResultCode.OrgLogoLenError")
	OrgUserUnabled            = AddResultCodeInfo(201317, "您已被当前组织禁止访问，请联系管理员解除限制", "ResultCode.OrgUserUnabled")
	OrgUserDeleted            = AddResultCodeInfo(201319, "您已被该组织移除", "ResultCode.OrgUserDeleted")
	OrgUserCheckStatusUnabled = AddResultCodeInfo(201320, "您未通过该组织的审核", "ResultCode.OrgUserCheckStatusUnabled")
	OrgFunctionInvalid        = AddResultCodeInfo(201321, "存在无效功能，请确认功能项", "ResultCode.OrgFunctionInvalid")
	ExportFieldIsNull         = AddResultCodeInfo(201322, "请选择导出字段", "ResultCode.ExportFieldIsNull")
	UserConfigUpdateError     = AddResultCodeInfo(201323, "用户设置更新失败", "Result.UserConfigUpdateError")
	OrgNameLenError           = AddResultCodeInfo(201324, "组织名称包含非法字符或长度超出20个字符", "Result.OrgNameLenError")

	GetUserInfoError      = AddResultCodeInfo(201400, "获取用户信息异常", "Result.GetUserInfoError")
	TargetNotExist        = AddResultCodeInfo(201401, "操作对象不存在", "Result.TargetNotExist")
	InviteCodeInvalid     = AddResultCodeInfo(201402, "邀请链接失效", "Result.InviteCodeInvalid")
	UnSupportLoginType    = AddResultCodeInfo(201403, "不支持的登录方式", "Result.UnSupportLoginType")
	PasswordEmptyError    = AddResultCodeInfo(201404, "请输入密码", "Result.PasswordEmptyError")
	PasswordNotSetError   = AddResultCodeInfo(201405, "密码未设置", "Result.PasswordNotSetError")
	PasswordNotMatchError = AddResultCodeInfo(201406, "密码验证错误", "Result.PasswordNotMatchError")
	PasswordEditTimeError = AddResultCodeInfo(201407, "三个月内没有修改过密码，请修改密码", "Result.PasswordEditTimeError")

	DepartmentNotExist        = AddResultCodeInfo(201501, "部门不存在", "Result.DepartmentNotExist")
	ParentDeptNotExist        = AddResultCodeInfo(201502, "父部门不存在", "Result.ParentDeptNotExist")
	RootDeptNotExistErr       = AddResultCodeInfo(201503, "根部门不存在", "Result.RootDeptNotExistErr")
	DepartmentNameInvalid     = AddResultCodeInfo(201506, "部门名称应为1~20个字符", "Result.DepartmentNameInvalid")
	DeptNameConflictErr       = AddResultCodeInfo(201507, "部门名称已使用", "Result.DeptNameConflictErr")
	DeptHaveSubDeptError      = AddResultCodeInfo(201508, "部门下存在有效子部门", "Result.DeptHaveSubDeptError")
	RootDeptCannotDeleteError = AddResultCodeInfo(201509, "根部门不可删除", "Result.RootDeptCannotDeleteError")
	UserIsLeaderAndUserDeny   = AddResultCodeInfo(201510, "抱歉，有用户被设置成既是主管又是普通成员。", "Result.UserIsLeaderAndUserDeny")
	Top2DeptMustHaveOne       = AddResultCodeInfo(201511, "至少有一个顶级部门。", "Result.Top2DeptMustHaveOne")

	UserOrgNotRelation          = AddResultCodeInfo(201600, "用户不是该组织成员", "Result.UserOrgNotRelation")
	UserDisabledError           = AddResultCodeInfo(201601, "已经被组织禁用", "Result.UserDisabledError")
	UpdateMemberIdsIsEmptyError = AddResultCodeInfo(201602, "变动的成员列表为空", "Result.UpdateMemberIdsIsEmptyError")
	UpdateMemberStatusFail      = AddResultCodeInfo(201603, "修改成员状态失败", "Result.UpdateMemberStatusFail")
	CannotChangeOwnerStatus     = AddResultCodeInfo(201604, "无权变更组织拥有者的状态", "Result.CannotChangeOwnerStatus")
	CannotChangeAdminStatus     = AddResultCodeInfo(201605, "无权变更系统管理员的状态", "Result.CannotChangeAdminStatus")
	CannotChangeSelfStatus      = AddResultCodeInfo(201606, "不允许变更自己的状态", "Result.CannotChangeSelfStatus")
	CannotEditOrgOwnerMainField = AddResultCodeInfo(201607, "无权变更组织拥有者的关键信息", "Result.CannotEditOrgOwnerMainField")
	CannotEditAdminMainField    = AddResultCodeInfo(201608, "无权变更系统管理员的关键信息", "Result.CannotEditAdminMainField")
	NoAdminInAdminGroup         = AddResultCodeInfo(201609, "该管理组没有管理员", "Result.NoAdminInAdminGroup")
	UsersNotInDept              = AddResultCodeInfo(201610, "这些成员属于当前部门下的下级部门，请在其所属部门移除。", "Result.UsersNotInDept")
	SysAdminGroupMustHasMember  = AddResultCodeInfo(201611, "系统角色至少有一名成员", "Result.SysAdminGroupMustHasMember")
	UserMustHasAdminGroup       = AddResultCodeInfo(201612, "用户至少需要关联一个角色", "Result.UserMustHasAdminGroup")
	CannotRemoveUser            = AddResultCodeInfo(201613, "不支持移除超管和自己，请重新选择成员", "Result.CannotRemoveUser")

	SourceChannelNotDefinedError = AddResultCodeInfo(201700, "来源通道未定义", "Result.SourceChannelNotDefinedError")
	OrgNotNeedInitError          = AddResultCodeInfo(201701, "组织已存在，不需要初始化", "Result.OrgNotNeedInitError")
	AuthCodeIsNull               = AddResultCodeInfo(201704, "验证码不得为空", "Result.AuthCodeIsNull")

	OrgConfigNotExist = AddResultCodeInfo(201802, "组织配置不存在", "Result.OrgConfigNotExist")
	IndustryNotExist  = AddResultCodeInfo(201803, "所选行业不存在", "Result.IndustryNotExist")

	SetUserPasswordError      = AddResultCodeInfo(201901, "设置密码失败", "Result.SetUserPasswordError")
	UnBindLoginNameFail       = AddResultCodeInfo(201902, "解绑登录方式失败", "Result.UnBindLoginNameFail")
	BindLoginNameFail         = AddResultCodeInfo(201903, "绑定登录方式失败", "Result.BindLoginNameFail")
	NotBindAccountError       = AddResultCodeInfo(201904, "该登录方式未绑定任何账号", "Result.NotBindAccountError")
	AccountAlreadyBindError   = AddResultCodeInfo(201905, "该登录方式已绑定其它账号", "Result.AccountAlreadyBindError")
	EmailNotBindAccountError  = AddResultCodeInfo(201906, "该邮箱未绑定任何账号，请重新输入或使用手机验证码登录", "Result.EmailNotBindAccountError")
	MobileNotBindAccountError = AddResultCodeInfo(201907, "该手机号未绑定任何账号", "Result.MobileNotBindAccountError")
	MobileAlreadyRegister     = AddResultCodeInfo(201908, "此手机号码已注册，请返回登录页进行登录", "Result.MobileAlreadyRegister")
	MobileAlreadyBind         = AddResultCodeInfo(201909, "此手机号码已绑定其它账号", "Result.MobileAlreadyBind")
	EmailAlreadyRegister      = AddResultCodeInfo(201910, "该邮箱已绑定手机号", "Result.EmailAlreadyRegister")
	AccountRegisterConflict   = AddResultCodeInfo(201911, "此用户名已注册，请返回登录页进行登录", "Result.AccountRegisterConflict")
	PhoneNumberFormatError    = AddResultCodeInfo(201912, "手机号格式错误，请重新输入", "ResultCode.PhoneNumberFormatError")

	//User
	UserInitError                       = AddResultCodeInfo(202000, "用户初始化失败", "Result.UserInitError")
	UserNotInitError                    = AddResultCodeInfo(202002, "用户未初始化", "Result.UserNotInitError")
	UserNotExist                        = AddResultCodeInfo(202003, "用户不存在", "Result.UserNotExist")
	UserRegisterError                   = AddResultCodeInfo(202005, "用户注册失败", "Result.UserRegisterError")
	LarkInitError                       = AddResultCodeInfo(202006, "示例数据已初始化", "Result.LarkInitError")
	UserSexFail                         = AddResultCodeInfo(202007, "用户性别错误", "Result.UserSexFail")
	UserNameEmpty                       = AddResultCodeInfo(202008, "用户姓名不能为空串", "Result.UserNameEmpty")
	EmailNotRegisterError               = AddResultCodeInfo(202009, "当前邮箱未注册", "Result.EmailNotRegisterError")
	EmailNotBindError                   = AddResultCodeInfo(202010, "邮箱未绑定", "Result.EmailNotBindError")
	MobileNotBindError                  = AddResultCodeInfo(202011, "手机号未绑定", "Result.MobileNotBindError")
	EmailAlreadyBindError               = AddResultCodeInfo(202012, "邮箱已绑定, 请先解绑", "Result.EmailAlreadyBindError")
	MobileAlreadyBindError              = AddResultCodeInfo(202013, "手机号已绑定， 请先解绑", "Result.MobileAlreadyBindError")
	EmailAlreadyBindByOtherAccountError = AddResultCodeInfo(202014, "该邮箱已被其他账号绑定", "Result.EmailAlreadyBindByOtherAccountError")
	MobileAlreadyBindOtherAccountError  = AddResultCodeInfo(202015, "该手机号已被其他账号绑定", "Result.MobileAlreadyBindOtherAccountError")
	AccountNotRegister                  = AddResultCodeInfo(202016, "用户未注册，请先注册或填写其他账号", "Result.AccountNotRegister")
	OrgMemberNotExistOrDisable          = AddResultCodeInfo(202017, "组织成员不存在或已禁用", "Result.OrgMemberNotExistOrDisable")
	OrgMemberNotExist                   = AddResultCodeInfo(202018, "组织成员不存在", "Result.OrgMemberNotExist")
	UserRebindCodeError                 = AddResultCodeInfo(202019, "抱歉，登录名更换绑定时，验证码错误。", "Result.UserRebindCodeError")
	OrgEmpNoConflictErr                 = AddResultCodeInfo(202020, "该工号已被其他成员使用", "Result.OrgUserEmpNoConflictErr")
	NicknameLenError                    = AddResultCodeInfo(202021, "昵称不超过20字符", "Result.NicknameLenError")
	UsernameLenError                    = AddResultCodeInfo(202022, "用户名应为2-20位的字母、数字组合", "Result.UsernameLenError")
	PwdAlreadySettingsErr               = AddResultCodeInfo(202023, "密码已设置过", "Result.PwdAlreadySettingsErr")
	PwdLengthError                      = AddResultCodeInfo(202024, "密码长度为8~16位，并为数字、字母、符号的组合形式", "Result.PwdLengthError")

	//domain
	ProjectDomainError    = AddResultCodeInfo(202100, "项目领域出错", "Result.ProjectDomainError")
	IssueDomainError      = AddResultCodeInfo(202101, "任务领域出错", "Result.IssueDomainError")
	UserDomainError       = AddResultCodeInfo(202102, "用户领域出错", "Result.UserDomainError")
	BaseDomainError       = AddResultCodeInfo(202103, "领域出错", "Result.BaseDomainError")
	TrendDomainError      = AddResultCodeInfo(202104, "动态领域出错", "Result.TrendDomainError")
	IterationDomainError  = AddResultCodeInfo(202105, "迭代领域出错", "Result.IterationDomainError")
	ObjectTypeDomainError = AddResultCodeInfo(202106, "对象类型领域出错", "Result.ObjectTypeDomainError")
	ResourceDomainError   = AddResultCodeInfo(202107, "资源领域出错", "Result.ResourceDomainError")
	ProcessDomainError    = AddResultCodeInfo(202108, "流程领域出错", "Result.ProcessDomainError")
	DepartmentDomainError = AddResultCodeInfo(202109, "部门领域出错", "Result.DepartmentDomainError")

	//权限验证领域
	IllegalityRoleOperation         = AddResultCodeInfo(202200, "非法的操作code", "Result.IllegalityRoleOperation")
	UserRoleNotDefinition           = AddResultCodeInfo(202201, "用户角色未定义", "Result.UserRoleNotDefinition")
	NoOperationPermissions          = AddResultCodeInfo(202202, "没有操作权限", "Result.NoOperationPermissions")
	PermissionNotExist              = AddResultCodeInfo(202203, "权限项不存在", "Result.PermissionNotExist")
	NoOperationPermissionForProject = AddResultCodeInfo(202204, "暂无权限操作，请联系项目负责人", "Result.NoOperationPermissionForProject")
	NoOperationPermissionForIssue   = AddResultCodeInfo(202205, "暂无权限操作，请联系任务负责人", "Result.NoOperationPermissionForIssue")

	//>>>dingtalk open_api api error
	SuiteTicketNotExistError      = AddResultCodeInfo(202301, "suiteTicket失效或不存在", "ResultCode.SuiteTicketNotExistError")
	DingTalkOpenApiCallError      = AddResultCodeInfo(202302, "钉钉OpenApi调用异常", "ResultCode.DingTalkOpenApiCallError")
	DingTalkAvoidCodeInvalidError = AddResultCodeInfo(202303, "钉钉免登code失效", "ResultCode.DingTalkAvoidCodeInvalidError")
	DingTalkClientError           = AddResultCodeInfo(202304, "钉钉Client获取失败", "ResultCode.DingTalkClientError")
	DingTalkGetUserInfoError      = AddResultCodeInfo(202305, "钉钉获取用户信息失败", "ResultCode.DingTalkGetUserInfoError")
	DingTalkOrgInitError          = AddResultCodeInfo(202306, "钉钉企业初始化失败", "ResultCode.DingTalkOrgInitError")
	DingTalkConfigError           = AddResultCodeInfo(202307, "钉钉配置错误", "ResultCode.DingTalkConfigError")

	//>>> 飞书 open_api api err
	FeiShuOpenApiCallError                = AddResultCodeInfo(202400, "飞书OpenApi调用异常", "ResultCode.FeiShuOpenApiCallError")
	FeiShuAppTicketNotExistError          = AddResultCodeInfo(202401, "AppTicket不存在", "ResultCode.FeiShuAppTicketNotExistError")
	FeiShuConfigNotExistError             = AddResultCodeInfo(202402, "飞书配置不存在", "ResultCode.FeiShuConfigNotExistError")
	FeiShuClientTenantError               = AddResultCodeInfo(202403, "飞书客户端获取失败", "ResultCode.FeiShuClientTenantError")
	FeiShuGetAppAccessTokenError          = AddResultCodeInfo(202404, "飞书获取AppAccessToken失败", "ResultCode.FeiShuGetAppAccessTokenError")
	FeiShuGetTenantAccessTokenError       = AddResultCodeInfo(202405, "飞书获取TenantAccessToken失败", "ResultCode.FeiShuGetTenantAccessTokenError")
	FeiShuAuthCodeInvalid                 = AddResultCodeInfo(202406, "飞书用户授权失败", "ResultCode.FeiShuAuthCodeInvalid")
	FeiShuCardCallSignVerifyError         = AddResultCodeInfo(202407, "飞书卡片回调签名校验失败", "ResultCode.FeiShuCardCallSigVerifyError")
	FeiShuCardCallMsgRepetError           = AddResultCodeInfo(202408, "飞书卡片消息重复推送", "ResultCode.FeiShuCardCallMsgRepetError")
	FeiShuUserNotInAppUseScopeOfAuthority = AddResultCodeInfo(202409, "不在通讯录权限范围内，请联系组织管理员!", "ResultCode.FeiShuUserNotInAppUseScopeOfAuthority")

	//Login Error
	SMSLoginCodeSendError                = AddResultCodeInfo(202500, "登录验证码发送失败", "ResultCode.SMSLoginCodeSendError")
	SMSSendLimitError                    = AddResultCodeInfo(202502, "发送过于频繁（服务商）", "ResultCode.SMSSendLimitError")
	SMSSendTimeLimitError                = AddResultCodeInfo(202503, "发送过于频繁", "ResultCode.SMSSendTimeLimitError")
	SMSLoginCodeInvalid                  = AddResultCodeInfo(202504, "验证码已失效，请重新获取", "ResultCode.SMSLoginCodeInvalid")
	SMSLoginCodeNotMatch                 = AddResultCodeInfo(202505, "验证码错误，请重新获取", "ResultCode.SMSLoginCodeNotMatch")
	SMSLoginCodeVerifyFailTimesOverLimit = AddResultCodeInfo(202506, "验证码错误，失败次数过多，请重新发送", "ResultCode.SMSLoginCodeVerifyFailTimesOverLimit")
	PwdLoginCodeNotMatch                 = AddResultCodeInfo(202507, "图形验证码错误", "ResultCode.PwdLoginCodeNotMatch")
	PwdLoginUsrOrPwdNotMatch             = AddResultCodeInfo(202508, "用户名或密码错误", "ResultCode.PwdLoginUsrOrPwdNotMatch")
	ChangeLoginNameInvalid               = AddResultCodeInfo(202509, "换绑操作已过期，请重新进行换绑操作", "ResultCode.ChangeLoginNameInvalid")
	AccountNotSetPwdErr                  = AddResultCodeInfo(202510, "不支持该登陆方式，请使用其他方式登陆", "ResultCode.AccountNotSetPwdErr")

	EmailFormatErr       = AddResultCodeInfo(202600, "邮箱格式错误", "ResultCode.EmailFormatErr")
	EmailSubjectEmptyErr = AddResultCodeInfo(202601, "邮箱标题不能为空", "ResultCode.EmailSubjectEmptyErr")
	EmailSendErr         = AddResultCodeInfo(202602, "邮件发送失败", "ResultCode.EmailSendErr")

	NotSupportedContactAddressType = AddResultCodeInfo(202700, "不支持的联系方式类型", "ResultCode.NotSupportedContactAddressType")
	NotSupportedAuthCodeType       = AddResultCodeInfo(202701, "不支持的验证码类型", "ResultCode.NotSupportedAuthCodeType")
	NotSupportedRegisterType       = AddResultCodeInfo(202702, "暂时不支持该注册方式", "ResultCode.NotSupportedRegisterType")
	HaveNoContract                 = AddResultCodeInfo(202703, "请保证至少一种联系方式", "ResultCode.HaveNoContract")

	RoleNotExist             = AddResultCodeInfo(202800, "角色不存在", "Result.RoleNotExist")
	RoleNameLenErr           = AddResultCodeInfo(202801, "角色名包含非法字符或长度超出10个字符", "Result.RoleNameLenErr")
	DefaultRoleCantModify    = AddResultCodeInfo(202802, "默认角色不允许编辑", "Result.DefaultRoleCantModify")
	RoleModifyBusy           = AddResultCodeInfo(202803, "角色更新繁忙", "Result.RoleEditBusy")
	RoleNameRepeatErr        = AddResultCodeInfo(202804, "角色名称重复", "Result.RoleNameRepeatErr")
	CannotRemoveProjectOwner = AddResultCodeInfo(202805, "项目负责人不能被移除", "Result.CannotRemoveProjectOwner")
	DefaultRoleNameErr       = AddResultCodeInfo(202806, "与系统角色名称冲突", "Result.DefaultRoleNameErr")

	RoleGroupNameLenErr    = AddResultCodeInfo(202900, "角色组名称包含非法字符或长度超出10个字符", "Result.RoleGroupNameLenErr")
	RoleGroupModifyBusy    = AddResultCodeInfo(202901, "角色组更新繁忙", "Result.RoleGroupModifyBusy")
	RoleGroupNameRepeatErr = AddResultCodeInfo(202902, "角色组名称重复", "Result.RoleGroupNameRepeatErr")
	RoleGroupHaveRoleErr   = AddResultCodeInfo(202903, "角色组下存在有效角色", "Result.RoleGroupHaveRoleErr")
	OrgRoleGroupNotExist   = AddResultCodeInfo(202904, "角色组不存在", "ResultCode.OrgRoleGroupNotExist")

	RoleUserRefModifyBusy = AddResultCodeInfo(203000, "用户角色更新繁忙", "Result.RoleUserRefModifyBusy")
	RoleUsersIsEmpty      = AddResultCodeInfo(203001, "用户列表为空", "Result.RoleUsersIsEmpty")

	ManageGroupNotExist               = AddResultCodeInfo(203100, "管理组不存在或已删除", "Result.ManageGroupNotExist")
	CannotRemoveDefaultManageGroupErr = AddResultCodeInfo(203101, "系统管理组不可删除", "Result.CannotRemoveDefaultManageGroupErr")
	DefaultManageGroupCantModify      = AddResultCodeInfo(203102, "系统管理组不可编辑", "Result.DefaultManageGroupCantModify")
	DefaultManageGroupErr             = AddResultCodeInfo(203103, "与默认管理组名称冲突", "Result.DefaultManageGroupErr")
	ManageGroupRefModifyBusy          = AddResultCodeInfo(203104, "管理组更新繁忙", "Result.ManageGroupRefModifyBusy")
	ManageGroupNameRepeatErr          = AddResultCodeInfo(203105, "管理组名称重复", "Result.ManageGroupNameRepeatErr")
	ManageGroupNameLenErr             = AddResultCodeInfo(203106, "角色名包含非法字符或长度超出10个字符", "Result.ManageGroupNameLenErr")
	SysManageGroupRepeatErr           = AddResultCodeInfo(203107, "系统管理组已存在", "Result.SysManageGroupRepeatErr")
	ManageGroupOptionsNotExist        = AddResultCodeInfo(203108, "修改项非法", "Result.ManageGroupOptionsNotExist")
	ManageGroupMemberConflict         = AddResultCodeInfo(203109, "成员已存在于其他管理组", "Result.ManageGroupMemberConflict")
	CannotDeleteSelf                  = AddResultCodeInfo(203110, "不可移除自身", "Result.CannotDeleteSelf")
	ManageGroupMemberCountLimitErr    = AddResultCodeInfo(203111, "只可设置一个子管理员", "Result.ManageGroupMemberCountLimitErr")
	DenyChangeSysAdminGroupOfUser     = AddResultCodeInfo(203112, "不能更改组织超级管理员的管理组。", "Result.DenyChangeSysAdminGroupOfUser")

	AppPackageNotExist = AddResultCodeInfo(203200, "应用包不存在或已删除", "Result.AppPackageNotExist")
	AppNotExist        = AddResultCodeInfo(203201, "应用不存在或已删除", "Result.AppNotExist")

	// OpenAPI
	ApiKeyAuthErr  = AddResultCodeInfo(204000, "OpenAPI签名校验失败", "ResultCode.ApiKeyAuthErr")
	ApiKeyIsOpened = AddResultCodeInfo(204001, "OpenAPI已开启，请勿重复操作", "ResultCode.ApiKeyIsOpened")
	ApiKeyIsClosed = AddResultCodeInfo(204002, "OpenAPI未开启", "ResultCode.ApiKeyClosed")

	// 职级
	PositionCreateErr              = AddResultCodeInfo(205000, "职级创建失败", "ResultCode.PositionCreateErr")
	PositionUpdateErr              = AddResultCodeInfo(205001, "职级修改失败", "ResultCode.PositionUpdateErr")
	PositionDeleteErr              = AddResultCodeInfo(205002, "职级删除失败", "ResultCode.PositionDeleteErr")
	PositionIsDisabledErr          = AddResultCodeInfo(205003, "职级已停用", "ResultCode.PositionIsDisabledErr")
	PositionNotExistErr            = AddResultCodeInfo(205004, "职级不存在或已删除", "ResultCode.PositionNotExistErr")
	PositionHaveUserExistErr       = AddResultCodeInfo(205005, "仍有成员时，不可操作", "ResultCode.PositionHaveUserExistErr")
	PositionStatusUpdateErr        = AddResultCodeInfo(205006, "修改职级状态失败", "ResultCode.PositionStatusUpdateErr")
	PositionModifyBusy             = AddResultCodeInfo(205007, "职级修改繁忙，请稍后重试", "ResultCode.PositionModifyBusy")
	PositionNameLenErr             = AddResultCodeInfo(205008, "职级名称包含非法字符或长度超出10个字符", "ResultCode.PositionNameLenErr")
	PositionLevelErr               = AddResultCodeInfo(205009, "职级等级最小1级，最大20级", "ResultCode.PositionLevelErr")
	DefaultPositionCannotDeleteErr = AddResultCodeInfo(205010, "默认职级不可删除", "ResultCode.PositionLevelErr")
	DefaultPositionCannotEditErr   = AddResultCodeInfo(205011, "默认职级不可修改", "ResultCode.DefaultPositionCannotEditErr")

	SynAddUserToCloudreveErr    = AddResultCodeInfo(205100, "新增用户同步云盘失败", "ResultCode.SynAddUserToCloudreveErr")
	SynEditUserToCloudreveErr   = AddResultCodeInfo(205101, "修改用户同步云盘失败", "ResultCode.SynEditUserToCloudreveErr")
	SynRemoveUserToCloudreveErr = AddResultCodeInfo(205102, "删除用户同步云盘失败", "ResultCode.SynRemoveUserToCloudreveErr")

	OrgUserUnabledErr            = AddResultCodeInfo(400024, "您已被当前组织禁止访问，请联系管理员解除限制", "ResultCode.OrgUserUnabled")
	OrgUserDeletedErr            = AddResultCodeInfo(400026, "您已被该组织移除", "ResultCode.OrgUserDeleted")
	OrgUserCheckStatusUnabledErr = AddResultCodeInfo(400027, "您未通过该组织的审核", "ResultCode.OrgUserCheckStatusUnabled")
	OrgUserInvalidErr            = AddResultCodeInfo(400034, "抱歉，您不在应用付费授权范围内，请联系管理员开通", "ResultCode.OrgUserInvalid")
)
