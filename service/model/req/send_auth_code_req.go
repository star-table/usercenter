package req

// 发送各种验证码请求结构体
type SendAuthCodeReq struct {
	// 验证方式: 1: 登录验证码，2：注册验证码，3：修改密码验证码，4：找回密码验证码，5：绑定验证码, 6：解绑验证码
	AuthType int `json:"authType"`
	// 地址类型: 1：手机号，2：邮箱
	AddressType int `json:"addressType"`
	// 联系地址，根据地址类型区分手机号或者邮箱
	Address string `json:"address"`
	// 验证码id
	CaptchaID *string `json:"captchaId"`
	// 输入的验证码
	CaptchaPassword *string `json:"captchaPassword"`
}

type AuthSmsCodeReq struct {
	// 手机号
	PhoneNumber string `json:"phoneNumber"`
	// 验证方式: 1: 登录验证码，2：注册验证码，3：修改密码验证码，4：找回密码验证码，5：绑定验证码, 6：解绑验证码，7更换管理员发送短信。
	AuthType int `json:"authType"`
	// 输入的验证码
	CaptchaPassword string `json:"captchaPassword"`
}

