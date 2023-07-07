package req

// 用户登录请求结构体
type UserLoginReq struct {
	// 登录类型: 1、短信验证码登录，2、账号密码登录，3、邮箱验证码登录（必填）
	LoginType int `json:"loginType"`
	// 登录类型为1时，loginName为手机号； 登录类型为3时，loginName为邮箱（必填）
	LoginName string `json:"loginName"`
	// 登录类型为2时，密码必传
	Password *string `json:"password"`
	// 验证码
	AuthCode *string `json:"authCode"`
	// 注册时可以带上名字
	Name *string `json:"name"`
	// 邀请码, 邀请注册时必填
	InviteCode *string `json:"inviteCode"`
	// 来源通道
	SourceChannel string `json:"sourceChannel"`
	// 平台
	SourcePlatform string `json:"sourcePlatform"`
}

// 用户注册请求结构体
type UserRegisterReq struct {
	// 手机号
	PhoneNumber string `json:"phoneNumber"`
	// 短信验证码
	AuthCode string `json:"authCode"`
	// 邮箱
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
	// 姓名（选填）
	Name string `json:"name"`
	// 邀请码, 邀请注册时必填
	InviteCode string `json:"inviteCode"`
	// 来源通道
	SourceChannel string `json:"sourceChannel"`
	// 平台
	SourcePlatform string `json:"sourcePlatform"`
}
